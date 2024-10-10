//go:build cli

package internal

func Run(cfg LaunchConfig) error {
	if len(cfg.PdfPath) == 0 && len(cfg.JsonPath) == 0 {
		cfg.JsonPath = "out.json"
	}

	return readAndSave(cfg)
}
