package ui

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/saltyui/service"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/image/colornames"
	"image/color"
	"time"
)

// AppManager Always call NewAppManager function to create AppManager instance
type AppManager struct {
	*app.Window
	currentPage Page
	pushPage    Page
	popupPage   bool
	history     []Page
	*NavDrawer
	*ModalNavDrawer
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

func (a *AppManager) PushPage(page Page) {
	a.pushPage = page
}

func (a *AppManager) PopUp() {
	a.popupPage = true
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
	settingsNavItem := NewNavItem(settingsPage, a.NavDrawer, "Settings", settingsIcon, make([]*NavItem, 0), a.Theme, a.AlphaPalette)
	newChatNavItem := NewNavItem(newChatPage, a.NavDrawer, "New Chat", newChatIcon, make([]*NavItem, 0), a.Theme, a.AlphaPalette)
	a.AddNavItem(settingsNavItem)
	a.AddNavItem(newChatNavItem)
	a.history = make([]Page, 1)
	a.currentPage = settingsPage
	a.history[0] = a.currentPage
	a.Theme = material.NewTheme(gofont.Collection())
	go func() {
		// Wait for the window to be ready
		time.Sleep(time.Millisecond * 500)
		if a.UseNonModalDrawer() {
			a.NavAnim.Appear(time.Now())
			a.Window.Invalidate()
		}
	}()
}

func (a *AppManager) Layout(gtx Gtx) Dim {
	if a.Theme == nil {
		a.Theme = material.NewTheme(gofont.Collection())
	}
	th := a.Theme
	if a.pushPage != nil {
		if len(a.history) == 0 {
			a.history = make([]Page, 0, 1)
		}
		a.history = append(a.history, a.pushPage)
		a.currentPage = a.history[len(a.history)-1]
		a.pushPage = nil
		a.popupPage = false
	} else if a.popupPage && len(a.history) > 1 {
		i := len(a.history) - 1
		a.history = a.history[:i]
		a.currentPage = a.history[i-1]
		a.NavDrawer.SetNavDestination(a.currentPage)
		a.pushPage = nil
		a.popupPage = false
	}
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
				if 350 > a.WindowWidth {
					gtx.Constraints.Max.X = a.WindowWidth - 100
					gtx.Constraints.Min.X = a.WindowWidth - 100
				} else {
					gtx.Constraints.Min.X = 350
					if gtx.Constraints.Max.X < gtx.Constraints.Min.X {
						gtx.Constraints.Max.X = gtx.Constraints.Min.X
					}
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
