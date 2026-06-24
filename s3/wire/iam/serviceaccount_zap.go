// Code generated from iam.proto; DO NOT EDIT.

package iamwire

import (
	zap "github.com/zap-proto/go"
)

// Service-account management request/response messages.

// --- CreateServiceAccountRequest ---

const (
	createServiceAccountRequestServiceAccountOff = 0
	createServiceAccountRequestSize              = 4
)

// CreateServiceAccountRequest is a zero-copy view into a ZAP-encoded
// CreateServiceAccountRequest message.
type CreateServiceAccountRequest struct{ o zap.Object }

// WrapCreateServiceAccountRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapCreateServiceAccountRequest(b []byte) (CreateServiceAccountRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CreateServiceAccountRequest{}, err
	}
	return CreateServiceAccountRequest{o: m.Root()}, nil
}

// ServiceAccount reads the service_account field (proto field 1, ServiceAccount).
func (t CreateServiceAccountRequest) ServiceAccount() ServiceAccount {
	return ServiceAccount{o: t.o.Object(createServiceAccountRequestServiceAccountOff)}
}

// CreateServiceAccountRequestInput collects the field values for
// NewCreateServiceAccountRequest.
type CreateServiceAccountRequestInput struct {
	ServiceAccount *ServiceAccountInput
}

