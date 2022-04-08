package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/mearaj/saltyui/alog"
	"go.mills.io/saltyim"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"time"
)

// Settings Always call NewSettingsPage function to create Settings page
type Settings struct {
	widget.List
	*AppManager
	Theme                  *material.Theme
	title                  string
	iconCreateNewID        *widget.Icon
	inputNewID             component.TextField
	inputNewIDStr          string
	buttonNewID            widget.Clickable
	buttonRegistration     widget.Clickable
	buttonNavigation       widget.Clickable
	navigationIcon         *widget.Icon
	iDDetailsAccordion     Accordion
	errorNewIDAccordion    Accordion
	errorRegisterAccordion Accordion
	iDDetailsView          IDDetailsView
	errorCreateNewID       error
	errorRegister          error
	errorParseAddr         error
	registerLoading        bool
}

// NewSettingsPage Always call this function to create Settings page
func NewSettingsPage(manager *AppManager, th *material.Theme) *Settings {
	navIcon, _ := widget.NewIcon(icons.NavigationMenu)
	iconCreateNewID, _ := widget.NewIcon(icons.ContentCreate)
	if th == nil {
		th = material.NewTheme(gofont.Collection())
	}
	errorTh := *th
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	return &Settings{
		AppManager:      manager,
		Theme:           th,
		title:           "Settings",
		navigationIcon:  navIcon,
		iconCreateNewID: iconCreateNewID,
		iDDetailsView: IDDetailsView{
			Theme:      th,
			AppManager: manager,
		},
		iDDetailsAccordion: Accordion{
			Theme: th,
			Title: "View Details",
			Animation: component.VisibilityAnimation{
				State:    component.Visible,
				Duration: time.Millisecond * 250,
			},
		},
		errorNewIDAccordion: Accordion{
			Theme: &errorTh,
			Title: "Create New ID Error",
			Animation: component.VisibilityAnimation{
				State:    component.Visible,
				Duration: time.Millisecond * 250,
			},
		},
		errorRegisterAccordion: Accordion{
			Theme: &errorTh,
			Title: "Register Error",
			Animation: component.VisibilityAnimation{
				State:    component.Visible,
				Duration: time.Millisecond * 250,
			},
		},
	}
}

func (s *Settings) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (s *Settings) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (s *Settings) Layout(gtx Gtx) Dim {
	if s.Theme == nil {
		s.Theme = material.NewTheme(gofont.Collection())
	}
	if s.inputNewID.Text() != s.inputNewIDStr {
		s.errorRegister = nil
		s.errorCreateNewID = nil
		s.errorParseAddr = nil
	}
	_, s.errorParseAddr = saltyim.ParseAddr(s.inputNewID.Text())
	s.inputNewIDStr = s.inputNewID.Text()
	th := s.Theme
	return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1.0, func(gtx Gtx) Dim {
				s.List.Axis = layout.Vertical
				return material.List(th, &s.List).Layout(gtx, 1, func(gtx Gtx, _ int) Dim {
					return layout.Flex{
						Alignment: layout.Middle,
						Axis:      layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx Gtx) Dim {
							return s.drawNewIDTextField(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return layout.Spacer{Height: unit.Dp(32)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return s.drawIDDetailsAccordion(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return s.drawErrorNewIDAccordion(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return layout.Spacer{Height: unit.Dp(32)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return s.drawRegistrationButton(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return layout.Spacer{Height: unit.Dp(32)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return s.drawErrorRegisterAccordion(gtx)
						}),
					)
				})
			}),
		)
	})
}

func (s *Settings) DrawAppBar(gtx Gtx) Dim {
	gtx.Constraints.Max.Y = gtx.Px(unit.Dp(56))
	th := s.Theme
	if s.buttonNavigation.Clicked() {
		if s.AppManager.UseNonModalDrawer() {
			s.NavAnim.ToggleVisibility(time.Now())
		} else {
			s.AppManager.ModalNavDrawer.Appear(gtx.Now)
			s.NavAnim.Disappear(gtx.Now)
		}
	}
	component.Rect{Size: gtx.Constraints.Max, Color: th.Palette.ContrastBg}.Layout(gtx)
	layout.Flex{
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			if s.navigationIcon == nil {
				return Dim{}
			}
			navigationIcon := s.navigationIcon
			button := material.IconButton(th, &s.buttonNavigation, navigationIcon, "Nav Icon Button")
			button.Size = unit.Dp(24)
			button.Background = th.Palette.ContrastBg
			button.Color = th.Palette.ContrastFg
			button.Inset = layout.UniformInset(unit.Dp(16))
			return button.Layout(gtx)
		}),
		layout.Rigid(func(gtx Gtx) Dim {
			return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx Gtx) Dim {
				titleText := s.title
				title := material.Body1(th, titleText)
				title.Color = th.Palette.ContrastFg
				title.TextSize = unit.Dp(18)
				return title.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx Gtx) Dim {
			return Dim{
				Size:     image.Point{X: 0},
				Baseline: 0,
			}
		}),
	)
	return Dim{Size: gtx.Constraints.Max}
}

