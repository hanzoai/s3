// Code generated from iam.proto; DO NOT EDIT.

package iamwire

import (
	zap "github.com/zap-proto/go"
)

// This file holds the shared IAM data types that the request/response
// messages embed: the leaf records (UserTag, Credential, Account, Policy),
// the composite records (Group, Identity, ServiceAccount), and the
// recursive configuration tree (S3ApiConfiguration).
//
// Each type T exposes:
//   - a zero-copy view  type T struct{ o zap.Object }  + WrapT(b) (T, error)
//   - typed field accessors reading off fixed byte offsets
//   - a TInput struct + NewT(in) []byte standalone builder
//   - an internal appendT(b, in) *zap.ObjectBuilder that builds T into an
//     existing buffer (children/lists laid down first) and returns its
//     unfinished ObjectBuilder, so a parent can embed T as a nested object
//     (childOff := appendT(b, in).Finish()) or NewT can root it.

// --- UserTag ---

const (
	userTagKeyOff   = 0
	userTagValueOff = 8
	userTagSize     = 16
)

// UserTag is a zero-copy view into a ZAP-encoded UserTag message.
type UserTag struct{ o zap.Object }

// WrapUserTag parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapUserTag(b []byte) (UserTag, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return UserTag{}, err
	}
	return UserTag{o: m.Root()}, nil
}

// Key reads the key field (proto field 1, string).
func (t UserTag) Key() string { return t.o.Text(userTagKeyOff) }

// Value reads the value field (proto field 2, string).
func (t UserTag) Value() string { return t.o.Text(userTagValueOff) }

// UserTagInput collects the field values for NewUserTag.
type UserTagInput struct {
	Key   string
	Value string
}

func appendUserTag(b *zap.Builder, in UserTagInput) *zap.ObjectBuilder {
	ob := b.StartObject(userTagSize)
	ob.SetText(userTagKeyOff, in.Key)
	ob.SetText(userTagValueOff, in.Value)
	return ob
}

// NewUserTag builds a ZAP-encoded UserTag message from in and returns the bytes.
func NewUserTag(in UserTagInput) []byte {
	b := zap.NewBuilder(256)
	appendUserTag(b, in).FinishAsRoot()
	return b.Finish()
}

// --- Credential ---

const (
	credentialAccessKeyOff = 0
	credentialSecretKeyOff = 8
	credentialStatusOff    = 16
	credentialSize         = 24
)

// Credential is a zero-copy view into a ZAP-encoded Credential message.
type Credential struct{ o zap.Object }

// WrapCredential parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapCredential(b []byte) (Credential, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Credential{}, err
	}
	return Credential{o: m.Root()}, nil
}

// AccessKey reads the access_key field (proto field 1, string).
func (t Credential) AccessKey() string { return t.o.Text(credentialAccessKeyOff) }

// SecretKey reads the secret_key field (proto field 2, string).
func (t Credential) SecretKey() string { return t.o.Text(credentialSecretKeyOff) }

// Status reads the status field (proto field 3, string).
func (t Credential) Status() string { return t.o.Text(credentialStatusOff) }

// CredentialInput collects the field values for NewCredential.
type CredentialInput struct {
	AccessKey string
	SecretKey string
	Status    string
}

func appendCredential(b *zap.Builder, in CredentialInput) *zap.ObjectBuilder {
	ob := b.StartObject(credentialSize)
	ob.SetText(credentialAccessKeyOff, in.AccessKey)
	ob.SetText(credentialSecretKeyOff, in.SecretKey)
	ob.SetText(credentialStatusOff, in.Status)
	return ob
}

// NewCredential builds a ZAP-encoded Credential message from in and returns the
// bytes.
func NewCredential(in CredentialInput) []byte {
	b := zap.NewBuilder(256)
	appendCredential(b, in).FinishAsRoot()
	return b.Finish()
}

// --- Account ---

const (
	accountIDOff           = 0
	accountDisplayNameOff  = 8
	accountEmailAddressOff = 16
	accountSize            = 24
)

// Account is a zero-copy view into a ZAP-encoded Account message.
type Account struct{ o zap.Object }

// WrapAccount parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapAccount(b []byte) (Account, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Account{}, err
	}
	return Account{o: m.Root()}, nil
}

