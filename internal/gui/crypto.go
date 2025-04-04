package gui

import (
	"encoding/pem"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/card"
	"github.com/ubavic/bas-celik/v2/internal/gui/reader"
	"github.com/ubavic/bas-celik/v2/internal/gui/widgets"
	"github.com/ubavic/bas-celik/v2/internal/logger"
)

func cryptoList() {
	if state.cryptoUiContainer.Visible() {
		return
	}

	state.mu.Lock()

	createCryptoUI()

	if state.cryptoUi != nil {
		state.startPage.Hide()
		state.documentUiMainContainer.Hide()
		state.cryptoUiContainer.Add(state.cryptoUi)
		state.cryptoUiContainer.Show()
	}

	state.mu.Unlock()
}

func createCryptoUI() {
	state.certs = nil
	state.selectedCert = -1

	gemaltoCard, ok := state.cardDocument.(*card.Gemalto)
	if !ok {
		state.cryptoUi = nil
		setStatus("crypto.wrongCard", fmt.Errorf("card could not be casted to Gemalto card"))
		return
	}

	reader.CancelReaderPoler()

	err := gemaltoCard.LoadCertificates()
	if err != nil {
		logger.Error(err)
	}

	state.certs = gemaltoCard.GetCertificates()
	state.selectedCert = 0

	logger.Info(fmt.Sprintf("loaded %d certificates", len(state.certs)))

	reader.RestartReaderPoler()

	canvasObjects := []fyne.CanvasObject{}

	if len(state.certs) == 0 {
		spacer := `                                                     `
		spacer += spacer + spacer
		noCertFound := widget.NewLabel(spacer + "\n\n\n\n\n" + t("crypto.noCertFound") + "\n\n\n\n\n\n\n")
		noCertFound.Alignment = fyne.TextAlignCenter
		labelContainer := container.New(layout.NewCenterLayout(), noCertFound)

		canvasObjects = append(canvasObjects, labelContainer)
	} else {
		// TODO: display multiple certificates

		fieldSerNo := widgets.NewField(t("crypto.serialNumber"), "", 280)
		fieldSig := widgets.NewField(t("crypto.type"), "", 290)
		generalRow1 := container.New(layout.NewHBoxLayout(), fieldSerNo, fieldSig)
		fieldNotBefore := widgets.NewField(t("crypto.notBefore"), "", 280)
		fieldNotAfter := widgets.NewField(t("crypto.notAfter"), "", 290)
		generalRow2 := container.New(layout.NewHBoxLayout(), fieldNotBefore, fieldNotAfter)
		generalGroup := widgets.NewGroup(t("crypto.general"), generalRow1, generalRow2)

		issuerSnField := widgets.NewField("SN", "", 280)
		issuerCnField := widgets.NewField("CN", "", 290)
		issuerRow1 := container.New(layout.NewHBoxLayout(), issuerSnField, issuerCnField)
		issuerOUField := widgets.NewField("OU", "", 280)
		issuerOField := widgets.NewField("O", "", 290)
		issuerRow2 := container.New(layout.NewHBoxLayout(), issuerOUField, issuerOField)
		issuerCField := widgets.NewField("C", "", 280)
		issuerLField := widgets.NewField("L", "", 290)
		issuerRow3 := container.New(layout.NewHBoxLayout(), issuerCField, issuerLField)
		issuerGroup := widgets.NewGroup(t("crypto.issuer"), issuerRow1, issuerRow2, issuerRow3)

		subjectSnField := widgets.NewField("SN", "", 280)
		subjectCnField := widgets.NewField("CN", "", 290)
		subjectRow1 := container.New(layout.NewHBoxLayout(), subjectSnField, subjectCnField)
		subjectOUField := widgets.NewField("OU", "", 280)
		subjectOField := widgets.NewField("O", "", 290)
		subjectRow2 := container.New(layout.NewHBoxLayout(), subjectOUField, subjectOField)
		subjectCField := widgets.NewField("C", "", 280)
		subjectLField := widgets.NewField("L", "", 290)
		subjectRow3 := container.New(layout.NewHBoxLayout(), subjectCField, subjectLField)
		subjectGroup := widgets.NewGroup(t("crypto.subject"), subjectRow1, subjectRow2, subjectRow3)

		selectCert := func(index int) {
			state.selectedCert = index

			if state.selectedCert >= len(state.certs) || state.selectedCert < 0 {
				return
			}

			cert := state.certs[state.selectedCert]
			fieldSerNo.SetValue(cert.SerialNumber.String())
			fieldSig.SetValue(cert.SignatureAlgorithm.String())
			fieldNotAfter.SetValue(cert.NotAfter.Format("02.01.2006"))
			fieldNotBefore.SetValue(cert.NotBefore.Format("02.01.2006"))

			issuerSnField.SetValue(cert.Issuer.SerialNumber)
			issuerCnField.SetValue(cert.Issuer.CommonName)
			issuerOUField.SetValue(strings.Join(cert.Issuer.OrganizationalUnit, ","))
			issuerOField.SetValue(strings.Join(cert.Issuer.Organization, ","))
			issuerLField.SetValue(strings.Join(cert.Issuer.Locality, ","))
			issuerCField.SetValue(strings.Join(cert.Issuer.Country, ","))

			subjectSnField.SetValue(cert.Subject.SerialNumber)
			subjectCnField.SetValue(cert.Subject.CommonName)
			subjectOUField.SetValue(strings.Join(cert.Subject.OrganizationalUnit, ","))
			subjectOField.SetValue(strings.Join(cert.Subject.Organization, ","))
			subjectLField.SetValue(strings.Join(cert.Subject.Locality, ","))
			subjectCField.SetValue(strings.Join(cert.Subject.Country, ","))

			if len(cert.Subject.OrganizationalUnit)+len(cert.Subject.Organization) == 0 {
				subjectOField.Hide()
				subjectOUField.Hide()
			} else {
				subjectOField.Show()
				subjectOField.Show()
			}
		}

		canvasObjects = append(canvasObjects, generalGroup, issuerGroup, subjectGroup)

		selectCert(0)
	}

	buttons := []fyne.CanvasObject{}
	exitButton := widget.NewButtonWithIcon(t("crypto.return"), theme.NavigateBackIcon(), closeCryptoUi)
	changePinButton := widget.NewButton(t("crypto.changePin"), pinChange())
	buttons = append(buttons, exitButton, layout.NewSpacer(), changePinButton)

	if len(state.certs) > 0 {
		saveCertButton := widget.NewButton(t("crypto.saveCert"), saveCert)
		buttons = append(buttons, saveCertButton)
	}

	buttonBar := container.New(layout.NewHBoxLayout(), buttons...)

	canvasObjects = append(canvasObjects, layout.NewSpacer(), buttonBar)

	state.cryptoUi = container.New(layout.NewVBoxLayout(), canvasObjects...)
}

func saveCert() {
	dialog := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
		if err != nil {
			setStatus("error.writingCert", fmt.Errorf("writing certificate: %w", err))
			return
		}

		if w == nil {
			return
		}

		if state.selectedCert >= len(state.certs) || state.selectedCert < 0 {
			return
		}

		cert := state.certs[state.selectedCert]
		pemBlock := pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		}

		err = pem.Encode(w, &pemBlock)
		if err != nil {
			setStatus("error.writingCert", fmt.Errorf("encoding certificate: %w", err))
			return
		}

		err = w.Close()
		if err != nil {
			setStatus("error.writingCert", fmt.Errorf("writing certificate: %w", err))
			return
		}

		setStatus("ui.certSaved", nil)
	}, state.window)

	dialog.SetFilter(storage.NewExtensionFileFilter([]string{".pem", ".crt", ".cer"}))

	lastUsedDirectoryURI := getLastUsedDirectory()
	if lastUsedDirectoryURI != nil {
		dialog.SetLocation(lastUsedDirectoryURI)
	}

	dialog.Show()
}

func closeCryptoUi() {
	state.cryptoUiContainer.Hide()
	state.documentUiMainContainer.Show()
	state.cryptoUiContainer.RemoveAll()
}
