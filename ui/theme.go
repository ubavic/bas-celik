package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct{}

func (MyTheme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0xD9, G: 0xDF, B: 0xFE, A: 0xFF}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0xE5, G: 0xE5, B: 0xE5, A: 0xFF}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x42}
	case theme.ColorNameError:
		return color.NRGBA{R: 0xF0, G: 0x47, B: 0x3B, A: 0xFF}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0xD9, G: 0xDF, B: 0xFE, A: 0xFF}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xFF}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x0F}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x19}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x19}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0x4B, G: 0x5E, B: 0xC3, A: 0xFF}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x99}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x0F}
	default:
		return theme.LightTheme().Color(c, v)
	}
}

func (MyTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.LightTheme().Font(s)
}

func (MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.LightTheme().Icon(n)
}

func (MyTheme) Size(s fyne.ThemeSizeName) float32 {
	return theme.LightTheme().Size(s)
}
