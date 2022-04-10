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

// AddClientPage Always call NewNewChatPage function to create AddClientPage page
type AddClientPage struct {
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

// NewNewChatPage Always call this function to create AddClientPage page
func NewNewChatPage(manager *AppManager, th *material.Theme) *AddClientPage {
	navIcon, _ := widget.NewIcon(icons.NavigationMenu)
	iconNewChat, _ := widget.NewIcon(icons.ContentCreate)
	if th == nil {
		th = material.NewTheme(gofont.Collection())
	}
	errorTh := *th
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	return &AddClientPage{
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
	}
}

func (nc *AddClientPage) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (nc *AddClientPage) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (nc *AddClientPage) Layout(gtx Gtx) Dim {
	if nc.Theme == nil {
		nc.Theme = material.NewTheme(gofont.Collection())
	}
	th := nc.Theme
	if nc.inputNewChat.Text() != nc.inputNewChatStr {
		nc.errorNewChat = nil
		nc.errorParseAddr = nil
	}
	_, nc.errorParseAddr = saltyim.ParseAddr(nc.inputNewChat.Text())
	nc.inputNewChatStr = nc.inputNewChat.Text()
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
}

func (nc *AddClientPage) DrawAppBar(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.Y = gtx.Px(unit.Dp(56))
	th := nc.Theme
	if nc.buttonNavigation.Clicked() {
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

func (nc *AddClientPage) drawNewChatTextField(gtx Gtx) Dim {
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
			if nc.Service.CurrentIdentity() != nil {
				nc.errorNewChat = nc.Service.NewChat(nc.inputNewChat.Text())
				if nc.errorNewChat != nil {
					alog.Logger().Println(nc.errorNewChat)
					nc.errorNewChatAccordion.Animation.Appear(gtx.Now)
				} else if clients := nc.Service.Clients(); len(clients) != 0 {
					client := nc.Service.GetClient(nc.inputNewChat.Text())
					for _, navItem := range nc.AppManager.NavItems() {
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
								client.String(),
								avatarIcon,
								make([]*NavItem, 0),
								nc.Theme,
								nc.AlphaPalette,
							)
							navItem.AddChild(newChatNavItem)
							navItem.SetSelectedItem(newChatNavItem)
							nc.Window.Invalidate()
							break
						}
					}
				}
			}
			nc.addingNewClient = false
		}()
	}
	return drawFormFieldRowWithLabel(gtx, nc.Theme, labelText, labelHintText, &nc.inputNewChat, &ib)
}

func (nc *AddClientPage) drawErrorNewIDAccordion(gtx Gtx) (d Dim) {
	if nc.errorNewChat != nil {
		errView := ErrorView{}
		nc.errorNewChatAccordion.Child = errView.Layout
		errView.Error = nc.errorNewChat.Error()
		return nc.errorNewChatAccordion.Layout(gtx)
	}
	return d
}
