package gui

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/internal/gui/widgets"
)

func setupCertificateView(t *testing.T) {
	t.Helper()
	a := test.NewApp()
	t.Cleanup(func() { a.Quit(); state = State{} })
	state = State{
		app: a, window: a.NewWindow("Certificates"),
		startPage: widgets.NewStartPage(), unknownCardPage: widgets.NewUnknownCardPage(),
		documentUiMainContainer: container.NewVBox(), cryptoUiContainer: container.NewVBox(),
		statusBar: widgets.NewStatusBar(),
		certs: []x509.Certificate{
			{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "First synthetic certificate"}},
			{SerialNumber: big.NewInt(2), Subject: pkix.Name{CommonName: "Second synthetic certificate"}},
		},
	}
	state.documentUi = container.NewVBox(state.startPage, state.unknownCardPage, state.documentUiMainContainer, state.cryptoUiContainer)
	state.mainContainer = container.NewVBox(state.documentUi)
	state.window.SetContent(state.mainContainer)
}

func certificateActionButtons() []*widget.Button {
	objects := state.cryptoUi.Objects
	bar := objects[len(objects)-1].(*fyne.Container)
	var buttons []*widget.Button
	for _, object := range bar.Objects {
		if b, ok := object.(*widget.Button); ok {
			buttons = append(buttons, b)
		}
	}
	return buttons
}

func TestCertificateOnlyView(t *testing.T) {
	setupCertificateView(t)
	setCertificateUI(state.certs, nil)
	if !state.cryptoUiContainer.Visible() || state.documentUiMainContainer.Visible() || state.startPage.Visible() || state.unknownCardPage.Visible() {
		t.Fatal("certificate-only view did not replace the document/unknown pages")
	}
	if len(certificateActionButtons()) != 1 {
		t.Fatal("certificate-only cards must offer export without PIN or document navigation")
	}
	if len(state.certsSelectorButtons) != 2 || state.selectedCert != 0 {
		t.Fatal("certificate selection was not initialized")
	}
	test.Tap(state.certsSelectorButtons[1])
	if state.selectedCert != 1 {
		t.Fatal("second certificate was not selected")
	}
	resetCardUI()
	if state.cryptoUiContainer.Visible() || len(state.cryptoUiContainer.Objects) != 0 || len(state.certs) != 0 || state.selectedCert != -1 || len(state.certsSelectorButtons) != 0 {
		t.Fatal("card removal left certificate data in the view")
	}
	// Insertion after removal must start with a fresh selection.
	setCertificateUI([]x509.Certificate{{SerialNumber: big.NewInt(3)}}, nil)
	if state.selectedCert != 0 || len(state.certsSelectorButtons) != 0 || len(state.certs) != 1 {
		t.Fatal("reinsertion retained the previous card's selection")
	}
}

func TestIdentityCertificateViewKeepsNavigation(t *testing.T) {
	setupCertificateView(t)
	showCryptoUI(true, nil)
	buttons := certificateActionButtons()
	if len(buttons) != 3 {
		t.Fatal("identity cards must retain back, PIN change and export")
	}
	test.Tap(buttons[0])
	if state.cryptoUiContainer.Visible() || !state.documentUiMainContainer.Visible() {
		t.Fatal("back did not restore the identity document")
	}
}

func TestPartialCertificateViewShowsWarning(t *testing.T) {
	setupCertificateView(t)
	setCertificateUI(state.certs[:1], errors.New("synthetic read failure"))
	if _, ok := state.cryptoUi.Objects[0].(*widget.Label); !ok {
		t.Fatal("partial read error is not visible")
	}
	if state.selectedCert != 0 || len(certificateActionButtons()) != 1 {
		t.Fatal("partial read must still allow export of the readable certificate")
	}
}
