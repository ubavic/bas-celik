package gui

import "github.com/ubavic/bas-celik/v2/internal/gui/translation"

func t(id string) string {
	return translation.Translate(id)
}
