package shell

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/hanzoai/s3/s3/pb/iam_pb"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

const anonymousUserName = "anonymous"

func init() {
	Commands = append(Commands, &commandS3AnonymousSet{})
}

type commandS3AnonymousSet struct {
}

func (c *commandS3AnonymousSet) Name() string {
	return "s3.anonymous.set"
}

func (c *commandS3AnonymousSet) Help() string {
	return `set anonymous (public) access on a bucket

	s3.anonymous.set -bucket <bucket_name> -access Read,List
	s3.anonymous.set -bucket <bucket_name> -access none

	Supported actions: Read, Write, List, Tagging, Admin
	Use "none" to remove all anonymous access for the bucket.

	This manages the special "anonymous" user's actions. It does not
	use IAM policies — it sets legacy per-bucket actions directly.
`
}

func (c *commandS3AnonymousSet) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3AnonymousSet) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	f := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	bucket := f.String("bucket", "", "bucket name")
	access := f.String("access", "", "comma-separated actions: Read,Write,List,Tagging,Admin or none")
	if err := f.Parse(args); err != nil {
		return err
	}

	if *bucket == "" {
		return fmt.Errorf("-bucket is required")
	}
	if *access == "" {
		return fmt.Errorf("-access is required")
	}

	return commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		// Get or create anonymous user
		identity, isNew, err := getOrCreateAnonymousUser(client)
		if err != nil {
			return err
		}

		// Remove existing actions for this bucket
		var kept []string
		for _, a := range identity.Actions {
			parts := strings.SplitN(a, ":", 2)
			if len(parts) != 2 || parts[1] != *bucket {
				kept = append(kept, a)
			}
		}

		// Add new actions unless "none"
		canonicalActions := map[string]string{
			"read": "Read", "write": "Write", "list": "List",
			"tagging": "Tagging", "admin": "Admin",
		}
		if strings.ToLower(strings.TrimSpace(*access)) != "none" {
			seen := make(map[string]struct{})
			for _, action := range strings.Split(*access, ",") {
				action = strings.TrimSpace(action)
				if action != "" {
					canonical, ok := canonicalActions[strings.ToLower(action)]
					if !ok {
						return fmt.Errorf("invalid action %q: supported actions are Read, Write, List, Tagging, Admin", action)
					}
					if _, dup := seen[canonical]; dup {
						continue
					}
					seen[canonical] = struct{}{}
					kept = append(kept, canonical+":"+*bucket)
				}
			}
		}

		identity.Actions = kept

		if isNew {
			_, _, err = client.CreateUser(context.Background(), iamwire.NewCreateUserRequest(iamwire.CreateUserRequestInput{Identity: iamwire.IdentityInputFromPB(identity)}))
		} else {
			_, _, err = client.UpdateUser(context.Background(), iamwire.NewUpdateUserRequest(iamwire.UpdateUserRequestInput{
				Username: anonymousUserName,
				Identity: iamwire.IdentityInputFromPB(identity),
			}))
		}
		if err != nil {
			return err
		}

		fmt.Fprintf(writer, "Set anonymous access on bucket %q to: %s\n", *bucket, *access)
		return nil
	})
}

func getOrCreateAnonymousUser(client *iamwire.HanzoIdentityAccessManagementClient) (*iam_pb.Identity, bool, error) {
	_, body, err := client.GetUser(context.Background(), iamwire.NewGetUserRequest(iamwire.GetUserRequestInput{Username: anonymousUserName}))
	if err == nil {
		identity, derr := iamwire.GetUserResp(body)
		if derr != nil {
			return nil, false, derr
		}
		if identity == nil {
			// Existing call succeeded but carried no identity: create a fresh one.
			return &iam_pb.Identity{Name: anonymousUserName, Actions: []string{}}, true, nil
		}
		return identity, false, nil
	}

	// Over the ZAP transport a missing anonymous user surfaces as a generic call
	// error (the gRPC NotFound path); create a fresh identity.
	return &iam_pb.Identity{
		Name:    anonymousUserName,
		Actions: []string{},
	}, true, nil
}
