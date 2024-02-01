package document

import (
	"os"
	"testing"
)

// Should be used only for testing
func SetDataFromLocalFiles(t *testing.T) {
	var err1, err2, err3 error
	fontRegular, err1 = os.ReadFile("../embed/liberationSansRegular.ttf")
	fontBold, err2 = os.ReadFile("../embed/liberationSansBold.ttf")
	rfzoLogo, err3 = os.ReadFile("../embed/rfzo.png")
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("Cant open file")
	}
}

func UnsetData(t *testing.T) {
	fontRegular = nil
	fontBold = nil
	rfzoLogo = nil
}
