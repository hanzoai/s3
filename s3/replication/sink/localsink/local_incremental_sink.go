package localsink

import (
	"github.com/hanzoai/s3/s3/replication/sink"
)

type LocalIncSink struct {
	LocalSink
}

func (localincsink *LocalIncSink) GetName() string {
	return "local_incremental"
}

func init() {
	sink.Sinks = append(sink.Sinks, &LocalIncSink{})
}
