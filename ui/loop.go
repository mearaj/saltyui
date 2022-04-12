package ui

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

func Loop(w *app.Window) error {
	var ops op.Ops
	am := NewAppManager(w)
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				am.Constraints = gtx.Constraints
				am.Metric = gtx.Metric
				am.Layout(gtx)
				e.Frame(gtx.Ops)
				if !am.isWindowLoaded {
					am.isWindowLoaded = true
				}
			}
		}
	}
}
