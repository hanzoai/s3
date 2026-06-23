package zaprpc

import (
	"context"
	"sync"

	"github.com/luxfi/zap"
)

// Direct-addressed ZAP RPC transport for Hanzo S3, replacing gRPC. It matches
// the master↔volume↔filer topology exactly: a server listens on a fixed port,
// a client dials a known host:port (no mDNS discovery). The whole application
// message rides as a single opaque bytes field inside the zap object; routing
// is by the message's flags high byte (the procedure code). So ZAP stays pure
// transport and our TLV codec (codec.go) owns the payload contents.
//
// Streaming RPCs are handled separately (stream.go, F4); this file is unary
// request/response.

const (
	payloadField   = 0 // the single bytes slot carrying the TLV payload
	payloadObjSize = 8 // one {relOffset u32, length u32} bytes pointer
)

// frame wraps a TLV payload into a zap.Message tagged with code for routing.
func frame(code uint8, payload []byte) (*zap.Message, error) {
	b := zap.NewBuilder(len(payload) + 32)
	ob := b.StartObject(payloadObjSize)
	ob.SetBytes(payloadField, payload)
	ob.FinishAsRoot()
	return zap.Parse(b.FinishWithFlags(uint16(code) << 8))
}

// unframe extracts the TLV payload from a received zap.Message.
func unframe(msg *zap.Message) []byte { return msg.Root().Bytes(payloadField) }

// Handler processes a request payload and returns a response payload. The
// generated per-service server stubs decode req via UnmarshalZap, dispatch to
// the service implementation, and encode the response via MarshalZap.
type Handler func(ctx context.Context, req []byte) ([]byte, error)

// Server is a direct-addressed ZAP RPC server (discovery disabled).
type Server struct{ node *zap.Node }

// NewServer creates a server listening on port (0 = OS-assigned).
func NewServer(nodeID string, port int) *Server {
	return &Server{node: zap.NewNode(zap.NodeConfig{NodeID: nodeID, Port: port, NoDiscovery: true})}
}

// Handle registers h for procedure code (1..255). Last registration wins.
func (s *Server) Handle(code uint8, h Handler) {
	s.node.Handle(uint16(code), func(ctx context.Context, _ string, msg *zap.Message) (*zap.Message, error) {
		resp, err := h(ctx, unframe(msg))
		if err != nil {
			return nil, err
		}
		return frame(code, resp)
	})
}

// Start binds the listener and begins accepting. Stop releases it.
func (s *Server) Start() error { return s.node.Start() }
func (s *Server) Stop()        { s.node.Stop() }

// Client is a direct-addressed ZAP RPC client. It memoizes the NodeID learned
// from each address's handshake so repeated calls skip the dial.
type Client struct {
	node    *zap.Node
	mu      sync.Mutex
	peerIDs map[string]string // addr -> learned NodeID
}

// NewClient creates and starts a client node.
func NewClient(nodeID string) (*Client, error) {
	n := zap.NewNode(zap.NodeConfig{NodeID: nodeID, Port: 0, NoDiscovery: true})
	if err := n.Start(); err != nil {
		return nil, err
	}
	return &Client{node: n, peerIDs: map[string]string{}}, nil
}

func (c *Client) resolve(addr string) (string, error) {
	c.mu.Lock()
	id, ok := c.peerIDs[addr]
	c.mu.Unlock()
	if ok {
		return id, nil
	}
	id, err := c.node.ConnectDirectID(addr)
	if err != nil {
		return "", err
	}
	c.mu.Lock()
	c.peerIDs[addr] = id
	c.mu.Unlock()
	return id, nil
}

// Call issues a unary RPC to addr for procedure code, returning the response
// payload. req/resp are TLV bytes produced/consumed by the codec.
func (c *Client) Call(ctx context.Context, addr string, code uint8, req []byte) ([]byte, error) {
	id, err := c.resolve(addr)
	if err != nil {
		return nil, err
	}
	reqMsg, err := frame(code, req)
	if err != nil {
		return nil, err
	}
	resp, err := c.node.Call(ctx, id, reqMsg)
	if err != nil {
		return nil, err
	}
	return unframe(resp), nil
}

// Stop releases the client node.
func (c *Client) Stop() { c.node.Stop() }
