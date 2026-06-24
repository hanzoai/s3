package iamwire

// Bridge between the protobuf IAM data model (iam_pb, still the in-process
// representation used by the credential manager and the rest of S3) and the
// ZAP wire Input/view types. The RPC envelope rides on ZAP; the records it
// carries are translated here in ONE place so client and server share the
// exact same mapping.
//
// Direction:
//   *iam_pb.T  -> TInput   (pb<name>Input): build an outgoing ZAP message.
//   T (view)   -> *iam_pb.T (<name>FromView): decode an incoming ZAP message.

import (
	"github.com/hanzoai/s3/s3/pb/iam_pb"
)

// ---- pb -> Input ----

func pbCredentialInput(c *iam_pb.Credential) *CredentialInput {
	if c == nil {
		return nil
	}
	return &CredentialInput{AccessKey: c.AccessKey, SecretKey: c.SecretKey, Status: c.Status}
}

func pbAccountInput(a *iam_pb.Account) *AccountInput {
	if a == nil {
		return nil
	}
	return &AccountInput{ID: a.Id, DisplayName: a.DisplayName, EmailAddress: a.EmailAddress}
}

func pbUserTagInputs(tags []*iam_pb.UserTag) []UserTagInput {
	if len(tags) == 0 {
		return nil
	}
	out := make([]UserTagInput, 0, len(tags))
	for _, t := range tags {
		if t == nil {
			continue
		}
		out = append(out, UserTagInput{Key: t.Key, Value: t.Value})
	}
	return out
}

func pbCredentialInputs(creds []*iam_pb.Credential) []CredentialInput {
	if len(creds) == 0 {
		return nil
	}
	out := make([]CredentialInput, 0, len(creds))
	for _, c := range creds {
		if c == nil {
			continue
		}
		out = append(out, CredentialInput{AccessKey: c.AccessKey, SecretKey: c.SecretKey, Status: c.Status})
	}
	return out
}

// pbIdentityInput converts an *iam_pb.Identity into an IdentityInput.
func pbIdentityInput(id *iam_pb.Identity) *IdentityInput {
	if id == nil {
		return nil
	}
	return &IdentityInput{
		Name:              id.Name,
		Credentials:       pbCredentialInputs(id.Credentials),
		Actions:           id.Actions,
		Account:           pbAccountInput(id.Account),
		Disabled:          id.Disabled,
		ServiceAccountIDs: id.ServiceAccountIds,
		PolicyNames:       id.PolicyNames,
		IsStatic:          id.IsStatic,
		Tags:              pbUserTagInputs(id.Tags),
	}
}

// pbServiceAccountInput converts an *iam_pb.ServiceAccount into a ServiceAccountInput.
func pbServiceAccountInput(sa *iam_pb.ServiceAccount) *ServiceAccountInput {
	if sa == nil {
		return nil
	}
	return &ServiceAccountInput{
		ID:          sa.Id,
		ParentUser:  sa.ParentUser,
		Description: sa.Description,
		Credential:  pbCredentialInput(sa.Credential),
		Actions:     sa.Actions,
		Expiration:  sa.Expiration,
		Disabled:    sa.Disabled,
		CreatedAt:   sa.CreatedAt,
		CreatedBy:   sa.CreatedBy,
	}
}

func pbGroupInput(g *iam_pb.Group) *GroupInput {
	if g == nil {
		return nil
	}
	return &GroupInput{Name: g.Name, Members: g.Members, PolicyNames: g.PolicyNames, Disabled: g.Disabled}
}

func pbPolicyInput(p *iam_pb.Policy) *PolicyInput {
	if p == nil {
		return nil
	}
	return &PolicyInput{Name: p.Name, Content: p.Content}
}

// pbConfigurationInput converts an *iam_pb.S3ApiConfiguration into an
// S3ApiConfigurationInput.
func pbConfigurationInput(c *iam_pb.S3ApiConfiguration) *S3ApiConfigurationInput {
	if c == nil {
		return nil
	}
	in := &S3ApiConfigurationInput{}
	for _, id := range c.Identities {
		if v := pbIdentityInput(id); v != nil {
			in.Identities = append(in.Identities, *v)
		}
	}
	for _, a := range c.Accounts {
		if v := pbAccountInput(a); v != nil {
			in.Accounts = append(in.Accounts, *v)
		}
	}
	for _, sa := range c.ServiceAccounts {
		if v := pbServiceAccountInput(sa); v != nil {
			in.ServiceAccounts = append(in.ServiceAccounts, *v)
		}
	}
	for _, p := range c.Policies {
		if v := pbPolicyInput(p); v != nil {
			in.Policies = append(in.Policies, *v)
		}
	}
	for _, g := range c.Groups {
		if v := pbGroupInput(g); v != nil {
			in.Groups = append(in.Groups, *v)
		}
	}
	return in
}

// ---- view -> pb ----

func credentialToPB(c Credential) *iam_pb.Credential {
	return &iam_pb.Credential{AccessKey: c.AccessKey(), SecretKey: c.SecretKey(), Status: c.Status()}
}

