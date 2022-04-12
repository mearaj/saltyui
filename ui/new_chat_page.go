package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/saltyui/alog"
	"go.mills.io/saltyim"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"time"
)

// NewChatPage Always call NewNewChatPage function to create NewChatPage page
type NewChatPage struct {
	widget.List
	*AppManager
	Theme                 *material.Theme
	title                 string
	iconNewChat           *widget.Icon
	inputNewChat          component.TextField
	inputNewChatStr       string
	buttonNewChat         widget.Clickable
	buttonNavigation      widget.Clickable
	navigationIcon        *widget.Icon
	errorNewChatAccordion Accordion
	errorNewChat          error
	errorParseAddr        error
	addingNewClient       bool
}

// NewNewChatPage Always call this function to create NewChatPage page
func NewNewChatPage(manager *AppManager, th *material.Theme) *NewChatPage {
	navIcon, _ := widget.NewIcon(icons.NavigationMenu)
	iconNewChat, _ := widget.NewIcon(icons.ContentCreate)
	if th == nil {
		th = material.NewTheme(gofont.Collection())
	}
	errorTh := *th
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	return &NewChatPage{
		AppManager:     manager,
		Theme:          th,
		title:          "New Chat",
		navigationIcon: navIcon,
		iconNewChat:    iconNewChat,
		errorNewChatAccordion: Accordion{
			Theme: &errorTh,
			Title: "Create New Chat Error",
			Animation: component.VisibilityAnimation{
				State:    component.Visible,
				Duration: time.Millisecond * 250,
			},
		},
		inputNewChat: component.TextField{Editor: widget.Editor{Submit: true, SingleLine: true}},
	}
}

func (nc *NewChatPage) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (nc *NewChatPage) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (nc *NewChatPage) Layout(gtx Gtx) Dim {
	if nc.Theme == nil {
		nc.Theme = material.NewTheme(gofont.Collection())
	}
	th := nc.Theme
	if nc.addingNewClient {
		loader := Loader{nc.Theme}
		return loader.Layout(gtx)
	}
	if nc.inputNewChat.Text() != nc.inputNewChatStr {
		nc.errorNewChat = nil
		nc.errorParseAddr = nil
	}
	_, nc.errorParseAddr = saltyim.ParseAddr(nc.inputNewChat.Text())
	nc.inputNewChatStr = nc.inputNewChat.Text()
	d := layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1.0, func(gtx layout.Context) layout.Dimensions {
				nc.List.Axis = layout.Vertical
				return material.List(th, &nc.List).Layout(gtx, 1, func(gtx Gtx, _ int) Dim {
					return layout.Flex{
						Alignment: layout.Middle,
						Axis:      layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx Gtx) Dim {
							return nc.drawNewChatTextField(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return layout.Spacer{Height: unit.Dp(32)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return nc.drawErrorNewIDAccordion(gtx)
						}),
					)
				})
			}),
		)
	})
	return d
}

func (nc *NewChatPage) DrawAppBar(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.Y = gtx.Px(unit.Dp(56))
	th := nc.Theme
	if nc.buttonNavigation.Clicked() {
		nc.NavDrawer.ToggleVisibility(time.Now())
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
			button := material.IconButton(th, &nc.buttonNavigation, navigationIcon, "Nav Icon Button")
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

func (nc *NewChatPage) drawNewChatTextField(gtx Gtx) Dim {
	labelText := "New Chat"
	labelHintText := "Start Chat with user@domain"
	buttonText := "New Chat"
	var button *widget.Clickable
	var th *material.Theme
	if nc.errorParseAddr != nil {
		button = &widget.Clickable{}
		th = material.NewTheme(gofont.Collection())
		th.ContrastBg = color.NRGBA(colornames.Grey500)
	} else {
		button = &nc.buttonNewChat
		th = nc.Theme
	}
	ib := IconButton{
		Theme:  th,
		Button: button,
		Icon:   nc.iconNewChat,
		Text:   buttonText,
	}

	if button.Clicked() && !nc.addingNewClient {
		nc.addingNewClient = true
		go func() {
			if nc.Service.CurrentIdentity() == nil {
				nc.addingNewClient = false
				nc.AppManager.PopUp()
				return
			}
			nc.errorNewChat = nc.Service.NewChat(nc.inputNewChat.Text())
			if nc.errorNewChat != nil {
				alog.Logger().Println(nc.errorNewChat)
				nc.errorNewChatAccordion.Animation.Appear(gtx.Now)
			} else if addrs := nc.Service.Addresses(); len(addrs) != 0 {
				addr := nc.Service.GetAddr(nc.inputNewChat.Text())
				for _, navItem := range nc.AppManager.DrawerItems() {
					if navItem != nil && navItem.Page() == nc {
						var page Page
						if len(navItem.Children()) == 0 {
							page = NewClientPage(nc.AppManager, nc.Theme)
						} else {
							page = navItem.Children()[0].Page()
						}
						avatarIcon, _ := widget.NewIcon(icons.SocialPerson)
						newChatNavItem := NewNavItem(page,
							nc.NavDrawer,
							addr.String(),
							avatarIcon,
							make([]*NavItem, 0),
							nc.Theme,
							ChatPageUrl,
						)
						navItem.AddChild(newChatNavItem)
						navItem.SetSelectedItem(newChatNavItem)
						if nc.NavDrawer.CurrentPage() == nc {
							nc.AppManager.PushPage(newChatNavItem.Page())
						}
						nc.Window.Invalidate()
						break
					}
				}
			}
			nc.addingNewClient = false
		}()
	}
	return drawFormFieldRowWithLabel(gtx, nc.Theme, labelText, labelHintText, &nc.inputNewChat, &ib)
}

func (nc *NewChatPage) drawErrorNewIDAccordion(gtx Gtx) (d Dim) {
	if nc.errorNewChat != nil {
		errView := ErrorView{}
		nc.errorNewChatAccordion.Child = errView.Layout
		errView.Error = nc.errorNewChat.Error()
		return nc.errorNewChatAccordion.Layout(gtx)
	}
	return d
}
