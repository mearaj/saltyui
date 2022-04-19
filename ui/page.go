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

type PageURL string

const (
	SettingsPageURL   PageURL = "/settings"
	IdentitiesPageURL PageURL = "/identities"
	StartChatPageURL  PageURL = "/new-chat"
	ChatRoomPageURL   PageURL = "/chat"
)