// NewCreateServiceAccountRequest builds a ZAP-encoded CreateServiceAccountRequest
// message from in and returns the bytes.
func NewCreateServiceAccountRequest(in CreateServiceAccountRequestInput) []byte {
	b := zap.NewBuilder(512)
	serviceAccountOff := 0
	if in.ServiceAccount != nil {
		serviceAccountOff = appendServiceAccount(b, *in.ServiceAccount).Finish()
	}
	ob := b.StartObject(createServiceAccountRequestSize)
	if serviceAccountOff != 0 {
		ob.SetObject(createServiceAccountRequestServiceAccountOff, serviceAccountOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CreateServiceAccountResponse ---

const createServiceAccountResponseSize = 0

// CreateServiceAccountResponse is a zero-copy view into a ZAP-encoded
// CreateServiceAccountResponse message. It carries no fields.
type CreateServiceAccountResponse struct{ o zap.Object }

// WrapCreateServiceAccountResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapCreateServiceAccountResponse(b []byte) (CreateServiceAccountResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CreateServiceAccountResponse{}, err
	}
	return CreateServiceAccountResponse{o: m.Root()}, nil
}

// CreateServiceAccountResponseInput collects the field values for
// NewCreateServiceAccountResponse. CreateServiceAccountResponse has no fields.
type CreateServiceAccountResponseInput struct{}

// NewCreateServiceAccountResponse builds a ZAP-encoded
// CreateServiceAccountResponse message and returns the bytes.
func NewCreateServiceAccountResponse(in CreateServiceAccountResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(createServiceAccountResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- UpdateServiceAccountRequest ---

const (
	updateServiceAccountRequestIDOff             = 0
	updateServiceAccountRequestServiceAccountOff = 8
	updateServiceAccountRequestSize              = 16
)

// UpdateServiceAccountRequest is a zero-copy view into a ZAP-encoded
// UpdateServiceAccountRequest message.
type UpdateServiceAccountRequest struct{ o zap.Object }

// WrapUpdateServiceAccountRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapUpdateServiceAccountRequest(b []byte) (UpdateServiceAccountRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return UpdateServiceAccountRequest{}, err
	}
	return UpdateServiceAccountRequest{o: m.Root()}, nil
}

// ID reads the id field (proto field 1, string).
func (t UpdateServiceAccountRequest) ID() string {
	return t.o.Text(updateServiceAccountRequestIDOff)
}

// ServiceAccount reads the service_account field (proto field 2, ServiceAccount).
func (t UpdateServiceAccountRequest) ServiceAccount() ServiceAccount {
	return ServiceAccount{o: t.o.Object(updateServiceAccountRequestServiceAccountOff)}
}

// UpdateServiceAccountRequestInput collects the field values for
// NewUpdateServiceAccountRequest.
type UpdateServiceAccountRequestInput struct {
	ID             string
	ServiceAccount *ServiceAccountInput
}

// NewUpdateServiceAccountRequest builds a ZAP-encoded UpdateServiceAccountRequest
// message from in and returns the bytes.
func NewUpdateServiceAccountRequest(in UpdateServiceAccountRequestInput) []byte {
	b := zap.NewBuilder(512)
	serviceAccountOff := 0
	if in.ServiceAccount != nil {
		serviceAccountOff = appendServiceAccount(b, *in.ServiceAccount).Finish()
	}
	ob := b.StartObject(updateServiceAccountRequestSize)
	ob.SetText(updateServiceAccountRequestIDOff, in.ID)
	if serviceAccountOff != 0 {
		ob.SetObject(updateServiceAccountRequestServiceAccountOff, serviceAccountOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- UpdateServiceAccountResponse ---

const updateServiceAccountResponseSize = 0

// UpdateServiceAccountResponse is a zero-copy view into a ZAP-encoded
// UpdateServiceAccountResponse message. It carries no fields.
type UpdateServiceAccountResponse struct{ o zap.Object }

// WrapUpdateServiceAccountResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapUpdateServiceAccountResponse(b []byte) (UpdateServiceAccountResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return UpdateServiceAccountResponse{}, err
	}
	return UpdateServiceAccountResponse{o: m.Root()}, nil
}

// UpdateServiceAccountResponseInput collects the field values for
// NewUpdateServiceAccountResponse. UpdateServiceAccountResponse has no fields.
type UpdateServiceAccountResponseInput struct{}

// NewUpdateServiceAccountResponse builds a ZAP-encoded
// UpdateServiceAccountResponse message and returns the bytes.
func NewUpdateServiceAccountResponse(in UpdateServiceAccountResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(updateServiceAccountResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteServiceAccountRequest ---

const (
	deleteServiceAccountRequestIDOff = 0
	deleteServiceAccountRequestSize  = 8
)

// DeleteServiceAccountRequest is a zero-copy view into a ZAP-encoded
// DeleteServiceAccountRequest message.
type DeleteServiceAccountRequest struct{ o zap.Object }

// WrapDeleteServiceAccountRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapDeleteServiceAccountRequest(b []byte) (DeleteServiceAccountRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteServiceAccountRequest{}, err
	}
	return DeleteServiceAccountRequest{o: m.Root()}, nil
}

// ID reads the id field (proto field 1, string).
func (t DeleteServiceAccountRequest) ID() string {
	return t.o.Text(deleteServiceAccountRequestIDOff)
}

// DeleteServiceAccountRequestInput collects the field values for
// NewDeleteServiceAccountRequest.
type DeleteServiceAccountRequestInput struct {
	ID string
}

// NewDeleteServiceAccountRequest builds a ZAP-encoded DeleteServiceAccountRequest
// message from in and returns the bytes.
func NewDeleteServiceAccountRequest(in DeleteServiceAccountRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteServiceAccountRequestSize)
	ob.SetText(deleteServiceAccountRequestIDOff, in.ID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteServiceAccountResponse ---

const deleteServiceAccountResponseSize = 0

// DeleteServiceAccountResponse is a zero-copy view into a ZAP-encoded
// DeleteServiceAccountResponse message. It carries no fields.
type DeleteServiceAccountResponse struct{ o zap.Object }

// WrapDeleteServiceAccountResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapDeleteServiceAccountResponse(b []byte) (DeleteServiceAccountResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteServiceAccountResponse{}, err
	}
	return DeleteServiceAccountResponse{o: m.Root()}, nil
}

// DeleteServiceAccountResponseInput collects the field values for
// NewDeleteServiceAccountResponse. DeleteServiceAccountResponse has no fields.
type DeleteServiceAccountResponseInput struct{}

// NewDeleteServiceAccountResponse builds a ZAP-encoded
// DeleteServiceAccountResponse message and returns the bytes.
func NewDeleteServiceAccountResponse(in DeleteServiceAccountResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteServiceAccountResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetServiceAccountRequest ---

const (
	getServiceAccountRequestIDOff = 0
	getServiceAccountRequestSize  = 8
)

// GetServiceAccountRequest is a zero-copy view into a ZAP-encoded
// GetServiceAccountRequest message.
type GetServiceAccountRequest struct{ o zap.Object }

// WrapGetServiceAccountRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetServiceAccountRequest(b []byte) (GetServiceAccountRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetServiceAccountRequest{}, err
	}
	return GetServiceAccountRequest{o: m.Root()}, nil
}

// ID reads the id field (proto field 1, string).
func (t GetServiceAccountRequest) ID() string { return t.o.Text(getServiceAccountRequestIDOff) }

