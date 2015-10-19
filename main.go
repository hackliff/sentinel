package main

import (
	"fmt"
	"os"

	"github.com/azer/logger"
	"github.com/voxelbrain/goptions"
)

const (
	VersionNumber = "0.1.0"
)

var log = logger.New("sentinel")

// Options declares command line flags
type Options struct {
	Version       bool `goptions:"-v, --version, description='Print version'"`
	goptions.Help `goptions:"-h, --help, description='Show this help'"`
	// NOTE use built-in support for *os.File ?
	Conf string `goptions:"-c, --conf, description='configuration file path'"`

	goptions.Verbs

	Agent struct {
		AdvertiseAddr string `goptions:"-a, --advertise, description='address to advertise to cluster'"`
	} `goptions:"agent"`
}

func run(opts *Options) int {
	var err error

	// Print version number and exit if the version flag is set
	if opts.Version {
		fmt.Printf("sentinel v%s\n", VersionNumber)
		return 0
	}

	switch opts.Verbs {
	case "agent":
		conf, err := LoadConfiguration(opts)
		if err != nil {
			log.Error("loading configuration: %v", err)
			return 1
		}
		return RunAgent(conf)
	default:
		goptions.PrintHelp()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI %s\n", err.Error())
		return 1
	}

	return 0
}

func main() {
	// TODO set defaults
	opts := &Options{}
	goptions.ParseAndFail(opts)

	os.Exit(run(opts))
}
