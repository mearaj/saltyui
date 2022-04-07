package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type IDDetailsView struct {
	*material.Theme
	widget.Clickable
	*AppManager
}

func (i *IDDetailsView) Layout(gtx Gtx) (d Dim) {
	if i.Theme == nil {
		i.Theme = material.NewTheme(gofont.Collection())
	}
	if i.AppManager.Service.CurrentIdentity() != nil {
		var contents = string(i.AppManager.Service.CurrentIdentity().Contents())
		if i.Clickable.Clicked() {
			i.AppManager.Window.WriteClipboard(contents)
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body1(i.Theme, contents).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Button(i.Theme, &i.Clickable, "Copy to Clipboard").Layout(gtx)
			}),
		)
	}
	return d
}
