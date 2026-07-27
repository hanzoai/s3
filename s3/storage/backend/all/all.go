// Package all registers every tiered-storage backend factory. Binaries that
// serve or manage tiered volumes blank-import it from their composition root;
// keeping the registrations out of s3/storage lets library consumers avoid
// pulling the backend SDKs into their dependency graph.
package all

import (
	_ "github.com/hanzoai/s3/s3/storage/backend/s3_backend"
)
