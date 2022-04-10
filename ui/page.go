package ui

import (
	"gioui.org/layout"
)

type (
	Gtx = layout.Context
	Dim = layout.Dimensions
)

type Page interface {
	DrawAppBar(gtx Gtx) Dim
	View
}
