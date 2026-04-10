// Copyright (c) 2026 Hanzo AI, Inc.
//
// This file is part of Hanzo S3 Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/hanzoai/s3/internal/hash"
	"github.com/hanzoai/s3/internal/logger"
	"github.com/luxfi/zap"
	"github.com/minio/pkg/v3/env"
)

// ZAP opcodes for S3 operations.
const (
	zapOpPut         uint16 = 0x01
	zapOpGet         uint16 = 0x02
	zapOpDelete      uint16 = 0x03
	zapOpList        uint16 = 0x04
	zapOpHead        uint16 = 0x05
	zapOpCreateBucket uint16 = 0x06
	zapOpListBuckets uint16 = 0x07
	zapOpHealth      uint16 = 0x08
)

// ZAP request field offsets.
// All requests share: bucket (text @0), key (text @8).
// Put additionally carries: data (bytes @16).
// List additionally carries: prefix (text @16).
const (
	zapFieldBucket = 0  // text: offset(4)+len(4) = 8 bytes
	zapFieldKey    = 8  // text: offset(4)+len(4) = 8 bytes
	zapFieldData   = 16 // bytes: offset(4)+len(4) = 8 bytes
	zapFieldPrefix = 16 // text: offset(4)+len(4) = 8 bytes (alias for list)
	zapReqSize     = 24 // minimum object size for requests with 3 fields
)

// ZAP response field offsets.
const (
	zapRespStatus  = 0  // uint32: 0=ok, non-zero=error
	zapRespMessage = 4  // text: error message or metadata
	zapRespData    = 12 // bytes: object data (get) or JSON list
	zapRespSize    = 20
)

const defaultZAPPort = "9002"

// startZAPServer starts the ZAP listener for inter-service object operations.
func startZAPServer() {
	port := env.Get("S3_ZAP_PORT", defaultZAPPort)

	portNum := 9002
	if _, err := fmt.Sscanf(port, "%d", &portNum); err != nil {
		logger.LogIf(GlobalContext, "zap", fmt.Errorf("invalid S3_ZAP_PORT=%q, using default %s", port, defaultZAPPort))
		portNum = 9002
	}

	node := zap.NewNode(zap.NodeConfig{
		NodeID:      fmt.Sprintf("s3-%s", globalLocalNodeName),
		ServiceType: "_s3._tcp",
		Port:        portNum,
		NoDiscovery: true, // Internal services connect directly
	})

	node.Handle(zapOpPut, zapHandlePut)
	node.Handle(zapOpGet, zapHandleGet)
	node.Handle(zapOpDelete, zapHandleDelete)
	node.Handle(zapOpList, zapHandleList)
	node.Handle(zapOpHead, zapHandleHead)
	node.Handle(zapOpCreateBucket, zapHandleCreateBucket)
	node.Handle(zapOpListBuckets, zapHandleListBuckets)
	node.Handle(zapOpHealth, zapHandleHealth)

	if err := node.Start(); err != nil {
		logger.LogIf(GlobalContext, "zap", fmt.Errorf("failed to start ZAP server on port %d: %w", portNum, err))
		return
	}

	logger.Info("ZAP transport started on :%d", portNum)
}

func zapHandlePut(ctx context.Context, from string, msg *zap.Message) (*zap.Message, error) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		return zapErrorResp("object layer not ready"), nil
	}

	root := msg.Root()
	bucket := root.Text(zapFieldBucket)
	key := root.Text(zapFieldKey)
	data := root.Bytes(zapFieldData)

	if bucket == "" || key == "" {
		return zapErrorResp("bucket and key required"), nil
	}

	reader := bytes.NewReader(data)
	hashReader, err := hash.NewReader(ctx, reader, int64(len(data)), "", "", int64(len(data)))
	if err != nil {
		return zapErrorResp(err.Error()), nil
	}

	_, err = objAPI.PutObject(ctx, bucket, key, NewPutObjReader(hashReader), ObjectOptions{
		UserDefined: map[string]string{
			"X-S3-Source": "zap/" + from,
		},
	})
	if err != nil {
		return zapErrorResp(err.Error()), nil
	}

	return zapOKResp(nil), nil
}

func zapHandleGet(ctx context.Context, from string, msg *zap.Message) (*zap.Message, error) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		return zapErrorResp("object layer not ready"), nil
	}

	root := msg.Root()
	bucket := root.Text(zapFieldBucket)
	key := root.Text(zapFieldKey)

	if bucket == "" || key == "" {
		return zapErrorResp("bucket and key required"), nil
	}

	gr, err := objAPI.GetObjectNInfo(ctx, bucket, key, nil, nil, ObjectOptions{})
	if err != nil {
		return zapErrorResp(err.Error()), nil
	}
	defer gr.Close()

	data, err := io.ReadAll(gr)
	if err != nil {
		return zapErrorResp(err.Error()), nil
	}

	return zapOKResp(data), nil
}

