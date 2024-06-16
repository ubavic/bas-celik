//go:build !cli

package internal

import (
	"github.com/ubavic/bas-celik/gui"
)

func Run(pdfPath, jsonPath string, verbose, getValidUntilFromRfzo bool, reader uint) error {
	if len(pdfPath) == 0 && len(jsonPath) == 0 {
		gui.StartGui(verbose, version)
		return nil
	} else {
		return readAndSave(pdfPath, jsonPath, reader, getValidUntilFromRfzo)
	}
}
