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

type StartChatNavItem struct {
	page         Page
	chatRoomPage *ChatRoomPage
	*AppManager
	Name     string
	Icon     *widget.Icon
	children []NavItem
	widget.Clickable
	*material.Theme
	ThemeAlt *material.Theme
	Accordion
	url               PageURL
	updateAndNavigate bool
}

func (n *StartChatNavItem) NavTitle() string {
	return n.Name
}

func NewStartChatItem(manager *AppManager, theme *material.Theme) *StartChatNavItem {
	newChatPage := NewStartChatPage(manager, theme)
	chatRoomPage := NewChatRoomPage(manager, theme)
	icon, _ := widget.NewIcon(icons.ContentAddBox)
	return &StartChatNavItem{
		page:         newChatPage,
		chatRoomPage: chatRoomPage,
		AppManager:   manager,
		Name:         "New Chat",
		Icon:         icon,
		children:     []NavItem{},
		Theme:        theme,
		Accordion: Accordion{
			Theme: theme,
			Title: "New Chat",
			Animation: component.VisibilityAnimation{
				State:    component.Visible,
				Duration: time.Millisecond * 250,
			},
			TitleIcon: icon,
		},
		url: StartChatPageURL,
	}
}

func (n *StartChatNavItem) ClickCallback() {
	if n.CurrentPage() != n.Page() {
		n.Accordion.NoToggleOnClick = true
	} else {
		n.Accordion.NoToggleOnClick = false
	}
	n.SetSelectedItem(n)
	n.AppManager.PushPage(n.Page())
}

func (n *StartChatNavItem) IsSelected() bool {
	ok := n.NavDrawer.selectedItem == n
	return ok
}

func (n *StartChatNavItem) Layout(gtx Gtx) Dim {
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
		n.Accordion.ClickCallback = n.ClickCallback
	}
	contacts := n.Service.ContactsAddresses()

	if len(contacts) != 0 {
		children := make([]layout.FlexChild, 0, len(contacts))
		for i, _ := range contacts {
			if len(n.Children()) < i+1 {
				title := contacts[i].String()
				n.children = append(n.children, NewChatRoomItem(n.AppManager, n.Theme, title))
			}
			children = append(children, layout.Rigid(n.Children()[i].Layout))
		}
		n.Child = func(gtx Gtx) Dim {
			inset := layout.Inset{Left: unit.Dp(20), Top: unit.Dp(4), Right: unit.Dp(6)}
			return inset.Layout(gtx, func(gtx Gtx) Dim {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
			})
		}
		if n.updateAndNavigate {
			n.SetSelectedItem(n.Children()[len(n.Children())-1])
			n.Window.Invalidate()
		}
	} else {
		n.Accordion.Child = nil
	}
	if n.selectedItem == n || n.Accordion.Hovered() {
		n.Accordion.Theme = n.ThemeAlt
	} else {
		n.Accordion.Theme = n.Theme
	}
	inset := layout.Inset{Top: unit.Dp(8)}
	return inset.Layout(gtx, func(gtx Gtx) Dim {
		return n.Accordion.Layout(gtx)
	})
}

func (n *StartChatNavItem) Page() Page {
	return n.page
}

func (n *StartChatNavItem) Children() []NavItem {
	return n.children
}

func (n *StartChatNavItem) AddChild(item *ChatRoomNavItem) {
	if n.children == nil {
		n.children = make([]NavItem, 0, 1)
	}
	n.Accordion.Animation.State = component.Visible
	n.children = append(n.children, item)
}

func (n *StartChatNavItem) ReplaceChildren(children []NavItem) {
	n.Child = nil
	n.children = children
}
func (n *StartChatNavItem) URL() PageURL {
	return n.url
}

// UpdateAndNavigate when a contact is created or deleted,
//  it should be called for setting selectedNavItem
func (n *StartChatNavItem) UpdateAndNavigate() {
	n.updateAndNavigate = true
}
