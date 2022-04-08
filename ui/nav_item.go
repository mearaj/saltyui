package ui

import (
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"math"
	"time"
)

type NavItem struct {
	*component.AlphaPalette
	page Page
	*NavDrawer
	Name     string
	Icon     *widget.Icon
	Children []*NavItem
	hovering bool
	widget.Clickable
	Animation component.VisibilityAnimation
	*material.Theme
}

func (n *NavItem) Layout(gtx Gtx) Dim {
	if n.Theme == nil {
		n.Theme = material.NewTheme(gofont.Collection())
	}
	th := n.Theme
	if n.Animation.Duration == time.Duration(0) {
		n.Animation.Duration = time.Millisecond * 100
		n.Animation.State = component.Invisible
	}
	if n.Clicked() && n.NavDrawer != nil {
		n.SetSelectedItem(n)
		n.AppManager.PushPage(gtx, n.page)
	}
	events := gtx.Events(n)
	for _, event := range events {
		switch event := event.(type) {
		case pointer.Event:
			switch event.Type {
			case pointer.Enter:
				n.hovering = true
			case pointer.Leave:
				n.hovering = false
			case pointer.Cancel:
				n.hovering = false
			}
		}
	}

	d := layout.Inset{Left: unit.Dp(12), Bottom: unit.Dp(6)}.Layout(gtx, func(gtx Gtx) Dim {
		d := material.Clickable(gtx, &n.Clickable, func(gtx Gtx) Dim {
			macro := op.Record(gtx.Ops)
			d := n.layoutContent(gtx, th)
			call := macro.Stop()
			pushOps := pointer.PassOp{}.Push(gtx.Ops)
			defer pushOps.Pop()
			defer clip.Rect(image.Rectangle{
				Max: d.Size,
			}).Push(gtx.Ops).Pop()
			pointer.InputOp{
				Tag:   n,
				Types: pointer.Enter | pointer.Leave,
			}.Add(gtx.Ops)
			d = layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx Gtx) Dim { return n.layoutBackground(gtx, d.Size) }),
				layout.Stacked(func(gtx Gtx) Dim {
					call.Add(gtx.Ops)
					return d
				}),
			)
			return d
		})
		if len(n.Children) == 0 {
			return d
		}
		if !n.Animation.Visible() {
			return d
		}

		children := make([]layout.FlexChild, 0, len(n.Children))
		children = append(children, layout.Rigid(
			func(gtx Gtx) Dim {
				progress := n.Animation.Revealed(gtx)
				height := int(math.Round(float64(float32(d.Size.Y) * progress)))
				d.Size.Y = height
				return d
			}))
		for _, child := range n.Children {
			children = append(children, layout.Rigid(
				func(gtx Gtx) Dim {
					return child.Layout(gtx)
				}))
		}

		d = layout.Inset{Top: unit.Dp(4), Left: unit.Dp(24)}.Layout(gtx, func(gtx Gtx) Dim {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		})
		return d
	})

	return d
}

func (n *NavItem) layoutContent(gtx Gtx, th *material.Theme) Dim {
	contentColor := th.Palette.Bg
	contentColor.A = 200
	if n == n.SelectedNavItem() || n.Hovered() {
		contentColor.A = 255
	}
	d := layout.Inset{
		Top:    unit.Dp(6),
		Right:  unit.Dp(24),
		Bottom: unit.Dp(6),
		Left:   unit.Dp(24),
	}.Layout(gtx, func(gtx Gtx) Dim {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1.0, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx Gtx) Dim {
						if n.Icon == nil {
							return layout.Dimensions{}
						}
						return layout.Inset{Right: unit.Dp(16)}.Layout(gtx,
							func(gtx Gtx) Dim {
								iconSize := gtx.Px(unit.Dp(24))
								gtx.Constraints = layout.Exact(image.Pt(iconSize, iconSize))
								return n.Icon.Layout(gtx, contentColor)
							})
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						label := material.Label(th, unit.Dp(14), n.Name)
						label.Color = contentColor
						//label.Font.Weight = text.Bold
						return layout.Center.Layout(gtx, component.TruncatingLabelStyle(label).Layout)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if len(n.Children) != 0 {
					affine := f32.Affine2D{}
					ic, _ := widget.NewIcon(icons.NavigationChevronRight)
					cl := th.Fg
					cl.A = 200
					origin := f32.Pt(12, 12)
					rotation := float32(0)
					if n.Animation.Visible() {
						rotation = float32(math.Pi * 0.5)
					}
					if n.Animation.Animating() {
						rotation *= n.Animation.Revealed(gtx)
					}
					affine = affine.Rotate(origin, rotation)
					defer op.Affine(affine).Push(gtx.Ops).Pop()
					return ic.Layout(gtx, cl)
				}
				return Dim{}
			}),
		)
	})
	return d
}

func (n *NavItem) layoutBackground(gtx Gtx, size image.Point) Dim {
	th := n.Theme
	if n != n.SelectedNavItem() && !n.hovering {
		return layout.Dimensions{}
	}
	var fill color.NRGBA
	if n.hovering {
		fill = component.WithAlpha(th.Palette.Bg, n.AlphaPalette.Hover)
	} else if n == n.SelectedNavItem() {
		fill = component.WithAlpha(th.Palette.Bg, n.AlphaPalette.Selected)
	}
	rr := float32(gtx.Px(unit.Dp(4)))
	defer clip.RRect{
		Rect: f32.Rectangle{
			Max: layout.FPt(size),
		},
		NE: rr,
		SE: rr,
		NW: rr,
		SW: rr,
	}.Push(gtx.Ops).Pop()
	PaintRect(gtx, size, fill)
	return layout.Dimensions{Size: size}
}

func (n *NavItem) Page() Page {
	return n.page
}
