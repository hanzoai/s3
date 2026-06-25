package shell

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/hanzoai/s3/s3/iam"
	"github.com/hanzoai/s3/s3/pb/iam_pb"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3UserProvision{})
}

type commandS3UserProvision struct {
}

func (c *commandS3UserProvision) Name() string {
	return "s3.user.provision"
}

func (c *commandS3UserProvision) Help() string {
	return `create a user with a bucket policy in one step

	s3.user.provision -name <username> -bucket <bucket_name> -role readwrite
	s3.user.provision -name <username> -bucket <bucket_name> -role readonly

	Convenience wrapper that performs these steps:
	  1. Creates an IAM policy for the bucket and role
	  2. Creates the user with auto-generated credentials
	  3. Attaches the policy to the user

	Roles:
	  readonly   - s3:GetObject, s3:ListBucket
	  readwrite  - s3:GetObject, s3:PutObject, s3:DeleteObject, s3:ListBucket
	  admin      - s3:* (full access to the bucket)
`
}

func (c *commandS3UserProvision) HasTag(CommandTag) bool {
	return false
}

var rolePolicies = map[string][]string{
	"readonly":  {"s3:GetObject"},
	"readwrite": {"s3:GetObject", "s3:PutObject", "s3:DeleteObject"},
	"admin":     {"s3:*"},
}

func (c *commandS3UserProvision) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	f := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	name := f.String("name", "", "user name")
	bucket := f.String("bucket", "", "bucket name")
	role := f.String("role", "", "role: readonly, readwrite, or admin")
	if err := f.Parse(args); err != nil {
		return err
	}

	if *name == "" {
		return fmt.Errorf("-name is required")
	}
	if *bucket == "" {
		return fmt.Errorf("-bucket is required")
	}
	if strings.ContainsAny(*bucket, "*?") {
		return fmt.Errorf("-bucket must be a literal bucket name, not a wildcard pattern")
	}
	if *role == "" {
		return fmt.Errorf("-role is required (readonly, readwrite, admin)")
	}

	actions, ok := rolePolicies[*role]
	if !ok {
		return fmt.Errorf("unknown role %q: must be readonly, readwrite, or admin", *role)
	}

	policyName := fmt.Sprintf("%s-%s-%s", *bucket, *name, *role)

	// Build the policy document
	bucketActions := []string{"s3:ListBucket"}
	if *role == "admin" {
		bucketActions = []string{"s3:*"}
	}
	policyDoc := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect":   "Allow",
				"Action":   actions,
				"Resource": []string{fmt.Sprintf("arn:aws:s3:::%s/*", *bucket)},
			},
			{
				"Effect":   "Allow",
				"Action":   bucketActions,
				"Resource": []string{fmt.Sprintf("arn:aws:s3:::%s", *bucket)},
			},
		},
	}
	policyJSON, err := json.Marshal(policyDoc)
	if err != nil {
		return fmt.Errorf("marshal policy: %v", err)
	}

	var ak, sk string
	var userCreated bool

	err = commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		// Step 0: Check if user already exists. Over the ZAP transport a missing
		// user surfaces as a generic call error (the gRPC NotFound path), so any
		// GetUser error is treated as "user does not exist".
		var existingIdentity *iam_pb.Identity
		if _, body, getErr := client.GetUser(context.Background(), iamwire.NewGetUserRequest(iamwire.GetUserRequestInput{Username: *name})); getErr == nil {
			identity, derr := iamwire.GetUserResp(body)
			if derr != nil {
				return derr
			}
			if identity != nil {
				existingIdentity = identity
				fmt.Fprintf(writer, "User %q already exists, adding policy\n", *name)
			}
		}

		// Step 1: Create policy
		_, _, err := client.PutPolicy(context.Background(), iamwire.NewPutPolicyRequest(iamwire.PutPolicyRequestInput{
			Name:    policyName,
			Content: string(policyJSON),
		}))
		if err != nil {
			return fmt.Errorf("create policy: %v", err)
		}
		fmt.Fprintf(writer, "Created policy %q\n", policyName)

		// rollbackPolicy removes the policy we just created. Used when a later
		// step fails, to avoid leaving the policy orphaned.
		rollbackPolicy := func() {
			if _, _, delErr := client.DeletePolicy(context.Background(), iamwire.NewDeletePolicyRequest(iamwire.DeletePolicyRequestInput{Name: policyName})); delErr != nil {
				fmt.Fprintf(writer, "Warning: failed to rollback policy %q: %v\n", policyName, delErr)
			}
		}

		if existingIdentity != nil {
			// User exists: attach the new policy if not already present
			for _, pn := range existingIdentity.PolicyNames {
				if pn == policyName {
					fmt.Fprintf(writer, "Policy %q already attached to user %q\n", policyName, *name)
					return nil
				}
			}
			existingIdentity.PolicyNames = append(existingIdentity.PolicyNames, policyName)
			_, _, err = client.UpdateUser(context.Background(), iamwire.NewUpdateUserRequest(iamwire.UpdateUserRequestInput{Username: *name, Identity: iamwire.IdentityInputFromPB(existingIdentity)}))
			if err != nil {
				rollbackPolicy()
				return fmt.Errorf("attach policy to existing user: %w", err)
			}
			fmt.Fprintf(writer, "Attached policy %q to existing user %q\n", policyName, *name)
		} else {
			// Step 2: Create new user with credentials
			ak, err = iam.GenerateRandomString(iam.AccessKeyIdLength, iam.CharsetUpper)
			if err != nil {
				rollbackPolicy()
				return fmt.Errorf("generate access key: %v", err)
			}
			sk, err = iam.GenerateSecretAccessKey()
			if err != nil {
				rollbackPolicy()
				return fmt.Errorf("generate secret key: %v", err)
			}

			identity := &iam_pb.Identity{
				Name: *name,
				Credentials: []*iam_pb.Credential{
					{
						AccessKey: ak,
						SecretKey: sk,
						Status:    iam.AccessKeyStatusActive,
					},
				},
				PolicyNames: []string{policyName},
			}
			_, _, err = client.CreateUser(context.Background(), iamwire.NewCreateUserRequest(iamwire.CreateUserRequestInput{Identity: iamwire.IdentityInputFromPB(identity)}))
			if err != nil {
				rollbackPolicy()
				return fmt.Errorf("create user: %w", err)
			}
			userCreated = true
			fmt.Fprintf(writer, "Created user %q with policy %q attached\n", *name, policyName)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if userCreated {
		fmt.Fprintln(writer)
		fmt.Fprintf(writer, "Access Key: %s\n", ak)
		fmt.Fprintf(writer, "Secret Key: %s\n", sk)
		fmt.Fprintln(writer)
		fmt.Fprintln(writer, "Save these credentials - the secret key cannot be retrieved later.")
	}
	return nil
}
