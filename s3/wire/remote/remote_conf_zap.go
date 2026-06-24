// Code generated from remote.proto; DO NOT EDIT.

package remotewire

import (
	zap "github.com/zap-proto/go"
)

// --- RemoteConf ---
//
// RemoteConf mirrors remote_pb.RemoteConf: the credential/endpoint bundle for
// every supported remote object store (S3, GCS, Azure, Backblaze, Aliyun,
// Tencent, Baidu, Wasabi, Filebase, Storj, Contabo). Every proto field is a
// string except the three S3 booleans.
//
// Wire layout: the 40 string fields occupy 8-byte tail pointers first (in
// proto field-number order), then the 3 bool fields occupy 1 byte each. The
// reader addresses every field by its fixed offset, so the grouping is a
// layout convenience and does not change field identity.
const (
	remoteConfTypeOff                            = 0
	remoteConfNameOff                            = 8
	remoteConfS3AccessKeyOff                     = 16
	remoteConfS3SecretKeyOff                     = 24
	remoteConfS3RegionOff                        = 32
	remoteConfS3EndpointOff                      = 40
	remoteConfS3StorageClassOff                  = 48
	remoteConfGcsGoogleApplicationCredentialsOff = 56
	remoteConfGcsProjectIDOff                    = 64
	remoteConfAzureAccountNameOff                = 72
	remoteConfAzureAccountKeyOff                 = 80
	remoteConfBackblazeKeyIDOff                  = 88
	remoteConfBackblazeApplicationKeyOff         = 96
	remoteConfBackblazeEndpointOff               = 104
	remoteConfBackblazeRegionOff                 = 112
	remoteConfAliyunAccessKeyOff                 = 120
	remoteConfAliyunSecretKeyOff                 = 128
	remoteConfAliyunEndpointOff                  = 136
	remoteConfAliyunRegionOff                    = 144
	remoteConfTencentSecretIDOff                 = 152
	remoteConfTencentSecretKeyOff                = 160
	remoteConfTencentEndpointOff                 = 168
	remoteConfBaiduAccessKeyOff                  = 176
	remoteConfBaiduSecretKeyOff                  = 184
	remoteConfBaiduEndpointOff                   = 192
	remoteConfBaiduRegionOff                     = 200
	remoteConfWasabiAccessKeyOff                 = 208
	remoteConfWasabiSecretKeyOff                 = 216
	remoteConfWasabiEndpointOff                  = 224
	remoteConfWasabiRegionOff                    = 232
	remoteConfFilebaseAccessKeyOff               = 240
	remoteConfFilebaseSecretKeyOff               = 248
	remoteConfFilebaseEndpointOff                = 256
	remoteConfStorjAccessKeyOff                  = 264
	remoteConfStorjSecretKeyOff                  = 272
	remoteConfStorjEndpointOff                   = 280
	remoteConfContaboAccessKeyOff                = 288
	remoteConfContaboSecretKeyOff                = 296
	remoteConfContaboEndpointOff                 = 304
	remoteConfContaboRegionOff                   = 312
	remoteConfS3ForcePathStyleOff                = 320
	remoteConfS3SupportTaggingOff                = 321
	remoteConfS3V4SignatureOff                   = 322
	remoteConfSize                               = 323
)

// RemoteConf is a zero-copy view into a ZAP-encoded RemoteConf message.
type RemoteConf struct{ o zap.Object }

// WrapRemoteConf parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapRemoteConf(b []byte) (RemoteConf, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RemoteConf{}, err
	}
	return RemoteConf{o: m.Root()}, nil
}

