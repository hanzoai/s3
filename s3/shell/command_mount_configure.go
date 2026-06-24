package shell

import (
	"flag"
	"fmt"
	"io"

	"github.com/hanzoai/s3/s3/util"
	mountwire "github.com/hanzoai/s3/s3/wire/mount"
	"github.com/zap-proto/go/transport"
)

func init() {
	Commands = append(Commands, &commandMountConfigure{})
}

type commandMountConfigure struct {
}

func (c *commandMountConfigure) Name() string {
	return "mount.configure"
}

func (c *commandMountConfigure) Help() string {
	return `configure the mount on current server

	mount.configure -dir=<mount_directory>

	This command connects with local mount via unix socket, so it can only run locally.
	The "mount_directory" value needs to be exactly the same as how mount was started in "s3 mount -dir=<mount_directory>"

`
}

func (c *commandMountConfigure) HasTag(CommandTag) bool {
	return false
}

func (c *commandMountConfigure) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {

	mountConfigureCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	mountDir := mountConfigureCommand.String("dir", "", "the mount directory same as how \"s3 mount -dir=<mount_directory>\" was started")
	mountQuota := mountConfigureCommand.Int("quotaMB", 0, "the quota in MB")
	if err = mountConfigureCommand.Parse(args); err != nil {
		return nil
	}

	mountDirHash := util.HashToInt32([]byte(*mountDir))
	if mountDirHash < 0 {
		mountDirHash = -mountDirHash
	}
	localSocket := fmt.Sprintf("/tmp/hanzo-mount-%d.sock", mountDirHash)

	conn, err := transport.Dial("unix", localSocket)
	if err != nil {
		return
	}
	defer conn.Close()

	client := mountwire.NewHanzoMountClient(conn, nil)
	_, _, err = client.Configure(mountwire.NewConfigureRequest(mountwire.ConfigureRequestInput{
		CollectionCapacity: int64(*mountQuota) * 1024 * 1024,
	}))

	return
}
