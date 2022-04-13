package main

import (
	"gioui.org/app"
	"github.com/mearaj/saltyui/ui"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		return
	}
	// parse string, this is built-in feature of logrus
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		return
	}
	// set global log level
	log.SetLevel(ll)
}

func main() {
	go func() {
		w := app.NewWindow(app.Title("Salty UI"))
		if err := ui.Loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
