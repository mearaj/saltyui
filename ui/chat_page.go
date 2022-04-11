package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"log"
	"strings"
)

// ChatPage Always call NewClientPage function to create ChatPage page
type ChatPage struct {
	layout.List
	*AppManager
	Theme            *material.Theme
	iconSendMessage  *widget.Icon
	inputMsgField    component.TextField
	inputNewChatStr  string
	buttonNavigation widget.Clickable
	submitButton     widget.Clickable
	navigationIcon   *widget.Icon
}

// NewClientPage Always call this function to create ChatPage page
func NewClientPage(manager *AppManager, th *material.Theme) *ChatPage {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	iconSendMessage, _ := widget.NewIcon(icons.ContentSend)
	if th == nil {
		th = material.NewTheme(gofont.Collection())
	}
	return &ChatPage{
		AppManager:      manager,
		Theme:           th,
		navigationIcon:  navIcon,
		iconSendMessage: iconSendMessage,
		inputMsgField: component.TextField{
			Editor: widget.Editor{},
		},
	}
}

func (cp *ChatPage) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (cp *ChatPage) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (cp *ChatPage) Layout(gtx Gtx) Dim {
	if cp.Theme == nil {
		cp.Theme = material.NewTheme(gofont.Collection())
	}
	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}
	return flex.Layout(gtx,
		layout.Flexed(1, cp.drawChatRoomList),
		layout.Rigid(cp.drawSendMsgField),
	)
}

func (cp *ChatPage) DrawAppBar(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.Y = gtx.Px(unit.Dp(56))
	th := cp.Theme
	if cp.buttonNavigation.Clicked() {
		cp.AppManager.PopUp()
	}
	component.Rect{Size: gtx.Constraints.Max, Color: th.Palette.ContrastBg}.Layout(gtx)
	layout.Flex{
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			if cp.navigationIcon == nil {
				return layout.Dimensions{}
			}
			navigationIcon := cp.navigationIcon
			button := material.IconButton(th, &cp.buttonNavigation, navigationIcon, "Nav Icon Button")
			button.Size = unit.Dp(24)
			button.Background = th.Palette.ContrastBg
			button.Color = th.Palette.ContrastFg
			button.Inset = layout.UniformInset(unit.Dp(16))
			return button.Layout(gtx)
		}),
		layout.Rigid(func(gtx Gtx) Dim {
			return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx Gtx) Dim {
				titleText := cp.AppManager.SelectedNavItem().Name
				title := material.Body1(th, titleText)
				title.Color = th.Palette.ContrastFg
				title.TextSize = unit.Dp(18)
				return title.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return Dim{
				Size:     image.Point{X: 0},
				Baseline: 0,
			}
		}),
	)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
func (cp *ChatPage) drawChatRoomList(gtx Gtx) Dim {
	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			inset := layout.UniformInset(unit.Dp(8))
			return inset.Layout(gtx, func(gtx Gtx) Dim {
				return cp.List.Layout(gtx, 1, cp.drawChatRoomListItem)
			})
		}))
}
func (cp *ChatPage) drawChatRoomListItem(gtx Gtx, index int) Dim {
	return material.Body1(cp.Theme, "Message here").Layout(gtx)
}
func (cp *ChatPage) drawSendMsgField(gtx Gtx) Dim {
	if cp.submitButton.Clicked() {
		canSend := strings.Trim(cp.inputMsgField.Text(), " ") != ""
		currNavItem := cp.SelectedNavItem()
		canSend = canSend && currNavItem != nil
		if canSend {
			err := cp.Service.SendMessage(currNavItem.Name, cp.inputMsgField.Text())
			if err != nil {
				log.Println(err)
			}
		}
	}
	fl := layout.Flex{
		Axis:      layout.Horizontal,
		Spacing:   layout.SpaceBetween,
		Alignment: layout.End,
		WeightSum: 1,
	}
	gtx.Constraints.Max.Y = 200
	inset := layout.Inset{
		Top:    unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(16),
		Left:   unit.Dp(16),
	}
	return inset.Layout(gtx, func(gtx Gtx) Dim {
		return fl.Layout(gtx,
			layout.Flexed(1.0, func(gtx Gtx) Dim {
				return cp.inputMsgField.Layout(gtx, cp.Theme,
					"Enter message here...")
			}),
			layout.Rigid(func(gtx Gtx) Dim {
				inset := layout.Inset{Left: unit.Dp(8.0)}
				return inset.Layout(
					gtx,
					func(gtx Gtx) Dim {
						return material.IconButtonStyle{
							Background: cp.Theme.ContrastBg,
							Color:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
							Icon:       cp.iconSendMessage,
							Size:       unit.Dp(24.0),
							Button:     &cp.submitButton,
							Inset:      layout.UniformInset(unit.Dp(9)),
						}.Layout(gtx)
					},
				)
			}),
		)
	})
}