// IsNull reports whether this view backs an absent account sub-message. Field
// reads on a null view dereference a nil backing message and panic, so callers
// that may observe an absent account must gate on IsNull first.
func (t Account) IsNull() bool { return t.o.IsNull() }

// ID reads the id field (proto field 1, string).
func (t Account) ID() string { return t.o.Text(accountIDOff) }

// DisplayName reads the display_name field (proto field 2, string).
func (t Account) DisplayName() string { return t.o.Text(accountDisplayNameOff) }

// EmailAddress reads the email_address field (proto field 3, string).
func (t Account) EmailAddress() string { return t.o.Text(accountEmailAddressOff) }

// AccountInput collects the field values for NewAccount.
type AccountInput struct {
	ID           string
	DisplayName  string
	EmailAddress string
}

func appendAccount(b *zap.Builder, in AccountInput) *zap.ObjectBuilder {
	ob := b.StartObject(accountSize)
	ob.SetText(accountIDOff, in.ID)
	ob.SetText(accountDisplayNameOff, in.DisplayName)
	ob.SetText(accountEmailAddressOff, in.EmailAddress)
	return ob
}

// NewAccount builds a ZAP-encoded Account message from in and returns the bytes.
func NewAccount(in AccountInput) []byte {
	b := zap.NewBuilder(256)
	appendAccount(b, in).FinishAsRoot()
	return b.Finish()
}

// --- Policy ---

const (
	policyNameOff    = 0
	policyContentOff = 8
	policySize       = 16
)

// Policy is a zero-copy view into a ZAP-encoded Policy message.
type Policy struct{ o zap.Object }

// WrapPolicy parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapPolicy(b []byte) (Policy, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Policy{}, err
	}
	return Policy{o: m.Root()}, nil
}

// Name reads the name field (proto field 1, string).
func (t Policy) Name() string { return t.o.Text(policyNameOff) }

// Content reads the content field (proto field 2, string, JSON policy document).
func (t Policy) Content() string { return t.o.Text(policyContentOff) }

// PolicyInput collects the field values for NewPolicy.
type PolicyInput struct {
	Name    string
	Content string
}

func appendPolicy(b *zap.Builder, in PolicyInput) *zap.ObjectBuilder {
	ob := b.StartObject(policySize)
	ob.SetText(policyNameOff, in.Name)
	ob.SetText(policyContentOff, in.Content)
	return ob
}

// NewPolicy builds a ZAP-encoded Policy message from in and returns the bytes.
func NewPolicy(in PolicyInput) []byte {
	b := zap.NewBuilder(256)
	appendPolicy(b, in).FinishAsRoot()
	return b.Finish()
}

// --- Group ---

const (
	groupNameOff        = 0
	groupMembersOff     = 8
	groupPolicyNamesOff = 16
	groupDisabledOff    = 24
	groupSize           = 32
)

// Group is a zero-copy view into a ZAP-encoded Group message.
type Group struct{ o zap.Object }

// WrapGroup parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapGroup(b []byte) (Group, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Group{}, err
	}
	return Group{o: m.Root()}, nil
}

// Name reads the name field (proto field 1, string).
func (t Group) Name() string { return t.o.Text(groupNameOff) }

// MembersLen returns the number of members (proto field 2, repeated string).
func (t Group) MembersLen() int { return t.o.List(groupMembersOff).Len() }

// MembersAt returns member i (proto field 2, repeated string, usernames).
func (t Group) MembersAt(i int) string { return string(t.o.List(groupMembersOff).BytesAt(i)) }

// PolicyNamesLen returns the number of attached managed policy names
// (proto field 3, repeated string).
func (t Group) PolicyNamesLen() int { return t.o.List(groupPolicyNamesOff).Len() }

// PolicyNamesAt returns policy name i (proto field 3, repeated string).
func (t Group) PolicyNamesAt(i int) string {
	return string(t.o.List(groupPolicyNamesOff).BytesAt(i))
}

// Disabled reads the disabled field (proto field 4, bool).
func (t Group) Disabled() bool { return t.o.Bool(groupDisabledOff) }

// GroupInput collects the field values for NewGroup.
type GroupInput struct {
	Name        string
	Members     []string
	PolicyNames []string
	Disabled    bool
}

