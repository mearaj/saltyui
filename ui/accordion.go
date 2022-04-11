package ui

import (
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"math"
	"time"
)

type Accordion struct {
	Animation component.VisibilityAnimation
	component.AlphaPalette
	widget.Clickable
	Child layout.Widget
	*material.Theme
	icon          *widget.Icon
	Title         string
	TitleIcon     *widget.Icon
	ClickCallback func()
}

func (a *Accordion) Layout(gtx Gtx) (d Dim) {
	if a.icon == nil {
		a.icon, _ = widget.NewIcon(icons.NavigationChevronRight)
	}
	if a.Theme == nil {
		a.Theme = material.NewTheme(gofont.Collection())
	}
	th := a.Theme
	if a.Animation.Duration == time.Duration(0) {
		a.Animation.Duration = time.Millisecond * 100
		a.Animation.State = component.Invisible
	}
	if a.Clicked() {
		a.Animation.ToggleVisibility(gtx.Now)
		if a.ClickCallback != nil {
			a.ClickCallback()
		}
	}

	d = layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			d = material.ButtonLayout(th, &a.Clickable).Layout(gtx,
				func(gtx Gtx) Dim {
					return a.layoutHeader(gtx)
				},
			)
			return d
		}),
		layout.Rigid(func(gtx Gtx) (d Dim) {
			if a.Child != nil {
				progress := a.Animation.Revealed(gtx)
				macro := op.Record(gtx.Ops)
				d = layout.Flex{}.Layout(gtx, layout.Flexed(1.0, func(gtx Gtx) Dim {
					return layout.Inset{
						Top:    unit.Dp(0),
						Bottom: unit.Dp(6),
						Left:   unit.Dp(12),
						Right:  unit.Dp(12),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return a.Child(gtx)
					})
				}))
				call := macro.Stop()
				height := int(math.Round(float64(float32(d.Size.Y) * progress)))
				d.Size.Y = height
				defer clip.Rect(image.Rectangle{
					Max: d.Size,
				}).Push(gtx.Ops).Pop()
				call.Add(gtx.Ops)
			}
			return d
		}),
	)

	return d
}

func (a *Accordion) layoutHeader(gtx Gtx) Dim {
	th := a.Theme

	d := layout.Inset{
		Top:    unit.Dp(6),
		Right:  unit.Dp(12),
		Bottom: unit.Dp(6),
		Left:   unit.Dp(12),
	}.Layout(gtx, func(gtx Gtx) Dim {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1.0, func(gtx Gtx) Dim {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if a.TitleIcon != nil {
							return layout.Flex{}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return a.TitleIcon.Layout(gtx, th.ContrastFg)
							}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
							}))
						}
						return Dim{}
					}),
					layout.Rigid(func(gtx Gtx) Dim {
						label := material.Label(th, unit.Dp(14), a.Title)
						label.Color = th.ContrastFg
						//label.Font.Weight = text.Bold
						return layout.Center.Layout(gtx, component.TruncatingLabelStyle(label).Layout)
					}),
				)
			}),
			layout.Rigid(func(gtx Gtx) (d Dim) {
				if a.Child != nil {
					affine := f32.Affine2D{}
					ic, _ := widget.NewIcon(icons.NavigationChevronRight)
					cl := th.ContrastFg
					origin := f32.Pt(12, 12)
					rotation := float32(0)
					if a.Animation.Visible() {
						rotation = float32(math.Pi * 0.5)
					}
					if a.Animation.Animating() {
						rotation *= a.Animation.Revealed(gtx)
					}
					affine = affine.Rotate(origin, rotation)
					defer op.Affine(affine).Push(gtx.Ops).Pop()
					return ic.Layout(gtx, cl)
				}
				return d
			}),
		)
	})
	return d
}
