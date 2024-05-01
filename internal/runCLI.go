//go:build cli

package internal

func Run(pdfPath, jsonPath string, verbose bool, reader uint) error {
	if len(pdfPath) == 0 && len(jsonPath) == 0 {
		jsonPath = "out.json"
	}

	return readAndSave(pdfPath, jsonPath, reader)
}
