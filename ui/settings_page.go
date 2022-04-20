package ui

import (
	"bytes"
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
	"io"
	"time"
)

// SettingsPage Always call NewSettingsPage function to create SettingsPage page
type SettingsPage struct {
	widget.List
	*AppManager
	Theme              *material.Theme
	title              string
	iconCreateNewID    *widget.Icon
	iconImportFile     *widget.Icon
	inputNewID         component.TextField
	inputImportFile    component.TextField
	inputNewIDStr      string
	inputImportFileStr string
	buttonNewID        widget.Clickable
	buttonRegistration widget.Clickable
	buttonNavigation   widget.Clickable
	buttonImport       widget.Clickable
	navigationIcon     *widget.Icon
	iDDetailsAccordion Accordion
	iDConfigAccordion  Accordion
	iDDetailsView      IDDetailsView
	errorCreateNewID   error
	errorRegister      error
	errorImportFile    error
	registerLoading    bool
	creatingNewID      bool
	importingFile      bool
	idLoadedFromDB     bool
}

// NewSettingsPage Always call this function to create SettingsPage page
func NewSettingsPage(manager *AppManager, th *material.Theme) *SettingsPage {
	navIcon, _ := widget.NewIcon(icons.NavigationMenu)
	iconCreateNewID, _ := widget.NewIcon(icons.ContentCreate)
	iconImportFile, _ := widget.NewIcon(icons.FileFileUpload)
	if th == nil {
		th = material.NewTheme(gofont.Collection())
	}
	errorTh := *th
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	return &SettingsPage{
		AppManager:         manager,
		Theme:              th,
		title:              "Settings",
		navigationIcon:     navIcon,
		iconCreateNewID:    iconCreateNewID,
		iconImportFile:     iconImportFile,
		inputImportFileStr: "Import ID key file by clicking button",
		iDDetailsView: IDDetailsView{
			Theme:      th,
			AppManager: manager,
		},
		iDDetailsAccordion: Accordion{
			Theme: th,
			Title: "View Current ID Details",
			Animation: component.VisibilityAnimation{
				State:    component.Visible,
				Duration: time.Millisecond * 250,
			},
		},
		iDConfigAccordion: Accordion{
			Theme: th,
			Title: "View Current ID Config",
			Animation: component.VisibilityAnimation{
				State:    component.Visible,
				Duration: time.Millisecond * 250,
			},
		},
	}
}

