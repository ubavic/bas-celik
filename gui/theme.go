package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct{}

func (MyTheme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if v == theme.VariantLight || v == 2 {
		switch c {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
		case theme.ColorNameButton:
			return color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF}
		case theme.ColorNameDisabledButton:
			return color.NRGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF}
		case theme.ColorNameDisabled:
			return color.NRGBA{R: 0x30, G: 0x30, B: 0x30, A: 0xFF}
		case theme.ColorNameError:
			return color.NRGBA{R: 0xF0, G: 0x47, B: 0x3B, A: 0xFF}
		case theme.ColorNameFocus:
			return color.NRGBA{R: 0xDE, G: 0xEB, B: 0xFA, A: 0xFF}
		case theme.ColorNameForeground:
			return color.NRGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xFF}
		case theme.ColorNameHeaderBackground:
			return color.NRGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xFF}
		case theme.ColorNameHover:
			return color.NRGBA{R: 0x00, G: 0x00, B: 0x40, A: 0x10}
		case theme.ColorNameHyperlink:
			return color.NRGBA{R: 0x50, G: 0x50, B: 0xA0, A: 0xFF}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF}
		case theme.ColorNameInputBorder:
			return color.NRGBA{R: 0xDA, G: 0xDA, B: 0xDA, A: 0xFF}
		case theme.ColorNameMenuBackground:
			return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 0xF9, G: 0xF9, B: 0xF9, A: 0xFF}
		case theme.ColorNamePlaceHolder:
			return color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
		case theme.ColorNamePressed:
			return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 0x5A, G: 0x73, B: 0x8F, A: 0xFF}
		case theme.ColorNameScrollBar:
			return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x99}
		case theme.ColorNameSelection:
			return color.NRGBA{R: 0xDE, G: 0xEB, B: 0xFA, A: 0xFF}
		case theme.ColorNameShadow:
			return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x10}
		default:
			return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x10}
		}
	} else {
		switch c {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 0x10, G: 0x10, B: 0x13, A: 0xFF}
		case theme.ColorNameButton:
			return color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF}
		case theme.ColorNameDisabledButton:
			return color.NRGBA{R: 0x12, G: 0x12, B: 0x12, A: 0xFF}
		case theme.ColorNameDisabled:
			return color.NRGBA{R: 0x40, G: 0x40, B: 0x40, A: 0xFF}
		case theme.ColorNameError:
			return color.NRGBA{R: 0xF0, G: 0x47, B: 0x3B, A: 0xFF}
		case theme.ColorNameFocus:
			return color.NRGBA{R: 0x23, G: 0x20, B: 0x24, A: 0xFF}
		case theme.ColorNameForeground:
			return color.NRGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF}
		case theme.ColorNameHeaderBackground:
			return color.NRGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xFF}
		case theme.ColorNameHover:
			return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x10}
		case theme.ColorNameHyperlink:
			return color.NRGBA{R: 0x80, G: 0x90, B: 0xF0, A: 0xFF}
		case theme.ColorNameInputBackground:
			return color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF}
		case theme.ColorNameInputBorder:
			return color.NRGBA{R: 0xDA, G: 0xDA, B: 0xDA, A: 0x00}
		case theme.ColorNameMenuBackground:
			return color.NRGBA{R: 0x15, G: 0x15, B: 0x15, A: 0xFF}
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 0x15, G: 0x15, B: 0x17, A: 0xFF}
		case theme.ColorNamePlaceHolder:
			return color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xFF}
		case theme.ColorNamePressed:
			return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 0x41, G: 0x4D, B: 0x7A, A: 0xFF}
		case theme.ColorNameScrollBar:
			return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x50}
		case theme.ColorNameSelection:
			return color.NRGBA{R: 0x23, G: 0x20, B: 0x24, A: 0xFF}
		case theme.ColorNameShadow:
			return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x40}
		default:
			return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x10}
		}
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