// CredentialToPBExported decodes a Credential view into an *iam_pb.Credential.
func CredentialToPBExported(c Credential) *iam_pb.Credential { return credentialToPB(c) }

func accountToPB(a Account) *iam_pb.Account {
	return &iam_pb.Account{Id: a.ID(), DisplayName: a.DisplayName(), EmailAddress: a.EmailAddress()}
}

// IdentityToPB decodes an Identity view into an *iam_pb.Identity.
func IdentityToPB(id Identity) *iam_pb.Identity {
	out := &iam_pb.Identity{
		Name:     id.Name(),
		Disabled: id.Disabled(),
		IsStatic: id.IsStatic(),
	}
	if n := id.CredentialsLen(); n > 0 {
		out.Credentials = make([]*iam_pb.Credential, n)
		for i := 0; i < n; i++ {
			out.Credentials[i] = credentialToPB(id.CredentialsAt(i))
		}
	}
	if n := id.ActionsLen(); n > 0 {
		out.Actions = make([]string, n)
		for i := 0; i < n; i++ {
			out.Actions[i] = id.ActionsAt(i)
		}
	}
	// Account is a message: an absent account decodes to the zero view whose
	// fields read empty; only attach when something is set.
	if acc := id.Account(); acc.ID() != "" || acc.DisplayName() != "" || acc.EmailAddress() != "" {
		out.Account = accountToPB(acc)
	}
	if n := id.ServiceAccountIDsLen(); n > 0 {
		out.ServiceAccountIds = make([]string, n)
		for i := 0; i < n; i++ {
			out.ServiceAccountIds[i] = id.ServiceAccountIDsAt(i)
		}
	}
	if n := id.PolicyNamesLen(); n > 0 {
		out.PolicyNames = make([]string, n)
		for i := 0; i < n; i++ {
			out.PolicyNames[i] = id.PolicyNamesAt(i)
		}
	}
	if n := id.TagsLen(); n > 0 {
		out.Tags = make([]*iam_pb.UserTag, n)
		for i := 0; i < n; i++ {
			t := id.TagsAt(i)
			out.Tags[i] = &iam_pb.UserTag{Key: t.Key(), Value: t.Value()}
		}
	}
	return out
}

// ServiceAccountToPB decodes a ServiceAccount view into an *iam_pb.ServiceAccount.
func ServiceAccountToPB(sa ServiceAccount) *iam_pb.ServiceAccount {
	out := &iam_pb.ServiceAccount{
		Id:          sa.ID(),
		ParentUser:  sa.ParentUser(),
		Description: sa.Description(),
		Expiration:  sa.Expiration(),
		Disabled:    sa.Disabled(),
		CreatedAt:   sa.CreatedAt(),
		CreatedBy:   sa.CreatedBy(),
	}
	if cred := sa.Credential(); cred.AccessKey() != "" || cred.SecretKey() != "" || cred.Status() != "" {
		out.Credential = credentialToPB(cred)
	}
	if n := sa.ActionsLen(); n > 0 {
		out.Actions = make([]string, n)
		for i := 0; i < n; i++ {
			out.Actions[i] = sa.ActionsAt(i)
		}
	}
	return out
}

func groupToPB(g Group) *iam_pb.Group {
	out := &iam_pb.Group{Name: g.Name(), Disabled: g.Disabled()}
	if n := g.MembersLen(); n > 0 {
		out.Members = make([]string, n)
		for i := 0; i < n; i++ {
			out.Members[i] = g.MembersAt(i)
		}
	}
	if n := g.PolicyNamesLen(); n > 0 {
		out.PolicyNames = make([]string, n)
		for i := 0; i < n; i++ {
			out.PolicyNames[i] = g.PolicyNamesAt(i)
		}
	}
	return out
}

func policyToPB(p Policy) *iam_pb.Policy {
	return &iam_pb.Policy{Name: p.Name(), Content: p.Content()}
}

// ConfigurationToPB decodes an S3ApiConfiguration view into an
// *iam_pb.S3ApiConfiguration.
func ConfigurationToPB(c S3ApiConfiguration) *iam_pb.S3ApiConfiguration {
	out := &iam_pb.S3ApiConfiguration{}
	if n := c.IdentitiesLen(); n > 0 {
		out.Identities = make([]*iam_pb.Identity, n)
		for i := 0; i < n; i++ {
			out.Identities[i] = IdentityToPB(c.IdentitiesAt(i))
		}
	}
	if n := c.AccountsLen(); n > 0 {
		out.Accounts = make([]*iam_pb.Account, n)
		for i := 0; i < n; i++ {
			out.Accounts[i] = accountToPB(c.AccountsAt(i))
		}
	}
	if n := c.ServiceAccountsLen(); n > 0 {
		out.ServiceAccounts = make([]*iam_pb.ServiceAccount, n)
		for i := 0; i < n; i++ {
			out.ServiceAccounts[i] = ServiceAccountToPB(c.ServiceAccountsAt(i))
		}
	}
	if n := c.PoliciesLen(); n > 0 {
		out.Policies = make([]*iam_pb.Policy, n)
		for i := 0; i < n; i++ {
			out.Policies[i] = policyToPB(c.PoliciesAt(i))
		}
	}
	if n := c.GroupsLen(); n > 0 {
		out.Groups = make([]*iam_pb.Group, n)
		for i := 0; i < n; i++ {
			out.Groups[i] = groupToPB(c.GroupsAt(i))
		}
	}
	return out
}

