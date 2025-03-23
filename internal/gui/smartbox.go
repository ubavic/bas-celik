package gui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/ubavic/bas-celik/v2/internal/smartbox/pkcs11"
	"github.com/ubavic/bas-celik/v2/internal/smartbox/server"
)

func startSmartboxUI() {
	preferences := state.app.Preferences()

	modulePaths := []server.ModulePath{}

	vendors := []pkcs11.CardVendor{
		pkcs11.CardVendorMup,
		pkcs11.CardVendorPks,
		pkcs11.CardVendorPosta,
		pkcs11.CardVendorHalcom,
		pkcs11.CardVendorEsmart,
	}

	for _, vendor := range vendors {
		preferenceKey := getVendorPreferenceKey(vendor)

		path := preferences.String(preferenceKey)

		if path != "" {
			modulePaths = append(modulePaths, server.ModulePath{
				Vendor: vendor,
				Path:   path,
			})
		}
	}

	rows := container.New(layout.NewVBoxLayout(), state.toolbar, state.startPage, state.statusBar)
	state.documentUi = rows

	state.mainContainer.Add(state.documentUi)

	address, err := server.StartServer(modulePaths)

	if err != nil {
		setStartPage("smartbox.error", "", err)
		state.statusBar.SetStatus(err.Error(), true)
		state.statusBar.Refresh()
	} else {
		setStartPage("smartbox.running", "", nil)
		state.statusBar.SetStatus(address, false)
		state.statusBar.Refresh()
	}

	resizeWindow(true)

	state.window.ShowAndRun()
}

func getVendorPreferenceKey(vendor pkcs11.CardVendor) string {

	switch vendor {
	case pkcs11.CardVendorHalcom:
		return halcomPkcsPathKey
	case pkcs11.CardVendorPosta:
		return postaPkcsPathKey
	case pkcs11.CardVendorEsmart:
		return esmartPkcsPathKey
	case pkcs11.CardVendorMup:
		return mupPkcsPathKey
	case pkcs11.CardVendorPks:
		return pksPkcsPathKey
	default:
		panic("Invalid card vendor")
	}
}