func (s *Settings) drawNewIDTextField(gtx Gtx) Dim {
	labelText := "Enter New ID"
	labelHintText := "User in the form user@domain"
	buttonText := "Create New ID"
	var button *widget.Clickable
	var th *material.Theme
	if s.errorParseAddr != nil {
		button = &widget.Clickable{}
		th = material.NewTheme(gofont.Collection())
		th.ContrastBg = color.NRGBA(colornames.Grey500)
	} else {
		button = &s.buttonNewID
		th = s.Theme
	}
	ib := IconButton{
		Theme:  th,
		Button: button,
		Icon:   s.iconCreateNewID,
		Text:   buttonText,
	}
	if button.Clicked() {
		s.errorCreateNewID = s.Service.CreateIdentity(s.inputNewID.Text())
		if s.errorCreateNewID != nil {
			s.errorNewIDAccordion.Animation.Appear(gtx.Now)
		}
	}
	return drawFormFieldRowWithLabel(gtx, th, labelText, labelHintText, &s.inputNewID, &ib)
}

func (s *Settings) drawIDDetailsAccordion(gtx Gtx) (d Dim) {
	if s.Service.CurrentIdentity() != nil {
		if s.iDDetailsAccordion.Child == nil {
			s.iDDetailsAccordion.Child = &s.iDDetailsView
		}
		return s.iDDetailsAccordion.Layout(gtx)
	} else {
		s.iDDetailsAccordion.Child = nil
	}
	return d
}

func (s *Settings) drawErrorNewIDAccordion(gtx Gtx) (d Dim) {
	if s.Service.CurrentIdentity() == nil && s.errorCreateNewID != nil {
		errView := ErrorView{}
		s.errorNewIDAccordion.Child = &errView
		errView.Error = s.errorCreateNewID.Error()
		return s.errorNewIDAccordion.Layout(gtx)
	}
	return d
}

func (s *Settings) drawRegistrationButton(gtx Gtx) Dim {
	buttonText := "Register with salty@domain for above id"
	var button *widget.Clickable
	var th *material.Theme
	if s.errorParseAddr != nil {
		button = &widget.Clickable{}
		th = material.NewTheme(gofont.Collection())
		th.ContrastBg = color.NRGBA(colornames.Grey500)
	} else {
		button = &s.buttonRegistration
		th = s.Theme
	}
	ib := IconButton{
		Theme:  th,
		Button: button,
		Icon:   s.iconCreateNewID,
		Text:   buttonText,
	}
	if button.Clicked() && !s.registerLoading {
		s.registerLoading = true
		go func() {
			if s.Service.CurrentIdentity() == nil {
				s.errorCreateNewID = s.Service.CreateIdentity(s.inputNewID.Text())
				if s.errorCreateNewID == nil {
					s.errorRegister = s.Service.Register()
					if s.errorRegister != nil {
						alog.Println(s.errorRegister)
					}
				}
			} else {
				s.errorRegister = s.Service.Register()
				if s.errorRegister != nil {
					alog.Println(s.errorRegister)
				}
			}
			s.registerLoading = false
			s.Window.Invalidate()
		}()
	}
	return ib.Layout(gtx)
}

func (s *Settings) drawErrorRegisterAccordion(gtx Gtx) (d Dim) {
	if s.errorRegister != nil {
		errView := ErrorView{}
		s.errorRegisterAccordion.Child = &errView
		errView.Error = s.errorRegister.Error()
		return s.errorRegisterAccordion.Layout(gtx)
	}
	return d
}
