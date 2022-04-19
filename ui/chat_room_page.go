package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/saltyui/alog"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"strings"
)

// ChatRoomPage Always call NewChatRoomPage function to create ChatRoomPage page
type ChatRoomPage struct {
	layout.List
	*AppManager
	Theme            *material.Theme
	iconSendMessage  *widget.Icon
	inputMsgField    component.TextField
	inputNewChatStr  string
	buttonNavigation widget.Clickable
	submitButton     widget.Clickable
	navigationIcon   *widget.Icon
	Title            string
}

// NewChatRoomPage Always call this function to create ChatRoomPage page
func NewChatRoomPage(manager *AppManager, th *material.Theme) *ChatRoomPage {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	iconSendMessage, _ := widget.NewIcon(icons.ContentSend)
	if th == nil {
		th = material.NewTheme(gofont.Collection())
	}
	return &ChatRoomPage{
		AppManager:      manager,
		Theme:           th,
		navigationIcon:  navIcon,
		iconSendMessage: iconSendMessage,
		inputMsgField: component.TextField{
			Editor: widget.Editor{},
		},
	}
}

func (cp *ChatRoomPage) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (cp *ChatRoomPage) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (cp *ChatRoomPage) Layout(gtx Gtx) Dim {
	if cp.Theme == nil {
		cp.Theme = material.NewTheme(gofont.Collection())
	}
	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}
	return flex.Layout(gtx,
		layout.Flexed(1, cp.drawChatRoomList),
		layout.Rigid(cp.drawSendMsgField),
	)
}

func (cp *ChatRoomPage) DrawAppBar(gtx Gtx) Dim {
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
				return Dim{}
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
				titleText := cp.AppManager.SelectedItem().NavTitle()
				title := material.Body1(th, titleText)
				title.Color = th.Palette.ContrastFg
				title.TextSize = unit.Dp(18)
				return title.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx Gtx) Dim {
			return Dim{
				Size:     image.Point{X: 0},
				Baseline: 0,
			}
		}),
	)
	return Dim{Size: gtx.Constraints.Max}
}
func (cp *ChatRoomPage) drawChatRoomList(gtx Gtx) Dim {
	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			inset := layout.UniformInset(unit.Dp(8))
			return inset.Layout(gtx, func(gtx Gtx) Dim {
				msgs := cp.Service.Messages(cp.SelectedItem().NavTitle())
				cp.List.Axis = layout.Vertical
				cp.List.ScrollToEnd = true
				return cp.List.Layout(gtx, len(msgs), cp.drawChatRoomListItem)
			})
		}))
}
func (cp *ChatRoomPage) drawChatRoomListItem(gtx Gtx, index int) Dim {
	msgs := cp.Service.Messages(cp.SelectedItem().NavTitle())
	return material.Body1(cp.Theme, msgs[index].Text).Layout(gtx)
}
func (cp *ChatRoomPage) drawSendMsgField(gtx Gtx) Dim {
	if cp.submitButton.Clicked() {
		canSend := strings.Trim(cp.inputMsgField.Text(), " ") != ""
		currNavItem := cp.SelectedItem()
		canSend = canSend && currNavItem != nil
		if canSend {
			msg := cp.inputMsgField.Text()
			cp.inputMsgField.Clear()
			go func() {
				err := <-cp.Service.SendMessage(currNavItem.NavTitle(), msg)
				cp.inputMsgField.Clear()
				if err != nil {
					alog.Logger().Errorln(err)
				} else {
					alog.Logger().Println("successfully sent msg...")
				}
			}()
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
