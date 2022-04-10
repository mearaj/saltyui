package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
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

			flex := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.End}
			return flex.Layout(gtx,
				layout.Flexed(1, func(gtx Gtx) Dim {
					return textField.Layout(gtx,
						th, labelHintText)
				}),
				layout.Rigid(func(gtx Gtx) Dim {
					return layout.Spacer{
						Width: unit.Dp(16),
					}.Layout(gtx)
				}),
				layout.Rigid(func(gtx Gtx) Dim {
					gtx.Constraints.Min.X = 180
					return button.Layout(gtx)
				}),
			)
		}),
	)
}
