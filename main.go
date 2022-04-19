package main

import (
	"gioui.org/app"
	"github.com/mearaj/saltyui/ui"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	//log.SetLevel(log.DebugLevel)
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
