package command

import (
	"errors"
	"fmt"
	"path"
	"sync/atomic"
	"time"

	"github.com/hanzoai/s3/s3/filer/archive"
	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/replication/repl_util"
	"github.com/hanzoai/s3/s3/replication/source"
	"github.com/hanzoai/s3/s3/security"
	"github.com/hanzoai/s3/s3/util"
)

type FilerArchiveCLIOptions struct {
	filer         *string
	bucket        *string
	remoteBucket  *string
	remotePrefix  *string
	schedule      *string
	retentionDays *int
	once          *bool
	proxyByFiler  *bool
}

var filerArchiveOptions FilerArchiveCLIOptions

func init() {
	cmdFilerArchive.Run = runFilerArchive
	filerArchiveOptions.filer = cmdFilerArchive.Flag.String("filer", "localhost:8888", "filer of one Hanzo S3 cluster")
	filerArchiveOptions.bucket = cmdFilerArchive.Flag.String("bucket", "", "source S3 bucket name under /buckets/<bucket>")
	filerArchiveOptions.remoteBucket = cmdFilerArchive.Flag.String("remoteBucket", "", "destination remote S3 archive bucket")
	filerArchiveOptions.remotePrefix = cmdFilerArchive.Flag.String("remotePrefix", "", "optional key prefix inside destination remote bucket")
	filerArchiveOptions.schedule = cmdFilerArchive.Flag.String("schedule", "daily", "per-bucket archive schedule: off, hourly, daily, or weekly")
	filerArchiveOptions.retentionDays = cmdFilerArchive.Flag.Int("retentionDays", 30, "archive log/data partition retention in days (0 disables pruning)")
	filerArchiveOptions.once = cmdFilerArchive.Flag.Bool("once", false, "run a single bounded incremental pass from the stored metadata offset and exit")
	filerArchiveOptions.proxyByFiler = cmdFilerArchive.Flag.Bool("filerProxy", false, "read file chunks via filer instead of volume servers")
}

var cmdFilerArchive = &Command{
	UsageLine: "filer.archive -filer=<filerHost>:<filerPort> -bucket=<bucket> -remoteBucket=<remoteBucket> -schedule=daily",
	Short:     "scheduled per-bucket incremental archive to a configured remote S3 bucket with per-bucket byte/request metering",
	Long: `Run scheduled incremental backup passes (off, hourly, daily, weekly) for a bucket under /buckets/<bucket>
from its stored filer metadata offset into a configured remote S3 sink defined in replication.toml.

Each pass advances the per-bucket checkpoint only after all events up to the pass watermark are archived,
meters archived bytes and remote requests per bucket, and enforces partition retention.`,
}

func runFilerArchive(cmd *Command, args []string) bool {
	sched, err := archive.ParseSchedule(*filerArchiveOptions.schedule)
	if err != nil {
		glog.Fatalf("filer.archive: %v", err)
	}
	pol := archive.BucketPolicy{
		Bucket:        *filerArchiveOptions.bucket,
		RemoteBucket:  *filerArchiveOptions.remoteBucket,
		RemotePrefix:  *filerArchiveOptions.remotePrefix,
		Schedule:      sched,
		RetentionDays: *filerArchiveOptions.retentionDays,
		PropagateDel:  true,
	}
	if err := pol.Validate(); err != nil {
		glog.Fatalf("filer.archive invalid policy: %v", err)
	}
	if sched == archive.ScheduleOff && !*filerArchiveOptions.once {
		glog.V(0).Infof("filer.archive: schedule is off for bucket %s; exiting", pol.Bucket)
		return true
	}

	util.LoadSecurityConfiguration()
	util.LoadConfiguration("replication", true)
	grpcDialOption := security.LoadClientTLS(util.GetViper(), "grpc.client")

	for {
		if err := runLiveFilerArchivePass(grpcDialOption, pol, *filerArchiveOptions.filer, *filerArchiveOptions.proxyByFiler); err != nil {
			glog.Errorf("filer.archive pass for bucket %s failed: %v", pol.Bucket, err)
		}
		if *filerArchiveOptions.once || sched == archive.ScheduleOff {
			return true
		}
		time.Sleep(sched.Interval())
	}
}

