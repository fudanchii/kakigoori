package main

import (
	"log"
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

	config, err := parseConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	event.StartListening()

	kakigoorifs := fs.NewKakigooriFileSystem(config.Root)
	finalFs = kakigoorifs

	opts := &nodefs.Options{
		NegativeTimeout: time.Second,
		AttrTimeout:     time.Second,
		EntryTimeout:    time.Second,
	}

	pathFs := pathfs.NewPathNodeFs(finalFs, nil)
	conn := nodefs.NewFileSystemConnector(pathFs.Root(), opts)

	mOpts := &fuse.MountOptions{
		AllowOther: false,
	}
	state, err := fuse.NewServer(conn.RawFS(), config.MountPoint, mOpts)
	if err != nil {
		log.Printf("Mount fail: %v\n", err)
		os.Exit(1)
	}

	state.SetDebug(false)

	log.Println("Mounted!")
	state.Serve()

}
