package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Aterocana/httpsrv"
)

const dateLayout = "2006/01/02-15:04:05"

var (
	version   = "unknown"
	buildDate = dateLayout
	gitCommit = "unknown"
)

// flags parses cmd line flags and it returns the path and the port to use.
func flags() []httpsrv.Options {
	var path string
	var port int
	var askHelp bool
	var askVersion bool

	flag.StringVar(&path, "path", "./", "the folder you want to serve. Default is current (./)")
	flag.IntVar(&port, "port", 0, "the port the server should listen on. Default is a random one")
	flag.BoolVar(&askHelp, "help", false, "show help")
	flag.BoolVar(&askVersion, "version", false, "show version")

	flag.Parse()
	if askHelp {
		flag.Usage()
		os.Exit(0)
	}
	if askVersion {
		buildDateTime, err := time.Parse(dateLayout, buildDate)
		if err != nil {
			buildDateTime = time.Time{}
		}
		fmt.Printf("%s (#%s) build on %s\n", version, gitCommit, buildDateTime)
		os.Exit(0)
	}
	return []httpsrv.Options{
		httpsrv.WithPort(port),
		httpsrv.WithPath(path),
		httpsrv.WithLogger(os.Stdout),
	}
}