func runLiveFilerArchivePass(grpcDialOption pb.DialOption, pol archive.BucketPolicy, filerAddr string, proxyByFiler bool) error {
	dataSink := findSink(util.GetViper())
	if dataSink == nil {
		return errors.New("no data sink configured in replication.toml")
	}
	sourceFiler := pb.ServerAddress(filerAddr)
	sourcePath := pol.FilerBucketPath()
	targetPath := path.Join(dataSink.GetSinkToDirectory(), pol.Bucket)

	sinkID := pol.SinkCheckpointID()
	lastOffsetTsNs, err := getOffset(grpcDialOption, sourceFiler, archive.KeyPrefix, sinkID)
	if err != nil {
		lastOffsetTsNs = 0
	}
	stopTsNs := time.Now().UnixNano()
	if stopTsNs <= lastOffsetTsNs {
		return nil
	}

	filerSource := &source.FilerSource{}
	if err := filerSource.DoInitialize(
		sourceFiler.ToHttpAddress(),
		sourceFiler.ToGrpcAddress(),
		sourcePath,
		proxyByFiler,
	); err != nil {
		return fmt.Errorf("filer source init: %w", err)
	}
	if err := repl_util.InitializeSSEForReplication(filerSource); err != nil {
		return fmt.Errorf("SSE init: %w", err)
	}
	dataSink.SetSourceFiler(filerSource)

	var archivedBytes atomic.Int64
	var archivedReqs atomic.Int64
	var maxObservedTsNs atomic.Int64
	maxObservedTsNs.Store(lastOffsetTsNs)
	baseFn := genProcessFunction(sourcePath, targetPath, nil, nil, nil, nil, dataSink, pol.PropagateDel, false)
	processFn := func(resp *filer_pb.SubscribeMetadataResponse) error {
		if err := baseFn(resp); err != nil {
			if isIgnorable404(err) {
				return nil
			}
			return err
		}
		archivedReqs.Add(1)
		if resp.TsNs > maxObservedTsNs.Load() {
			maxObservedTsNs.Store(resp.TsNs)
		}
		if resp.EventNotification != nil && resp.EventNotification.NewEntry != nil && resp.EventNotification.NewEntry.Attributes != nil {
			archivedBytes.Add(int64(resp.EventNotification.NewEntry.Attributes.FileSize))
		}
		return nil
	}

	followOpt := &pb.MetadataFollowOption{
		ClientName:     "archive_" + pol.Bucket + "_" + dataSink.GetName(),
		ClientId:       util.RandomInt32(),
		ClientEpoch:    1,
		PathPrefix:     sourcePath,
		StartTsNs:      lastOffsetTsNs,
		StopTsNs:       stopTsNs,
		EventErrorType: pb.TrivialOnError,
	}
	if err := pb.FollowMetadata(sourceFiler, grpcDialOption, followOpt, processFn); err != nil {
		return err
	}
	nextOffset := maxObservedTsNs.Load()
	if nextOffset > lastOffsetTsNs {
		if err := setOffset(grpcDialOption, sourceFiler, archive.KeyPrefix, sinkID, nextOffset); err != nil {
			return err
		}
	}
	if dataSink.IsIncremental() && pol.RetentionDays > 0 {
		expiredDir := util.Join(targetPath, time.Now().UTC().Add(-time.Duration(pol.RetentionDays)*24*time.Hour).Format("2006-01-02"))
		_ = dataSink.DeleteEntry(expiredDir, true, true, nil)
	}
	glog.V(0).Infof("filer.archive bucket=%s offset=%d bytes=%d reqs=%d", pol.Bucket, nextOffset, archivedBytes.Load(), archivedReqs.Load())
	return nil
}
