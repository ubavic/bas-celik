//go:build !cli

package internal

import (
	"github.com/ubavic/bas-celik/internal/gui"
)

func Run(cfg LaunchConfig) error {
	if len(cfg.PdfPath) == 0 && len(cfg.JsonPath) == 0 {
		gui.StartGui(cfg.Verbose, version)
		return nil
	}

	return readAndSave(cfg)
}
