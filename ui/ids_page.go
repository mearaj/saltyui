package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
)

// IDPage Always call NewIDPage function to create IDPage page
type IDPage struct {
	layout.List
	*AppManager
	Theme            *material.Theme
	title            string
	iconNewChat      *widget.Icon
	inputNewChat     component.TextField
	inputNewChatStr  string
	buttonNewChat    widget.Clickable
	buttonNavigation widget.Clickable
	navigationIcon   *widget.Icon
	identitiesViews  []*IDListItem
}

// NewIDPage Always call this function to create IDPage page
func NewIDPage(manager *AppManager, th *material.Theme) *IDPage {
	navIcon, _ := widget.NewIcon(icons.NavigationArrowBack)
	iconNewChat, _ := widget.NewIcon(icons.ContentCreate)
	if th == nil {
		th = material.NewTheme(gofont.Collection())
	}
	errorTh := *th
	errorTh.ContrastBg = color.NRGBA(colornames.Red500)
	return &IDPage{
		AppManager:      manager,
		Theme:           th,
		title:           "Identities",
		navigationIcon:  navIcon,
		iconNewChat:     iconNewChat,
		List:            layout.List{Axis: layout.Vertical},
		identitiesViews: []*IDListItem{},
	}
}

func (ids *IDPage) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (ids *IDPage) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (ids *IDPage) Layout(gtx Gtx) Dim {
	if ids.Theme == nil {
		ids.Theme = material.NewTheme(gofont.Collection())
	}
	inset := layout.UniformInset(unit.Dp(16))
	d := inset.Layout(gtx, func(gtx Gtx) Dim {
		flex := layout.Flex{Axis: layout.Vertical}
		return flex.Layout(gtx,
			layout.Rigid(func(gtx Gtx) Dim {
				return Dim{}
			}),
			layout.Rigid(func(gtx Gtx) Dim {
				if len(ids.Service.Identities()) == 0 {
					return inset.Layout(gtx, func(gtx Gtx) Dim {
						return ids.drawNoIdentitiesCreated(gtx)
					})
				}
				return inset.Layout(gtx, ids.drawIdentitiesItems)
			}))
	})
	return d
}

func (ids *IDPage) DrawAppBar(gtx Gtx) Dim {
	gtx.Constraints.Max.Y = gtx.Px(unit.Dp(56))
	th := ids.Theme
	if ids.buttonNavigation.Clicked() {
		ids.AppManager.PopUp()
	}
	component.Rect{Size: gtx.Constraints.Max, Color: th.Palette.ContrastBg}.Layout(gtx)
	layout.Flex{
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx Gtx) Dim {
			if ids.navigationIcon == nil {
				return Dim{}
			}
			navigationIcon := ids.navigationIcon
			button := material.IconButton(th, &ids.buttonNavigation, navigationIcon, "Nav Icon Button")
			button.Size = unit.Dp(24)
			button.Background = th.Palette.ContrastBg
			button.Color = th.Palette.ContrastFg
			button.Inset = layout.UniformInset(unit.Dp(16))
			return button.Layout(gtx)
		}),
		layout.Rigid(func(gtx Gtx) Dim {
			return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx Gtx) Dim {
				titleText := ids.title
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

func (ids *IDPage) drawIdentitiesItems(gtx Gtx) Dim {
	userIDs := ids.Service.Identities()
	return ids.List.Layout(gtx, len(userIDs), func(gtx Gtx, index int) (d Dim) {
		if len(ids.identitiesViews) < index+1 {
			ids.identitiesViews = append(ids.identitiesViews, &IDListItem{
				Theme:      ids.Theme,
				AppManager: ids.AppManager,
			})
		}
		inset := layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}
		return inset.Layout(gtx, func(gtx2 Gtx) Dim {
			return ids.identitiesViews[index].Layout(gtx, index)
		})
	})
}

func (ids *IDPage) drawNoIdentitiesCreated(gtx Gtx) Dim {
	flex := layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides}
	return flex.Layout(gtx, layout.Rigid(func(gtx Gtx) Dim {
		return layout.Center.Layout(gtx, func(gtx Gtx) Dim {
			return material.Body1(ids.Theme, "You haven't created any identities").Layout(gtx)
		})
	}))
}
