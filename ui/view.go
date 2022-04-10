package ui

type View interface {
	Layout(gtx Gtx) Dim
}
