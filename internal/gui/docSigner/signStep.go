package docSigning

import (
	"io"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/ubavic/bas-celik/v2/internal/xmlsign"
)

type signStep struct {
	signBtn *widget.Button
	wrapper *fyne.Container
}

func (s *signStep) activate() *fyne.Container {
	if s.signBtn == nil {
		s.signBtn = widget.NewButton(gWizard.t("xmlSign.signFile"), s.openFile)
		s.signBtn.Importance = widget.HighImportance
	} else {
		s.signBtn.Enable()
	}

	if s.wrapper == nil {
		s.wrapper = container.New(layout.NewVBoxLayout(), s.signBtn)
	}

	return s.wrapper
}

func (s *signStep) openFile() {
	s.signBtn.Disable()

	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			s.signBtn.Enable()
			return
		}
		defer reader.Close()

		data, readErr := io.ReadAll(reader)
		if readErr != nil {
			if gWizard.onStatusMsg != nil {
				gWizard.onStatusMsg(gWizard.t("xmlSign.signError") + ": " + readErr.Error())
			}
			s.signBtn.Enable()
			return
		}

		s.signAsync(data, filepath.Base(reader.URI().Path()))
	}, gWizard.window)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".xml"}))
	fd.Show()
}

func (s *signStep) signAsync(xmlBytes []byte, xmlFileName string) {
	certIdx := gWizard.selectedCertIdx
	cert := gWizard.certs[certIdx]
	signer := &pkcs11Signer{
		session: gWizard.session,
		certId:  cert.Id,
		cert:    cert.Certificate,
	}

	go func() {
		result, err := xmlsign.SignXMLBytes(xmlBytes, signer, cert.Certificate)
		fyne.Do(func() {
			if err != nil {
				if gWizard.onStatusMsg != nil {
					gWizard.onStatusMsg(gWizard.t("xmlSign.signError") + ": " + err.Error())
				}
				s.signBtn.Enable()
				return
			}
			s.saveFile(result, xmlFileName)
		})
	}()
}

func (s *signStep) saveFile(signedXML, xmlFileName string) {
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		defer s.signBtn.Enable()

		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		if _, writeErr := writer.Write([]byte(signedXML)); writeErr != nil {
			if gWizard.onStatusMsg != nil {
				gWizard.onStatusMsg(gWizard.t("xmlSign.signError") + ": " + writeErr.Error())
			}
			return
		}

		if gWizard.onStatusMsg != nil {
			gWizard.onStatusMsg(gWizard.t("xmlSign.xmlSaved"))
		}
	}, gWizard.window)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".xml"}))
	fd.SetFileName("signed_" + xmlFileName)
	fd.Show()
}

func (s *signStep) complete() *stepError { return nil }
func (s *signStep) revert()              {}
