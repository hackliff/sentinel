package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hackliff/sentinel"
	"github.com/hackliff/sentinel/radio"
	"github.com/hackliff/sentinel/sensor"
)

const usage = `usage: %s [OPTIONS]
	sentinel %s
		Git commit: %s
		Built time: %s
		Go version: %s

sentinel - watch on your things, so you don't have to

OPTIONS:
`

// Options stores command line flags values
type Options struct {
	ApiKey  string
	LogFile string
}

func getOpts() *Options {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0], sentinel.Version, sentinel.GitCommit, sentinel.BuildTime, sentinel.GoVersion)
		flag.PrintDefaults()
	}

	var opts = new(Options)
	flag.StringVar(&opts.ApiKey, "api-key", os.Getenv("PUSHBULLET_API_KEY"), "Pushbullet private API key")
	flag.StringVar(&opts.LogFile, "log-file", "/tmp/sentinel.log", "Log output destination")
	flag.Parse()

	if opts.ApiKey == "" {
		log.Fatalln("no api key found for pushbullet radio, aborting")
	}

	log.SetPrefix("sentinel | ")

	return opts
}

func run(opts *Options) int {
	squad_ := sentinel.NewSquad()

	buddySentinel_ := BuddySentinel{
		Radio: radio.NewPushbulletRadio(opts.ApiKey),
	}

	webSentinel_ := WebSentinel{
		Radio:  radio.NewPushbulletRadio(opts.ApiKey),
		Sensor: &sensor.HTTPSensor{},
	}

	squad_.Register("health", webSentinel_)
	squad_.Register("member-leave", buddySentinel_)
	squad_.Register("member-failed", buddySentinel_)

	if err := squad_.Dispatch(); err != nil {
		// let serf to know something went wrong
		log.Printf("error dispatching event: %s\n", err.Error())
	}

	return 0
}

func main() {
	// parse command line
	opts := getOpts()
	// send logs in a file (otherwise serf will hide it
	fd, err := os.OpenFile(opts.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer fd.Close()
	log.SetOutput(fd)

	// unix like exit code
	exitCode := run(opts)
	log.Println("exiting with code", exitCode)
	os.Exit(exitCode)
}
