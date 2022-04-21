package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"time"
)

type ChatRoomNavItem struct {
	page *ChatRoomPage
	*AppManager
	Name     string
	Icon     *widget.Icon
	children []NavItem
	widget.Clickable
	*material.Theme
	ThemeAlt *material.Theme
	Accordion
	url PageURL
}

func (n *ChatRoomNavItem) NavTitle() string {
	return n.Name
}

func NewChatRoomItem(manager *AppManager, theme *material.Theme, navTitle string) *ChatRoomNavItem {
	page := NewChatRoomPage(manager, theme)
	page.Title = navTitle
	icon, _ := widget.NewIcon(icons.CommunicationChat)

	return &ChatRoomNavItem{
		page:       page,
		AppManager: manager,
		Icon:       icon,
		Name:       navTitle,
		children:   make([]NavItem, 0, 1),
		Theme:      theme,
		Accordion: Accordion{
			Theme: theme,
			Title: navTitle,
			Animation: component.VisibilityAnimation{
				State:    component.Visible,
				Duration: time.Millisecond * 250,
			},
			TitleIcon: icon,
		},
		url: ChatRoomPageURL,
	}
}

func (n *ChatRoomNavItem) OnClick() {
	n.SetSelectedItem(n)
	n.AppManager.PushPage(n.Page())
}

func (n *ChatRoomNavItem) IsSelected() bool {
	ok := n.NavDrawer.selectedItem == n
	return ok
}

func (n *ChatRoomNavItem) Layout(gtx Gtx) Dim {
	if n.Theme == nil {
		n.Theme = material.NewTheme(gofont.Collection())
	}
	if n.ThemeAlt == nil {
		newTheme := *n.Theme
		newTheme.ContrastBg = color.NRGBA{}
		newTheme.ContrastBg.A = 50
		n.ThemeAlt = &newTheme
	}
	if n.Accordion.ClickCallback == nil {
		n.Accordion.ClickCallback = n.OnClick
	}

	if len(n.Children()) != 0 {
		children := make([]layout.FlexChild, 0, len(n.Children()))
		for _, child := range n.Children() {
			children = append(children, layout.Rigid(child.Layout))
		}
		n.Child = func(gtx Gtx) Dim {
			inset := layout.Inset{Left: unit.Dp(20), Top: unit.Dp(4), Right: unit.Dp(6)}
			return inset.Layout(gtx, func(gtx Gtx) Dim {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
			})
		}
	}
	if n.SelectedItem() == n || n.Hovered() {
		n.Accordion.Background = color.NRGBA{A: 50}
	} else {
		n.Accordion.Background = n.Theme.ContrastBg
	}
	inset := layout.Inset{Top: unit.Dp(8)}
	return inset.Layout(gtx, func(gtx Gtx) Dim {
		return n.Accordion.Layout(gtx)
	})
}

func (n *ChatRoomNavItem) Page() Page {
	return n.page
}

func (n *ChatRoomNavItem) Children() []NavItem {
	return n.children
}

func (n *ChatRoomNavItem) AddChild(item NavItem) {
	if n.children == nil {
		n.children = make([]NavItem, 0, 1)
	}
	n.Accordion.Animation.State = component.Visible
	n.children = append(n.children, item)
}

func (n *ChatRoomNavItem) ReplaceChildren(children []NavItem) {
	n.Child = nil
	n.children = children
}
func (n *ChatRoomNavItem) URL() PageURL {
	return n.url
}
