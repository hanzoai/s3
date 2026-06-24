package shell

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3GroupDelete{})
}

type commandS3GroupDelete struct {
}

func (c *commandS3GroupDelete) Name() string {
	return "s3.group.delete"
}

func (c *commandS3GroupDelete) Help() string {
	return `delete an S3 IAM group

	s3.group.delete -name <groupname>

	The group must have no members and no attached policies.
`
}

func (c *commandS3GroupDelete) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3GroupDelete) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	f := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	name := f.String("name", "", "group name")
	if err := f.Parse(args); err != nil {
		return err
	}
	if *name == "" {
		return fmt.Errorf("-name is required")
	}

	return commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetConfiguration(iamwire.NewGetConfigurationRequest(iamwire.GetConfigurationRequestInput{}))
		if err != nil {
			return err
		}
		cfg, err := iamwire.GetConfigurationResp(body)
		if err != nil {
			return err
		}
		if cfg == nil {
			return fmt.Errorf("no IAM configuration found")
		}

		for i, g := range cfg.Groups {
			if g.Name == *name {
				if len(g.Members) > 0 {
					return fmt.Errorf("cannot delete group %s: has %d member(s)", *name, len(g.Members))
				}
				if len(g.PolicyNames) > 0 {
					return fmt.Errorf("cannot delete group %s: has %d attached policy(ies)", *name, len(g.PolicyNames))
				}
				cfg.Groups = append(cfg.Groups[:i], cfg.Groups[i+1:]...)
				if _, _, err := client.PutConfiguration(iamwire.NewPutConfigurationRequest(iamwire.PutConfigurationRequestInput{Configuration: iamwire.ConfigurationInputFromPB(cfg)})); err != nil {
					return err
				}
				return json.NewEncoder(writer).Encode(map[string]string{"deleted": *name})
			}
		}
		return fmt.Errorf("group %s not found", *name)
	})
}
