package integration

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/zap-proto/go/transport"

	"github.com/hanzoai/s3/s3/filer_client"
	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/mq"
	mqsvc "github.com/hanzoai/s3/s3/svc/mq"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
	"github.com/hanzoai/s3/s3/pb/schema_pb"
	"github.com/hanzoai/s3/s3/util"
)

const (
	brokerHealthCheckTimeout  = 5 * time.Second
	brokerHealthCheckAttempts = 3
	brokerHealthCheckBackoff  = 250 * time.Millisecond
)

// dialBrokerZap opens a ZAP connection to a broker address, PQ-secured mTLS when
// grpc.mq.cert/.key is configured and plaintext otherwise — the single place this
// package's broker dial-side security gate lives, mirroring pb.dialBrokerZapAddr.
func dialBrokerZap(brokerAddress string) (transport.Conn, error) {
	if cfg := pb.ClientTLSConfig(util.GetViper(), "grpc.mq"); cfg != nil {
		return transport.DialTLS("tcp", brokerAddress, transport.PQTLSConfig(cfg))
	}
	return transport.Dial("tcp", brokerAddress)
}

// NewBrokerClientWithFilerAccessor creates a client with a shared filer accessor
func NewBrokerClientWithFilerAccessor(brokerAddress string, filerClientAccessor *filer_client.FilerClientAccessor) (*BrokerClient, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Connect to the broker over the native ZAP transport. The broker fully cut over
	// to ZAP (no gRPC), so its address IS the ZAP endpoint. When grpc.mq.cert/.key is
	// configured the connection is PQ-secured mTLS (transport.PQTLSConfig pins the
	// X25519MLKEM768 hybrid); otherwise plaintext (loopback / dev), matching the
	// broker server's gating. The Conn outlives BrokerClient.Close of individual
	// streams: closing it (in Close) drops every multiplexed stream at once.
	util.LoadSecurityConfiguration()
	conn, err := dialBrokerZap(brokerAddress)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to broker %s: %v", brokerAddress, err)
	}

	client := mqsvc.New(conn, nil)

	return &BrokerClient{
		filerClientAccessor:         filerClientAccessor,
		brokerAddress:               brokerAddress,
		conn:                        conn,
		client:                      client,
		publishers:                  make(map[string]*BrokerPublisherSession),
		subscribers:                 make(map[string]*BrokerSubscriberSession),
		fetchRequests:               make(map[string]*FetchRequest),
		partitionAssignmentCache:    make(map[string]*partitionAssignmentCacheEntry),
		partitionAssignmentCacheTTL: 30 * time.Second, // Same as broker's cache TTL
		ctx:                         ctx,
		cancel:                      cancel,
	}, nil
}

// Close shuts down the broker client and all streams
func (bc *BrokerClient) Close() error {
	bc.cancel()

	// Close all publisher streams
	bc.publishersLock.Lock()
	for key, session := range bc.publishers {
		if session.Stream != nil {
			_ = session.Stream.CloseSend()
		}
		delete(bc.publishers, key)
	}
	bc.publishersLock.Unlock()

	// Close all subscriber streams
	bc.subscribersLock.Lock()
	for key, session := range bc.subscribers {
		if session.Stream != nil {
			_ = session.Stream.CloseSend()
		}
		if session.Cancel != nil {
			session.Cancel()
		}
		delete(bc.subscribers, key)
	}
	bc.subscribersLock.Unlock()

	return bc.conn.Close()
}

// HealthCheck verifies the broker connection is working
func (bc *BrokerClient) HealthCheck() error {
	if bc.client == nil {
		return fmt.Errorf("broker client not connected")
	}

	var lastErr error
	for attempt := 1; attempt <= brokerHealthCheckAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(bc.ctx, brokerHealthCheckTimeout)
		_, err := bc.client.FindBrokerLeader(ctx, &mq_pb.FindBrokerLeaderRequest{})
		cancel()
		if err == nil {
			return nil
		}
		lastErr = err
		if attempt == brokerHealthCheckAttempts {
			break
		}
		select {
		case <-bc.ctx.Done():
			return bc.ctx.Err()
		case <-time.After(time.Duration(attempt) * brokerHealthCheckBackoff):
		}
	}

	return fmt.Errorf("broker health check failed after %d attempts: %v", brokerHealthCheckAttempts, lastErr)
}

