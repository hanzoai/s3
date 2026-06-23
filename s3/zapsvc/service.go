// Package zapsvc is S3's native ZAP service — the interface internal Hanzo
// services use to talk to S3 over the mesh (zero-copy ZAP, never HTTP). The
// external HTTPS S3 API stays only for outside S3 clients (aws-sdk/rclone);
// in-cloud callers (team-go, cloud, commerce, …) use this.
//
// It is the reference implementation of the platform law (internal = ZAP only):
// requests/responses are zero-copy ZAP schemas (s3/wire), carried natively over
// the zaprpc transport, dispatched by schema Kind. Every other Hanzo service
// exposes itself with this same shape.
package zapsvc

import (
	"context"

	"github.com/hanzoai/s3/s3/wire"
	"github.com/hanzoai/s3/s3/zaprpc"
)

// ObjectStore is the backend the ZAP service delegates to. The real S3 engine
// (filer + volumes) implements this; tests use an in-memory map. Keeping the
// service layer behind this interface is what lets the ZAP surface land before
// the internal cluster-transport rip finishes.
type ObjectStore interface {
	Get(bucket, key string) (data []byte, contentType, etag string, err error)
	Put(bucket, key string, data []byte, contentType string) (etag string, err error)
}

// Register binds the S3 procedures onto a ZAP server. Handlers Wrap the request
// buffer into its typed View (zero-copy), call the store, and build the
// response schema directly — no struct, no marshal.
func Register(s *zaprpc.Server, store ObjectStore) {
	s.Handle(uint8(wire.KindGetObjectRequest), func(_ context.Context, reqBuf []byte) ([]byte, error) {
		req, err := wire.WrapGetObjectRequest(reqBuf)
		if err != nil {
			return nil, err
		}
		data, ct, etag, err := store.Get(wire.GetObjectRequestBucket(req), wire.GetObjectRequestKey(req))
		if err != nil {
			return nil, err
		}
		_, buf := wire.NewGetObjectResponse(data, ct, etag)
		return buf, nil
	})
	s.Handle(uint8(wire.KindPutObjectRequest), func(_ context.Context, reqBuf []byte) ([]byte, error) {
		req, err := wire.WrapPutObjectRequest(reqBuf)
		if err != nil {
			return nil, err
		}
		etag, err := store.Put(
			wire.PutObjectRequestBucket(req), wire.PutObjectRequestKey(req),
			wire.PutObjectRequestData(req), wire.PutObjectRequestContentType(req))
		if err != nil {
			return nil, err
		}
		_, buf := wire.NewPutObjectResponse(etag)
		return buf, nil
	})
}

// Client is the typed S3 ZAP service client internal services hold.
type Client struct {
	rpc  *zaprpc.Client
	addr string
}

// NewClient targets the S3 service at addr (e.g. "s3.hanzo.svc:<zapport>").
func NewClient(rpc *zaprpc.Client, addr string) *Client {
	return &Client{rpc: rpc, addr: addr}
}

// GetObject fetches an object over ZAP.
func (c *Client) GetObject(ctx context.Context, bucket, key string) (data []byte, contentType, etag string, err error) {
	_, reqBuf := wire.NewGetObjectRequest(bucket, key)
	respBuf, err := c.rpc.Call(ctx, c.addr, uint8(wire.KindGetObjectRequest), reqBuf)
	if err != nil {
		return nil, "", "", err
	}
	resp, err := wire.WrapGetObjectResponse(respBuf)
	if err != nil {
		return nil, "", "", err
	}
	return wire.GetObjectResponseData(resp), wire.GetObjectResponseContentType(resp), wire.GetObjectResponseEtag(resp), nil
}

// PutObject stores an object over ZAP, returning its etag.
func (c *Client) PutObject(ctx context.Context, bucket, key string, data []byte, contentType string) (etag string, err error) {
	_, reqBuf := wire.NewPutObjectRequest(bucket, key, data, contentType)
	respBuf, err := c.rpc.Call(ctx, c.addr, uint8(wire.KindPutObjectRequest), reqBuf)
	if err != nil {
		return "", err
	}
	resp, err := wire.WrapPutObjectResponse(respBuf)
	if err != nil {
		return "", err
	}
	return wire.PutObjectResponseEtag(resp), nil
}
