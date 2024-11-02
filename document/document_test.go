package document_test

import (
	"os"
	"testing"

	"github.com/ubavic/bas-celik/document"
)

// Should be used only for testing
func setDocumentConfigFromLocalFiles(t *testing.T) {
	documentConfig := document.DocumentConfig{}
	var err error

	documentConfig.FontRegular, err = os.ReadFile("../embed/liberationSansRegular.ttf")
	if err != nil {
		t.Fatal(err.Error())
	}

	documentConfig.FontBold, err = os.ReadFile("../embed/liberationSansBold.ttf")
	if err != nil {
		t.Fatal(err.Error())
	}

	documentConfig.RfzoLogo, err = os.ReadFile("../embed/rfzo.png")
	if err != nil {
		t.Fatal(err.Error())
	}

	document.Configure(documentConfig)
}

func unsetDocumentConfig() {
	document.Configure(document.DocumentConfig{})
}