// GetServiceAccountRequestInput collects the field values for
// NewGetServiceAccountRequest.
type GetServiceAccountRequestInput struct {
	ID string
}

// NewGetServiceAccountRequest builds a ZAP-encoded GetServiceAccountRequest
// message from in and returns the bytes.
func NewGetServiceAccountRequest(in GetServiceAccountRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getServiceAccountRequestSize)
	ob.SetText(getServiceAccountRequestIDOff, in.ID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetServiceAccountResponse ---

const (
	getServiceAccountResponseServiceAccountOff = 0
	getServiceAccountResponseSize              = 4
)

// GetServiceAccountResponse is a zero-copy view into a ZAP-encoded
// GetServiceAccountResponse message.
type GetServiceAccountResponse struct{ o zap.Object }

// WrapGetServiceAccountResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapGetServiceAccountResponse(b []byte) (GetServiceAccountResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetServiceAccountResponse{}, err
	}
	return GetServiceAccountResponse{o: m.Root()}, nil
}

// ServiceAccount reads the service_account field (proto field 1, ServiceAccount).
func (t GetServiceAccountResponse) ServiceAccount() ServiceAccount {
	return ServiceAccount{o: t.o.Object(getServiceAccountResponseServiceAccountOff)}
}

// GetServiceAccountResponseInput collects the field values for
// NewGetServiceAccountResponse.
type GetServiceAccountResponseInput struct {
	ServiceAccount *ServiceAccountInput
}