func appendGroup(b *zap.Builder, in GroupInput) *zap.ObjectBuilder {
	ob := b.StartObject(groupSize)
	ob.SetText(groupNameOff, in.Name)
	ob.SetBool(groupDisabledOff, in.Disabled)
	if len(in.Members) > 0 {
		lb := b.StartList(0)
		for _, s := range in.Members {
			lb.AddObjectBytes([]byte(s))
		}
		off, n := lb.Finish()
		ob.SetList(groupMembersOff, off, n)
	}
	if len(in.PolicyNames) > 0 {
		lb := b.StartList(0)
		for _, s := range in.PolicyNames {
			lb.AddObjectBytes([]byte(s))
		}
		off, n := lb.Finish()
		ob.SetList(groupPolicyNamesOff, off, n)
	}
	return ob
}

// NewGroup builds a ZAP-encoded Group message from in and returns the bytes.
func NewGroup(in GroupInput) []byte {
	b := zap.NewBuilder(256)
	appendGroup(b, in).FinishAsRoot()
	return b.Finish()
}

// --- Identity ---

const (
	identityNameOff              = 0
	identityCredentialsOff       = 8
	identityActionsOff           = 16
	identityAccountOff           = 24
	identityDisabledOff          = 28
	identityServiceAccountIDsOff = 32
	identityPolicyNamesOff       = 40
	identityIsStaticOff          = 48
	identityTagsOff              = 56
	identitySize                 = 64
)

// Identity is a zero-copy view into a ZAP-encoded Identity message.
type Identity struct{ o zap.Object }

// WrapIdentity parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapIdentity(b []byte) (Identity, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Identity{}, err
	}
	return Identity{o: m.Root()}, nil
}

// Name reads the name field (proto field 1, string).
func (t Identity) Name() string { return t.o.Text(identityNameOff) }

// CredentialsLen returns the number of credentials (proto field 2, repeated Credential).
func (t Identity) CredentialsLen() int { return t.o.List(identityCredentialsOff).Len() }

// CredentialsAt returns credential i (proto field 2, repeated Credential).
func (t Identity) CredentialsAt(i int) Credential {
	return Credential{o: t.o.List(identityCredentialsOff).ObjectAt(i)}
}

// ActionsLen returns the number of actions (proto field 3, repeated string).
func (t Identity) ActionsLen() int { return t.o.List(identityActionsOff).Len() }

// ActionsAt returns action i (proto field 3, repeated string).
func (t Identity) ActionsAt(i int) string { return string(t.o.List(identityActionsOff).BytesAt(i)) }

// Account reads the account field (proto field 4, Account).
func (t Identity) Account() Account { return Account{o: t.o.Object(identityAccountOff)} }

// Disabled reads the disabled field (proto field 5, bool). false = enabled (default).
func (t Identity) Disabled() bool { return t.o.Bool(identityDisabledOff) }

// ServiceAccountIDsLen returns the number of owned service-account IDs
// (proto field 6, repeated string).
func (t Identity) ServiceAccountIDsLen() int { return t.o.List(identityServiceAccountIDsOff).Len() }

// ServiceAccountIDsAt returns service-account ID i (proto field 6, repeated string).
func (t Identity) ServiceAccountIDsAt(i int) string {
	return string(t.o.List(identityServiceAccountIDsOff).BytesAt(i))
}

// PolicyNamesLen returns the number of policy names (proto field 7, repeated string).
func (t Identity) PolicyNamesLen() int { return t.o.List(identityPolicyNamesOff).Len() }

// PolicyNamesAt returns policy name i (proto field 7, repeated string).
func (t Identity) PolicyNamesAt(i int) string {
	return string(t.o.List(identityPolicyNamesOff).BytesAt(i))
}

// IsStatic reads the is_static field (proto field 8, bool). Loaded from static
// config file (read-only, not editable via API).
func (t Identity) IsStatic() bool { return t.o.Bool(identityIsStaticOff) }

// TagsLen returns the number of tags (proto field 9, repeated UserTag).
func (t Identity) TagsLen() int { return t.o.List(identityTagsOff).Len() }

// TagsAt returns tag i (proto field 9, repeated UserTag).
func (t Identity) TagsAt(i int) UserTag {
	return UserTag{o: t.o.List(identityTagsOff).ObjectAt(i)}
}