func (s *SettingsPage) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (s *SettingsPage) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (s *SettingsPage) Layout(gtx Gtx) Dim {
	if s.Theme == nil {
		s.Theme = material.NewTheme(gofont.Collection())
	}
	th := s.Theme
	if !s.idLoadedFromDB {
		s.idLoadedFromDB = true
		id := s.Service.Identity()
		if id != nil {
			s.inputNewIDStr = id.Addr().String()
			s.inputNewID.SetText(s.inputNewIDStr)
		}
	}
	s.inputImportFile.SetText(s.inputImportFileStr)
	if s.inputNewID.Text() != s.inputNewIDStr {
		s.errorRegister = nil
		s.errorCreateNewID = nil
	}
	_, s.errorCreateNewID = saltyim.ParseAddr(s.inputNewID.Text())
	s.inputNewIDStr = s.inputNewID.Text()
	return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx Gtx) Dim {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1.0, func(gtx Gtx) Dim {
				s.List.Axis = layout.Vertical
				return material.List(th, &s.List).Layout(gtx, 1, func(gtx Gtx, _ int) Dim {
					return layout.Flex{
						Alignment: layout.Middle,
						Axis:      layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx Gtx) Dim {
							return s.drawImportFileField(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return layout.Spacer{Height: unit.Dp(32)}.Layout(gtx)
						}),
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
							return layout.Spacer{Height: unit.Dp(32)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return s.drawRegistrationButton(gtx)
						}),
						layout.Rigid(func(gtx Gtx) Dim {
							return layout.Spacer{Height: unit.Dp(32)}.Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}

func (s *SettingsPage) DrawAppBar(gtx Gtx) Dim {
	gtx.Constraints.Max.Y = gtx.Px(unit.Dp(56))
	th := s.Theme
	if s.buttonNavigation.Clicked() {
		s.NavDrawer.ToggleVisibility(time.Now())
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

func (s *SettingsPage) drawImportFileField(gtx Gtx) Dim {
	labelText := "Import ID"
	labelHintText := "Import ID From File"
	buttonText := "Import key file"
	textTh := material.NewTheme(gofont.Collection())
	textTh.Fg.A = 200
	textTh.ContrastBg = textTh.Fg
	if s.errorImportFile != nil {
		s.inputImportFile.SetError(s.errorImportFile.Error())
	} else {
		s.inputImportFile.ClearError()
	}
	ib := IconButton{
		Theme:  s.Theme,
		Button: &s.buttonImport,
		Icon:   s.iconImportFile,
		Text:   buttonText,
	}

	if s.buttonImport.Clicked() && !s.importingFile {
		s.importingFile = true
		go func() {
			var cl io.ReadCloser
			cl, s.errorImportFile = s.Explorer.ChooseFile(".key")
			defer func() {
				if cl != nil {
					_ = cl.Close()
				}
			}()
			if s.errorImportFile == nil {
				buff := bytes.Buffer{}
				_, s.errorImportFile = buff.ReadFrom(cl)
				if s.errorImportFile == nil {
					s.errorImportFile = <-s.Service.CreateIDFromBytes(buff.Bytes())
					if s.errorImportFile == nil {
						if s.Service.Identity() != nil {
							addr := s.Service.Identity().Addr().String()
							s.inputNewID.SetText(addr)
							s.inputImportFile.SetText(addr)
							s.inputNewIDStr = addr
							s.inputImportFileStr = addr
						}
					}
				}
			}
			s.importingFile = false
		}()
	}
	return drawFormFieldRowWithLabel(gtx, textTh, labelText, labelHintText, &s.inputImportFile, &ib)
}

func (s *SettingsPage) drawNewIDTextField(gtx Gtx) Dim {
	labelText := "Enter New ID"
	labelHintText := "User in the form user@domain"
	buttonText := "Create New ID"
	var button *widget.Clickable
	var th *material.Theme
	if s.errorCreateNewID != nil && s.inputNewIDStr != "" {
		button = &widget.Clickable{}
		th = material.NewTheme(gofont.Collection())
		th.ContrastBg = color.NRGBA(colornames.Grey500)
		s.inputNewID.SetError(s.errorCreateNewID.Error())
	} else {
		button = &s.buttonNewID
		th = s.Theme
		s.inputNewID.ClearError()
	}

	ib := IconButton{
		Theme:  th,
		Button: button,
		Icon:   s.iconCreateNewID,
		Text:   buttonText,
	}
	if button.Clicked() && !s.creatingNewID {
		s.creatingNewID = true
		go func() {
			s.errorCreateNewID = <-s.Service.CreateIDFromAddrStr(s.inputNewID.Text())
			s.creatingNewID = false
		}()
	}
	return drawFormFieldRowWithLabel(gtx, th, labelText, labelHintText, &s.inputNewID, &ib)
}

func (s *SettingsPage) drawIDDetailsAccordion(gtx Gtx) (d Dim) {
	identity := s.Service.Identity()
	if identity != nil {
		if s.iDDetailsAccordion.Child == nil {
			s.iDDetailsAccordion.Child = s.iDDetailsView.Layout
		}
		return s.iDDetailsAccordion.Layout(gtx)
	} else {
		s.iDDetailsAccordion.Child = nil
	}
	return d
}

func (s *SettingsPage) drawRegistrationButton(gtx Gtx) Dim {
	buttonText := "Register with salty@domain for above id"
	var button *widget.Clickable
	var th *material.Theme
	if s.errorRegister != nil {
		button = &widget.Clickable{}
		th = material.NewTheme(gofont.Collection())
		th.ContrastBg = color.NRGBA(colornames.Red500)
		th.Fg = th.ContrastBg
	} else {
		button = &s.buttonRegistration
		th = s.Theme
	}
	ib := IconButton{
		Theme:  s.Theme,
		Button: button,
		Icon:   s.iconCreateNewID,
		Text:   buttonText,
	}
	if button.Clicked() && !s.registerLoading {
		s.registerLoading = true
		go func() {
			s.errorRegister = <-s.Service.Register(s.inputNewID.Text())
			s.registerLoading = false
			if s.errorRegister != nil {
				alog.Logger().Println(s.errorRegister)
			}
			s.Window.Invalidate()
		}()
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			return ib.Layout(gtx)
		}),
		layout.Rigid(func(gtx Gtx) Dim {
			if s.errorRegister != nil {
				return layout.Inset{Bottom: unit.Dp(8.0)}.Layout(gtx,
					func(gtx Gtx) Dim {
						return material.Body1(th, s.errorRegister.Error()).Layout(gtx)
					})
			}
			return Dim{}
		}),
	)
}
