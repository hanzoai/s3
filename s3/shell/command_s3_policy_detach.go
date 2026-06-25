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
	Commands = append(Commands, &commandS3PolicyDetach{})
}

type commandS3PolicyDetach struct {
}

func (c *commandS3PolicyDetach) Name() string {
	return "s3.policy.detach"
}

func (c *commandS3PolicyDetach) Help() string {
	return `detach a policy from an S3 IAM user

	s3.policy.detach -policy <policy_name> -user <username>
`
}

func (c *commandS3PolicyDetach) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3PolicyDetach) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
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

		found := false
		var kept []string
		for _, p := range identity.PolicyNames {
			if p == *policy {
				found = true
			} else {
				kept = append(kept, p)
			}
		}
		if !found {
			return fmt.Errorf("policy %q is not attached to user %q", *policy, *user)
		}

		identity.PolicyNames = kept
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
