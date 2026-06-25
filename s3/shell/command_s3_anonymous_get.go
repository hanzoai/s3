package shell

import (
	"context"
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"

	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3AnonymousGet{})
}

type commandS3AnonymousGet struct {
}

func (c *commandS3AnonymousGet) Name() string {
	return "s3.anonymous.get"
}

func (c *commandS3AnonymousGet) Help() string {
	return `show anonymous access for a bucket

	s3.anonymous.get -bucket <bucket_name>
`
}

func (c *commandS3AnonymousGet) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3AnonymousGet) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	f := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	bucket := f.String("bucket", "", "bucket name")
	if err := f.Parse(args); err != nil {
		return err
	}

	if *bucket == "" {
		return fmt.Errorf("-bucket is required")
	}

	return commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetUser(context.Background(), iamwire.NewGetUserRequest(iamwire.GetUserRequestInput{Username: anonymousUserName}))
		if err != nil {
			// Over the ZAP transport a missing anonymous identity surfaces as a
			// generic call error; treat it as "no anonymous access" as the gRPC
			// NotFound path did.
			fmt.Fprintf(writer, "Bucket: %s\nAccess: none\n", *bucket)
			return nil
		}
		identity, err := iamwire.GetUserResp(body)
		if err != nil {
			return err
		}
		if identity == nil {
			fmt.Fprintf(writer, "Bucket: %s\nAccess: none\n", *bucket)
			return nil
		}

		var actions []string
		for _, a := range identity.Actions {
			parts := strings.SplitN(a, ":", 2)
			if len(parts) == 2 && parts[1] == *bucket {
				actions = append(actions, parts[0])
			}
		}

		fmt.Fprintf(writer, "Bucket: %s\n", *bucket)
		if len(actions) == 0 {
			fmt.Fprintln(writer, "Access: none")
		} else {
			sort.Strings(actions)
			fmt.Fprintf(writer, "Access: %s\n", strings.Join(actions, ", "))
		}

		return nil
	})
}
