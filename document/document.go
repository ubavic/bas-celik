// Package provides type definitions fot different types of documents
// present on smart cards, as well methods for exporting those documents to different file formats.
// In order for `BuildPDF` methods to work properly,
// this package must be (pre)configured with `Configure` function.
package document

// Represents any document handled by Bas Celik
type Document interface {
	BuildPdf() ([]byte, string, error)   // Renders document to pdf
	BuildJson() ([]byte, error)          // Renders document to json
	BuildExcel() ([]byte, string, error) // Renders document to xlsx
}

var (
	fontRegular []byte
	fontBold    []byte
	rfzoLogo    []byte
)

type DocumentConfig struct {
	FontRegular []byte // regular font used for PDF render
	FontBold    []byte // bold font used for PDF render
	RfzoLogo    []byte // logo used in PDF render of medical cards
}

// Sets fonts and graphics used for rendering PDF
func Configure(config DocumentConfig) error {
	fontRegular = config.FontRegular
	fontBold = config.FontBold
	rfzoLogo = config.RfzoLogo

	return nil
}
