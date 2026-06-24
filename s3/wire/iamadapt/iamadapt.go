// Package iamadapt bridges the protobuf iam_pb domain structs and the
// zero-copy iamwire ZAP messages for the HanzoS3IamCache RPC. It is the ONE
// place that knows how an iam_pb.Identity / iam_pb.Group maps onto its
// iamwire input (build side) and back out of an iamwire view (read side), so
// the cache client (filer) and the cache server (s3 gateway) share a single
// definition instead of each carrying its own copy.
package iamadapt

import (
	"net"
	"strconv"

	"github.com/hanzoai/s3/s3/pb/iam_pb"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

// CacheZapPortOffset is the fixed offset added to an S3 server's gRPC port to
// derive the port of its HanzoS3IamCache ZAP listener. The cache RPC moved off
// the shared gRPC server onto its own ZAP transport; deriving the port from the
// already-published gRPC port keeps the address discoverable from the cluster
// node list without a new flag.
const CacheZapPortOffset = 1

// CacheZapAddress maps an S3 server's gRPC address (host:port) to the address of
// its HanzoS3IamCache ZAP listener. A malformed address is returned unchanged.
func CacheZapAddress(grpcAddress string) string {
	host, port, err := net.SplitHostPort(grpcAddress)
	if err != nil {
		return grpcAddress
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return grpcAddress
	}
	return net.JoinHostPort(host, strconv.Itoa(p+CacheZapPortOffset))
}

// IdentityInput maps an iam_pb.Identity onto its iamwire build input.
func IdentityInput(id *iam_pb.Identity) *iamwire.IdentityInput {
	if id == nil {
		return nil
	}
	in := &iamwire.IdentityInput{
		Name:              id.Name,
		Actions:           id.Actions,
		Disabled:          id.Disabled,
		ServiceAccountIDs: id.ServiceAccountIds,
		PolicyNames:       id.PolicyNames,
		IsStatic:          id.IsStatic,
	}
	if id.Account != nil {
		in.Account = &iamwire.AccountInput{
			ID:           id.Account.Id,
			DisplayName:  id.Account.DisplayName,
			EmailAddress: id.Account.EmailAddress,
		}
	}
	for _, c := range id.Credentials {
		if c == nil {
			continue
		}
		in.Credentials = append(in.Credentials, iamwire.CredentialInput{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
			Status:    c.Status,
		})
	}
	for _, t := range id.Tags {
		if t == nil {
			continue
		}
		in.Tags = append(in.Tags, iamwire.UserTagInput{Key: t.Key, Value: t.Value})
	}
	return in
}

// IdentityFromWire materializes an iam_pb.Identity from an iamwire view.
func IdentityFromWire(v iamwire.Identity) *iam_pb.Identity {
	id := &iam_pb.Identity{
		Name:     v.Name(),
		Disabled: v.Disabled(),
		IsStatic: v.IsStatic(),
	}
	if n := v.CredentialsLen(); n > 0 {
		id.Credentials = make([]*iam_pb.Credential, 0, n)
		for i := 0; i < n; i++ {
			c := v.CredentialsAt(i)
			id.Credentials = append(id.Credentials, &iam_pb.Credential{
				AccessKey: c.AccessKey(),
				SecretKey: c.SecretKey(),
				Status:    c.Status(),
			})
		}
	}
	if n := v.ActionsLen(); n > 0 {
		id.Actions = make([]string, 0, n)
		for i := 0; i < n; i++ {
			id.Actions = append(id.Actions, v.ActionsAt(i))
		}
	}
	if acc := v.Account(); acc.ID() != "" || acc.DisplayName() != "" || acc.EmailAddress() != "" {
		id.Account = &iam_pb.Account{
			Id:           acc.ID(),
			DisplayName:  acc.DisplayName(),
			EmailAddress: acc.EmailAddress(),
		}
	}
	if n := v.ServiceAccountIDsLen(); n > 0 {
		id.ServiceAccountIds = make([]string, 0, n)
		for i := 0; i < n; i++ {
			id.ServiceAccountIds = append(id.ServiceAccountIds, v.ServiceAccountIDsAt(i))
		}
	}
	if n := v.PolicyNamesLen(); n > 0 {
		id.PolicyNames = make([]string, 0, n)
		for i := 0; i < n; i++ {
			id.PolicyNames = append(id.PolicyNames, v.PolicyNamesAt(i))
		}
	}
	if n := v.TagsLen(); n > 0 {
		id.Tags = make([]*iam_pb.UserTag, 0, n)
		for i := 0; i < n; i++ {
			t := v.TagsAt(i)
			id.Tags = append(id.Tags, &iam_pb.UserTag{Key: t.Key(), Value: t.Value()})
		}
	}
	return id
}

// GroupInput maps an iam_pb.Group onto its iamwire build input.
func GroupInput(g *iam_pb.Group) *iamwire.GroupInput {
	if g == nil {
		return nil
	}
	return &iamwire.GroupInput{
		Name:        g.Name,
		Members:     g.Members,
		PolicyNames: g.PolicyNames,
		Disabled:    g.Disabled,
	}
}

// GroupFromWire materializes an iam_pb.Group from an iamwire view.
func GroupFromWire(v iamwire.Group) *iam_pb.Group {
	g := &iam_pb.Group{
		Name:     v.Name(),
		Disabled: v.Disabled(),
	}
	if n := v.MembersLen(); n > 0 {
		g.Members = make([]string, 0, n)
		for i := 0; i < n; i++ {
			g.Members = append(g.Members, v.MembersAt(i))
		}
	}
	if n := v.PolicyNamesLen(); n > 0 {
		g.PolicyNames = make([]string, 0, n)
		for i := 0; i < n; i++ {
			g.PolicyNames = append(g.PolicyNames, v.PolicyNamesAt(i))
		}
	}
	return g
}
