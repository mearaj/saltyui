package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/saltyui/service"
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
	var isMe bool
	isMe = c.Message.UserAddr == c.Message.From
	if isMe {
		return c.commonLayout(gtx, layout.SpaceStart)
	}
	return c.commonLayout(gtx, layout.SpaceEnd)
}

func (c *ChatRoomPageItem) commonLayout(gtx Gtx, spacing layout.Spacing) Dim {
	macro := op.Record(gtx.Ops)
	inset := layout.UniformInset(unit.Dp(12))
	d := inset.Layout(gtx, func(gtx Gtx) Dim {
		flex := layout.Flex{Spacing: spacing}
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
}
