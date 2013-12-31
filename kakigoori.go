package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	"github.com/fudanchii/kakigoori/event"
	"github.com/fudanchii/kakigoori/fs"
)

func main() {

	var finalFs pathfs.FileSystem

	other := flag.Bool("allow-other", false, "mount with -o allowother.")

	flag.Parse()
	mountPoint := flag.Arg(0)
	orig := flag.Arg(1)

	event.StartListening()

	kakigoorifs := fs.NewKakigooriFileSystem(orig)
	finalFs = kakigoorifs

	opts := &nodefs.Options{
		NegativeTimeout: time.Second,
		AttrTimeout:     time.Second,
		EntryTimeout:    time.Second,
	}

	pathFs := pathfs.NewPathNodeFs(finalFs, nil)
	conn := nodefs.NewFileSystemConnector(pathFs, opts)

	mOpts := &fuse.MountOptions{
		AllowOther: *other,
	}
	state, err := fuse.NewServer(conn.RawFS(), mountPoint, mOpts)
	if err != nil {
		fmt.Printf("Mount fail: %v\n", err)
		os.Exit(1)
	}

	state.SetDebug(false)

	fmt.Println("Mounted!")
	state.Serve()

}
