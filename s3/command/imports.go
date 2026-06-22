package command

import (
	_ "net/http/pprof"

	_ "github.com/hanzoai/s3/s3/remote_storage/azure"
	_ "github.com/hanzoai/s3/s3/remote_storage/gcs"
	_ "github.com/hanzoai/s3/s3/remote_storage/s3"

	_ "github.com/hanzoai/s3/s3/replication/sink/azuresink"
	_ "github.com/hanzoai/s3/s3/replication/sink/b2sink"
	_ "github.com/hanzoai/s3/s3/replication/sink/filersink"
	_ "github.com/hanzoai/s3/s3/replication/sink/gcssink"
	_ "github.com/hanzoai/s3/s3/replication/sink/localsink"
	_ "github.com/hanzoai/s3/s3/replication/sink/s3sink"

	_ "github.com/hanzoai/s3/s3/filer/arangodb"
	_ "github.com/hanzoai/s3/s3/filer/cassandra"
	_ "github.com/hanzoai/s3/s3/filer/elastic/v7"
	_ "github.com/hanzoai/s3/s3/filer/etcd"
	_ "github.com/hanzoai/s3/s3/filer/hbase"
	_ "github.com/hanzoai/s3/s3/filer/leveldb"
	_ "github.com/hanzoai/s3/s3/filer/leveldb2"
	_ "github.com/hanzoai/s3/s3/filer/leveldb3"
	_ "github.com/hanzoai/s3/s3/filer/mongodb"
	_ "github.com/hanzoai/s3/s3/filer/mysql"
	_ "github.com/hanzoai/s3/s3/filer/mysql2"
	_ "github.com/hanzoai/s3/s3/filer/postgres"
	_ "github.com/hanzoai/s3/s3/filer/postgres2"
	_ "github.com/hanzoai/s3/s3/filer/redis"
	_ "github.com/hanzoai/s3/s3/filer/redis2"
	_ "github.com/hanzoai/s3/s3/filer/redis3"
	_ "github.com/hanzoai/s3/s3/filer/sqlite"
	_ "github.com/hanzoai/s3/s3/filer/tarantool"
	_ "github.com/hanzoai/s3/s3/filer/tikv"
	_ "github.com/hanzoai/s3/s3/filer/ydb"

	_ "github.com/hanzoai/s3/s3/credential/filer_etc"
)
