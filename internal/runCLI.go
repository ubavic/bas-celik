//go:build cli

package internal

import "github.com/ubavic/bas-celik/internal/logger"

func Run(cfg LaunchConfig) error {
	if len(cfg.PdfPath) == 0 && len(cfg.JsonPath) == 0 && len(cfg.ExcelPath) == 0 {
		logger.Info("no output file path detected, using default value")
		cfg.JsonPath = "out.json"
	}

	return readAndSave(cfg)
}
