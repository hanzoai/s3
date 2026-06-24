package grpc

import (
	"context"

	"github.com/hanzoai/s3/s3/pb/iam_pb"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func (store *IamGrpcStore) CreateServiceAccount(ctx context.Context, sa *iam_pb.ServiceAccount) error {
	return store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.CreateServiceAccount(iamwire.NewCreateServiceAccountRequest(iamwire.CreateServiceAccountRequestInput{
			ServiceAccount: iamwire.ServiceAccountInputFromPB(sa),
		}))
		return err
	})
}

func (store *IamGrpcStore) UpdateServiceAccount(ctx context.Context, id string, sa *iam_pb.ServiceAccount) error {
	return store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.UpdateServiceAccount(iamwire.NewUpdateServiceAccountRequest(iamwire.UpdateServiceAccountRequestInput{
			ID:             id,
			ServiceAccount: iamwire.ServiceAccountInputFromPB(sa),
		}))
		return err
	})
}

func (store *IamGrpcStore) DeleteServiceAccount(ctx context.Context, id string) error {
	return store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.DeleteServiceAccount(iamwire.NewDeleteServiceAccountRequest(iamwire.DeleteServiceAccountRequestInput{
			ID: id,
		}))
		return err
	})
}

func (store *IamGrpcStore) GetServiceAccount(ctx context.Context, id string) (*iam_pb.ServiceAccount, error) {
	var sa *iam_pb.ServiceAccount
	err := store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetServiceAccount(iamwire.NewGetServiceAccountRequest(iamwire.GetServiceAccountRequestInput{
			ID: id,
		}))
		if err != nil {
			return err
		}
		sa, err = iamwire.GetServiceAccountResp(body)
		return err
	})
	return sa, err
}

func (store *IamGrpcStore) ListServiceAccounts(ctx context.Context) ([]*iam_pb.ServiceAccount, error) {
	var accounts []*iam_pb.ServiceAccount
	err := store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.ListServiceAccounts(iamwire.NewListServiceAccountsRequest(iamwire.ListServiceAccountsRequestInput{}))
		if err != nil {
			return err
		}
		accounts, err = iamwire.ListServiceAccountsResp(body)
		return err
	})
	return accounts, err
}

func (store *IamGrpcStore) GetServiceAccountByAccessKey(ctx context.Context, accessKey string) (*iam_pb.ServiceAccount, error) {
	var sa *iam_pb.ServiceAccount
	err := store.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetServiceAccountByAccessKey(iamwire.NewGetServiceAccountByAccessKeyRequest(iamwire.GetServiceAccountByAccessKeyRequestInput{
			AccessKey: accessKey,
		}))
		if err != nil {
			return err
		}
		sa, err = iamwire.GetServiceAccountByAccessKeyResp(body)
		return err
	})
	return sa, err
}
