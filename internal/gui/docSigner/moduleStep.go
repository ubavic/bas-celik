package docSigning

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
)

type moduleStep struct {
	moduleSelect *widget.Select
}

func (m *moduleStep) activate() *fyne.Container {
	vendorNames := make([]string, len(gWizard.loadedVendors))
	for i, v := range gWizard.loadedVendors {
		vendorNames[i] = v.String()
	}

	if m.moduleSelect == nil {
		sel := widget.NewSelect(vendorNames, nil)
		sel.PlaceHolder = gWizard.t("xmlSign.step.selectModule")
		sel.OnChanged = func(_ string) {
			gWizard.selectedVendorIdx = sel.SelectedIndex()
			gWizard.nextBtn.Enable()
		}
		m.moduleSelect = sel
	}
	m.moduleSelect.DisableableWidget.Enable()

	if len(gWizard.loadedVendors) > 0 {
		m.moduleSelect.SetSelectedIndex(0)
	}

	return container.New(layout.NewVBoxLayout(), m.moduleSelect)
}

func (m *moduleStep) complete() *stepError {
	if gWizard.selectedVendorIdx < 0 {
		return &stepError{"", nil}
	}

	vendor := gWizard.loadedVendors[gWizard.selectedVendorIdx]

	session, err := pkcs11.GetPkcsSession(vendor)
	if err != nil {
		return &stepError{"xmlSign.signError", fmt.Errorf("getting PKCS#11 session: %w", err)}
	}

	gWizard.session = &session
	m.moduleSelect.DisableableWidget.Disable()

	return nil
}

func (m *moduleStep) revert() {
	gWizard.selectedVendorIdx = -1
	m.activate()
}
