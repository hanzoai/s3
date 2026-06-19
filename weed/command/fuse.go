package command

func init() {
	cmdFuse.Run = runFuse // break init cycle
}

var cmdFuse = &Command{
	UsageLine: "fuse /mnt/mount/point -o \"filer=localhost:8888,filer.path=/\"",
	Short:     "Allow use s3 with linux's mount command",
	Long: `Allow use s3 with linux's mount command

  You can use -t s3 on mount command:
  mv s3 /sbin/mount.weed
  mount -t s3 fuse /mnt -o "filer=localhost:8888,filer.path=/"

  Or you can use -t fuse on mount command:
  mv s3 /sbin/weed
  mount -t fuse.s3 fuse /mnt -o "filer=localhost:8888,filer.path=/"
  mount -t fuse "weed#fuse" /mnt -o "filer=localhost:8888,filer.path=/"

  To use without mess with your /sbin:
  mount -t fuse./home/user/bin/s3 fuse /mnt -o "filer=localhost:8888,filer.path=/"
  mount -t fuse "/home/user/bin/weed#fuse" /mnt -o "filer=localhost:8888,filer.path=/"

  To pass more than one parameter use quotes, example:
  mount -t s3 fuse /mnt -o "filer='192.168.0.1:8888,192.168.0.2:8888',filer.path=/"

  To check valid options look "s3 mount --help"
  `,
}

func GetFuseCommandName() string {
	return cmdFuse.Name()
}
