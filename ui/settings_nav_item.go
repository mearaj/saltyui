package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"time"
)

type SettingsNavItem struct {
	page Page
	*AppManager
	Name     string
	Icon     *widget.Icon
	children []NavItem
	widget.Clickable
	*material.Theme
	//ThemeAlt *material.Theme
	Accordion
	url PageURL
}

func (n *SettingsNavItem) NavTitle() string {
	return n.Title
}

func NewSettingsItem(manager *AppManager, theme *material.Theme) *SettingsNavItem {
	settingsPage := NewSettingsPage(manager, theme)
	identitiesItem := NewIDsNavItem(manager, theme)
	icon, _ := widget.NewIcon(icons.ActionSettings)
	return &SettingsNavItem{
		page:       settingsPage,
		AppManager: manager,
		Name:       "Settings",
		Icon:       icon,
		children:   []NavItem{identitiesItem},
		Theme:      theme,
		Accordion: Accordion{
			Theme: theme,
			Title: "Settings",
			Animation: component.VisibilityAnimation{
				State:    component.Visible,
				Duration: time.Millisecond * 250,
			},
			TitleIcon: icon,
		},
		url: SettingsPageURL,
	}
}

func (n *SettingsNavItem) ClickCallback() {
	if n.CurrentPage() != n.Page() {
		n.Accordion.NoToggleOnClick = true
	} else {
		n.Accordion.NoToggleOnClick = false
	}
	n.SetSelectedItem(n)
	n.AppManager.PushPage(n.Page())
}

func (n *SettingsNavItem) IsSelected() bool {
	ok := n.NavDrawer.selectedItem == n
	return ok
}

func (n *SettingsNavItem) Layout(gtx Gtx) Dim {
	if n.Theme == nil {
		n.Theme = material.NewTheme(gofont.Collection())
	}

	if n.Accordion.ClickCallback == nil {
		n.Accordion.ClickCallback = n.ClickCallback
	}

	if len(n.Children()) != 0 {
		children := make([]layout.FlexChild, 0, len(n.Children()))
		for _, child := range n.Children() {
			children = append(children, layout.Rigid(child.Layout))
		}
		n.Child = func(gtx Gtx) Dim {
			inset := layout.Inset{Left: unit.Dp(20), Top: unit.Dp(4)}
			return inset.Layout(gtx, func(gtx Gtx) Dim {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
			})
		}
	}
	if n.selectedItem == n || n.Hovered() {
		n.Accordion.ButtonLayoutStyle.Background = n.Theme.ContrastBg
	}
	inset := layout.Inset{Top: unit.Dp(8)}
	return inset.Layout(gtx, func(gtx Gtx) Dim {
		return n.Accordion.Layout(gtx)
	})
}

func (n *SettingsNavItem) Page() Page {
	return n.page
}

func (n *SettingsNavItem) Children() []NavItem {
	return n.children
}

func (n *SettingsNavItem) URL() PageURL {
	return n.url
}
