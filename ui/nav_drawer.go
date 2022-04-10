package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"time"
)

var (
	hoverOverlayAlpha    uint8 = 96
	selectedOverlayAlpha uint8 = 48
)

type NavDrawer struct {
	component.AlphaPalette
	Title        string
	Subtitle     string
	Anchor       component.VerticalAnchorPosition
	selectedItem *NavItem
	drawerItems  []*NavItem
	navList      widget.List
	NavAnim      component.VisibilityAnimation
	*material.Theme
	*AppManager
	scrollList widget.List
}

func NewNav(title, subtitle string, manager *AppManager, th *material.Theme) *NavDrawer {
	m := NavDrawer{
		Title:    title,
		Subtitle: subtitle,
		AlphaPalette: component.AlphaPalette{
			Hover:    hoverOverlayAlpha,
			Selected: selectedOverlayAlpha,
		},
		AppManager: manager,
		Theme:      th,
	}
	return &m
}

func (n *NavDrawer) AddNavItem(item *NavItem) {
	item.AlphaPalette = n.AlphaPalette
	n.drawerItems = append(n.drawerItems, item)
	if len(n.drawerItems) == 1 {
		n.drawerItems[0].selectedItem = n.drawerItems[0]
	}
}

func (n *NavDrawer) Layout(gtx Gtx, anim *component.VisibilityAnimation) Dim {
	if n.Theme == nil {
		n.Theme = material.NewTheme(gofont.Collection())
	}
	th := n.Theme
	sheet := component.NewSheet()
	return sheet.Layout(gtx, th, anim, func(gtx Gtx) Dim {
		return n.LayoutContents(gtx, anim)
	})
}

func (n *NavDrawer) LayoutContents(gtx Gtx, anim *component.VisibilityAnimation) Dim {
	th := n.Theme
	if !anim.Visible() {
		return Dim{}
	}
	spacing := layout.SpaceEnd
	if n.Anchor == component.Bottom {
		spacing = layout.SpaceStart
	}

	layout.Flex{
		Spacing: spacing,
		Axis:    layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			return layout.Inset{
				Left: unit.Dp(16),
			}.Layout(gtx, func(gtx Gtx) Dim {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx Gtx) Dim {
						gtx.Constraints.Max.Y = gtx.Px(unit.Dp(36))
						gtx.Constraints.Min = gtx.Constraints.Max
						title := material.Label(th, unit.Dp(18), n.Title)
						title.Font.Weight = text.Bold
						return layout.SW.Layout(gtx, title.Layout)
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						gtx.Constraints.Max.Y = gtx.Px(unit.Dp(20))
						gtx.Constraints.Min = gtx.Constraints.Max
						return layout.SW.Layout(gtx, material.Label(th, unit.Dp(12), n.Subtitle).Layout)
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx Gtx) Dim {
			return n.layoutNavList(gtx, th, anim)
		}),
	)
	return Dim{Size: gtx.Constraints.Max}
}

// layoutNavList draw root nav items here
func (n *NavDrawer) layoutNavList(gtx Gtx, th *material.Theme, anim *component.VisibilityAnimation) Dim {
	n.navList.Axis = layout.Vertical
	inset := layout.Inset{Left: unit.Dp(16), Top: unit.Dp(16), Right: unit.Dp(6)}
	xConstraints := gtx.Constraints.Max.X
	dim := material.List(th, &n.scrollList).Layout(gtx, 1, func(gtx Gtx, index int) Dim {
		return inset.Layout(gtx, func(gtx Gtx) Dim {
			gtx.Constraints.Max.X = xConstraints - 24.0
			dim := material.List(th, &n.navList).Layout(gtx, len(n.drawerItems), func(gtx Gtx, index int) Dim {
				inset := layout.Inset{Top: unit.Dp(8)}
				dim := inset.Layout(gtx, func(gtx Gtx) Dim {
					return n.drawerItems[index].Layout(gtx)
				})
				return dim
			})
			return dim
		})
	})
	return dim
}
func (n *NavDrawer) SetSelectedItem(item *NavItem) {
	n.selectedItem = item
	item.Animation.ToggleVisibility(time.Now())
	n.AppManager.PushPage(item.Page())
}
func (n *NavDrawer) SelectedNavItem() *NavItem {
	return n.selectedItem
}

func (n *NavDrawer) SetNavDestination(page Page) {
	for _, item := range n.drawerItems {
		if item.Page() == page {
			n.SetSelectedItem(item)
			break
		}
	}
}

func (n *NavDrawer) NavItems() []*NavItem {
	return n.drawerItems
}
