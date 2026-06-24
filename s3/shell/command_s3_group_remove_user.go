package shell

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3GroupRemoveUser{})
}

type commandS3GroupRemoveUser struct {
}

func (c *commandS3GroupRemoveUser) Name() string {
	return "s3.group.remove-user"
}

func (c *commandS3GroupRemoveUser) Help() string {
	return `remove a user from an S3 IAM group

	s3.group.remove-user -group <groupname> -user <username>
`
}

func (c *commandS3GroupRemoveUser) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3GroupRemoveUser) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	f := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	group := f.String("group", "", "group name")
	user := f.String("user", "", "user name")
	if err := f.Parse(args); err != nil {
		return err
	}
	if *group == "" {
		return fmt.Errorf("-group is required")
	}
	if *user == "" {
		return fmt.Errorf("-user is required")
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

		for _, g := range cfg.Groups {
			if g.Name == *group {
				for i, m := range g.Members {
					if m == *user {
						g.Members = append(g.Members[:i], g.Members[i+1:]...)
						if _, _, err := client.PutConfiguration(iamwire.NewPutConfigurationRequest(iamwire.PutConfigurationRequestInput{Configuration: iamwire.ConfigurationInputFromPB(cfg)})); err != nil {
							return err
						}
						return json.NewEncoder(writer).Encode(map[string]string{"group": *group, "removed": *user})
					}
				}
				return fmt.Errorf("user %s is not a member of group %s", *user, *group)
			}
		}
		return fmt.Errorf("group %s not found", *group)
	})
}
