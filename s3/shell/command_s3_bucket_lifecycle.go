package shell

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"io"

	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/s3api/lifecycle_xml"
	"github.com/hanzoai/s3/s3/s3api/s3_constants"
	"github.com/hanzoai/s3/s3/s3api/s3bucket"
)

func init() {
	Commands = append(Commands, &commandS3BucketLifecycle{})
}

type commandS3BucketLifecycle struct {
}

func (c *commandS3BucketLifecycle) Name() string {
	return "s3.bucket.lifecycle"
}

func (c *commandS3BucketLifecycle) Help() string {
	return `view, set, or clear the bucket lifecycle policy

	The policy is the same XML PutBucketLifecycleConfiguration stores, so what
	is set here is what GetBucketLifecycleConfiguration returns. Setting writes
	one rule covering the whole bucket and replaces the entire document —
	whatever the S3 API put there before is gone.

	Both rules act on versions only: -noncurrent-expire-days needs noncurrent
	versions and -expire-delete-marker needs delete markers, neither of which
	an unversioned bucket ever has. Such a policy stores fine and then never
	fires. Turn versioning on first with
	s3.bucket.versioning -name <bucket_name> -enable.

	Example:
		# Show the current policy
		s3.bucket.lifecycle -name <bucket_name>

		# Delete a version 7 days after it stops being current
		s3.bucket.lifecycle -name <bucket_name> -noncurrent-expire-days 7

		# ... and drop a delete marker once it is the only version left
		s3.bucket.lifecycle -name <bucket_name> -noncurrent-expire-days 7 -expire-delete-marker

		# Clear the policy
		s3.bucket.lifecycle -name <bucket_name> -remove
`
}

func (c *commandS3BucketLifecycle) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3BucketLifecycle) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {
	bucketCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	bucketName := bucketCommand.String("name", "", "bucket name")
	noncurrentDays := bucketCommand.Int("noncurrent-expire-days", 0, "delete a version this many days after it stops being current")
	expireDeleteMarker := bucketCommand.Bool("expire-delete-marker", false, "delete a delete marker once it is an object's only remaining version")
	remove := bucketCommand.Bool("remove", false, "clear the lifecycle policy")
	if err = bucketCommand.Parse(args); err != nil {
		return err
	}

	if *bucketName == "" {
		return fmt.Errorf("empty bucket name")
	}
	if err := s3bucket.VerifyS3BucketName(*bucketName); err != nil {
		return fmt.Errorf("invalid bucket name %q: %w", *bucketName, err)
	}
	if *noncurrentDays < 0 {
		return fmt.Errorf("-noncurrent-expire-days cannot be negative")
	}
	if *remove && (*noncurrentDays > 0 || *expireDeleteMarker) {
		return fmt.Errorf("-remove cannot be combined with a rule flag")
	}

	return commandEnv.WithFilerClient(false, func(client filer_pb.HanzoFilerClient) error {
		resp, err := client.GetFilerConfiguration(context.Background(), &filer_pb.GetFilerConfigurationRequest{})
		if err != nil {
			return fmt.Errorf("get filer configuration: %w", err)
		}
		filerBucketsPath := resp.DirBuckets

		lookupResp, err := client.LookupDirectoryEntry(context.Background(), &filer_pb.LookupDirectoryEntryRequest{
			Directory: filerBucketsPath,
			Name:      *bucketName,
		})
		if err != nil {
			return fmt.Errorf("lookup bucket %s: %w", *bucketName, err)
		}
		entry := lookupResp.Entry

		if !*remove && *noncurrentDays == 0 && !*expireDeleteMarker {
			fmt.Fprintf(writer, "Bucket: %s\n", *bucketName)
			policy := entry.Extended[s3_constants.ExtLifecycleConfigurationXMLKey]
			if len(policy) == 0 {
				fmt.Fprintf(writer, "Lifecycle: none\n")
				return nil
			}
			fmt.Fprintf(writer, "Lifecycle:\n%s\n", policy)
			return nil
		}

		if entry.Extended == nil {
			entry.Extended = make(map[string][]byte)
		}
		if *remove {
			delete(entry.Extended, s3_constants.ExtLifecycleConfigurationXMLKey)
		} else {
			policy, err := lifecyclePolicy(*noncurrentDays, *expireDeleteMarker)
			if err != nil {
				return fmt.Errorf("marshal lifecycle policy: %w", err)
			}
			entry.Extended[s3_constants.ExtLifecycleConfigurationXMLKey] = policy
		}

		if _, err := client.UpdateEntry(context.Background(), &filer_pb.UpdateEntryRequest{
			Directory: filerBucketsPath,
			Entry:     entry,
		}); err != nil {
			return fmt.Errorf("failed to update bucket: %w", err)
		}

		if *remove {
			fmt.Fprintf(writer, "Bucket %s lifecycle policy cleared\n", *bucketName)
			return nil
		}
		fmt.Fprintf(writer, "Bucket %s lifecycle policy set\n", *bucketName)
		if !isBucketVersioned(entry) {
			fmt.Fprintf(writer, "warning: bucket %s is not versioned, so this policy has nothing to expire\n", *bucketName)
		}
		return nil
	})
}

// lifecyclePolicy renders the one-rule document the set path stores, using
// the same types PutBucketLifecycleConfiguration unmarshals so the shell and
// the S3 API cannot disagree about the format. An empty Prefix is the whole
// bucket, matching what GetBucketLifecycleConfiguration emits.
func lifecyclePolicy(noncurrentDays int, expireDeleteMarker bool) ([]byte, error) {
	rule := lifecycle_xml.Rule{
		Status: lifecycle_xml.Enabled,
		Prefix: lifecycle_xml.NewPrefix(""),
	}
	if noncurrentDays > 0 {
		rule.NoncurrentVersionExpiration = lifecycle_xml.NewNoncurrentDays(noncurrentDays)
	}
	if expireDeleteMarker {
		rule.Expiration = lifecycle_xml.NewExpirationDeleteMarker()
	}
	return xml.Marshal(lifecycle_xml.Lifecycle{Rules: []lifecycle_xml.Rule{rule}})
}

// isBucketVersioned reports whether the bucket holds versions at all.
// Suspended counts: it stops making new ones but keeps those already there,
// which these rules still expire.
func isBucketVersioned(entry *filer_pb.Entry) bool {
	status, lockEnabled := getBucketVersioningState(entry)
	return lockEnabled || status == s3_constants.VersioningEnabled || status == s3_constants.VersioningSuspended
}