// GetPartitionRangeInfo gets comprehensive range information from HanzoMQ broker's native range manager
func (bc *BrokerClient) GetPartitionRangeInfo(topic string, partition int32) (*PartitionRangeInfo, error) {

	if bc.client == nil {
		return nil, fmt.Errorf("broker client not connected")
	}

	// Get the actual partition assignment from the broker instead of hardcoding
	pbTopic := &schema_pb.Topic{
		Namespace: "kafka",
		Name:      topic,
	}

	// Get the actual partition assignment for this Kafka partition
	actualPartition, err := bc.getActualPartitionAssignment(topic, partition)
	if err != nil {
		return nil, fmt.Errorf("failed to get actual partition assignment: %v", err)
	}

	// Call the broker's gRPC method
	resp, err := bc.client.GetPartitionRangeInfo(context.Background(), &mq_pb.GetPartitionRangeInfoRequest{
		Topic:     pbTopic,
		Partition: actualPartition,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get partition range info from broker: %v", err)
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("broker error: %s", resp.Error)
	}

	// Extract offset range information
	var earliestOffset, latestOffset, highWaterMark int64
	if resp.OffsetRange != nil {
		earliestOffset = resp.OffsetRange.EarliestOffset
		latestOffset = resp.OffsetRange.LatestOffset
		highWaterMark = resp.OffsetRange.HighWaterMark
	}

	// Extract timestamp range information
	var earliestTimestampNs, latestTimestampNs int64
	if resp.TimestampRange != nil {
		earliestTimestampNs = resp.TimestampRange.EarliestTimestampNs
		latestTimestampNs = resp.TimestampRange.LatestTimestampNs
	}

	info := &PartitionRangeInfo{
		EarliestOffset:      earliestOffset,
		LatestOffset:        latestOffset,
		HighWaterMark:       highWaterMark,
		EarliestTimestampNs: earliestTimestampNs,
		LatestTimestampNs:   latestTimestampNs,
		RecordCount:         resp.RecordCount,
		ActiveSubscriptions: resp.ActiveSubscriptions,
	}

	return info, nil
}

// GetHighWaterMark gets the high water mark for a topic partition
func (bc *BrokerClient) GetHighWaterMark(topic string, partition int32) (int64, error) {

	// Primary approach: Use HanzoMQ's native range manager via gRPC
	info, err := bc.GetPartitionRangeInfo(topic, partition)
	if err != nil {
		// Fallback to chunk metadata approach
		highWaterMark, err := bc.getHighWaterMarkFromChunkMetadata(topic, partition)
		if err != nil {
			return 0, err
		}
		return highWaterMark, nil
	}

	return info.HighWaterMark, nil
}

// GetEarliestOffset gets the earliest offset from HanzoMQ broker's native offset manager
func (bc *BrokerClient) GetEarliestOffset(topic string, partition int32) (int64, error) {

	// Primary approach: Use HanzoMQ's native range manager via gRPC
	info, err := bc.GetPartitionRangeInfo(topic, partition)
	if err != nil {
		// Fallback to chunk metadata approach
		earliestOffset, err := bc.getEarliestOffsetFromChunkMetadata(topic, partition)
		if err != nil {
			return 0, err
		}
		return earliestOffset, nil
	}

	return info.EarliestOffset, nil
}

// getOffsetRangeFromChunkMetadata reads chunk metadata to find both earliest and latest offsets
func (bc *BrokerClient) getOffsetRangeFromChunkMetadata(topic string, partition int32) (earliestOffset int64, highWaterMark int64, err error) {
	if bc.filerClientAccessor == nil {
		return 0, 0, fmt.Errorf("filer client not available")
	}

	// Get the topic path and find the latest version
	topicPath := fmt.Sprintf("/topics/kafka/%s", topic)

	// First, list the topic versions to find the latest
	var latestVersion string
	err = bc.filerClientAccessor.WithFilerClient(false, func(client filer_pb.HanzoFilerClient) error {
		stream, err := client.ListEntries(context.Background(), &filer_pb.ListEntriesRequest{
			Directory: topicPath,
		})
		if err != nil {
			return err
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if resp.Entry != nil && resp.Entry.IsDirectory && strings.HasPrefix(resp.Entry.Name, "v") {
				if latestVersion == "" || resp.Entry.Name > latestVersion {
					latestVersion = resp.Entry.Name
				}
			}
		}
		return nil
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to list topic versions: %v", err)
	}

	if latestVersion == "" {
		return 0, 0, nil
	}

	// Find the partition directory
	versionPath := fmt.Sprintf("%s/%s", topicPath, latestVersion)
	var partitionDir string
	err = bc.filerClientAccessor.WithFilerClient(false, func(client filer_pb.HanzoFilerClient) error {
		stream, err := client.ListEntries(context.Background(), &filer_pb.ListEntriesRequest{
			Directory: versionPath,
		})
		if err != nil {
			return err
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if resp.Entry != nil && resp.Entry.IsDirectory && strings.Contains(resp.Entry.Name, "-") {
				// Parse partition range (format: NNNN-NNNN) and match requested partition
				var pStart, pStop int32
				if n, scanErr := fmt.Sscanf(resp.Entry.Name, "%04d-%04d", &pStart, &pStop); n == 2 && scanErr == nil {
					if partition >= pStart && partition < pStop {
						partitionDir = resp.Entry.Name
						break
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to list partition directories: %v", err)
	}

	if partitionDir == "" {
		return 0, 0, nil
	}

	// Scan all message files to find the highest offset_max and lowest offset_min
	partitionPath := fmt.Sprintf("%s/%s", versionPath, partitionDir)
	highWaterMark = 0
	earliestOffset = -1 // -1 indicates no data found yet

	err = bc.filerClientAccessor.WithFilerClient(false, func(client filer_pb.HanzoFilerClient) error {
		stream, err := client.ListEntries(context.Background(), &filer_pb.ListEntriesRequest{
			Directory: partitionPath,
		})
		if err != nil {
			return err
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if resp.Entry != nil && !resp.Entry.IsDirectory && resp.Entry.Name != "checkpoint.offset" {
				// Check for offset ranges in Extended attributes (both log files and parquet files)
				if resp.Entry.Extended != nil {
					// Track maximum offset for high water mark
					if maxOffsetBytes, exists := resp.Entry.Extended[mq.ExtendedAttrOffsetMax]; exists && len(maxOffsetBytes) == 8 {
						maxOffset := int64(binary.BigEndian.Uint64(maxOffsetBytes))
						if maxOffset > highWaterMark {
							highWaterMark = maxOffset
						}
					}

					// Track minimum offset for earliest offset
					if minOffsetBytes, exists := resp.Entry.Extended[mq.ExtendedAttrOffsetMin]; exists && len(minOffsetBytes) == 8 {
						minOffset := int64(binary.BigEndian.Uint64(minOffsetBytes))
						if earliestOffset == -1 || minOffset < earliestOffset {
							earliestOffset = minOffset
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to scan message files: %v", err)
	}

	// High water mark is the next offset after the highest written offset
	if highWaterMark > 0 {
		highWaterMark++
	}

	// If no data found, set earliest offset to 0
	if earliestOffset == -1 {
		earliestOffset = 0
	}

	return earliestOffset, highWaterMark, nil
}

// getHighWaterMarkFromChunkMetadata is a wrapper for backward compatibility
func (bc *BrokerClient) getHighWaterMarkFromChunkMetadata(topic string, partition int32) (int64, error) {
	_, highWaterMark, err := bc.getOffsetRangeFromChunkMetadata(topic, partition)
	return highWaterMark, err
}

// getEarliestOffsetFromChunkMetadata gets the earliest offset from chunk metadata (fallback)
func (bc *BrokerClient) getEarliestOffsetFromChunkMetadata(topic string, partition int32) (int64, error) {
	earliestOffset, _, err := bc.getOffsetRangeFromChunkMetadata(topic, partition)
	return earliestOffset, err
}

// GetFilerAddress returns the first filer address used by this broker client (for backward compatibility)
func (bc *BrokerClient) GetFilerAddress() string {
	if bc.filerClientAccessor != nil && bc.filerClientAccessor.GetFilers != nil {
		filers := bc.filerClientAccessor.GetFilers()
		if len(filers) > 0 {
			return string(filers[0])
		}
	}
	return ""
}

// Delegate methods to the shared filer client accessor
func (bc *BrokerClient) WithFilerClient(streamingMode bool, fn func(client filer_pb.HanzoFilerClient) error) error {
	return bc.filerClientAccessor.WithFilerClient(streamingMode, fn)
}

func (bc *BrokerClient) GetFilers() []pb.ServerAddress {
	return bc.filerClientAccessor.GetFilers()
}

func (bc *BrokerClient) GetGrpcDialOption() pb.DialOption {
	return bc.filerClientAccessor.GetGrpcDialOption()
}

// ListTopics gets all topics from HanzoMQ broker (includes in-memory topics)
func (bc *BrokerClient) ListTopics() ([]string, error) {
	if bc.client == nil {
		return nil, fmt.Errorf("broker client not connected")
	}

	ctx, cancel := context.WithTimeout(bc.ctx, 5*time.Second)
	defer cancel()

	resp, err := bc.client.ListTopics(ctx, &mq_pb.ListTopicsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list topics from broker: %v", err)
	}

	var topics []string
	for _, topic := range resp.Topics {
		// Filter for kafka namespace topics
		if topic.Namespace == "kafka" {
			topics = append(topics, topic.Name)
		}
	}

	return topics, nil
}

// GetTopicConfiguration gets topic configuration including partition count from the broker
func (bc *BrokerClient) GetTopicConfiguration(topicName string) (*mq_pb.GetTopicConfigurationResponse, error) {
	if bc.client == nil {
		return nil, fmt.Errorf("broker client not connected")
	}

	ctx, cancel := context.WithTimeout(bc.ctx, 5*time.Second)
	defer cancel()

	resp, err := bc.client.GetTopicConfiguration(ctx, &mq_pb.GetTopicConfigurationRequest{
		Topic: &schema_pb.Topic{
			Namespace: "kafka",
			Name:      topicName,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get topic configuration from broker: %v", err)
	}

	return resp, nil
}

// TopicExists checks if a topic exists in HanzoMQ broker (includes in-memory topics)
func (bc *BrokerClient) TopicExists(topicName string) (bool, error) {
	if bc.client == nil {
		return false, fmt.Errorf("broker client not connected")
	}

	ctx, cancel := context.WithTimeout(bc.ctx, 5*time.Second)
	defer cancel()

	glog.V(2).Infof("[BrokerClient] TopicExists: Querying broker for topic %s", topicName)
	resp, err := bc.client.TopicExists(ctx, &mq_pb.TopicExistsRequest{
		Topic: &schema_pb.Topic{
			Namespace: "kafka",
			Name:      topicName,
		},
	})
	if err != nil {
		glog.V(1).Infof("[BrokerClient] TopicExists: ERROR for topic %s: %v", topicName, err)
		return false, fmt.Errorf("failed to check topic existence: %v", err)
	}

	glog.V(2).Infof("[BrokerClient] TopicExists: Topic %s exists=%v", topicName, resp.Exists)
	return resp.Exists, nil
}