func (t RemoteConf) Type() string           { return t.o.Text(remoteConfTypeOff) }
func (t RemoteConf) Name() string           { return t.o.Text(remoteConfNameOff) }
func (t RemoteConf) S3AccessKey() string    { return t.o.Text(remoteConfS3AccessKeyOff) }
func (t RemoteConf) S3SecretKey() string    { return t.o.Text(remoteConfS3SecretKeyOff) }
func (t RemoteConf) S3Region() string       { return t.o.Text(remoteConfS3RegionOff) }
func (t RemoteConf) S3Endpoint() string     { return t.o.Text(remoteConfS3EndpointOff) }
func (t RemoteConf) S3StorageClass() string { return t.o.Text(remoteConfS3StorageClassOff) }
func (t RemoteConf) GcsGoogleApplicationCredentials() string {
	return t.o.Text(remoteConfGcsGoogleApplicationCredentialsOff)
}
func (t RemoteConf) GcsProjectID() string     { return t.o.Text(remoteConfGcsProjectIDOff) }
func (t RemoteConf) AzureAccountName() string { return t.o.Text(remoteConfAzureAccountNameOff) }
func (t RemoteConf) AzureAccountKey() string  { return t.o.Text(remoteConfAzureAccountKeyOff) }
func (t RemoteConf) BackblazeKeyID() string   { return t.o.Text(remoteConfBackblazeKeyIDOff) }
func (t RemoteConf) BackblazeApplicationKey() string {
	return t.o.Text(remoteConfBackblazeApplicationKeyOff)
}
func (t RemoteConf) BackblazeEndpoint() string { return t.o.Text(remoteConfBackblazeEndpointOff) }
func (t RemoteConf) BackblazeRegion() string   { return t.o.Text(remoteConfBackblazeRegionOff) }
func (t RemoteConf) AliyunAccessKey() string   { return t.o.Text(remoteConfAliyunAccessKeyOff) }
func (t RemoteConf) AliyunSecretKey() string   { return t.o.Text(remoteConfAliyunSecretKeyOff) }
func (t RemoteConf) AliyunEndpoint() string    { return t.o.Text(remoteConfAliyunEndpointOff) }
func (t RemoteConf) AliyunRegion() string      { return t.o.Text(remoteConfAliyunRegionOff) }
func (t RemoteConf) TencentSecretID() string   { return t.o.Text(remoteConfTencentSecretIDOff) }
func (t RemoteConf) TencentSecretKey() string  { return t.o.Text(remoteConfTencentSecretKeyOff) }
func (t RemoteConf) TencentEndpoint() string   { return t.o.Text(remoteConfTencentEndpointOff) }
func (t RemoteConf) BaiduAccessKey() string    { return t.o.Text(remoteConfBaiduAccessKeyOff) }
func (t RemoteConf) BaiduSecretKey() string    { return t.o.Text(remoteConfBaiduSecretKeyOff) }
func (t RemoteConf) BaiduEndpoint() string     { return t.o.Text(remoteConfBaiduEndpointOff) }
func (t RemoteConf) BaiduRegion() string       { return t.o.Text(remoteConfBaiduRegionOff) }
func (t RemoteConf) WasabiAccessKey() string   { return t.o.Text(remoteConfWasabiAccessKeyOff) }
func (t RemoteConf) WasabiSecretKey() string   { return t.o.Text(remoteConfWasabiSecretKeyOff) }
func (t RemoteConf) WasabiEndpoint() string    { return t.o.Text(remoteConfWasabiEndpointOff) }
func (t RemoteConf) WasabiRegion() string      { return t.o.Text(remoteConfWasabiRegionOff) }
func (t RemoteConf) FilebaseAccessKey() string { return t.o.Text(remoteConfFilebaseAccessKeyOff) }
func (t RemoteConf) FilebaseSecretKey() string { return t.o.Text(remoteConfFilebaseSecretKeyOff) }
func (t RemoteConf) FilebaseEndpoint() string  { return t.o.Text(remoteConfFilebaseEndpointOff) }
func (t RemoteConf) StorjAccessKey() string    { return t.o.Text(remoteConfStorjAccessKeyOff) }
func (t RemoteConf) StorjSecretKey() string    { return t.o.Text(remoteConfStorjSecretKeyOff) }
func (t RemoteConf) StorjEndpoint() string     { return t.o.Text(remoteConfStorjEndpointOff) }
func (t RemoteConf) ContaboAccessKey() string  { return t.o.Text(remoteConfContaboAccessKeyOff) }
func (t RemoteConf) ContaboSecretKey() string  { return t.o.Text(remoteConfContaboSecretKeyOff) }
func (t RemoteConf) ContaboEndpoint() string   { return t.o.Text(remoteConfContaboEndpointOff) }
func (t RemoteConf) ContaboRegion() string     { return t.o.Text(remoteConfContaboRegionOff) }
func (t RemoteConf) S3ForcePathStyle() bool    { return t.o.Bool(remoteConfS3ForcePathStyleOff) }
func (t RemoteConf) S3SupportTagging() bool    { return t.o.Bool(remoteConfS3SupportTaggingOff) }
func (t RemoteConf) S3V4Signature() bool       { return t.o.Bool(remoteConfS3V4SignatureOff) }

// RemoteConfInput collects the field values for NewRemoteConf.
type RemoteConfInput struct {
	Type                            string
	Name                            string
	S3AccessKey                     string
	S3SecretKey                     string
	S3Region                        string
	S3Endpoint                      string
	S3StorageClass                  string
	GcsGoogleApplicationCredentials string
	GcsProjectID                    string
	AzureAccountName                string
	AzureAccountKey                 string
	BackblazeKeyID                  string
	BackblazeApplicationKey         string
	BackblazeEndpoint               string
	BackblazeRegion                 string
	AliyunAccessKey                 string
	AliyunSecretKey                 string
	AliyunEndpoint                  string
	AliyunRegion                    string
	TencentSecretID                 string
	TencentSecretKey                string
	TencentEndpoint                 string
	BaiduAccessKey                  string
	BaiduSecretKey                  string
	BaiduEndpoint                   string
	BaiduRegion                     string
	WasabiAccessKey                 string
	WasabiSecretKey                 string
	WasabiEndpoint                  string
	WasabiRegion                    string
	FilebaseAccessKey               string
	FilebaseSecretKey               string
	FilebaseEndpoint                string
	StorjAccessKey                  string
	StorjSecretKey                  string
	StorjEndpoint                   string
	ContaboAccessKey                string
	ContaboSecretKey                string
	ContaboEndpoint                 string
	ContaboRegion                   string
	S3ForcePathStyle                bool
	S3SupportTagging                bool
	S3V4Signature                   bool
}

