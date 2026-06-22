//go:build !elastic

package command

import (
	"fmt"

	"github.com/hanzoai/s3/s3/pb/filer_pb"
)

func sendToElasticSearchFunc(servers string, esIndex string) (func(resp *filer_pb.SubscribeMetadataResponse) error, error) {
	return nil, fmt.Errorf("this s3 binary was built without Elasticsearch support. Rebuild with `-tags elastic` (see docker/Dockerfile.go_build) to use -es/-es.index")
}
