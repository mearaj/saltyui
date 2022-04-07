package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// DetailRow lays out two widgets in a horizontal row, with the left
// widget considered the "Primary" widget.
type DetailRow struct {
	// PrimaryWidth is the fraction of the available width that should
	// be allocated to the primary widget. It should be in the range
	// (0,1.0]. Defaults to 0.3 if not set.
	PrimaryWidth float32
	// Inset is automatically applied to both widgets. This inset is
	// required, and will default to a uniform 8DP inset if not set.
	layout.Inset
}

var DefaultInset = layout.UniformInset(unit.Dp(8))

// Layout the DetailRow with the provided widgets.
func (d DetailRow) Layout(gtx Gtx, primary, detail layout.Widget) Dim {
	if d.PrimaryWidth == 0 {
		d.PrimaryWidth = 0.3
	}
	if d.Inset == (layout.Inset{}) {
		d.Inset = DefaultInset
	}
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(d.PrimaryWidth, func(gtx Gtx) Dim {
			return d.Inset.Layout(gtx, primary)
		}),
		layout.Flexed(1-d.PrimaryWidth, func(gtx Gtx) Dim {
			return d.Inset.Layout(gtx, detail)
		}),
	)
}
