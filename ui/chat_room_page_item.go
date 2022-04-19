package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
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
		return c.meLayout(gtx)
	}
	return c.youLayout(gtx)
}

func (c ChatRoomPageItem) meLayout(gtx Gtx) Dim {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Spacing: layout.SpaceStart}.Layout(gtx,
		layout.Rigid(func(gtx2 Gtx) Dim {
			gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) / 1.5)
			bd := material.Body1(c.Theme, c.Message.Text)
			bd.Alignment = text.End
			return bd.Layout(gtx)
		}))
}
func (c ChatRoomPageItem) youLayout(gtx Gtx) Dim {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Spacing: layout.SpaceEnd}.Layout(gtx,
		layout.Rigid(func(gtx2 Gtx) Dim {
			gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) / 1.5)
			return material.Body1(c.Theme, c.Message.Text).Layout(gtx)
		}))
}