// IdentityInput collects the field values for NewIdentity.
type IdentityInput struct {
	Name              string
	Credentials       []CredentialInput
	Actions           []string
	Account           *AccountInput
	Disabled          bool
	ServiceAccountIDs []string
	PolicyNames       []string
	IsStatic          bool
	Tags              []UserTagInput
}

func appendIdentity(b *zap.Builder, in IdentityInput) *zap.ObjectBuilder {
	// Nested singular message laid down before the parent object so its
	// pointer can resolve.
	accountOff := 0
	if in.Account != nil {
		accountOff = appendAccount(b, *in.Account).Finish()
	}
	ob := b.StartObject(identitySize)
	ob.SetText(identityNameOff, in.Name)
	ob.SetBool(identityDisabledOff, in.Disabled)
	ob.SetBool(identityIsStaticOff, in.IsStatic)
	if accountOff != 0 {
		ob.SetObject(identityAccountOff, accountOff)
	}
	if len(in.Credentials) > 0 {
		lb := b.StartList(0)
		for i := range in.Credentials {
			lb.AddObjectBytes(NewCredential(in.Credentials[i]))
		}
		off, n := lb.Finish()
		ob.SetList(identityCredentialsOff, off, n)
	}
	if len(in.Actions) > 0 {
		lb := b.StartList(0)
		for _, s := range in.Actions {
			lb.AddObjectBytes([]byte(s))
		}
		off, n := lb.Finish()
		ob.SetList(identityActionsOff, off, n)
	}
	if len(in.ServiceAccountIDs) > 0 {
		lb := b.StartList(0)
		for _, s := range in.ServiceAccountIDs {
			lb.AddObjectBytes([]byte(s))
		}
		off, n := lb.Finish()
		ob.SetList(identityServiceAccountIDsOff, off, n)
	}
	if len(in.PolicyNames) > 0 {
		lb := b.StartList(0)
		for _, s := range in.PolicyNames {
			lb.AddObjectBytes([]byte(s))
		}
		off, n := lb.Finish()
		ob.SetList(identityPolicyNamesOff, off, n)
	}
	if len(in.Tags) > 0 {
		lb := b.StartList(0)
		for i := range in.Tags {
			lb.AddObjectBytes(NewUserTag(in.Tags[i]))
		}
		off, n := lb.Finish()
		ob.SetList(identityTagsOff, off, n)
	}
	return ob
}

// NewIdentity builds a ZAP-encoded Identity message from in and returns the bytes.
func NewIdentity(in IdentityInput) []byte {
	b := zap.NewBuilder(512)
	appendIdentity(b, in).FinishAsRoot()
	return b.Finish()
}

// --- ServiceAccount ---

const (
	serviceAccountIDOff          = 0
	serviceAccountParentUserOff  = 8
	serviceAccountDescriptionOff = 16
	serviceAccountCredentialOff  = 24
	serviceAccountActionsOff     = 32
	serviceAccountExpirationOff  = 40
	serviceAccountDisabledOff    = 48
	serviceAccountCreatedAtOff   = 56
	serviceAccountCreatedByOff   = 64
	serviceAccountSize           = 72
)

// ServiceAccount is a zero-copy view into a ZAP-encoded ServiceAccount message.
// A service account holds special credentials for applications, linked to a
// parent user with restricted permissions.
type ServiceAccount struct{ o zap.Object }

// WrapServiceAccount parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapServiceAccount(b []byte) (ServiceAccount, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ServiceAccount{}, err
	}
	return ServiceAccount{o: m.Root()}, nil
}

// ID reads the id field (proto field 1, string). Unique identifier (e.g. "sa-xxxxx").
func (t ServiceAccount) ID() string { return t.o.Text(serviceAccountIDOff) }

// ParentUser reads the parent_user field (proto field 2, string). Parent identity name.
func (t ServiceAccount) ParentUser() string { return t.o.Text(serviceAccountParentUserOff) }

// Description reads the description field (proto field 3, string).
func (t ServiceAccount) Description() string { return t.o.Text(serviceAccountDescriptionOff) }

// Credential reads the credential field (proto field 4, Credential).
func (t ServiceAccount) Credential() Credential {
	return Credential{o: t.o.Object(serviceAccountCredentialOff)}
}

