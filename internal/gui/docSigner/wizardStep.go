package docSigning

import "fyne.io/fyne/v2"

type wizardStep interface {
	activate() *fyne.Container
	complete() *stepError
	revert()
}

type stepError struct {
	msg string
	err error
}
