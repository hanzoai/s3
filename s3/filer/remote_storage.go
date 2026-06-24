package filer

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/remote_pb"
	"github.com/hanzoai/s3/s3/remote_storage"
	"github.com/hanzoai/s3/s3/util"
	remotewire "github.com/hanzoai/s3/s3/wire/remote"
	"google.golang.org/grpc"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/viant/ptrie"
)

const REMOTE_STORAGE_CONF_SUFFIX = ".conf"
const REMOTE_STORAGE_MOUNT_FILE = "mount.mapping"

// The persisted remote-storage artifacts (per-storage `<name>.conf` RemoteConf
// entries and the `mount.mapping` RemoteStorageMapping) are serialized with the
// native ZAP wire format via the remotewire package. The in-memory data model
// stays remote_pb.* (threaded through the RemoteStorageClient interface and the
// ptrie/maps below); these bridge functions are the one and only seam between
// that model and its on-disk ZAP encoding.

// MarshalRemoteConf encodes a RemoteConf as ZAP wire bytes.
func MarshalRemoteConf(conf *remote_pb.RemoteConf) []byte {
	return remotewire.NewRemoteConf(remotewire.RemoteConfInput{
		Type:                            conf.Type,
		Name:                            conf.Name,
		S3AccessKey:                     conf.S3AccessKey,
		S3SecretKey:                     conf.S3SecretKey,
		S3Region:                        conf.S3Region,
		S3Endpoint:                      conf.S3Endpoint,
		S3StorageClass:                  conf.S3StorageClass,
		GcsGoogleApplicationCredentials: conf.GcsGoogleApplicationCredentials,
		GcsProjectID:                    conf.GcsProjectId,
		AzureAccountName:                conf.AzureAccountName,
		AzureAccountKey:                 conf.AzureAccountKey,
		BackblazeKeyID:                  conf.BackblazeKeyId,
		BackblazeApplicationKey:         conf.BackblazeApplicationKey,
		BackblazeEndpoint:               conf.BackblazeEndpoint,
		BackblazeRegion:                 conf.BackblazeRegion,
		AliyunAccessKey:                 conf.AliyunAccessKey,
		AliyunSecretKey:                 conf.AliyunSecretKey,
		AliyunEndpoint:                  conf.AliyunEndpoint,
		AliyunRegion:                    conf.AliyunRegion,
		TencentSecretID:                 conf.TencentSecretId,
		TencentSecretKey:                conf.TencentSecretKey,
		TencentEndpoint:                 conf.TencentEndpoint,
		BaiduAccessKey:                  conf.BaiduAccessKey,
		BaiduSecretKey:                  conf.BaiduSecretKey,
		BaiduEndpoint:                   conf.BaiduEndpoint,
		BaiduRegion:                     conf.BaiduRegion,
		WasabiAccessKey:                 conf.WasabiAccessKey,
		WasabiSecretKey:                 conf.WasabiSecretKey,
		WasabiEndpoint:                  conf.WasabiEndpoint,
		WasabiRegion:                    conf.WasabiRegion,
		FilebaseAccessKey:               conf.FilebaseAccessKey,
		FilebaseSecretKey:               conf.FilebaseSecretKey,
		FilebaseEndpoint:                conf.FilebaseEndpoint,
		StorjAccessKey:                  conf.StorjAccessKey,
		StorjSecretKey:                  conf.StorjSecretKey,
		StorjEndpoint:                   conf.StorjEndpoint,
		ContaboAccessKey:                conf.ContaboAccessKey,
		ContaboSecretKey:                conf.ContaboSecretKey,
		ContaboEndpoint:                 conf.ContaboEndpoint,
		ContaboRegion:                   conf.ContaboRegion,
		S3ForcePathStyle:                conf.S3ForcePathStyle,
		S3SupportTagging:                conf.S3SupportTagging,
		S3V4Signature:                   conf.S3V4Signature,
	})
}

