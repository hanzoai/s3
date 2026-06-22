package mount

import (
	"testing"

	"github.com/hanzoai/s3/s3/filer"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
)

func TestIsFilerConfUpdateEventMatchesRenameTarget(t *testing.T) {
	event := &filer_pb.SubscribeMetadataResponse{
		Directory: "/tmp",
		EventNotification: &filer_pb.EventNotification{
			OldEntry:      &filer_pb.Entry{Name: filer.FilerConfName},
			NewEntry:      &filer_pb.Entry{Name: filer.FilerConfName},
			NewParentPath: filer.DirectoryEtcHanzo,
		},
	}

	if !isFilerConfUpdateEvent(event, filer.DirectoryEtcHanzo, filer.FilerConfName) {
		t.Fatalf("expected rename target to match filer.conf watcher")
	}
}
