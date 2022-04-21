package ui

import (
	"encoding/json"
	"fmt"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/mearaj/saltyui/alog"
	"image/color"
)

type IDDetailsView struct {
	*material.Theme
	copyButton   widget.Clickable
	configButton widget.Clickable
	*AppManager
	Contents        string
	exportingConfig bool
}

func (i *IDDetailsView) Layout(gtx Gtx) (d Dim) {
	if i.Theme == nil {
		i.Theme = material.NewTheme(gofont.Collection())
	}
	identity := i.AppManager.Service.Identity()
	if identity != nil {
		var contents = string(identity.Contents())
		if i.copyButton.Clicked() {
			i.AppManager.Window.WriteClipboard(contents)
		}
		if i.configButton.Clicked() && !i.exportingConfig {
			i.exportingConfig = true
			go func() {
				cfg, err := i.Service.ConfigJSON()
				if err != nil {
					alog.Logger().Errorln(err)
				}
				if cfg != nil {
					cfgPrint := map[string]interface{}{
						"config": map[string]interface{}{
							"key":      cfg.Config.Key,
							"endpoint": cfg.Config.Endpoint,
						},
						"hash": cfg.Hash,
					}
					cfgMar, _ := json.MarshalIndent(cfgPrint, "", "  ")
					fmt.Println(string(cfgMar))
				}
				i.exportingConfig = false
			}()
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx Gtx) Dim {
				return material.Body1(i.Theme, contents).Layout(gtx)
			}),
			layout.Rigid(func(gtx Gtx) Dim {
				maxWidth := gtx.Constraints.Max.X
				return layout.Flex{Spacing: layout.SpaceSides}.Layout(gtx,
					layout.Rigid(func(gtx2 Gtx) Dim {
						gtx.Constraints.Max.X = (maxWidth / 3) - 32
						return material.Button(i.Theme, &i.copyButton, "Copy to Clipboard").Layout(gtx)
					}),
					layout.Rigid(func(gtx2 Gtx) Dim {
						return layout.Spacer{Width: unit.Dp(32)}.Layout(gtx)
					}),
					layout.Rigid(func(gtx2 Gtx) Dim {
						gtx.Constraints.Max.X = (maxWidth / 3) - 32
						button := &i.configButton
						if i.exportingConfig {
							button = &widget.Clickable{}
						}
						return material.Button(i.Theme, button, "Export config").Layout(gtx)
					}),
					layout.Rigid(func(gtx2 Gtx) Dim {
						return layout.Spacer{Width: unit.Dp(32)}.Layout(gtx)
					}),
				)
			}),
		)
	}
	return d
}

type IDListItem struct {
	*material.Theme
	widget   widget.Clickable
	Selected bool
	*AppManager
}

func (i *IDListItem) Layout(gtx Gtx, index int) Dim {
	if i.widget.Clicked() {
		i.Selected = !i.Selected
	}
	iconUnChecked := i.Theme.Icon.CheckBoxUnchecked
	iconChecked := i.Theme.Icon.CheckBoxChecked
	userIds := i.Service.Identities()
	if index >= len(userIds) {
		return Dim{}
	}
	identity := i.Service.Identities()[index]
	if identity == nil {
		return Dim{}
	}

	ins := layout.UniformInset(unit.Dp(8))
	flx := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}
	btnStyle := material.ButtonLayoutStyle{
		Background:   color.NRGBA{},
		CornerRadius: unit.Dp(8),
		Button:       &i.widget,
	}
	if i.Selected {
		btnStyle.Background = color.NRGBA{A: 0x44, R: 0x88, G: 0x88, B: 0x88}
	} else {
		btnStyle.Background = color.NRGBA{}
	}

	return btnStyle.Layout(gtx, func(gtx Gtx) Dim {
		return ins.Layout(gtx, func(gtx Gtx) Dim {
			return flx.Layout(gtx,
				layout.Flexed(1.0, func(gtx Gtx) Dim {
					body := material.Body1(i.Theme, string(identity.Contents()))
					//gtx.Constraints.Max.X = gtx.Constraints.Max.X - gtx.Px(unit.Dp(64))
					d := body.Layout(gtx)
					return d
				}),
				layout.Rigid(func(gtx Gtx) Dim {
					return layout.Inset{Left: unit.Dp(16.0), Right: unit.Dp(8.0)}.Layout(gtx,
						func(gtx Gtx) Dim {
							return layout.Flex{Alignment: layout.Middle,
								Spacing: layout.SpaceSides}.Layout(gtx,
								layout.Rigid(func(gtx Gtx) Dim {
									gtx.Constraints.Max.X = gtx.Px(unit.Dp(32))
									gtx.Constraints.Min.X = gtx.Px(unit.Dp(32))
									gtx.Constraints.Max.Y = gtx.Px(unit.Dp(32))
									gtx.Constraints.Min.Y = gtx.Px(unit.Dp(32))
									if i.Selected {
										return iconChecked.Layout(gtx, i.Theme.ContrastBg)
									}
									return iconUnChecked.Layout(gtx, i.Theme.Fg)
								}),
							)

						})
				}),
			)
		})
	})
}
