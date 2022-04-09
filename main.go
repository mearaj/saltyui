package main

import (
	"gioui.org/app"
	"github.com/mearaj/saltyui/alog"
	"github.com/mearaj/saltyui/ui"
	"os"
)

func main() {
	go func() {
		w := app.NewWindow(app.Title("Salty UI"))
		if err := ui.Loop(w); err != nil {
			alog.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
