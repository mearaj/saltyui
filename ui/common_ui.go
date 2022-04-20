package ui

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"image"
	"image/color"
)

func drawFormFieldRowWithLabel(gtx Gtx, th *material.Theme, labelText string, labelHintText string, textField *component.TextField, button *IconButton) Dim {
	return layout.Flex{
		Axis:      layout.Vertical,
		Spacing:   layout.SpaceStart,
		Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1.0, func(gtx Gtx) Dim {
					return layout.Inset{
						Top:    unit.Dp(0),
						Right:  unit.Dp(0),
						Bottom: unit.Dp(8.0),
						Left:   unit.Dp(0),
					}.Layout(gtx, func(gtx Gtx) Dim {
						return material.Label(th, unit.Dp(16.0), labelText).Layout(gtx)
					})
				}),
			)
		}),
		layout.Rigid(func(gtx Gtx) Dim {
			if gtx.Constraints.Max.X < 600 {
				flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd, Alignment: layout.Start}
				return flex.Layout(gtx,
					layout.Rigid(func(gtx Gtx) Dim {
						return textField.Layout(gtx,
							th, labelHintText)
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						return layout.Spacer{
							Height: unit.Dp(16),
						}.Layout(gtx)
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						gtx.Constraints.Min.X = 180
						return button.Layout(gtx)
					}),
				)
			}

			flex := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Start}
			return flex.Layout(gtx,
				layout.Flexed(1, func(gtx Gtx) Dim {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx2 Gtx) Dim {
							return textField.Layout(gtx,
								th, labelHintText)
						}),
					)
				}),
				layout.Rigid(func(gtx Gtx) Dim {
					return layout.Spacer{
						Width: unit.Dp(16),
					}.Layout(gtx)
				}),
				layout.Rigid(func(gtx Gtx) Dim {
					gtx.Constraints.Min.X = 180
					return layout.Inset{Top: unit.Dp(7)}.Layout(gtx, func(gtx Gtx) Dim {
						return button.Layout(gtx)
					})
				}),
			)
		}),
	)
}
func drawAvatar(gtx Gtx, initials string, bgColor color.NRGBA, textTheme *material.Theme) Dim {
	d := component.Rect{
		Color: bgColor,
		Size:  image.Point{X: gtx.Px(unit.Dp(48)), Y: gtx.Px(unit.Dp(48))},
		Radii: float32(gtx.Px(unit.Dp(48)) / 2),
	}.Layout(gtx)
	macro2 := op.Record(gtx.Ops)
	d2 := material.Label(textTheme, unit.Dp(20), initials).Layout(gtx)
	macro2.Stop()
	op.Offset(f32.Point{
		X: float32(d.Size.X-d2.Size.X) / 2,
		Y: float32(d.Size.Y-d2.Size.Y) / 2,
	}).Add(gtx.Ops)
	material.Label(textTheme, unit.Dp(20), initials).Layout(gtx)
	return d
}
