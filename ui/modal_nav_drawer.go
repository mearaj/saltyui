package ui

import (
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"image"
	"image/color"
	"time"
)

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
