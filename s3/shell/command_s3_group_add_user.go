package shell

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3GroupAddUser{})
}

type commandS3GroupAddUser struct {
}

func (c *commandS3GroupAddUser) Name() string {
	return "s3.group.add-user"
}

func (c *commandS3GroupAddUser) Help() string {
	return `add a user to an S3 IAM group

	s3.group.add-user -group <groupname> -user <username>
`
}

func (c *commandS3GroupAddUser) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3GroupAddUser) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
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

		// Verify user exists
		userFound := false
		for _, id := range cfg.Identities {
			if id.Name == *user {
				userFound = true
				break
			}
		}
		if !userFound {
			return fmt.Errorf("user %s not found", *user)
		}

		for _, g := range cfg.Groups {
			if g.Name == *group {
				// Check if already a member
				for _, m := range g.Members {
					if m == *user {
						return fmt.Errorf("user %s is already a member of group %s", *user, *group)
					}
				}
				g.Members = append(g.Members, *user)
				if _, _, err := client.PutConfiguration(iamwire.NewPutConfigurationRequest(iamwire.PutConfigurationRequestInput{Configuration: iamwire.ConfigurationInputFromPB(cfg)})); err != nil {
					return err
				}
				return json.NewEncoder(writer).Encode(map[string]string{"group": *group, "user": *user})
			}
		}
		return fmt.Errorf("group %s not found", *group)
	})
}
