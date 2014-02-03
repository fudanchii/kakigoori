package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	"github.com/fudanchii/kakigoori/event"
	"github.com/fudanchii/kakigoori/event/handler"
	"github.com/fudanchii/kakigoori/fs"
)

var (
	APPNAME    = "kakigoori"
	APPVERSION = "0.1-norev"
	configPath = flag.String("c", "config.json", "Config file to use.")
	flVersion  = flag.Bool("v", false, "Display application version.")
)

func main() {

	flag.Parse()
	check_intro()

	var finalFs pathfs.FileSystem

	config, err := parseConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	hnd := event.StartListening(config.Handlers)
	hnd.RegisterHandler(event.Write | event.Close, handler.Spawner)

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

func check_intro() {
	if *flVersion {
		show_version()
		os.Exit(1)
	}
}

func show_version() {
	fmt.Printf("%s %s\n", APPNAME, APPVERSION)
}
