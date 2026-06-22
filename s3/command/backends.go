package command

// Register the tiered-storage backends here, at the binary's composition root,
// so library consumers of s3/storage don't pull the backend SDKs.
import _ "github.com/hanzoai/s3/s3/storage/backend/all"
