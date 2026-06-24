package shell

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3UserEnable{})
}

type commandS3UserEnable struct {
}

func (c *commandS3UserEnable) Name() string {
	return "s3.user.enable"
}

func (c *commandS3UserEnable) Help() string {
	return `enable a disabled S3 IAM user

	s3.user.enable -name <username>
`
}

func (c *commandS3UserEnable) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3UserEnable) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	f := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	name := f.String("name", "", "user name")
	if err := f.Parse(args); err != nil {
		return err
	}

	if *name == "" {
		return fmt.Errorf("-name is required")
	}

	err := commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetUser(iamwire.NewGetUserRequest(iamwire.GetUserRequestInput{Username: *name}))
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
			identity.Disabled = false
			_, _, err = client.UpdateUser(iamwire.NewUpdateUserRequest(iamwire.UpdateUserRequestInput{
				Username: *name,
				Identity: iamwire.IdentityInputFromPB(identity),
			}))
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return json.NewEncoder(writer).Encode(map[string]string{"name": *name, "status": "enabled"})
}