func zapHandleDelete(ctx context.Context, from string, msg *zap.Message) (*zap.Message, error) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		return zapErrorResp("object layer not ready"), nil
	}

	root := msg.Root()
	bucket := root.Text(zapFieldBucket)
	key := root.Text(zapFieldKey)

	if bucket == "" || key == "" {
		return zapErrorResp("bucket and key required"), nil
	}

	_, err := objAPI.DeleteObject(ctx, bucket, key, ObjectOptions{})
	if err != nil {
		return zapErrorResp(err.Error()), nil
	}

	return zapOKResp(nil), nil
}

func zapHandleList(ctx context.Context, from string, msg *zap.Message) (*zap.Message, error) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		return zapErrorResp("object layer not ready"), nil
	}

	root := msg.Root()
	bucket := root.Text(zapFieldBucket)
	prefix := root.Text(zapFieldPrefix)

	if bucket == "" {
		return zapErrorResp("bucket required"), nil
	}

	result, err := objAPI.ListObjects(ctx, bucket, prefix, "", "", 1000)
	if err != nil {
		return zapErrorResp(err.Error()), nil
	}

	// Encode object names as newline-separated text.
	var sb strings.Builder
	for _, obj := range result.Objects {
		sb.WriteString(obj.Name)
		sb.WriteByte('\n')
	}

	return zapOKResp([]byte(sb.String())), nil
}

func zapHandleHead(ctx context.Context, from string, msg *zap.Message) (*zap.Message, error) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		return zapErrorResp("object layer not ready"), nil
	}

	root := msg.Root()
	bucket := root.Text(zapFieldBucket)
	key := root.Text(zapFieldKey)

	if bucket == "" || key == "" {
		return zapErrorResp("bucket and key required"), nil
	}

	info, err := objAPI.GetObjectInfo(ctx, bucket, key, ObjectOptions{})
	if err != nil {
		return zapErrorResp(err.Error()), nil
	}

	// Return size and content-type as text metadata.
	meta := fmt.Sprintf("size=%d\ncontent-type=%s\nmodified=%s",
		info.Size, info.ContentType, info.ModTime.UTC().Format("2006-01-02T15:04:05Z"))

	return zapOKRespText(meta), nil
}

func zapHandleCreateBucket(ctx context.Context, from string, msg *zap.Message) (*zap.Message, error) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		return zapErrorResp("object layer not ready"), nil
	}

	root := msg.Root()
	bucket := root.Text(zapFieldBucket)

	if bucket == "" {
		return zapErrorResp("bucket required"), nil
	}

	err := objAPI.MakeBucket(ctx, bucket, MakeBucketOptions{})
	if err != nil {
		return zapErrorResp(err.Error()), nil
	}

	return zapOKResp(nil), nil
}

func zapHandleListBuckets(ctx context.Context, from string, msg *zap.Message) (*zap.Message, error) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		return zapErrorResp("object layer not ready"), nil
	}

	buckets, err := objAPI.ListBuckets(ctx, BucketOptions{})
	if err != nil {
		return zapErrorResp(err.Error()), nil
	}

	var sb strings.Builder
	for _, b := range buckets {
		sb.WriteString(b.Name)
		sb.WriteByte('\n')
	}

	return zapOKResp([]byte(sb.String())), nil
}

func zapHandleHealth(ctx context.Context, from string, msg *zap.Message) (*zap.Message, error) {
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		return zapErrorResp("not ready"), nil
	}

	result := objAPI.Health(ctx, HealthOptions{})
	if !result.Healthy {
		return zapErrorResp("unhealthy"), nil
	}

	return zapOKResp(nil), nil
}

// zapOKResp builds a success response with optional data payload.
func zapOKResp(data []byte) *zap.Message {
	b := zap.NewBuilder(zap.HeaderSize + zapRespSize + len(data) + 64)
	obj := b.StartObject(zapRespSize)
	obj.SetUint32(zapRespStatus, 0)
	if len(data) > 0 {
		obj.SetBytes(zapRespData, data)
	}
	obj.FinishAsRoot()
	msg, _ := zap.Parse(b.Finish())
	return msg
}

// zapOKRespText builds a success response with a text message.
func zapOKRespText(text string) *zap.Message {
	b := zap.NewBuilder(zap.HeaderSize + zapRespSize + len(text) + 64)
	obj := b.StartObject(zapRespSize)
	obj.SetUint32(zapRespStatus, 0)
	obj.SetText(zapRespMessage, text)
	obj.FinishAsRoot()
	msg, _ := zap.Parse(b.Finish())
	return msg
}

// zapErrorResp builds an error response.
func zapErrorResp(errMsg string) *zap.Message {
	b := zap.NewBuilder(zap.HeaderSize + zapRespSize + len(errMsg) + 64)
	obj := b.StartObject(zapRespSize)
	obj.SetUint32(zapRespStatus, 1)
	obj.SetText(zapRespMessage, errMsg)
	obj.FinishAsRoot()
	msg, _ := zap.Parse(b.Finish())
	return msg
}
