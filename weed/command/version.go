package command

import (
	"fmt"
	"runtime"

	"github.com/seaweedfs/seaweedfs/weed/util/version"
)

var cmdVersion = &Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "print Hanzo S3 version",
	Long:      `Version prints the Hanzo S3 version`,
}

func runVersion(cmd *Command, args []string) bool {
	if len(args) != 0 {
		cmd.Usage()
	}

	fmt.Printf("Hanzo S3 (fork of SeaweedFS) version %s %s %s\n", version.Version(), runtime.GOOS, runtime.GOARCH)
	return true
}
