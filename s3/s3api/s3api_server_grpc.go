package s3api

import (
	"fmt"

	"github.com/hanzoai/s3/s3/glog"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
	"github.com/hanzoai/s3/s3/wire/iamadapt"
)

// HanzoS3IamCacheHandler implementation over the ZAP transport.
// This service is dedicated to UNIDIRECTIONAL updates from Filer to S3 Server.
// S3 Server acts purely as a cache. Each method takes/returns an opaque iamwire
// buffer; Wrap the request, mutate the cache, and build the reply with New*.

func (s3a *S3ApiServer) PutIdentity(req []byte) ([]byte, error) {
	request, err := iamwire.WrapPutIdentityRequest(req)
	if err != nil {
		return nil, err
	}
	identity := iamadapt.IdentityFromWire(request.Identity())
	if identity.Name == "" {
		return nil, fmt.Errorf("identity is required")
	}
	// Direct in-memory cache update
	glog.V(1).Infof("IAM: received identity update for %s", identity.Name)
	if err := s3a.iam.UpsertIdentity(identity); err != nil {
		glog.Errorf("failed to update identity cache for %s: %v", identity.Name, err)
		return nil, fmt.Errorf("failed to update identity cache: %v", err)
	}
	return iamwire.NewPutIdentityResponse(iamwire.PutIdentityResponseInput{}), nil
}

func (s3a *S3ApiServer) RemoveIdentity(req []byte) ([]byte, error) {
	request, err := iamwire.WrapRemoveIdentityRequest(req)
	if err != nil {
		return nil, err
	}
	username := request.Username()
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	// Direct in-memory cache update
	glog.V(1).Infof("IAM: received identity removal for %s", username)
	s3a.iam.RemoveIdentity(username)
	return iamwire.NewRemoveIdentityResponse(iamwire.RemoveIdentityResponseInput{}), nil
}

func (s3a *S3ApiServer) PutPolicy(req []byte) ([]byte, error) {
	request, err := iamwire.WrapPutPolicyRequest(req)
	if err != nil {
		return nil, err
	}
	name, content := request.Name(), request.Content()
	if name == "" {
		return nil, fmt.Errorf("policy name is required")
	}
	if content == "" {
		return nil, fmt.Errorf("policy content is required")
	}

	// Update IAM policy cache
	glog.V(1).Infof("IAM: received policy update for %s", name)
	if s3a.iam == nil {
		return nil, fmt.Errorf("IAM not initialized")
	}

	if err := s3a.iam.PutPolicy(name, content); err != nil {
		glog.Errorf("failed to update policy cache for %s: %v", name, err)
		return nil, fmt.Errorf("failed to update policy cache: %v", err)
	}
	return iamwire.NewPutPolicyResponse(iamwire.PutPolicyResponseInput{}), nil
}

func (s3a *S3ApiServer) DeletePolicy(req []byte) ([]byte, error) {
	request, err := iamwire.WrapDeletePolicyRequest(req)
	if err != nil {
		return nil, err
	}
	name := request.Name()
	if name == "" {
		return nil, fmt.Errorf("policy name is required")
	}

	// Delete from IAM policy cache
	glog.V(1).Infof("IAM: received policy removal for %s", name)
	if s3a.iam == nil {
		return nil, fmt.Errorf("IAM not initialized")
	}

	if err := s3a.iam.DeletePolicy(name); err != nil {
		glog.Errorf("failed to delete policy cache for %s: %v", name, err)
		return nil, fmt.Errorf("failed to delete policy cache: %v", err)
	}
	return iamwire.NewDeletePolicyResponse(iamwire.DeletePolicyResponseInput{}), nil
}

func (s3a *S3ApiServer) GetPolicy(req []byte) ([]byte, error) {
	request, err := iamwire.WrapGetPolicyRequest(req)
	if err != nil {
		return nil, err
	}
	name := request.Name()
	if name == "" {
		return nil, fmt.Errorf("policy name is required")
	}
	if s3a.iam == nil {
		return nil, fmt.Errorf("IAM not initialized")
	}
	policy, err := s3a.iam.GetPolicy(name)
	if err != nil {
		return iamwire.NewGetPolicyResponse(iamwire.GetPolicyResponseInput{}), nil // Not found is fine for cache
	}
	return iamwire.NewGetPolicyResponse(iamwire.GetPolicyResponseInput{
		Name:    policy.Name,
		Content: policy.Content,
	}), nil
}

func (s3a *S3ApiServer) PutGroup(req []byte) ([]byte, error) {
	request, err := iamwire.WrapPutGroupRequest(req)
	if err != nil {
		return nil, err
	}
	group := iamadapt.GroupFromWire(request.Group())
	if group.Name == "" {
		return nil, fmt.Errorf("group is required")
	}
	glog.V(1).Infof("IAM: received group update for %s", group.Name)
	if s3a.iam == nil {
		return nil, fmt.Errorf("IAM not initialized")
	}
	if err := s3a.iam.PutGroup(group); err != nil {
		glog.Errorf("failed to update group cache for %s: %v", group.Name, err)
		return nil, fmt.Errorf("failed to update group cache: %v", err)
	}
	return iamwire.NewPutGroupResponse(iamwire.PutGroupResponseInput{}), nil
}

func (s3a *S3ApiServer) RemoveGroup(req []byte) ([]byte, error) {
	request, err := iamwire.WrapRemoveGroupRequest(req)
	if err != nil {
		return nil, err
	}
	groupName := request.GroupName()
	if groupName == "" {
		return nil, fmt.Errorf("group name is required")
	}
	glog.V(1).Infof("IAM: received group removal for %s", groupName)
	if s3a.iam == nil {
		return nil, fmt.Errorf("IAM not initialized")
	}
	s3a.iam.RemoveGroup(groupName)
	return iamwire.NewRemoveGroupResponse(iamwire.RemoveGroupResponseInput{}), nil
}

func (s3a *S3ApiServer) ListPolicies(req []byte) ([]byte, error) {
	if _, err := iamwire.WrapListPoliciesRequest(req); err != nil {
		return nil, err
	}
	if s3a.iam == nil {
		return nil, fmt.Errorf("IAM not initialized")
	}
	var policies []iamwire.PolicyInput
	for _, policy := range s3a.iam.ListPolicies() {
		if policy == nil {
			continue
		}
		policies = append(policies, iamwire.PolicyInput{Name: policy.Name, Content: policy.Content})
	}
	return iamwire.NewListPoliciesResponse(iamwire.ListPoliciesResponseInput{Policies: policies}), nil
}
