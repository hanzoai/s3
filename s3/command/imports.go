package command

import (
	_ "net/http/pprof"

	_ "github.com/hanzoai/s3/s3/remote_storage/s3"

	_ "github.com/hanzoai/s3/s3/replication/sink/filersink"
	_ "github.com/hanzoai/s3/s3/replication/sink/localsink"
	_ "github.com/hanzoai/s3/s3/replication/sink/s3sink"

	_ "github.com/hanzoai/s3/s3/filer/sqlite"

	_ "github.com/hanzoai/s3/s3/credential/filer_etc"
)
