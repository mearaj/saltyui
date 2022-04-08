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
	items        []*NavItem
	navList      widget.List
	NavAnim      component.VisibilityAnimation
	*material.Theme
	*AppManager
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

func (m *NavDrawer) AddNavItem(item *NavItem) {
	item.AlphaPalette = &m.AlphaPalette
	m.items = append(m.items, item)
	if len(m.items) == 1 {
		m.items[0].selectedItem = m.items[0]
	}
}

func (m *NavDrawer) Layout(gtx Gtx, anim *component.VisibilityAnimation) Dim {
	if m.Theme == nil {
		m.Theme = material.NewTheme(gofont.Collection())
	}
	th := m.Theme
	sheet := component.NewSheet()
	return sheet.Layout(gtx, th, anim, func(gtx Gtx) Dim {
		return m.LayoutContents(gtx, anim)
	})
}

func (m *NavDrawer) LayoutContents(gtx Gtx, anim *component.VisibilityAnimation) Dim {
	th := m.Theme
	if !anim.Visible() {
		return Dim{}
	}
	spacing := layout.SpaceEnd
	if m.Anchor == component.Bottom {
		spacing = layout.SpaceStart
	}

	layout.Flex{
		Spacing: spacing,
		Axis:    layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			return layout.Inset{
				Left:   unit.Dp(16),
				Bottom: unit.Dp(18),
			}.Layout(gtx, func(gtx Gtx) Dim {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx Gtx) Dim {
						gtx.Constraints.Max.Y = gtx.Px(unit.Dp(36))
						gtx.Constraints.Min = gtx.Constraints.Max
						title := material.Label(th, unit.Dp(18), m.Title)
						title.Font.Weight = text.Bold
						return layout.SW.Layout(gtx, title.Layout)
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						gtx.Constraints.Max.Y = gtx.Px(unit.Dp(20))
						gtx.Constraints.Min = gtx.Constraints.Max
						return layout.SW.Layout(gtx, material.Label(th, unit.Dp(12), m.Subtitle).Layout)
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx Gtx) Dim {
			return m.layoutNavList(gtx, th, anim)
		}),
	)
	return Dim{Size: gtx.Constraints.Max}
}

func (m *NavDrawer) layoutNavList(gtx Gtx, th *material.Theme, anim *component.VisibilityAnimation) Dim {
	m.navList.Axis = layout.Vertical
	return material.List(th, &m.navList).Layout(gtx, len(m.items), func(gtx Gtx, index int) Dim {
		dimensions := m.items[index].Layout(gtx)
		return dimensions
	})
}
func (m *NavDrawer) SetSelectedItem(item *NavItem) {
	m.selectedItem = item
}
func (m *NavDrawer) SelectedNavItem() *NavItem {
	return m.selectedItem
}

func (m *NavDrawer) SetNavDestination(page Page) {
	for _, item := range m.items {
		if item.Page() == page {
			m.SetSelectedItem(item)
			break
		}
	}
}

type ModalNavDrawer struct {
	*NavDrawer
	sheet *component.ModalSheet
}

// NewModalNav configures a modal navigation drawer that will render itself into the provided ModalLayer
func NewModalNav(modal *component.ModalLayer, title, subtitle string, manager *AppManager, theme *material.Theme) *ModalNavDrawer {
	nav := NewNav(title, subtitle, manager, theme)
	return ModalNavFrom(nav, modal)
}

func ModalNavFrom(nav *NavDrawer, modal *component.ModalLayer) *ModalNavDrawer {
	m := &ModalNavDrawer{}
	modalSheet := component.NewModalSheet(modal)
	m.NavDrawer = nav
	m.sheet = modalSheet
	return m
}

func (m *ModalNavDrawer) Layout() Dim {
	m.sheet.LayoutModal(func(gtx Gtx, th *material.Theme, anim *component.VisibilityAnimation) Dim {
		PaintRect(gtx, gtx.Constraints.Max, th.ContrastBg)
		dims := m.NavDrawer.LayoutContents(gtx, anim)
		// FixMe:
		//if m.selectedChanged {
		//	anim.Disappear(gtx.Now)
		//}
		return dims
	})
	return Dim{}
}

func (m *ModalNavDrawer) ToggleVisibility(when time.Time) {
	m.Layout()
	m.sheet.ToggleVisibility(when)
}

func (m *ModalNavDrawer) Appear(when time.Time) {
	m.Layout()
	m.sheet.Appear(when)
}

func (m *ModalNavDrawer) Disappear(when time.Time) {
	m.Layout()
	m.sheet.Disappear(when)
}

func PaintRect(gtx Gtx, size image.Point, fill color.NRGBA) {
	component.Rect{
		Color: fill,
		Size:  size,
	}.Layout(gtx)
}
