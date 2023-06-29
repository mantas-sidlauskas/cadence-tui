package main

import (
	config "github.com/mantas-sidlauskas/cadence-tui/config"
	"github.com/mantas-sidlauskas/cadence-tui/tui"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {

	cfg := new(config.Config)
	if _, err := flags.Parse(cfg); err != nil {
		os.Exit(1)
	}

	tui := tui.New(cfg)
	if err := tui.Run(); err != nil {
		log.Fatal(err)
	}

}