// UnmarshalRemoteConf decodes ZAP wire bytes into conf. Empty data leaves conf
// at its zero values (mirroring proto.Unmarshal of an empty buffer).
func UnmarshalRemoteConf(data []byte, conf *remote_pb.RemoteConf) error {
	if len(data) == 0 {
		return nil
	}
	v, err := remotewire.WrapRemoteConf(data)
	if err != nil {
		return err
	}
	conf.Type = v.Type()
	conf.Name = v.Name()
	conf.S3AccessKey = v.S3AccessKey()
	conf.S3SecretKey = v.S3SecretKey()
	conf.S3Region = v.S3Region()
	conf.S3Endpoint = v.S3Endpoint()
	conf.S3StorageClass = v.S3StorageClass()
	conf.GcsGoogleApplicationCredentials = v.GcsGoogleApplicationCredentials()
	conf.GcsProjectId = v.GcsProjectID()
	conf.AzureAccountName = v.AzureAccountName()
	conf.AzureAccountKey = v.AzureAccountKey()
	conf.BackblazeKeyId = v.BackblazeKeyID()
	conf.BackblazeApplicationKey = v.BackblazeApplicationKey()
	conf.BackblazeEndpoint = v.BackblazeEndpoint()
	conf.BackblazeRegion = v.BackblazeRegion()
	conf.AliyunAccessKey = v.AliyunAccessKey()
	conf.AliyunSecretKey = v.AliyunSecretKey()
	conf.AliyunEndpoint = v.AliyunEndpoint()
	conf.AliyunRegion = v.AliyunRegion()
	conf.TencentSecretId = v.TencentSecretID()
	conf.TencentSecretKey = v.TencentSecretKey()
	conf.TencentEndpoint = v.TencentEndpoint()
	conf.BaiduAccessKey = v.BaiduAccessKey()
	conf.BaiduSecretKey = v.BaiduSecretKey()
	conf.BaiduEndpoint = v.BaiduEndpoint()
	conf.BaiduRegion = v.BaiduRegion()
	conf.WasabiAccessKey = v.WasabiAccessKey()
	conf.WasabiSecretKey = v.WasabiSecretKey()
	conf.WasabiEndpoint = v.WasabiEndpoint()
	conf.WasabiRegion = v.WasabiRegion()
	conf.FilebaseAccessKey = v.FilebaseAccessKey()
	conf.FilebaseSecretKey = v.FilebaseSecretKey()
	conf.FilebaseEndpoint = v.FilebaseEndpoint()
	conf.StorjAccessKey = v.StorjAccessKey()
	conf.StorjSecretKey = v.StorjSecretKey()
	conf.StorjEndpoint = v.StorjEndpoint()
	conf.ContaboAccessKey = v.ContaboAccessKey()
	conf.ContaboSecretKey = v.ContaboSecretKey()
	conf.ContaboEndpoint = v.ContaboEndpoint()
	conf.ContaboRegion = v.ContaboRegion()
	conf.S3ForcePathStyle = v.S3ForcePathStyle()
	conf.S3SupportTagging = v.S3SupportTagging()
	conf.S3V4Signature = v.S3V4Signature()
	return nil
}

// MarshalRemoteStorageMapping encodes a RemoteStorageMapping as ZAP wire bytes.
// The proto map<string,RemoteStorageLocation> is emitted as the wire format's
// list of {key,value} entries.
func MarshalRemoteStorageMapping(mappings *remote_pb.RemoteStorageMapping) []byte {
	entries := make([]remotewire.RemoteStorageMappingEntryInput, 0, len(mappings.Mappings))
	for dir, loc := range mappings.Mappings {
		entries = append(entries, remotewire.RemoteStorageMappingEntryInput{
			Key:   dir,
			Value: remoteStorageLocationInput(loc),
		})
	}
	return remotewire.NewRemoteStorageMapping(remotewire.RemoteStorageMappingInput{
		Mappings:                 entries,
		PrimaryBucketStorageName: mappings.PrimaryBucketStorageName,
	})
}

// UnmarshalRemoteStorageMapping decodes ZAP wire bytes into mappings. Empty data
// leaves mappings.Mappings empty (mirroring proto.Unmarshal of an empty buffer).
func UnmarshalRemoteStorageMapping(data []byte, mappings *remote_pb.RemoteStorageMapping) error {
	if mappings.Mappings == nil {
		mappings.Mappings = make(map[string]*remote_pb.RemoteStorageLocation)
	}
	if len(data) == 0 {
		return nil
	}
	v, err := remotewire.WrapRemoteStorageMapping(data)
	if err != nil {
		return err
	}
	mappings.PrimaryBucketStorageName = v.PrimaryBucketStorageName()
	for i := 0; i < v.MappingsLen(); i++ {
		e := v.MappingsAt(i)
		loc, err := e.Value()
		if err != nil {
			return err
		}
		mappings.Mappings[e.Key()] = remoteStorageLocationFromView(loc)
	}
	return nil
}

