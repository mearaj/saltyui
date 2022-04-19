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

type IdentitiesNavItem struct {
	page Page
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

func (n *IdentitiesNavItem) NavTitle() string {
	return n.Name
}

func NewIdentitiesItem(manager *AppManager, theme *material.Theme) *IdentitiesNavItem {
	identitiesPage := NewIdentitiesPage(manager, theme)
	icon, _ := widget.NewIcon(icons.ActionAccountBox)

	return &IdentitiesNavItem{
		page:       identitiesPage,
		AppManager: manager,
		Name:       "Identities",
		Icon:       icon,
		children:   make([]NavItem, 0, 1),
		Theme:      theme,
		Accordion: Accordion{
			Theme: theme,
			Title: "Identities",
			Animation: component.VisibilityAnimation{
				State:    component.Visible,
				Duration: time.Millisecond * 250,
			},
			TitleIcon: icon,
		},
		url: IdentitiesPageURL,
	}
}

func (n *IdentitiesNavItem) OnClick() {
	n.SetSelectedItem(n)
	n.AppManager.PushPage(n.Page())
}

func (n *IdentitiesNavItem) IsSelected() bool {
	ok := n.NavDrawer.selectedItem == n
	return ok
}

func (n *IdentitiesNavItem) Layout(gtx Gtx) Dim {
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
			inset := layout.Inset{Left: unit.Dp(20), Top: unit.Dp(4)}
			return inset.Layout(gtx, func(gtx Gtx) Dim {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
			})
		}
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

func (n *IdentitiesNavItem) Page() Page {
	return n.page
}

func (n *IdentitiesNavItem) Children() []NavItem {
	return n.children
}

func (n *IdentitiesNavItem) AddChild(item *IdentitiesNavItem) {
	if n.children == nil {
		n.children = make([]NavItem, 0, 1)
	}
	n.Accordion.Animation.State = component.Visible
	n.children = append(n.children, item)
}

func (n *IdentitiesNavItem) ReplaceChildren(children []NavItem) {
	n.Child = nil
	n.children = children
}
func (n *IdentitiesNavItem) URL() PageURL {
	return n.url
}