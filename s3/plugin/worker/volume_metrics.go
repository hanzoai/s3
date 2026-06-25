package pluginworker

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/hanzoai/s3/s3/admin/topology"
	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	workertypes "github.com/hanzoai/s3/s3/worker/types"
)

// CollectVolumeMetricsFromMasters dials the provided master addresses in order
// until one returns a usable volume list, then converts that into per-volume
// health metrics, an active-topology view, and a replica-location map.
func CollectVolumeMetricsFromMasters(
	ctx context.Context,
	masterAddresses []string,
	collectionFilter string,
	grpcDialOption pb.DialOption,
) ([]*workertypes.VolumeHealthMetrics, *topology.ActiveTopology, map[uint32][]workertypes.ReplicaLocation, error) {
	if len(masterAddresses) == 0 {
		return nil, nil, nil, fmt.Errorf("no master addresses provided in cluster context")
	}

	for _, masterAddress := range masterAddresses {
		response, err := FetchVolumeList(ctx, masterAddress, grpcDialOption)
		if err != nil {
			glog.Warningf("Plugin worker failed master volume list at %s: %v", masterAddress, err)
			continue
		}

		metrics, activeTopology, replicaMap, buildErr := buildVolumeMetrics(response, collectionFilter)
		if buildErr != nil {
			// Configuration errors (e.g. invalid regex) will fail on every master,
			// so return immediately instead of masking them with retries.
			if isConfigError(buildErr) {
				return nil, nil, nil, buildErr
			}
			glog.Warningf("Plugin worker failed to build metrics from master %s: %v", masterAddress, buildErr)
			continue
		}
		return metrics, activeTopology, replicaMap, nil
	}

	return nil, nil, nil, fmt.Errorf("failed to load topology from all provided masters")
}

// FetchVolumeList dials the given master address (trying both the address as
// given and the gRPC port variant) and returns the master's volume list. Used
// by detection helpers that already know which master address to talk to.
// FetchDefaultReplicaPlacement returns the master's configured default replication
// (GetMasterConfiguration), used by detectors as the replica-placement fallback so
// the plugin path matches the shell. Returns "" if it cannot be fetched, so callers
// fall back to even spread rather than failing detection.
func FetchDefaultReplicaPlacement(ctx context.Context, masterAddresses []string, grpcDialOption pb.DialOption) string {
	for _, address := range masterAddresses {
		if ctx.Err() != nil {
			return ""
		}
		var replication string
		callErr := pb.WithMasterClient(ctx, false, pb.ServerAddress(address), grpcDialOption, false, func(client master_pb.HanzoClient) error {
			callCtx, cancelCall := context.WithTimeout(ctx, 10*time.Second)
			defer cancelCall()
			resp, err := client.GetMasterConfiguration(callCtx, &master_pb.GetMasterConfigurationRequest{})
			if err != nil {
				return err
			}
			replication = resp.DefaultReplication
			return nil
		})
		if callErr == nil {
			return replication
		}
	}
	return ""
}

