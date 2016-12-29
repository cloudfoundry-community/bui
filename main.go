package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/bui/api"
	"github.com/starkandwayne/goutils/log"
	"github.com/voxelbrain/goptions"
)

type BuiServerOpts struct {
	Help       bool   `goptions:"-h, --help, description='Show the help screen'"`
	ConfigFile string `goptions:"-c, --config, description='Path to the api configuration file'"`
	Log        string `goptions:"-l, --log-level, description='Set logging level to debug, info, notice, warn, error, crit, alert, or emerg'"`
	Version    bool   `goptions:"-v, --version, description='Display the api version'"`
}

var Version = ""

func main() {
	api.Version = Version
	var opts BuiServerOpts
	opts.Log = "Info"
	if err := goptions.Parse(&opts); err != nil {
		fmt.Printf("%s\n", err)
		goptions.PrintHelp()
		return
	}

	if opts.Help {
		goptions.PrintHelp()
		os.Exit(0)
	}
	if opts.Version {
		if Version == "" {
			fmt.Printf("bui (development)\n")
		} else {
			fmt.Printf("bui v%s\n", Version)
		}
		os.Exit(0)
	}

	if opts.ConfigFile == "" {
		fmt.Fprintf(os.Stderr, "No config specified. Please try again using the -c/--config argument\n")
		os.Exit(1)
	}

	log.SetupLogging(log.LogConfig{Type: "console", Level: opts.Log})
	log.Infof("starting bui api")

	a := api.NewApi()
	if err := a.ReadConfig(opts.ConfigFile); err != nil {
		log.Errorf("Failed to load config: %s", err)
		return
	}

	if err := a.Run(); err != nil {
		log.Errorf("bui failed: %s", err)
	}
	log.Infof("stopping bui api")
}