// ActionsLen returns the number of allowed actions (proto field 5, repeated string).
func (t ServiceAccount) ActionsLen() int { return t.o.List(serviceAccountActionsOff).Len() }

// ActionsAt returns action i (proto field 5, repeated string, subset of parent).
func (t ServiceAccount) ActionsAt(i int) string {
	return string(t.o.List(serviceAccountActionsOff).BytesAt(i))
}

// Expiration reads the expiration field (proto field 6, int64). Unix timestamp,
// 0 = no expiration.
func (t ServiceAccount) Expiration() int64 { return t.o.Int64(serviceAccountExpirationOff) }

// Disabled reads the disabled field (proto field 7, bool). false = enabled (default).
func (t ServiceAccount) Disabled() bool { return t.o.Bool(serviceAccountDisabledOff) }

// CreatedAt reads the created_at field (proto field 8, int64). Creation timestamp.
func (t ServiceAccount) CreatedAt() int64 { return t.o.Int64(serviceAccountCreatedAtOff) }

// CreatedBy reads the created_by field (proto field 9, string).
func (t ServiceAccount) CreatedBy() string { return t.o.Text(serviceAccountCreatedByOff) }

// ServiceAccountInput collects the field values for NewServiceAccount.
type ServiceAccountInput struct {
	ID          string
	ParentUser  string
	Description string
	Credential  *CredentialInput
	Actions     []string
	Expiration  int64
	Disabled    bool
	CreatedAt   int64
	CreatedBy   string
}

func appendServiceAccount(b *zap.Builder, in ServiceAccountInput) *zap.ObjectBuilder {
	credentialOff := 0
	if in.Credential != nil {
		credentialOff = appendCredential(b, *in.Credential).Finish()
	}
	ob := b.StartObject(serviceAccountSize)
	ob.SetText(serviceAccountIDOff, in.ID)
	ob.SetText(serviceAccountParentUserOff, in.ParentUser)
	ob.SetText(serviceAccountDescriptionOff, in.Description)
	ob.SetInt64(serviceAccountExpirationOff, in.Expiration)
	ob.SetBool(serviceAccountDisabledOff, in.Disabled)
	ob.SetInt64(serviceAccountCreatedAtOff, in.CreatedAt)
	ob.SetText(serviceAccountCreatedByOff, in.CreatedBy)
	if credentialOff != 0 {
		ob.SetObject(serviceAccountCredentialOff, credentialOff)
	}
	if len(in.Actions) > 0 {
		lb := b.StartList(0)
		for _, s := range in.Actions {
			lb.AddObjectBytes([]byte(s))
		}
		off, n := lb.Finish()
		ob.SetList(serviceAccountActionsOff, off, n)
	}
	return ob
}

// NewServiceAccount builds a ZAP-encoded ServiceAccount message from in and
// returns the bytes.
func NewServiceAccount(in ServiceAccountInput) []byte {
	b := zap.NewBuilder(512)
	appendServiceAccount(b, in).FinishAsRoot()
	return b.Finish()
}

// --- S3ApiConfiguration ---

const (
	s3APIConfigurationIdentitiesOff      = 0
	s3APIConfigurationAccountsOff        = 8
	s3APIConfigurationServiceAccountsOff = 16
	s3APIConfigurationPoliciesOff        = 24
	s3APIConfigurationGroupsOff          = 32
	s3APIConfigurationSize               = 40
)

// S3ApiConfiguration is a zero-copy view into a ZAP-encoded S3ApiConfiguration
// message: the recursive IAM configuration tree.
type S3ApiConfiguration struct{ o zap.Object }

// WrapS3ApiConfiguration parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapS3ApiConfiguration(b []byte) (S3ApiConfiguration, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return S3ApiConfiguration{}, err
	}
	return S3ApiConfiguration{o: m.Root()}, nil
}

// IdentitiesLen returns the number of identities (proto field 1, repeated Identity).
func (t S3ApiConfiguration) IdentitiesLen() int {
	return t.o.List(s3APIConfigurationIdentitiesOff).Len()
}

// IdentitiesAt returns identity i (proto field 1, repeated Identity).
func (t S3ApiConfiguration) IdentitiesAt(i int) Identity {
	return Identity{o: t.o.List(s3APIConfigurationIdentitiesOff).ObjectAt(i)}
}

