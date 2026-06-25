package shell

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/hanzoai/s3/s3/pb/iam_pb"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	Commands = append(Commands, &commandS3GroupCreate{})
}

type commandS3GroupCreate struct {
}

func (c *commandS3GroupCreate) Name() string {
	return "s3.group.create"
}

func (c *commandS3GroupCreate) Help() string {
	return `create an S3 IAM group

	s3.group.create -name <groupname>

	Creates a new empty group. Add users with s3.group.add-user and
	attach policies with s3.policy.attach or the IAM API.
`
}

func (c *commandS3GroupCreate) HasTag(CommandTag) bool {
	return false
}

func (c *commandS3GroupCreate) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	f := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	name := f.String("name", "", "group name")
	if err := f.Parse(args); err != nil {
		return err
	}
	if *name == "" {
		return fmt.Errorf("-name is required")
	}

	return commandEnv.withIamClient(func(client *iamwire.HanzoIdentityAccessManagementClient) error {
		_, body, err := client.GetConfiguration(context.Background(), iamwire.NewGetConfigurationRequest(iamwire.GetConfigurationRequestInput{}))
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
			if g.Name == *name {
				return fmt.Errorf("group %s already exists", *name)
			}
		}

		cfg.Groups = append(cfg.Groups, &iam_pb.Group{Name: *name})

		if _, _, err := client.PutConfiguration(context.Background(), iamwire.NewPutConfigurationRequest(iamwire.PutConfigurationRequestInput{Configuration: iamwire.ConfigurationInputFromPB(cfg)})); err != nil {
			return err
		}

		return json.NewEncoder(writer).Encode(map[string]string{"group": *name})
	})
}
