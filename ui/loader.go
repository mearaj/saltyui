package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type Loader struct {
	*material.Theme
}

func (l *Loader) Layout(gtx Gtx) Dim {
	var th *material.Theme
	if l.Theme == nil {
		l.Theme = material.NewTheme(gofont.Collection())
	}
	th = l.Theme
	return layout.Flex{Alignment: layout.Middle,
		Axis:    layout.Vertical,
		Spacing: layout.SpaceSides}.Layout(gtx,
		layout.Flexed(1.0, func(gtx Gtx) Dim {
			return layout.Center.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = 56
					gtx.Constraints.Min.Y = 56
					loader := material.Loader(th)
					loader.Color = th.ContrastBg
					return loader.Layout(gtx)
				},
			)
		}),
	)
}
