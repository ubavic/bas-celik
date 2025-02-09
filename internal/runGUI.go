//go:build !cli

package internal

import (
	"github.com/ubavic/bas-celik/v2/internal/gui"
	"github.com/ubavic/bas-celik/v2/internal/gui/icon"
	"github.com/ubavic/bas-celik/v2/internal/gui/translation"
	"github.com/ubavic/bas-celik/v2/internal/logger"
)

func Run(cfg LaunchConfig) error {
	if len(cfg.PdfPath) == 0 && len(cfg.JsonPath) == 0 && len(cfg.ExcelPath) == 0 {
		err := translation.SetTranslations(cfg.EmbedDirectory)
		if err != nil {
			return err
		}

		err = icon.LoadIcons(cfg.EmbedDirectory)
		if err != nil {
			return err
		}

		gui.StartGui(version)
		return nil
	}

	logger.Info("output file detected")
	return readAndSave(cfg)
}
