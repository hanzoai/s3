package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hanzoai/s3/s3/s3api/policy_engine"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func (store *IamGrpcStore) GetPolicies(ctx context.Context) (map[string]policy_engine.PolicyDocument, error) {
	policies := make(map[string]policy_engine.PolicyDocument)
	err := store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.ListPolicies(ctx, iamwire.NewListPoliciesRequest(iamwire.ListPoliciesRequestInput{}))
		if err != nil {
			return err
		}
		list, err := iamwire.ListPoliciesResp(body)
		if err != nil {
			return err
		}
		for _, p := range list {
			var doc policy_engine.PolicyDocument
			if err := json.Unmarshal([]byte(p.Content), &doc); err != nil {
				return fmt.Errorf("failed to unmarshal policy %s: %v", p.Name, err)
			}
			policies[p.Name] = doc
		}
		return nil
	})
	return policies, err
}

func (store *IamGrpcStore) PutPolicy(ctx context.Context, name string, document policy_engine.PolicyDocument) error {
	content, err := json.Marshal(document)
	if err != nil {
		return err
	}
	return store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.PutPolicy(ctx, iamwire.NewPutPolicyRequest(iamwire.PutPolicyRequestInput{
			Name:    name,
			Content: string(content),
		}))
		return err
	})
}

func (store *IamGrpcStore) DeletePolicy(ctx context.Context, name string) error {
	return store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.DeletePolicy(ctx, iamwire.NewDeletePolicyRequest(iamwire.DeletePolicyRequestInput{
			Name: name,
		}))
		return err
	})
}

func (store *IamGrpcStore) GetPolicy(ctx context.Context, name string) (*policy_engine.PolicyDocument, error) {
	var doc policy_engine.PolicyDocument
	err := store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetPolicy(ctx, iamwire.NewGetPolicyRequest(iamwire.GetPolicyRequestInput{
			Name: name,
		}))
		if err != nil {
			return err
		}
		_, content, err := iamwire.GetPolicyResp(body)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(content), &doc)
	})
	if err != nil {
		// A missing policy errors on the filer side; over the ZAP transport that
		// surfaces as a generic call error, so treat any failure as "not found"
		// and return nil (consistent with other stores).
		return nil, nil
	}
	return &doc, nil
}

// CreatePolicy creates a new policy (delegates to PutPolicy)
func (store *IamGrpcStore) CreatePolicy(ctx context.Context, name string, document policy_engine.PolicyDocument) error {
	return store.PutPolicy(ctx, name, document)
}

// ListPolicyNames retrieves names of all IAM policies via gRPC.
func (store *IamGrpcStore) ListPolicyNames(ctx context.Context) ([]string, error) {
	var names []string
	err := store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.ListPolicies(ctx, iamwire.NewListPoliciesRequest(iamwire.ListPoliciesRequestInput{}))
		if err != nil {
			return err
		}
		list, err := iamwire.ListPoliciesResp(body)
		if err != nil {
			return err
		}
		for _, policy := range list {
			names = append(names, policy.Name)
		}
		return nil
	})
	return names, err
}

// UpdatePolicy updates an existing policy (delegates to PutPolicy)
func (store *IamGrpcStore) UpdatePolicy(ctx context.Context, name string, document policy_engine.PolicyDocument) error {
	return store.PutPolicy(ctx, name, document)
}
