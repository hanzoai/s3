package shell

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3PolicyAttach{})
}

type commandS3PolicyAttach struct {
}

func (c *commandS3PolicyAttach) Name() string {
	return "s3.policy.attach"
}

func (c *commandS3PolicyAttach) Help() string {
	return `attach a policy to an S3 IAM user

	s3.policy.attach -policy <policy_name> -user <username>

	The policy must already exist (create it with s3.policy -put).
`
}

func (c *commandS3PolicyAttach) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3PolicyAttach) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	f := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	policy := f.String("policy", "", "policy name")
	user := f.String("user", "", "user name")
	if err := f.Parse(args); err != nil {
		return err
	}

	if *policy == "" {
		return fmt.Errorf("-policy is required")
	}
	if *user == "" {
		return fmt.Errorf("-user is required")
	}

	return commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		// Verify the policy exists
		if _, _, err := client.GetPolicy(context.Background(), iamwire.NewGetPolicyRequest(iamwire.GetPolicyRequestInput{Name: *policy})); err != nil {
			return fmt.Errorf("get policy %q: %w", *policy, err)
		}

		// Get the user
		_, body, err := client.GetUser(context.Background(), iamwire.NewGetUserRequest(iamwire.GetUserRequestInput{Username: *user}))
		if err != nil {
			return fmt.Errorf("get user %q: %w", *user, err)
		}
		identity, err := iamwire.GetUserResp(body)
		if err != nil {
			return err
		}
		if identity == nil {
			return fmt.Errorf("user %q returned empty identity", *user)
		}

		// Check if already attached
		for _, p := range identity.PolicyNames {
			if p == *policy {
				return json.NewEncoder(writer).Encode(map[string]string{"policy": *policy, "user": *user})
			}
		}

		identity.PolicyNames = append(identity.PolicyNames, *policy)
		_, _, err = client.UpdateUser(context.Background(), iamwire.NewUpdateUserRequest(iamwire.UpdateUserRequestInput{
			Username: *user,
			Identity: iamwire.IdentityInputFromPB(identity),
		}))
		if err != nil {
			return err
		}

		return json.NewEncoder(writer).Encode(map[string]string{"policy": *policy, "user": *user})
	})
}