// remoteStorageLocationInput projects a RemoteStorageLocation into its wire
// input form.
func remoteStorageLocationInput(loc *remote_pb.RemoteStorageLocation) remotewire.RemoteStorageLocationInput {
	return remotewire.RemoteStorageLocationInput{
		Name:                   loc.Name,
		Bucket:                 loc.Bucket,
		Path:                   loc.Path,
		ListingCacheTTLSeconds: loc.ListingCacheTtlSeconds,
	}
}

// remoteStorageLocationFromView materializes a RemoteStorageLocation from a
// zero-copy wire view.
func remoteStorageLocationFromView(v remotewire.RemoteStorageLocation) *remote_pb.RemoteStorageLocation {
	return &remote_pb.RemoteStorageLocation{
		Name:                   v.Name(),
		Bucket:                 v.Bucket(),
		Path:                   v.Path(),
		ListingCacheTtlSeconds: v.ListingCacheTTLSeconds(),
	}
}

type FilerRemoteStorage struct {
	rules             ptrie.Trie[*remote_pb.RemoteStorageLocation]
	storageNameToConf map[string]*remote_pb.RemoteConf
}

func NewFilerRemoteStorage() (rs *FilerRemoteStorage) {
	rs = &FilerRemoteStorage{
		rules:             ptrie.New[*remote_pb.RemoteStorageLocation](),
		storageNameToConf: make(map[string]*remote_pb.RemoteConf),
	}
	return rs
}

func (rs *FilerRemoteStorage) LoadRemoteStorageConfigurationsAndMapping(filer *Filer) (err error) {
	// execute this on filer

	limit := int64(math.MaxInt32)

	entries, _, err := filer.ListDirectoryEntries(context.Background(), DirectoryEtcRemote, "", false, limit, "", "", "")
	if err != nil {
		if err == filer_pb.ErrNotFound {
			return nil
		}
		glog.Errorf("read remote storage %s: %v", DirectoryEtcRemote, err)
		return
	}

	for _, entry := range entries {
		if entry.Name() == REMOTE_STORAGE_MOUNT_FILE {
			if err := rs.loadRemoteStorageMountMapping(entry.Content); err != nil {
				return err
			}
			continue
		}
		if !strings.HasSuffix(entry.Name(), REMOTE_STORAGE_CONF_SUFFIX) {
			return nil
		}
		conf := &remote_pb.RemoteConf{}
		if err := UnmarshalRemoteConf(entry.Content, conf); err != nil {
			return fmt.Errorf("unmarshal %s/%s: %v", DirectoryEtcRemote, entry.Name(), err)
		}
		rs.storageNameToConf[conf.Name] = conf
	}
	return nil
}

func (rs *FilerRemoteStorage) loadRemoteStorageMountMapping(data []byte) (err error) {
	mappings := &remote_pb.RemoteStorageMapping{}
	if err := UnmarshalRemoteStorageMapping(data, mappings); err != nil {
		return fmt.Errorf("unmarshal %s/%s: %v", DirectoryEtcRemote, REMOTE_STORAGE_MOUNT_FILE, err)
	}
	for dir, storageLocation := range mappings.Mappings {
		rs.mapDirectoryToRemoteStorage(util.FullPath(dir), storageLocation)
	}
	return nil
}

func (rs *FilerRemoteStorage) mapDirectoryToRemoteStorage(dir util.FullPath, loc *remote_pb.RemoteStorageLocation) {
	rs.rules.Put([]byte(dir+"/"), loc)
}

