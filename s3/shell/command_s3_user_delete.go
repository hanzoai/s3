package shell

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3UserDelete{})
}

type commandS3UserDelete struct {
}

func (c *commandS3UserDelete) Name() string {
	return "s3.user.delete"
}

func (c *commandS3UserDelete) Help() string {
	return `delete an S3 IAM user

	s3.user.delete -name <username>
`
}

func (c *commandS3UserDelete) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3UserDelete) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	f := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	name := f.String("name", "", "user name")
	if err := f.Parse(args); err != nil {
		return err
	}

	if *name == "" {
		return fmt.Errorf("-name is required")
	}

	err := commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, _, err := client.DeleteUser(iamwire.NewDeleteUserRequest(iamwire.DeleteUserRequestInput{Username: *name}))
		return err
	})
	if err != nil {
		return err
	}

	return json.NewEncoder(writer).Encode(map[string]string{"name": *name})
}
