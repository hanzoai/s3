package command

import (
	_ "net/http/pprof"

	_ "github.com/hanzoai/s3/weed/remote_storage/azure"
	_ "github.com/hanzoai/s3/weed/remote_storage/gcs"
	_ "github.com/hanzoai/s3/weed/remote_storage/s3"

	_ "github.com/hanzoai/s3/weed/replication/sink/azuresink"
	_ "github.com/hanzoai/s3/weed/replication/sink/b2sink"
	_ "github.com/hanzoai/s3/weed/replication/sink/filersink"
	_ "github.com/hanzoai/s3/weed/replication/sink/gcssink"
	_ "github.com/hanzoai/s3/weed/replication/sink/localsink"
	_ "github.com/hanzoai/s3/weed/replication/sink/s3sink"

	_ "github.com/hanzoai/s3/weed/filer/arangodb"
	_ "github.com/hanzoai/s3/weed/filer/cassandra"
	_ "github.com/hanzoai/s3/weed/filer/elastic/v7"
	_ "github.com/hanzoai/s3/weed/filer/etcd"
	_ "github.com/hanzoai/s3/weed/filer/hbase"
	_ "github.com/hanzoai/s3/weed/filer/leveldb"
	_ "github.com/hanzoai/s3/weed/filer/leveldb2"
	_ "github.com/hanzoai/s3/weed/filer/leveldb3"
	_ "github.com/hanzoai/s3/weed/filer/mongodb"
	_ "github.com/hanzoai/s3/weed/filer/mysql"
	_ "github.com/hanzoai/s3/weed/filer/mysql2"
	_ "github.com/hanzoai/s3/weed/filer/postgres"
	_ "github.com/hanzoai/s3/weed/filer/postgres2"
	_ "github.com/hanzoai/s3/weed/filer/redis"
	_ "github.com/hanzoai/s3/weed/filer/redis2"
	_ "github.com/hanzoai/s3/weed/filer/redis3"
	_ "github.com/hanzoai/s3/weed/filer/sqlite"
	_ "github.com/hanzoai/s3/weed/filer/tarantool"
	_ "github.com/hanzoai/s3/weed/filer/tikv"
	_ "github.com/hanzoai/s3/weed/filer/ydb"

	_ "github.com/hanzoai/s3/weed/credential/filer_etc"
)
