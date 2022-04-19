package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/saltyui/service"
	"golang.org/x/image/colornames"
	"image/color"
	"time"
)

type ChatRoomPageItem struct {
	service.Message
	*material.Theme
}

func (c *ChatRoomPageItem) Layout(gtx Gtx) (d Dim) {
	if c.Message.Text == "" {
		return d
	}
	if c.Theme == nil {
		c.Theme = material.NewTheme(gofont.Collection())
	}

	isMe := c.Message.UserAddr == c.Message.From
	inset := layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(16)}
	d = inset.Layout(gtx, func(gtx Gtx) Dim {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		flex := layout.Flex{Axis: layout.Vertical}
		if isMe {
			flex.Alignment = layout.End
		} else {
			flex.Alignment = layout.Start
		}
		d := flex.Layout(gtx,
			layout.Rigid(func(gtx Gtx) Dim {
				timeVal, _ := time.Parse(time.RFC3339, c.Message.Created)
				txtMsg := timeVal.Format("Mon Jan 2 15:04 2006")
				label := material.Label(c.Theme, unit.Sp(12), txtMsg)
				label.Color = color.NRGBA{
					R: colornames.Gray.R,
					G: colornames.Gray.G,
					B: colornames.Gray.B,
					A: colornames.Gray.A,
				}
				label.Font.Weight = text.Bold
				label.Font.Style = text.Italic
				inset = layout.Inset{Bottom: unit.Dp(8.0)}
				d := inset.Layout(gtx, func(gtx Gtx) Dim {
					flex := layout.Flex{}
					if isMe {
						flex.Spacing = layout.SpaceStart
					} else {
						flex.Spacing = layout.SpaceEnd
					}
					return flex.Layout(gtx, layout.Rigid(func(gtx Gtx) Dim {
						return component.TruncatingLabelStyle(label).Layout(gtx)
					}))
				})
				return d
			}),
			layout.Rigid(func(gtx2 Gtx) Dim {
				macro := op.Record(gtx.Ops)
				inset := layout.UniformInset(unit.Dp(12))
				d := inset.Layout(gtx, func(gtx Gtx) Dim {
					flex := layout.Flex{}
					gtx.Constraints.Min.X = 0
					return flex.Layout(gtx,
						layout.Rigid(func(gtx2 Gtx) Dim {
							gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) / 1.5)
							bd := material.Body1(c.Theme, c.Message.Text)
							return bd.Layout(gtx)
						}))
				})
				bgColor := c.Theme.ContrastBg
				bgColor.A = 50
				rect := component.Rect{Color: bgColor, Size: d.Size, Radii: unit.Dp(8).V}
				call := macro.Stop()
				defer call.Add(gtx.Ops)
				return rect.Layout(gtx)
			}))
		return d
	})
	return d
}
