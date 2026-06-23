package zaprpc

import (
	"context"
	"encoding/binary"
	"sync"

	"github.com/luxfi/zap"
)

// Direct-addressed ZAP RPC transport for the Hanzo service mesh — the east-west
// substrate for the platform law (internal service↔service = ZAP only; external
// = HTTPS/WS at the edge). It carries ZAP schema buffers NATIVELY: a request is
// a built schema (NewX → buffer); the buffer IS the wire message, sent as-is.
// No codec, no framing, no copy. Routing is by the schema's Kind, stamped into
// the message flags' high byte (the byte zap.Node dispatches on).
//
// Matches the master↔volume↔filer + service↔service topology: a server listens
// on a known port, a client dials a known host:port (no mDNS).

// routeByKind stamps the schema Kind into the message flags (data[6:8] hi byte)
// so zap.Node routes to the matching Handle(kind). The schema's own Kind byte
// lives at object offset 0; this mirrors it into the transport header.
func routeByKind(buf []byte, kind uint8) {
	if len(buf) >= 8 {
		binary.LittleEndian.PutUint16(buf[6:8], uint16(kind)<<8)
	}
}

// Handler processes a request schema buffer and returns a response schema
// buffer. Both are zero-copy ZAP buffers (Wrap on the way in, New on the way
// out) — never deserialized into structs.
type Handler func(ctx context.Context, reqBuf []byte) (respBuf []byte, err error)

// Server exposes a service over ZAP (discovery disabled; direct-addressed).
type Server struct{ node *zap.Node }

// NewServer listens on port (0 = OS-assigned).
func NewServer(nodeID string, port int) *Server {
	return &Server{node: zap.NewNode(zap.NodeConfig{NodeID: nodeID, Port: port, NoDiscovery: true})}
}

// Handle binds h to a schema Kind (the request type's discriminator).
func (s *Server) Handle(kind uint8, h Handler) {
	s.node.Handle(uint16(kind), func(ctx context.Context, _ string, msg *zap.Message) (*zap.Message, error) {
		respBuf, err := h(ctx, msg.Bytes())
		if err != nil {
			return nil, err
		}
		return zap.WrapBuffer(respBuf), nil
	})
}

func (s *Server) Start() error { return s.node.Start() }
func (s *Server) Stop()        { s.node.Stop() }

// Client invokes services over ZAP. It memoizes each address's handshake-
// learned NodeID so repeated calls skip the dial.
type Client struct {
	node  *zap.Node
	mu    sync.Mutex
	peers map[string]string // addr -> NodeID
}

// NewClient creates and starts a client node.
func NewClient(nodeID string) (*Client, error) {
	n := zap.NewNode(zap.NodeConfig{NodeID: nodeID, Port: 0, NoDiscovery: true})
	if err := n.Start(); err != nil {
		return nil, err
	}
	return &Client{node: n, peers: map[string]string{}}, nil
}

func (c *Client) resolve(addr string) (string, error) {
	c.mu.Lock()
	id, ok := c.peers[addr]
	c.mu.Unlock()
	if ok {
		return id, nil
	}
	id, err := c.node.ConnectDirectID(addr)
	if err != nil {
		return "", err
	}
	c.mu.Lock()
	c.peers[addr] = id
	c.mu.Unlock()
	return id, nil
}

// Call sends a request schema buffer (of type kind) to addr and returns the
// response schema buffer. The caller Wraps the response into its typed View.
func (c *Client) Call(ctx context.Context, addr string, kind uint8, reqBuf []byte) ([]byte, error) {
	id, err := c.resolve(addr)
	if err != nil {
		return nil, err
	}
	routeByKind(reqBuf, kind)
	resp, err := c.node.Call(ctx, id, zap.WrapBuffer(reqBuf))
	if err != nil {
		return nil, err
	}
	return resp.Bytes(), nil
}

func (c *Client) Stop() { c.node.Stop() }
