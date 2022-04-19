package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"image"
	"image/color"
	"runtime"
	"time"
)

type NavDrawer struct {
	DrawerTitle  string
	Subtitle     string
	Anchor       component.VerticalAnchorPosition
	selectedItem NavItem
	drawerItems  []NavItem
	navList      widget.List
	navAnim      component.VisibilityAnimation
	*material.Theme
	*AppManager
	scrollList widget.List
	*component.ModalLayer
	*component.ModalSheet
	maxWidth int
}

func NewNavDrawer(title, subtitle string, manager *AppManager, th *material.Theme, layer *component.ModalLayer) *NavDrawer {
	sheet := component.NewModalSheet(layer)
	m := NavDrawer{
		DrawerTitle: title,
		Subtitle:    subtitle,
		AppManager:  manager,
		Theme:       th,
		ModalLayer:  layer,
		ModalSheet:  sheet,
		navAnim: component.VisibilityAnimation{
			State:    component.Visible,
			Duration: time.Millisecond * 250,
		},
	}
	return &m
}

func (n *NavDrawer) SetSelectedItem(item NavItem) {
	n.selectedItem = item
}

func (n *NavDrawer) SelectedItem() NavItem {
	return n.selectedItem
}

func (n *NavDrawer) AddRootNavItem(item NavItem) {
	n.drawerItems = append(n.drawerItems, item)
	if len(n.drawerItems) == 1 {
		n.selectedItem = n.drawerItems[0]
	}
}

func (n *NavDrawer) MaxWidth() int {
	return n.maxWidth
}

func (n *NavDrawer) Layout(gtx Gtx) Dim {
	n.maxWidth = gtx.Constraints.Max.X
	if n.useModalDrawer() {
		n.navAnim.State = component.Invisible
		return n.DrawerModalLayout()
	} else {
		n.Modal.VisibilityAnimation.State = component.Invisible
		return n.DrawerLayout(gtx)
	}
}

func (n *NavDrawer) DrawerLayout(gtx Gtx) Dim {
	if n.Theme == nil {
		n.Theme = material.NewTheme(gofont.Collection())
	}
	th := n.Theme
	sheet := component.NewSheet()
	return sheet.Layout(gtx, th, &n.navAnim, func(gtx Gtx) Dim {
		return n.LayoutContents(gtx, &n.navAnim)
	})
}
func (n *NavDrawer) DrawerModalLayout() Dim {
	n.ModalSheet.LayoutModal(func(gtx Gtx, th *material.Theme, anim *component.VisibilityAnimation) Dim {
		PaintRect(gtx, gtx.Constraints.Max, th.ContrastBg)
		dims := n.LayoutContents(gtx, anim)
		return dims
	})
	return Dim{}
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
						title := material.Label(th, unit.Dp(18), n.DrawerTitle)
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
	dim := material.List(th, &n.scrollList).Layout(gtx, 1, func(gtx Gtx, index int) Dim {
		dim := inset.Layout(gtx, func(gtx Gtx) Dim {
			gtx.Constraints.Max.X = n.MaxWidth() - 24.0
			dim := material.List(th, &n.navList).Layout(gtx, len(n.drawerItems), func(gtx Gtx, index int) Dim {
				inset := layout.Inset{Top: unit.Dp(8)}
				dim := inset.Layout(gtx, func(gtx Gtx) Dim {
					return n.drawerItems[index].Layout(gtx)
				})
				return dim
			})
			return dim
		})
		return dim
	})
	return dim
}

func (n *NavDrawer) DrawerItems() []NavItem {
	return n.drawerItems
}

func (n *NavDrawer) setNavDestinationRecursively(navItem NavItem, page Page) {
	if navItem.Page() == page {
		n.SetSelectedItem(navItem)
		return
	}
	for _, item := range navItem.Children() {
		n.setNavDestinationRecursively(item, page)
	}
}

func (n *NavDrawer) SetNavDestination(page Page) {
	for _, item := range n.DrawerItems() {
		n.setNavDestinationRecursively(item, page)
	}
}

func (n *NavDrawer) ToggleVisibility(when time.Time) {
	if n.useModalDrawer() {
		n.DrawerModalLayout()
		n.ModalSheet.ToggleVisibility(when)
	} else {
		n.navAnim.ToggleVisibility(when)
	}
}

func (n *NavDrawer) Appear(when time.Time) {
	if n.useModalDrawer() {
		n.DrawerModalLayout()
		n.ModalSheet.Appear(when)
	} else {
		n.navAnim.Appear(when)
	}
}

func (n *NavDrawer) Disappear(when time.Time) {
	if n.useModalDrawer() {
		n.DrawerModalLayout()
		n.ModalSheet.Disappear(when)
	} else {
		n.navAnim.Appear(when)
	}
}

func (n *NavDrawer) useModalDrawer() bool {
	return n.GetWindowWidthInDp() < 800 || runtime.GOOS == "android" || runtime.GOOS == "ios"
}
func PaintRect(gtx Gtx, size image.Point, fill color.NRGBA) {
	component.Rect{
		Color: fill,
		Size:  size,
	}.Layout(gtx)
}