func FetchVolumeList(ctx context.Context, address string, grpcDialOption pb.DialOption) (*master_pb.VolumeListResponse, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	var response *master_pb.VolumeListResponse
	err := pb.WithMasterClient(ctx, false, pb.ServerAddress(address), grpcDialOption, false, func(client master_pb.HanzoClient) error {
		callCtx, cancelCall := context.WithTimeout(ctx, 10*time.Second)
		defer cancelCall()
		resp, callErr := client.VolumeList(callCtx, &master_pb.VolumeListRequest{})
		if callErr != nil {
			return callErr
		}
		response = resp
		return nil
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func buildVolumeMetrics(
	response *master_pb.VolumeListResponse,
	collectionFilter string,
) ([]*workertypes.VolumeHealthMetrics, *topology.ActiveTopology, map[uint32][]workertypes.ReplicaLocation, error) {
	if response == nil || response.TopologyInfo == nil {
		return nil, nil, nil, fmt.Errorf("volume list response has no topology info")
	}

	activeTopology := topology.NewActiveTopology(10)
	if err := activeTopology.UpdateTopology(response.TopologyInfo); err != nil {
		return nil, nil, nil, err
	}

	var collectionRegex *regexp.Regexp
	trimmedFilter := strings.TrimSpace(collectionFilter)
	filterMode := CollectionFilterMode(trimmedFilter)
	if trimmedFilter != "" && filterMode != CollectionFilterAll && filterMode != CollectionFilterEach && trimmedFilter != "*" {
		var err error
		collectionRegex, err = regexp.Compile(trimmedFilter)
		if err != nil {
			return nil, nil, nil, &configError{err: fmt.Errorf("invalid collection_filter regex %q: %w", trimmedFilter, err)}
		}
	}

	volumeSizeLimitBytes := uint64(response.VolumeSizeLimitMb) * 1024 * 1024
	now := time.Now()
	metrics := make([]*workertypes.VolumeHealthMetrics, 0, 256)
	replicaMap := make(map[uint32][]workertypes.ReplicaLocation)

	for _, dc := range response.TopologyInfo.DataCenterInfos {
		for _, rack := range dc.RackInfos {
			for _, node := range rack.DataNodeInfos {
				for diskType, diskInfo := range node.DiskInfos {
					for _, volume := range diskInfo.VolumeInfos {
						// Build replica map from ALL volumes BEFORE collection filtering,
						// since replicas may span filtered/unfiltered nodes.
						replicaMap[volume.Id] = append(replicaMap[volume.Id], workertypes.ReplicaLocation{
							DataCenter: dc.Id,
							Rack:       rack.Id,
							NodeID:     node.Id,
							Host:       pb.NewServerAddressFromDataNode(node).ToHost(),
						})

						if collectionRegex != nil && !collectionRegex.MatchString(volume.Collection) {
							continue
						}

						metric := &workertypes.VolumeHealthMetrics{
							VolumeID:         volume.Id,
							Server:           node.Id,
							ServerAddress:    string(pb.NewServerAddressFromDataNode(node)),
							DiskType:         diskType,
							DiskId:           volume.DiskId,
							DataCenter:       dc.Id,
							Rack:             rack.Id,
							Collection:       volume.Collection,
							Size:             volume.Size,
							DeletedBytes:     volume.DeletedByteCount,
							LastModified:     time.Unix(volume.ModifiedAtSecond, 0),
							ReplicaCount:     1,
							ExpectedReplicas: int(volume.ReplicaPlacement),
							IsReadOnly:       volume.ReadOnly,
							HasRemoteCopy:    volume.RemoteStorageName != "",
						}
						if metric.Size > 0 {
							metric.GarbageRatio = float64(metric.DeletedBytes) / float64(metric.Size)
						}
						if volumeSizeLimitBytes > 0 {
							metric.FullnessRatio = float64(metric.Size) / float64(volumeSizeLimitBytes)
						}
						metric.Age = now.Sub(metric.LastModified)
						metrics = append(metrics, metric)
					}
				}
			}
		}
	}

	replicaCounts := make(map[uint32]int)
	for _, metric := range metrics {
		replicaCounts[metric.VolumeID]++
	}
	for _, metric := range metrics {
		metric.ReplicaCount = replicaCounts[metric.VolumeID]
	}

	return metrics, activeTopology, replicaMap, nil
}

// configError wraps configuration errors that should not be retried across masters.
type configError struct {
	err error
}

func (e *configError) Error() string { return e.err.Error() }
func (e *configError) Unwrap() error { return e.err }

func isConfigError(err error) bool {
	var ce *configError
	return errors.As(err, &ce)
}

// MasterAddressCandidates returns address forms to try when dialing a master:
// the address as given plus the gRPC variant (port + 10000). Both are tried
// because callers may pass either an HTTP-port or gRPC-port address.
func MasterAddressCandidates(address string) []string {
	trimmed := strings.TrimSpace(address)
	if trimmed == "" {
		return nil
	}
	candidateSet := map[string]struct{}{
		trimmed: {},
	}
	converted := pb.ServerToGrpcAddress(trimmed)
	candidateSet[converted] = struct{}{}

	candidates := make([]string, 0, len(candidateSet))
	for candidate := range candidateSet {
		candidates = append(candidates, candidate)
	}
	sort.Strings(candidates)
	return candidates
}
