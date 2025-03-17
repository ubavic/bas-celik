package gui

import "github.com/ubavic/bas-celik/v2/internal/gui/translation"

func t(id string, vals ...any) string {
	return translation.Translate(id, vals...)
}