// FindMountDirectory returns the mount directory and location for p. When multiple
// mounts match (e.g. /buckets/b and /buckets/b/prefix), ptrie MatchPrefix visits
// shorter prefixes first, so the last match is the longest prefix.
func (rs *FilerRemoteStorage) FindMountDirectory(p util.FullPath) (mountDir util.FullPath, remoteLocation *remote_pb.RemoteStorageLocation) {
	rs.rules.MatchPrefix([]byte(p), func(key []byte, value *remote_pb.RemoteStorageLocation) bool {
		mountDir = util.FullPath(string(key[:len(key)-1]))
		remoteLocation = value
		return true
	})
	return
}

func (rs *FilerRemoteStorage) FindRemoteStorageClient(p util.FullPath) (client remote_storage.RemoteStorageClient, remoteConf *remote_pb.RemoteConf, found bool) {
	var storageLocation *remote_pb.RemoteStorageLocation
	rs.rules.MatchPrefix([]byte(p), func(key []byte, value *remote_pb.RemoteStorageLocation) bool {
		storageLocation = value
		return true
	})

	if storageLocation == nil {
		found = false
		return
	}

	return rs.GetRemoteStorageClient(storageLocation.Name)
}

func (rs *FilerRemoteStorage) GetRemoteStorageClient(storageName string) (client remote_storage.RemoteStorageClient, remoteConf *remote_pb.RemoteConf, found bool) {
	remoteConf, found = rs.storageNameToConf[storageName]
	if !found {
		return
	}

	var err error
	if client, err = remote_storage.GetRemoteStorage(remoteConf); err == nil {
		found = true
		return
	}
	return
}

func UnmarshalRemoteStorageMappings(oldContent []byte) (mappings *remote_pb.RemoteStorageMapping, err error) {
	mappings = &remote_pb.RemoteStorageMapping{
		Mappings: make(map[string]*remote_pb.RemoteStorageLocation),
	}
	if len(oldContent) > 0 {
		if err = UnmarshalRemoteStorageMapping(oldContent, mappings); err != nil {
			glog.Warningf("unmarshal existing mappings: %v", err)
		}
	}
	return
}

func ReadRemoteStorageConf(grpcDialOption grpc.DialOption, filerAddress pb.ServerAddress, storageName string) (conf *remote_pb.RemoteConf, readErr error) {
	var oldContent []byte
	if readErr = pb.WithFilerClient(false, 0, filerAddress, grpcDialOption, func(client filer_pb.HanzoFilerClient) error {
		oldContent, readErr = ReadInsideFiler(context.Background(), client, DirectoryEtcRemote, storageName+REMOTE_STORAGE_CONF_SUFFIX)
		return readErr
	}); readErr != nil {
		return nil, readErr
	}

	// unmarshal storage configuration
	conf = &remote_pb.RemoteConf{}
	if unMarshalErr := UnmarshalRemoteConf(oldContent, conf); unMarshalErr != nil {
		readErr = fmt.Errorf("unmarshal %s/%s: %v", DirectoryEtcRemote, storageName+REMOTE_STORAGE_CONF_SUFFIX, unMarshalErr)
		return
	}

	return
}

func DetectMountInfo(grpcDialOption grpc.DialOption, filerAddress pb.ServerAddress, dir string) (*remote_pb.RemoteStorageMapping, string, *remote_pb.RemoteStorageLocation, *remote_pb.RemoteConf, error) {

	mappings, listErr := ReadMountMappings(grpcDialOption, filerAddress)
	if listErr != nil {
		return nil, "", nil, nil, listErr
	}
	if dir == "" {
		return mappings, "", nil, nil, fmt.Errorf("need to specify '-dir' option")
	}

	var localMountedDir string
	var remoteStorageMountedLocation *remote_pb.RemoteStorageLocation
	for k, loc := range mappings.Mappings {
		if strings.HasPrefix(dir, k) {
			localMountedDir, remoteStorageMountedLocation = k, loc
		}
	}
	if localMountedDir == "" {
		return mappings, localMountedDir, remoteStorageMountedLocation, nil, fmt.Errorf("%s is not mounted", dir)
	}

	// find remote storage configuration
	remoteStorageConf, err := ReadRemoteStorageConf(grpcDialOption, filerAddress, remoteStorageMountedLocation.Name)
	if err != nil {
		return mappings, localMountedDir, remoteStorageMountedLocation, remoteStorageConf, err
	}

	return mappings, localMountedDir, remoteStorageMountedLocation, remoteStorageConf, nil
}
