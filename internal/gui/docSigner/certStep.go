package docSigning

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type certStep struct {
	certSelect *widget.Select
}

func (c *certStep) activate() *fyne.Container {
	certNames := make([]string, len(gWizard.certs))
	for i, cert := range gWizard.certs {
		name := cert.Certificate.Subject.CommonName
		if name == "" {
			name = cert.Certificate.Subject.String()
		}
		certNames[i] = name
	}

	if c.certSelect == nil {
		c.certSelect = widget.NewSelect(certNames, nil)
		c.certSelect.PlaceHolder = gWizard.t("xmlSign.step.selectCert")
		c.certSelect.OnChanged = func(_ string) {
			gWizard.selectedCertIdx = c.certSelect.SelectedIndex()
			gWizard.nextBtn.Enable()
		}
	} else {
		c.certSelect.Options = certNames
		c.certSelect.Enable()
		c.certSelect.Refresh()
	}

	if len(certNames) > 0 {
		c.certSelect.SetSelectedIndex(0)
	}

	return container.New(layout.NewVBoxLayout(), c.certSelect)
}

func (c *certStep) complete() *stepError {
	if gWizard.selectedCertIdx < 0 {
		return &stepError{"", nil}
	}
	c.certSelect.Disable()
	return nil
}

func (c *certStep) revert() {
	gWizard.selectedCertIdx = -1
}
