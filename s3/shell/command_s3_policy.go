package shell

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/hanzoai/s3/s3/s3api/policy_engine"
	"github.com/hanzoai/s3/s3/util"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3Policy{})
}

type commandS3Policy struct {
}

func (c *commandS3Policy) Name() string {
	return "s3.policy"
}

func (c *commandS3Policy) Help() string {
	return `manage s3 policies

	# create or update a policy
	s3.policy -put -name=mypolicy -file=policy.json

	# list all policies
	s3.policy -list

	# get a policy
	s3.policy -get -name=mypolicy

	# delete a policy
	s3.policy -delete -name=mypolicy
	`
}

func (c *commandS3Policy) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3Policy) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {

	s3PolicyCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	put := s3PolicyCommand.Bool("put", false, "create or update a policy")
	get := s3PolicyCommand.Bool("get", false, "get a policy")
	list := s3PolicyCommand.Bool("list", false, "list all policies")
	del := s3PolicyCommand.Bool("delete", false, "delete a policy")
	name := s3PolicyCommand.String("name", "", "policy name")
	file := s3PolicyCommand.String("file", "", "policy file (json)")

	if err = s3PolicyCommand.Parse(args); err != nil {
		return err
	}

	actionCount := 0
	for _, v := range []bool{*put, *get, *list, *del} {
		if v {
			actionCount++
		}
	}
	if actionCount == 0 {
		return fmt.Errorf("one of -put, -get, -list, -delete must be specified")
	}
	if actionCount > 1 {
		return fmt.Errorf("only one of -put, -get, -list, -delete can be specified")
	}

	return commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		if *put {
			if *name == "" {
				return fmt.Errorf("-name is required")
			}
			if *file == "" {
				return fmt.Errorf("-file is required")
			}
			data, err := os.ReadFile(util.ResolvePath(*file))
			if err != nil {
				return fmt.Errorf("failed to read policy file: %v", err)
			}

			// Validate JSON
			var policy policy_engine.PolicyDocument
			if err := json.Unmarshal(data, &policy); err != nil {
				return fmt.Errorf("invalid policy json: %v", err)
			}

			_, _, err = client.PutPolicy(iamwire.NewPutPolicyRequest(iamwire.PutPolicyRequestInput{
				Name:    *name,
				Content: string(data),
			}))
			return err
		}

		if *get {
			if *name == "" {
				return fmt.Errorf("-name is required")
			}
			_, body, err := client.GetPolicy(iamwire.NewGetPolicyRequest(iamwire.GetPolicyRequestInput{
				Name: *name,
			}))
			if err != nil {
				return err
			}
			_, content, err := iamwire.GetPolicyResp(body)
			if err != nil {
				return err
			}
			if content == "" {
				return fmt.Errorf("policy not found")
			}
			fmt.Fprintf(writer, "%s\n", content)
			return nil
		}

		if *list {
			_, body, err := client.ListPolicies(iamwire.NewListPoliciesRequest(iamwire.ListPoliciesRequestInput{}))
			if err != nil {
				return err
			}
			policies, err := iamwire.ListPoliciesResp(body)
			if err != nil {
				return err
			}
			for _, policy := range policies {
				fmt.Fprintf(writer, "Name: %s\n", policy.Name)
				fmt.Fprintf(writer, "Content: %s\n", policy.Content)
				fmt.Fprintf(writer, "---\n")
			}
			return nil
		}

		if *del {
			if *name == "" {
				return fmt.Errorf("-name is required")
			}
			_, _, err := client.DeletePolicy(iamwire.NewDeletePolicyRequest(iamwire.DeletePolicyRequestInput{
				Name: *name,
			}))
			return err
		}

		return nil
	})

}