// NewRemoteConf builds a ZAP-encoded RemoteConf message from in and returns the
// bytes.
func NewRemoteConf(in RemoteConfInput) []byte {
	b := zap.NewBuilder(512)
	ob := b.StartObject(remoteConfSize)
	ob.SetText(remoteConfTypeOff, in.Type)
	ob.SetText(remoteConfNameOff, in.Name)
	ob.SetText(remoteConfS3AccessKeyOff, in.S3AccessKey)
	ob.SetText(remoteConfS3SecretKeyOff, in.S3SecretKey)
	ob.SetText(remoteConfS3RegionOff, in.S3Region)
	ob.SetText(remoteConfS3EndpointOff, in.S3Endpoint)
	ob.SetText(remoteConfS3StorageClassOff, in.S3StorageClass)
	ob.SetText(remoteConfGcsGoogleApplicationCredentialsOff, in.GcsGoogleApplicationCredentials)
	ob.SetText(remoteConfGcsProjectIDOff, in.GcsProjectID)
	ob.SetText(remoteConfAzureAccountNameOff, in.AzureAccountName)
	ob.SetText(remoteConfAzureAccountKeyOff, in.AzureAccountKey)
	ob.SetText(remoteConfBackblazeKeyIDOff, in.BackblazeKeyID)
	ob.SetText(remoteConfBackblazeApplicationKeyOff, in.BackblazeApplicationKey)
	ob.SetText(remoteConfBackblazeEndpointOff, in.BackblazeEndpoint)
	ob.SetText(remoteConfBackblazeRegionOff, in.BackblazeRegion)
	ob.SetText(remoteConfAliyunAccessKeyOff, in.AliyunAccessKey)
	ob.SetText(remoteConfAliyunSecretKeyOff, in.AliyunSecretKey)
	ob.SetText(remoteConfAliyunEndpointOff, in.AliyunEndpoint)
	ob.SetText(remoteConfAliyunRegionOff, in.AliyunRegion)
	ob.SetText(remoteConfTencentSecretIDOff, in.TencentSecretID)
	ob.SetText(remoteConfTencentSecretKeyOff, in.TencentSecretKey)
	ob.SetText(remoteConfTencentEndpointOff, in.TencentEndpoint)
	ob.SetText(remoteConfBaiduAccessKeyOff, in.BaiduAccessKey)
	ob.SetText(remoteConfBaiduSecretKeyOff, in.BaiduSecretKey)
	ob.SetText(remoteConfBaiduEndpointOff, in.BaiduEndpoint)
	ob.SetText(remoteConfBaiduRegionOff, in.BaiduRegion)
	ob.SetText(remoteConfWasabiAccessKeyOff, in.WasabiAccessKey)
	ob.SetText(remoteConfWasabiSecretKeyOff, in.WasabiSecretKey)
	ob.SetText(remoteConfWasabiEndpointOff, in.WasabiEndpoint)
	ob.SetText(remoteConfWasabiRegionOff, in.WasabiRegion)
	ob.SetText(remoteConfFilebaseAccessKeyOff, in.FilebaseAccessKey)
	ob.SetText(remoteConfFilebaseSecretKeyOff, in.FilebaseSecretKey)
	ob.SetText(remoteConfFilebaseEndpointOff, in.FilebaseEndpoint)
	ob.SetText(remoteConfStorjAccessKeyOff, in.StorjAccessKey)
	ob.SetText(remoteConfStorjSecretKeyOff, in.StorjSecretKey)
	ob.SetText(remoteConfStorjEndpointOff, in.StorjEndpoint)
	ob.SetText(remoteConfContaboAccessKeyOff, in.ContaboAccessKey)
	ob.SetText(remoteConfContaboSecretKeyOff, in.ContaboSecretKey)
	ob.SetText(remoteConfContaboEndpointOff, in.ContaboEndpoint)
	ob.SetText(remoteConfContaboRegionOff, in.ContaboRegion)
	ob.SetBool(remoteConfS3ForcePathStyleOff, in.S3ForcePathStyle)
	ob.SetBool(remoteConfS3SupportTaggingOff, in.S3SupportTagging)
	ob.SetBool(remoteConfS3V4SignatureOff, in.S3V4Signature)
	ob.FinishAsRoot()
	return b.Finish()
}
