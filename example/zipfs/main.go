package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/zipfs"
	"fmt"
	"flag"
	"log"
	"os"
)

var _ = log.Printf

func main() {
	// Scans the arg list and sets up flags
	debug := flag.Bool("debug", false, "print debugging messages.")
	latencies := flag.Bool("latencies", false, "record operation latencies.")
	
	flag.Parse()
	if flag.NArg() < 2 {
		// TODO - where to get program name?
		fmt.Println("usage: main MOUNTPOINT ZIP-FILE")
		os.Exit(2)
	}

	var fs fuse.FileSystem
	fs = zipfs.NewZipArchiveFileSystem(flag.Arg(1))
	debugFs := fuse.NewFileSystemDebug()

	if *latencies {
		debugFs.Original = fs
		fs = debugFs
	}
	
	conn := fuse.NewFileSystemConnector(fs)
	state := fuse.NewMountState(conn)

	if *latencies {
		debugFs.AddFileSystemConnector(conn)
		debugFs.AddMountState(state)
	}
	
	mountPoint := flag.Arg(0)
	state.RecordStatistics = *latencies
	state.Debug = *debug
	state.Mount(mountPoint)

	fmt.Printf("Mounted %s - PID %s\n", mountPoint, fuse.MyPID())
	state.Loop(true)
}
