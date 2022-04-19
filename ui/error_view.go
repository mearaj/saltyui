package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type ErrorView struct {
	*material.Theme
	*AppManager
	Error string
}

func (i *ErrorView) Layout(gtx Gtx) (d Dim) {
	if i.Theme == nil {
		i.Theme = material.NewTheme(gofont.Collection())
	}
	if i.Error != "" {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx Gtx) Dim {
				return material.Body1(i.Theme, i.Error).Layout(gtx)
			}),
		)
	}
	return d
}
