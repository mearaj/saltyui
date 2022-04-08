package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
)

type IconButton struct {
	*material.Theme
	Button *widget.Clickable
	Icon   *widget.Icon
	Text   string
}

func (b *IconButton) Layout(gtx Gtx) Dim {
	button := b.Button
	if button == nil {
		button = &widget.Clickable{}
	}
	return material.ButtonLayout(b.Theme, button).Layout(gtx, func(gtx Gtx) Dim {
		return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx Gtx) Dim {
			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			textIconSpacer := unit.Dp(5)

			layIcon := layout.Rigid(func(gtx Gtx) Dim {
				return layout.Inset{Right: textIconSpacer}.Layout(gtx, func(gtx Gtx) Dim {
					var d Dim
					if b.Icon != nil {
						size := gtx.Px(unit.Dp(56)) / 3
						gtx.Constraints = layout.Exact(image.Pt(size, size))
						d = b.Icon.Layout(gtx, b.Theme.ContrastFg)
					}
					return d
				})
			})

			layLabel := layout.Rigid(func(gtx Gtx) Dim {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx Gtx) Dim {
					l := material.Body1(b.Theme, b.Text)
					l.Color = b.Theme.Palette.ContrastFg
					return l.Layout(gtx)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)
		})
	})
}
