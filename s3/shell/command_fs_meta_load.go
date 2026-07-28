package shell

import (
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/util"
)

func init() {
	Commands = append(Commands, &commandFsMetaLoad{})
}

type commandFsMetaLoad struct {
	dirPrefix *string
}

func (c *commandFsMetaLoad) Name() string {
	return "fs.meta.load"
}

func (c *commandFsMetaLoad) Help() string {
	return `load saved filer meta data to restore the directory and file structure

	fs.meta.load <filer_host>-<port>-<time>.meta
	fs.meta.load -v=false <filer_host>-<port>-<time>.meta // skip printing out the verbose output
 	fs.meta.load -concurrency=1 <filer_host>-<port>-<time>.meta // number of parallel meta load to filer
	fs.meta.load -dirPrefix=/buckets/important <filer_host>.meta // load any dirs with prefix "important"

`
}

func (c *commandFsMetaLoad) HasTag(CommandTag) bool {
	return false
}

func (c *commandFsMetaLoad) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {

	if len(args) == 0 {
		fmt.Fprintf(writer, "missing a metadata file\n")
		return nil
	}

	fileName := util.ResolvePath(args[len(args)-1])

	metaLoadCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	c.dirPrefix = metaLoadCommand.String("dirPrefix", "", "load entries only with directories matching prefix")
	concurrency := metaLoadCommand.Int("concurrency", 1, "number of parallel meta load to filer")
	verbose := metaLoadCommand.Bool("v", true, "verbose mode")
	if err = metaLoadCommand.Parse(args[0 : len(args)-1]); err != nil {
		return nil
	}

	var dst io.Reader

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", fileName, err)
	}
	defer f.Close()

	dst = f

	if strings.HasSuffix(fileName, ".gz") || strings.HasSuffix(fileName, ".gzip") {
		var gr *gzip.Reader
		gr, err = gzip.NewReader(dst)
		if err != nil {
			return err
		}
		defer func() {
			err1 := gr.Close()
			if err == nil {
				err = err1
			}
		}()

		dst = gr
	}

	var dirCount, fileCount uint64
	lastLogTime := time.Now()

	err = commandEnv.WithFilerClient(false, func(client filer_pb.HanzoFilerClient) (loadErr error) {

		sizeBuf := make([]byte, 4)
		waitChan := make(chan struct{}, *concurrency)
		var wg sync.WaitGroup

		// File entries are written on their own goroutines, so EVERY exit from
		// this function — including the EOF that ends a *successful* load — has
		// to wait for them first. The caller closes the filer connection as soon
		// as this returns, so returning early abandoned the in-flight writes:
		// their CreateEntry calls died against a torn-down connection, the
		// entries were silently lost, and the command still printed "is loaded".
		//
		// At the default concurrency of 1 that dropped the tail of every load.
		// It is why a restored tree came back with its directories present and
		// its files missing — directories are created synchronously just below,
		// files were not.
		var mu sync.Mutex
		var asyncErr error
		firstAsyncErr := func() error {
			mu.Lock()
			defer mu.Unlock()
			return asyncErr
		}
		defer func() {
			wg.Wait()
			close(waitChan)
			if loadErr == nil {
				loadErr = firstAsyncErr()
			}
		}()

		for {
			if _, err := io.ReadFull(dst, sizeBuf); err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}

			size := util.BytesToUint32(sizeBuf)

			data := make([]byte, int(size))

			if _, err := io.ReadFull(dst, data); err != nil {
				return err
			}

			fullEntry := &filer_pb.FullEntry{}
			if err = proto.Unmarshal(data, fullEntry); err != nil {
				return err
			}

			// check collection name pattern
			entryFullName := string(util.FullPath(fullEntry.Dir).Child(fullEntry.Entry.Name))
			if *c.dirPrefix != "" {
				if !strings.HasPrefix(fullEntry.Dir, *c.dirPrefix) {
					if *verbose {
						fmt.Fprintf(writer, "not match dir prefix %s\n", entryFullName)
					}
					continue
				}
			}

			if *verbose || lastLogTime.Add(time.Second).Before(time.Now()) {
				if !*verbose {
					lastLogTime = time.Now()
				}
				fmt.Fprintf(writer, "load %s\n", entryFullName)
			}

			fullEntry.Entry.Name = strings.ReplaceAll(fullEntry.Entry.Name, "/", "x")
			if fullEntry.Entry.IsDirectory {
				wg.Wait()
				if errEntry := filer_pb.CreateEntry(context.Background(), client, &filer_pb.CreateEntryRequest{
					Directory: fullEntry.Dir,
					Entry:     fullEntry.Entry,
				}); errEntry != nil {
					return errEntry
				}
				dirCount++
			} else {
				wg.Add(1)
				waitChan <- struct{}{}
				go func(entry *filer_pb.FullEntry) {
					defer wg.Done()
					defer func() { <-waitChan }()
					if errEntry := filer_pb.CreateEntry(context.Background(), client, &filer_pb.CreateEntryRequest{
						Directory: entry.Dir,
						Entry:     entry.Entry,
					}); errEntry != nil {
						// Keep the FIRST failure and keep it off the outer err:
						// several of these run at once, and assigning to the
						// enclosing error from each was a data race that could
						// also erase a real failure with a later success.
						mu.Lock()
						if asyncErr == nil {
							asyncErr = errEntry
						}
						mu.Unlock()
					}
				}(fullEntry)
				if e := firstAsyncErr(); e != nil {
					return e
				}
				fileCount++
			}
		}
	})

	if err == nil {
		fmt.Fprintf(writer, "\ntotal %d directories, %d files", dirCount, fileCount)
		fmt.Fprintf(writer, "\n%s is loaded.\n", fileName)
	}

	return err
}
