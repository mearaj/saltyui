package ui

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/widget"
	"github.com/mearaj/saltyui/service"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type (
	Gtx = layout.Context
	Dim = layout.Dimensions
)

type Page interface {
	DrawAppBar(gtx Gtx) Dim
	View
}

type View interface {
	Layout(gtx Gtx) Dim
}

// AppManager Always call NewAppManager function to create AppManager instance
type AppManager struct {
	*app.Window
	currentPage Page
	history     []Page
	*NavDrawer
	*ModalNavDrawer
	NavAnim component.VisibilityAnimation
	*component.ModalLayer
	BottomBar bool
	*material.Theme
	WindowHeight int
	WindowWidth  int
	Service      *service.Service
}

// NewAppManager Always call this function to create AppManager instance
func NewAppManager(w *app.Window) *AppManager {
	s := service.NewService()
	am := &AppManager{Window: w, Service: s}
	am.init()
	return am
}

func (a *AppManager) UseNonModalDrawer() bool {
	return a.WindowWidth >= 800
}

func (a *AppManager) PushPage(gtx Gtx, page Page) {
	if len(a.history) == 0 {
		a.history = make([]Page, 1)
	}
	a.history = append(a.history, page)
	a.currentPage = a.history[len(a.history)-1]
	op.InvalidateOp{}.Add(gtx.Ops)
}

func (a *AppManager) PopUp() {
	if len(a.history) > 1 {
		i := len(a.history) - 1
		a.history = a.history[:i]
		a.currentPage = a.history[i-1]
	}
}

func (a *AppManager) init() {
	a.Theme = material.NewTheme(gofont.Collection())
	a.ModalLayer = component.NewModal()
	navDrwTh := *a.Theme
	navDrwTh.Bg, navDrwTh.Fg, navDrwTh.ContrastBg, navDrwTh.ContrastFg =
		a.Theme.ContrastBg, a.Theme.ContrastFg, a.Theme.Bg, a.Theme.Fg
	a.NavDrawer = NewNav("Salty UI", "Decentralized Chat App", a, &navDrwTh)
	a.ModalNavDrawer = ModalNavFrom(a.NavDrawer, a.ModalLayer)
	a.NavAnim = component.VisibilityAnimation{
		State:    component.Invisible,
		Duration: time.Millisecond * 250,
	}
	settingsPage := NewSettingsPage(a, a.Theme)
	newChatPage := NewNewChatPage(a, a.Theme)
	settingsIcon, _ := widget.NewIcon(icons.ActionSettings)
	newChatIcon, _ := widget.NewIcon(icons.CommunicationChat)
	settingsNavItem := &NavItem{
		page:         settingsPage,
		NavDrawer:    a.NavDrawer,
		Name:         "Settings",
		Icon:         settingsIcon,
		Children:     make([]*NavItem, 0),
		Theme:        a.Theme,
		AlphaPalette: &a.AlphaPalette,
	}
	newChatNavItem := &NavItem{
		page:         newChatPage,
		NavDrawer:    a.NavDrawer,
		Name:         "New Chat",
		Icon:         newChatIcon,
		Children:     make([]*NavItem, 0),
		Theme:        a.Theme,
		AlphaPalette: &a.AlphaPalette,
	}
	a.AddNavItem(settingsNavItem)
	a.AddNavItem(newChatNavItem)
	a.history = make([]Page, 1)
	a.currentPage = settingsPage
	a.history[0] = a.currentPage
	a.Theme = material.NewTheme(gofont.Collection())
	go func() {
		if !a.UseNonModalDrawer() {
			a.NavAnim.Appear(time.Now())
		}
	}()
}

func (a *AppManager) Layout(gtx Gtx) Dim {
	if a.Theme == nil {
		a.Theme = material.NewTheme(gofont.Collection())
	}
	th := a.Theme
	//if a.ModalNavDrawer.NavDestinationChanged() {
	//	a.PushPage(gtx, a.ModalNavDrawer.CurrentNavDestination().(Page))
	//}
	paint.Fill(gtx.Ops, th.Palette.Bg)
	content := layout.Flexed(1, func(gtx Gtx) Dim {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx Gtx) Dim {
				th := *th
				th.Bg = th.ContrastBg
				th.ContrastBg = color.NRGBA(colornames.White)
				th.Fg = color.NRGBA(colornames.White)
				gtx.Constraints.Max.X /= 3
				if gtx.Constraints.Max.X > 350 {
					gtx.Constraints.Max.X = 350
				}
				return a.NavDrawer.Layout(gtx, &a.NavAnim)
			}),
			layout.Flexed(1, func(gtx Gtx) Dim {
				bar := layout.Rigid(func(gtx Gtx) Dim {
					return a.currentPage.DrawAppBar(gtx)
				})
				flex := layout.Flex{Axis: layout.Vertical}
				var currentView View
				if !a.Service.Loaded() {
					currentView = &Loader{Theme: a.Theme}
				} else {
					currentView = a.currentPage
				}
				if a.BottomBar {
					flex.Spacing = layout.SpaceBetween
					return flex.Layout(gtx, layout.Rigid(currentView.Layout), bar)
				} else {
					return flex.Layout(gtx, bar, layout.Rigid(currentView.Layout))
				}
			}),
		)
	})
	layout.Flex{}.Layout(gtx, content)
	if !a.UseNonModalDrawer() {
		a.ModalLayer.Layout(gtx, th)
	}
	return Dim{Size: gtx.Constraints.Max}
}
