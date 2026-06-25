// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mqzap

import (
	"github.com/zap-proto/go/transport"

	"github.com/hanzoai/s3/s3/security"
	"github.com/hanzoai/s3/s3/util"
)

// DialBroker dials a broker's ZAP endpoint with the broker-mesh transport
// security: PQ-mTLS (X25519MLKEM768, applied by transport.DialTLS) when
// grpc.msg_broker certs are configured, plaintext otherwise — the SAME gate as
// pb.dialBrokerZapAddr and kafka.dialBrokerZap. The broker→follower replication
// legs (PublishFollowMe / SubscribeFollowMe) dial through here so follower
// fan-out keeps the mTLS-when-configured contract instead of silently dialing
// plaintext (which, against a PQ-mTLS follower listener, would fail the
// handshake and break replication; against a plaintext one, would ship payloads
// in cleartext).
func DialBroker(addr string) (transport.Conn, error) {
	if cfg := security.ClientTLSConfig(util.GetViper(), "grpc.msg_broker"); cfg != nil {
		return transport.DialTLS("tcp", addr, cfg)
	}
	return transport.Dial("tcp", addr)
}
