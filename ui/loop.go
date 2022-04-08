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
				if am.WindowWidth != gtx.Constraints.Max.X {
					am.WindowWidth = gtx.Constraints.Max.X
					if !am.UseNonModalDrawer() && am.NavAnim.Visible() {
						am.NavAnim.Disappear(gtx.Now)
					}
				}
				if am.WindowHeight != gtx.Constraints.Max.Y {
					am.WindowHeight = gtx.Constraints.Max.Y
				}

				am.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}
}
