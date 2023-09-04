package gui

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
		return color.NRGBA{R: 0xC2, G: 0xEB, B: 0xF5, A: 0xFF}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0xE5, G: 0xE5, B: 0xE5, A: 0xFF}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x42}
	case theme.ColorNameError:
		return color.NRGBA{R: 0xF0, G: 0x47, B: 0x3B, A: 0xFF}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0xC5, G: 0xEB, B: 0xF4, A: 0xFF}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xFF}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x82, G: 0xA0, B: 0xE3, A: 0x30}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x19}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0xF0, G: 0x00, B: 0x00, A: 0x19}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0x6B, G: 0xB2, B: 0xC3, A: 0xFF}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x99}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0xDF, G: 0xF9, B: 0xFE, A: 0xFF}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x0F}
	default:
		return theme.LightTheme().Color(c, v)
	}
}

func (MyTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

func (MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (MyTheme) Size(s fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(s)
}

func (MyTheme) CornerRadius() float32 {
	return 3
}
