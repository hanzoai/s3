// Package object is S3's native ZAP service — the interface internal Hanzo
// services use to talk to S3 over the mesh (zero-copy ZAP, never HTTP). The
// external HTTPS S3 API stays only for outside S3 clients (aws-sdk/rclone);
// in-cloud callers (team-go, cloud, commerce, …) use this.
//
// It is the reference implementation of the platform law (internal = ZAP only):
// requests/responses are zero-copy ZAP schemas (s3/wire/object), carried over
// the canonical github.com/zap-proto/go transport and dispatched by the
// generated HanzoS3Object service — the SAME wire hanzo and lux share. No
// struct, no marshal, no protobuf, no per-service one-off transport.
package object

import (
	"crypto/tls"

	objectwire "github.com/hanzoai/s3/s3/wire/object"
	"github.com/zap-proto/go/transport"
)

// ObjectStore is the backend the ZAP service delegates to. The real S3 engine
// (filer + volumes) implements this (FilerStore); tests use an in-memory map.
type ObjectStore interface {
	Get(bucket, key string) (data []byte, contentType, etag string, err error)
	Put(bucket, key string, data []byte, contentType string) (etag string, err error)
}

// handler adapts an ObjectStore to the generated HanzoS3ObjectHandler: it Wraps
// each request buffer into its zero-copy view, calls the store, and builds the
// response buffer directly.
type handler struct{ store ObjectStore }

func (h handler) GetObject(req []byte) ([]byte, error) {
	v, err := objectwire.WrapGetObjectRequest(req)
	if err != nil {
		return nil, err
	}
	data, ct, etag, err := h.store.Get(v.Bucket(), v.Key())
	if err != nil {
		return nil, err
	}
	return objectwire.NewGetObjectResponse(objectwire.GetObjectResponseInput{
		Data: data, ContentType: ct, Etag: etag,
	}), nil
}

func (h handler) PutObject(req []byte) ([]byte, error) {
	v, err := objectwire.WrapPutObjectRequest(req)
	if err != nil {
		return nil, err
	}
	etag, err := h.store.Put(v.Bucket(), v.Key(), v.Data(), v.ContentType())
	if err != nil {
		return nil, err
	}
	return objectwire.NewPutObjectResponse(objectwire.PutObjectResponseInput{Etag: etag}), nil
}

// Dispatch is the HanzoS3Object server dispatch bound to store — pass it to
// transport.Listen. Exposed so callers can compose it (e.g. behind TLS).
func Dispatch(store ObjectStore) transport.Dispatch {
	h := handler{store: store}
	return func(envelope []byte) ([]byte, error) {
		return objectwire.DispatchHanzoS3Object(h, envelope)
	}
}

// ServeTLS starts the native ZAP S3 service over PQ-secured TLS — the default
// for the production mesh. conf must carry the service certificate (from KMS);
// it is wrapped with transport.PQTLSConfig so the X25519MLKEM768 hybrid (PQ
// X-Wing) is REQUIRED — a classical-only peer fails the handshake, no downgrade.
func ServeTLS(network, addr string, store ObjectStore, conf *tls.Config) (*transport.Server, error) {
	return transport.ListenTLS(network, addr, conf, Dispatch(store))
}

// Serve starts the native ZAP S3 service over plaintext on network/addr (e.g.
// "tcp", ":18901"), backed by store. This is the loopback/test path — the
// production mesh uses ServeTLS (PQ X-Wing). Returns the running server.
func Serve(network, addr string, store ObjectStore) (*transport.Server, error) {
	return transport.Listen(network, addr, Dispatch(store))
}

// Client is the typed S3 ZAP service client internal services hold. It owns the
// transport connection to one S3 endpoint.
type Client struct {
	conn transport.Conn
	rpc  *objectwire.HanzoS3ObjectClient
}

// DialTLS connects to the S3 ZAP service at addr (e.g. "s3.hanzo.svc:18901")
// over PQ-secured TLS — the default for the production mesh. conf carries the
// client trust (RootCAs from KMS); it is wrapped with transport.PQTLSConfig so
// the session key rides X25519MLKEM768 (PQ X-Wing) or the handshake fails.
func DialTLS(network, addr string, conf *tls.Config) (*Client, error) {
	conn, err := transport.DialTLS(network, addr, conf)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, rpc: objectwire.NewHanzoS3ObjectClient(conn, nil)}, nil
}

// Dial connects to the S3 ZAP service over plaintext TCP — the loopback/test
// path. The production mesh uses DialTLS (PQ X-Wing).
func Dial(network, addr string) (*Client, error) {
	conn, err := transport.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, rpc: objectwire.NewHanzoS3ObjectClient(conn, nil)}, nil
}

// NewClient wraps an already-established transport.Conn (TCP, Unix, or PQ-TLS).
func NewClient(conn transport.Conn) *Client {
	return &Client{conn: conn, rpc: objectwire.NewHanzoS3ObjectClient(conn, nil)}
}

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

// GetObject fetches an object over ZAP.
func (c *Client) GetObject(bucket, key string) (data []byte, contentType, etag string, err error) {
	req := objectwire.NewGetObjectRequest(objectwire.GetObjectRequestInput{Bucket: bucket, Key: key})
	_, body, err := c.rpc.GetObject(req)
	if err != nil {
		return nil, "", "", err
	}
	resp, err := objectwire.WrapGetObjectResponse(body)
	if err != nil {
		return nil, "", "", err
	}
	return resp.Data(), resp.ContentType(), resp.Etag(), nil
}

// PutObject stores an object over ZAP, returning its etag.
func (c *Client) PutObject(bucket, key string, data []byte, contentType string) (etag string, err error) {
	req := objectwire.NewPutObjectRequest(objectwire.PutObjectRequestInput{
		Bucket: bucket, Key: key, Data: data, ContentType: contentType,
	})
	_, body, err := c.rpc.PutObject(req)
	if err != nil {
		return "", err
	}
	resp, err := objectwire.WrapPutObjectResponse(body)
	if err != nil {
		return "", err
	}
	return resp.Etag(), nil
}