// NewGetServiceAccountResponse builds a ZAP-encoded GetServiceAccountResponse
// message from in and returns the bytes.
func NewGetServiceAccountResponse(in GetServiceAccountResponseInput) []byte {
	b := zap.NewBuilder(512)
	serviceAccountOff := 0
	if in.ServiceAccount != nil {
		serviceAccountOff = appendServiceAccount(b, *in.ServiceAccount).Finish()
	}
	ob := b.StartObject(getServiceAccountResponseSize)
	if serviceAccountOff != 0 {
		ob.SetObject(getServiceAccountResponseServiceAccountOff, serviceAccountOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListServiceAccountsRequest ---

const listServiceAccountsRequestSize = 0

// ListServiceAccountsRequest is a zero-copy view into a ZAP-encoded
// ListServiceAccountsRequest message. It carries no fields.
type ListServiceAccountsRequest struct{ o zap.Object }

// WrapListServiceAccountsRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapListServiceAccountsRequest(b []byte) (ListServiceAccountsRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListServiceAccountsRequest{}, err
	}
	return ListServiceAccountsRequest{o: m.Root()}, nil
}

// ListServiceAccountsRequestInput collects the field values for
// NewListServiceAccountsRequest. ListServiceAccountsRequest has no fields.
type ListServiceAccountsRequestInput struct{}

// NewListServiceAccountsRequest builds a ZAP-encoded ListServiceAccountsRequest
// message and returns the bytes.
func NewListServiceAccountsRequest(in ListServiceAccountsRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(listServiceAccountsRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListServiceAccountsResponse ---

const (
	listServiceAccountsResponseServiceAccountsOff = 0
	listServiceAccountsResponseSize               = 8
)

// ListServiceAccountsResponse is a zero-copy view into a ZAP-encoded
// ListServiceAccountsResponse message.
type ListServiceAccountsResponse struct{ o zap.Object }

// WrapListServiceAccountsResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapListServiceAccountsResponse(b []byte) (ListServiceAccountsResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListServiceAccountsResponse{}, err
	}
	return ListServiceAccountsResponse{o: m.Root()}, nil
}

// ServiceAccountsLen returns the number of service accounts (proto field 1,
// repeated ServiceAccount).
func (t ListServiceAccountsResponse) ServiceAccountsLen() int {
	return t.o.List(listServiceAccountsResponseServiceAccountsOff).Len()
}

// ServiceAccountsAt returns service account i (proto field 1, repeated ServiceAccount).
func (t ListServiceAccountsResponse) ServiceAccountsAt(i int) ServiceAccount {
	return ServiceAccount{o: t.o.List(listServiceAccountsResponseServiceAccountsOff).ObjectAt(i)}
}

// ListServiceAccountsResponseInput collects the field values for
// NewListServiceAccountsResponse.
type ListServiceAccountsResponseInput struct {
	ServiceAccounts []ServiceAccountInput
}

// NewListServiceAccountsResponse builds a ZAP-encoded ListServiceAccountsResponse
// message from in and returns the bytes.
func NewListServiceAccountsResponse(in ListServiceAccountsResponseInput) []byte {
	b := zap.NewBuilder(512)
	ob := b.StartObject(listServiceAccountsResponseSize)
	if len(in.ServiceAccounts) > 0 {
		lb := b.StartList(0)
		for i := range in.ServiceAccounts {
			lb.AddObjectBytes(NewServiceAccount(in.ServiceAccounts[i]))
		}
		off, n := lb.Finish()
		ob.SetList(listServiceAccountsResponseServiceAccountsOff, off, n)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetServiceAccountByAccessKeyRequest ---

const (
	getServiceAccountByAccessKeyRequestAccessKeyOff = 0
	getServiceAccountByAccessKeyRequestSize         = 8
)

// GetServiceAccountByAccessKeyRequest is a zero-copy view into a ZAP-encoded
// GetServiceAccountByAccessKeyRequest message.
type GetServiceAccountByAccessKeyRequest struct{ o zap.Object }

// WrapGetServiceAccountByAccessKeyRequest parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapGetServiceAccountByAccessKeyRequest(b []byte) (GetServiceAccountByAccessKeyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetServiceAccountByAccessKeyRequest{}, err
	}
	return GetServiceAccountByAccessKeyRequest{o: m.Root()}, nil
}

// AccessKey reads the access_key field (proto field 1, string).
func (t GetServiceAccountByAccessKeyRequest) AccessKey() string {
	return t.o.Text(getServiceAccountByAccessKeyRequestAccessKeyOff)
}

// GetServiceAccountByAccessKeyRequestInput collects the field values for
// NewGetServiceAccountByAccessKeyRequest.
type GetServiceAccountByAccessKeyRequestInput struct {
	AccessKey string
}

// NewGetServiceAccountByAccessKeyRequest builds a ZAP-encoded
// GetServiceAccountByAccessKeyRequest message from in and returns the bytes.
func NewGetServiceAccountByAccessKeyRequest(in GetServiceAccountByAccessKeyRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getServiceAccountByAccessKeyRequestSize)
	ob.SetText(getServiceAccountByAccessKeyRequestAccessKeyOff, in.AccessKey)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetServiceAccountByAccessKeyResponse ---

const (
	getServiceAccountByAccessKeyResponseServiceAccountOff = 0
	getServiceAccountByAccessKeyResponseSize              = 4
)

// GetServiceAccountByAccessKeyResponse is a zero-copy view into a ZAP-encoded
// GetServiceAccountByAccessKeyResponse message.
type GetServiceAccountByAccessKeyResponse struct{ o zap.Object }

// WrapGetServiceAccountByAccessKeyResponse parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapGetServiceAccountByAccessKeyResponse(b []byte) (GetServiceAccountByAccessKeyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetServiceAccountByAccessKeyResponse{}, err
	}
	return GetServiceAccountByAccessKeyResponse{o: m.Root()}, nil
}

// ServiceAccount reads the service_account field (proto field 1, ServiceAccount).
func (t GetServiceAccountByAccessKeyResponse) ServiceAccount() ServiceAccount {
	return ServiceAccount{o: t.o.Object(getServiceAccountByAccessKeyResponseServiceAccountOff)}
}

// GetServiceAccountByAccessKeyResponseInput collects the field values for
// NewGetServiceAccountByAccessKeyResponse.
type GetServiceAccountByAccessKeyResponseInput struct {
	ServiceAccount *ServiceAccountInput
}

// NewGetServiceAccountByAccessKeyResponse builds a ZAP-encoded
// GetServiceAccountByAccessKeyResponse message from in and returns the bytes.
func NewGetServiceAccountByAccessKeyResponse(in GetServiceAccountByAccessKeyResponseInput) []byte {
	b := zap.NewBuilder(512)
	serviceAccountOff := 0
	if in.ServiceAccount != nil {
		serviceAccountOff = appendServiceAccount(b, *in.ServiceAccount).Finish()
	}
	ob := b.StartObject(getServiceAccountByAccessKeyResponseSize)
	if serviceAccountOff != 0 {
		ob.SetObject(getServiceAccountByAccessKeyResponseServiceAccountOff, serviceAccountOff)
	}
	ob.FinishAsRoot()
	return b.Finish()
}
