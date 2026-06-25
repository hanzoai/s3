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
	Commands = append(Commands, &commandS3UserDisable{})
}

type commandS3UserDisable struct {
}

func (c *commandS3UserDisable) Name() string {
	return "s3.user.disable"
}

func (c *commandS3UserDisable) Help() string {
	return `disable an S3 IAM user

	s3.user.disable -name <username>

	Disabled users cannot authenticate. Their credentials and policies
	are preserved and will take effect again when the user is re-enabled.
`
}

func (c *commandS3UserDisable) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3UserDisable) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	f := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	name := f.String("name", "", "user name")
	if err := f.Parse(args); err != nil {
		return err
	}

	if *name == "" {
		return fmt.Errorf("-name is required")
	}

	err := commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetUser(context.Background(), iamwire.NewGetUserRequest(iamwire.GetUserRequestInput{Username: *name}))
		if err != nil {
			return fmt.Errorf("get user %q: %w", *name, err)
		}
		identity, err := iamwire.GetUserResp(body)
		if err != nil {
			return err
		}
		if identity == nil {
			return fmt.Errorf("user %q returned empty identity", *name)
		}

		if identity.Disabled {
			return nil
		}

		identity.Disabled = true
		_, _, err = client.UpdateUser(context.Background(), iamwire.NewUpdateUserRequest(iamwire.UpdateUserRequestInput{
			Username: *name,
			Identity: iamwire.IdentityInputFromPB(identity),
		}))
		return err
	})
	if err != nil {
		return err
	}

	return json.NewEncoder(writer).Encode(map[string]string{"name": *name, "status": "disabled"})
}
