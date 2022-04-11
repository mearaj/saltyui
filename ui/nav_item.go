package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"image/color"
	"time"
)

type NavItem struct {
	page Page
	*NavDrawer
	Name     string
	Icon     *widget.Icon
	children []*NavItem
	widget.Clickable
	*material.Theme
	ThemeAlt *material.Theme
	Accordion
}

func NewNavItem(page Page, drawer *NavDrawer, name string, icon *widget.Icon, children []*NavItem, th *material.Theme) *NavItem {
	return &NavItem{
		page:      page,
		NavDrawer: drawer,
		Name:      name,
		Icon:      icon,
		children:  children,
		Theme:     th,
		Accordion: Accordion{
			Theme: th,
			Title: name,
			Animation: component.VisibilityAnimation{
				State:    component.Visible,
				Duration: time.Millisecond * 250,
			},
			TitleIcon: icon,
		},
	}
}

func (n *NavItem) OnClick() {
	n.SetSelectedItem(n)
}

func (n *NavItem) IsSelected() bool {
	ok := n.NavDrawer.selectedItem == n
	return ok
}

func (n *NavItem) Layout(gtx Gtx) Dim {
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
		n.Child = func(gtx layout.Context) layout.Dimensions {
			inset := layout.Inset{Left: unit.Dp(20), Top: unit.Dp(4), Right: unit.Dp(6)}
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
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return n.Accordion.Layout(gtx)
	})
}

func (n *NavItem) Page() Page {
	return n.page
}

func (n *NavItem) Children() []*NavItem {
	return n.children
}

func (n *NavItem) AddChild(item *NavItem) {
	if n.children == nil {
		n.children = make([]*NavItem, 0, 1)
	}
	n.Accordion.Animation.State = component.Visible
	n.children = append(n.children, item)
}
