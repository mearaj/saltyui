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
	"time"
)

// NewChat Always call NewNewChatPage function to create NewChat page
type NewChat struct {
	widget.List
	*AppManager
	Theme            *material.Theme
	title            string
	navItemIcon      *widget.Icon
	iconNewChat      *widget.Icon
	newIDInput       component.TextField
	inputNewChat     component.TextField
	buttonNewChat    widget.Clickable
	NavigationButton widget.Clickable
	navigationIcon   *widget.Icon
}

// NewNewChatPage Always call this function to create NewChat page
func NewNewChatPage(manager *AppManager, th *material.Theme) *NewChat {
	navItemIcon, _ := widget.NewIcon(icons.CommunicationChat)
	navIcon, _ := widget.NewIcon(icons.NavigationMenu)
	iconNewChat, _ := widget.NewIcon(icons.ContentCreate)
	if th == nil {
		th = material.NewTheme(gofont.Collection())
	}
	return &NewChat{
		AppManager:     manager,
		Theme:          th,
		navItemIcon:    navItemIcon,
		title:          "New Chat",
		navigationIcon: navIcon,
		iconNewChat:    iconNewChat,
	}
}

func (nc *NewChat) NavItem() NavItem {
	return NavItem{
		Tag:  nc,
		Name: nc.title,
		Icon: nc.navItemIcon,
	}
}

func (nc *NewChat) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (nc *NewChat) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (nc *NewChat) Layout(gtx Gtx) Dim {
	th := nc.Theme
	return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1.0, func(gtx layout.Context) layout.Dimensions {
				nc.List.Axis = layout.Vertical
				return material.List(th, &nc.List).Layout(gtx, 1, func(gtx Gtx, _ int) Dim {
					return layout.Flex{
						Alignment: layout.Middle,
						Axis:      layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx Gtx) Dim {
							return nc.drawNewChatTextField(gtx, th)
						}),
					)
				})
			}),
		)
	})
}

func (nc *NewChat) DrawAppBar(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.Y = gtx.Px(unit.Dp(56))
	th := nc.Theme
	if nc.NavigationButton.Clicked() {
		if nc.AppManager.UseNonModalDrawer() {
			nc.NavAnim.ToggleVisibility(time.Now())
		} else {
			nc.AppManager.ModalNavDrawer.Appear(gtx.Now)
			nc.NavAnim.Disappear(gtx.Now)
		}
	}
	component.Rect{Size: gtx.Constraints.Max, Color: th.Palette.ContrastBg}.Layout(gtx)
	layout.Flex{
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			if nc.navigationIcon == nil {
				return layout.Dimensions{}
			}
			navigationIcon := nc.navigationIcon
			button := material.IconButton(th, &nc.NavigationButton, navigationIcon, "Nav Icon Button")
			button.Size = unit.Dp(24)
			button.Background = th.Palette.ContrastBg
			button.Color = th.Palette.ContrastFg
			button.Inset = layout.UniformInset(unit.Dp(16))
			return button.Layout(gtx)
		}),
		layout.Rigid(func(gtx Gtx) Dim {
			return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx Gtx) Dim {
				titleText := nc.title
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

func (nc *NewChat) drawNewChatTextField(gtx Gtx, th *material.Theme) Dim {
	labelText := "New Chat"
	labelHintText := "Start Chat with user@domain"
	buttonText := "New Chat"
	ib := IconButton{
		Theme:  nc.Theme,
		Button: &nc.buttonNewChat,
		Icon:   nc.iconNewChat,
		Text:   buttonText,
	}
	return drawFormFieldRowWithLabel(gtx, th, labelText, labelHintText, &nc.newIDInput, &ib)
}
