package sub

import (
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/util"
)

type NotificationInput interface {
	// GetName gets the name to locate the configuration in sync.toml file
	GetName() string
	// Initialize initializes the file store
	Initialize(configuration util.Configuration, prefix string) error
	ReceiveMessage() (key string, message *filer_pb.EventNotification, onSuccessFn func(), onFailureFn func(), err error)
}

var (
	NotificationInputs []NotificationInput
)