// ---- response decoders (used by client call sites) ----
//
// Each takes the wire response body returned by the typed client and yields the
// in-process pb model. A nil result for a singular message means the field was
// absent on the wire (the gRPC client's nil-pointer semantics).

// GetConfigurationResp decodes a GetConfiguration response body.
func GetConfigurationResp(body []byte) (*iam_pb.S3ApiConfiguration, error) {
	resp, err := WrapGetConfigurationResponse(body)
	if err != nil {
		return nil, err
	}
	if resp.o.Object(getConfigurationResponseConfigurationOff).IsNull() {
		return nil, nil
	}
	return ConfigurationToPB(resp.Configuration()), nil
}

// GetUserResp decodes a GetUser response body, returning a nil identity when
// the response carried none.
func GetUserResp(body []byte) (*iam_pb.Identity, error) {
	resp, err := WrapGetUserResponse(body)
	if err != nil {
		return nil, err
	}
	if resp.o.Object(getUserResponseIdentityOff).IsNull() {
		return nil, nil
	}
	return IdentityToPB(resp.Identity()), nil
}

// GetUserByAccessKeyResp decodes a GetUserByAccessKey response body.
func GetUserByAccessKeyResp(body []byte) (*iam_pb.Identity, error) {
	resp, err := WrapGetUserByAccessKeyResponse(body)
	if err != nil {
		return nil, err
	}
	if resp.o.Object(getUserByAccessKeyResponseIdentityOff).IsNull() {
		return nil, nil
	}
	return IdentityToPB(resp.Identity()), nil
}

// GetPolicyResp decodes a GetPolicy response body into (name, content).
func GetPolicyResp(body []byte) (string, string, error) {
	resp, err := WrapGetPolicyResponse(body)
	if err != nil {
		return "", "", err
	}
	return resp.Name(), resp.Content(), nil
}

// ListPoliciesResp decodes a ListPolicies response body.
func ListPoliciesResp(body []byte) ([]*iam_pb.Policy, error) {
	resp, err := WrapListPoliciesResponse(body)
	if err != nil {
		return nil, err
	}
	n := resp.PoliciesLen()
	out := make([]*iam_pb.Policy, n)
	for i := 0; i < n; i++ {
		out[i] = policyToPB(resp.PoliciesAt(i))
	}
	return out, nil
}

// GetServiceAccountResp decodes a GetServiceAccount response body.
func GetServiceAccountResp(body []byte) (*iam_pb.ServiceAccount, error) {
	resp, err := WrapGetServiceAccountResponse(body)
	if err != nil {
		return nil, err
	}
	if resp.o.Object(getServiceAccountResponseServiceAccountOff).IsNull() {
		return nil, nil
	}
	return ServiceAccountToPB(resp.ServiceAccount()), nil
}

// GetServiceAccountByAccessKeyResp decodes a GetServiceAccountByAccessKey response body.
func GetServiceAccountByAccessKeyResp(body []byte) (*iam_pb.ServiceAccount, error) {
	resp, err := WrapGetServiceAccountByAccessKeyResponse(body)
	if err != nil {
		return nil, err
	}
	if resp.o.Object(getServiceAccountByAccessKeyResponseServiceAccountOff).IsNull() {
		return nil, nil
	}
	return ServiceAccountToPB(resp.ServiceAccount()), nil
}

// ListServiceAccountsResp decodes a ListServiceAccounts response body.
func ListServiceAccountsResp(body []byte) ([]*iam_pb.ServiceAccount, error) {
	resp, err := WrapListServiceAccountsResponse(body)
	if err != nil {
		return nil, err
	}
	n := resp.ServiceAccountsLen()
	out := make([]*iam_pb.ServiceAccount, n)
	for i := 0; i < n; i++ {
		out[i] = ServiceAccountToPB(resp.ServiceAccountsAt(i))
	}
	return out, nil
}

// ---- exported pb -> Input wrappers (used by client call sites) ----

// IdentityInputFromPB exposes pbIdentityInput to call sites.
func IdentityInputFromPB(id *iam_pb.Identity) *IdentityInput { return pbIdentityInput(id) }

// ServiceAccountInputFromPB exposes pbServiceAccountInput to call sites.
func ServiceAccountInputFromPB(sa *iam_pb.ServiceAccount) *ServiceAccountInput {
	return pbServiceAccountInput(sa)
}

// CredentialInputFromPB exposes pbCredentialInput to call sites.
func CredentialInputFromPB(c *iam_pb.Credential) *CredentialInput { return pbCredentialInput(c) }

// ConfigurationInputFromPB exposes pbConfigurationInput to call sites.
func ConfigurationInputFromPB(c *iam_pb.S3ApiConfiguration) *S3ApiConfigurationInput {
	return pbConfigurationInput(c)
}
