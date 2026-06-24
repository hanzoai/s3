package shell

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3AnonymousList{})
}

type commandS3AnonymousList struct {
}

func (c *commandS3AnonymousList) Name() string {
	return "s3.anonymous.list"
}

func (c *commandS3AnonymousList) Help() string {
	return `list all buckets with anonymous access

	s3.anonymous.list
`
}

func (c *commandS3AnonymousList) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3AnonymousList) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	return commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetUser(iamwire.NewGetUserRequest(iamwire.GetUserRequestInput{Username: anonymousUserName}))
		if err != nil {
			// A missing anonymous identity surfaces as a generic call error over
			// the ZAP transport; treat it as no anonymous access (gRPC NotFound).
			fmt.Fprintln(writer, "No anonymous access configured.")
			return nil
		}
		identity, err := iamwire.GetUserResp(body)
		if err != nil {
			return err
		}
		if identity == nil {
			fmt.Fprintln(writer, "No anonymous access configured.")
			return nil
		}

		// Group actions by bucket
		bucketActions := map[string][]string{}
		for _, a := range identity.Actions {
			parts := strings.SplitN(a, ":", 2)
			if len(parts) == 2 {
				bucketActions[parts[1]] = append(bucketActions[parts[1]], parts[0])
			}
		}

		if len(bucketActions) == 0 {
			fmt.Fprintln(writer, "No anonymous access configured.")
			return nil
		}

		// Sort bucket names
		buckets := make([]string, 0, len(bucketActions))
		for b := range bucketActions {
			buckets = append(buckets, b)
		}
		sort.Strings(buckets)

		tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
		fmt.Fprintln(tw, "BUCKET\tACCESS")
		for _, b := range buckets {
			actions := bucketActions[b]
			sort.Strings(actions)
			fmt.Fprintf(tw, "%s\t%s\n", b, strings.Join(actions, ", "))
		}
		return tw.Flush()
	})
}