// AccountsLen returns the number of accounts (proto field 2, repeated Account).
func (t S3ApiConfiguration) AccountsLen() int { return t.o.List(s3APIConfigurationAccountsOff).Len() }

// AccountsAt returns account i (proto field 2, repeated Account).
func (t S3ApiConfiguration) AccountsAt(i int) Account {
	return Account{o: t.o.List(s3APIConfigurationAccountsOff).ObjectAt(i)}
}

// ServiceAccountsLen returns the number of service accounts (proto field 3,
// repeated ServiceAccount).
func (t S3ApiConfiguration) ServiceAccountsLen() int {
	return t.o.List(s3APIConfigurationServiceAccountsOff).Len()
}

// ServiceAccountsAt returns service account i (proto field 3, repeated ServiceAccount).
func (t S3ApiConfiguration) ServiceAccountsAt(i int) ServiceAccount {
	return ServiceAccount{o: t.o.List(s3APIConfigurationServiceAccountsOff).ObjectAt(i)}
}

// PoliciesLen returns the number of policies (proto field 4, repeated Policy).
func (t S3ApiConfiguration) PoliciesLen() int { return t.o.List(s3APIConfigurationPoliciesOff).Len() }

// PoliciesAt returns policy i (proto field 4, repeated Policy).
func (t S3ApiConfiguration) PoliciesAt(i int) Policy {
	return Policy{o: t.o.List(s3APIConfigurationPoliciesOff).ObjectAt(i)}
}

// GroupsLen returns the number of groups (proto field 5, repeated Group).
func (t S3ApiConfiguration) GroupsLen() int { return t.o.List(s3APIConfigurationGroupsOff).Len() }

// GroupsAt returns group i (proto field 5, repeated Group).
func (t S3ApiConfiguration) GroupsAt(i int) Group {
	return Group{o: t.o.List(s3APIConfigurationGroupsOff).ObjectAt(i)}
}

// S3ApiConfigurationInput collects the field values for NewS3ApiConfiguration.
type S3ApiConfigurationInput struct {
	Identities      []IdentityInput
	Accounts        []AccountInput
	ServiceAccounts []ServiceAccountInput
	Policies        []PolicyInput
	Groups          []GroupInput
}

func appendS3ApiConfiguration(b *zap.Builder, in S3ApiConfigurationInput) *zap.ObjectBuilder {
	ob := b.StartObject(s3APIConfigurationSize)
	if len(in.Identities) > 0 {
		lb := b.StartList(0)
		for i := range in.Identities {
			lb.AddObjectBytes(NewIdentity(in.Identities[i]))
		}
		off, n := lb.Finish()
		ob.SetList(s3APIConfigurationIdentitiesOff, off, n)
	}
	if len(in.Accounts) > 0 {
		lb := b.StartList(0)
		for i := range in.Accounts {
			lb.AddObjectBytes(NewAccount(in.Accounts[i]))
		}
		off, n := lb.Finish()
		ob.SetList(s3APIConfigurationAccountsOff, off, n)
	}
	if len(in.ServiceAccounts) > 0 {
		lb := b.StartList(0)
		for i := range in.ServiceAccounts {
			lb.AddObjectBytes(NewServiceAccount(in.ServiceAccounts[i]))
		}
		off, n := lb.Finish()
		ob.SetList(s3APIConfigurationServiceAccountsOff, off, n)
	}
	if len(in.Policies) > 0 {
		lb := b.StartList(0)
		for i := range in.Policies {
			lb.AddObjectBytes(NewPolicy(in.Policies[i]))
		}
		off, n := lb.Finish()
		ob.SetList(s3APIConfigurationPoliciesOff, off, n)
	}
	if len(in.Groups) > 0 {
		lb := b.StartList(0)
		for i := range in.Groups {
			lb.AddObjectBytes(NewGroup(in.Groups[i]))
		}
		off, n := lb.Finish()
		ob.SetList(s3APIConfigurationGroupsOff, off, n)
	}
	return ob
}

// NewS3ApiConfiguration builds a ZAP-encoded S3ApiConfiguration message from in
// and returns the bytes.
func NewS3ApiConfiguration(in S3ApiConfigurationInput) []byte {
	b := zap.NewBuilder(1024)
	appendS3ApiConfiguration(b, in).FinishAsRoot()
	return b.Finish()
}
