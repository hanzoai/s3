package shell

import (
	"flag"
	"fmt"
	"io"
	"text/tabwriter"

	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3AccessKeyList{})
}

type commandS3AccessKeyList struct {
}

func (c *commandS3AccessKeyList) Name() string {
	return "s3.accesskey.list"
}

func (c *commandS3AccessKeyList) Help() string {
	return `list access keys for an S3 IAM user

	s3.accesskey.list -user <username>
`
}

func (c *commandS3AccessKeyList) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3AccessKeyList) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	f := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	user := f.String("user", "", "user name")
	if err := f.Parse(args); err != nil {
		return err
	}

	if *user == "" {
		return fmt.Errorf("-user is required")
	}

	return commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetUser(iamwire.NewGetUserRequest(iamwire.GetUserRequestInput{Username: *user}))
		if err != nil {
			return err
		}
		identity, err := iamwire.GetUserResp(body)
		if err != nil {
			return err
		}

		if identity == nil || len(identity.Credentials) == 0 {
			fmt.Fprintf(writer, "No access keys for user %q.\n", *user)
			return nil
		}

		tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
		fmt.Fprintln(tw, "ACCESS KEY\tSTATUS")
		for _, cred := range identity.Credentials {
			st := cred.Status
			if st == "" {
				st = "Active"
			}
			fmt.Fprintf(tw, "%s\t%s\n", cred.AccessKey, st)
		}
		return tw.Flush()
	})
}
