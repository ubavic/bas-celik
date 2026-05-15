package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/ubavic/bas-celik/v2/internal/gui/widgets"
	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
	"github.com/ubavic/bas-celik/v2/internal/smartbox/server"
)

func startSmartboxUI() {
	preferences := state.app.Preferences()
	modulePaths := getValidModulePaths(preferences)

	rows := container.New(layout.NewVBoxLayout(), state.toolbar, state.startPage, state.documentUiMainContainer, state.statusBar)
	state.documentUi = rows

	state.mainContainer.Add(state.documentUi)

	port, loadedVendors, err := server.StartServer(modulePaths)
	if err != nil {
		setStartPage("smartbox.error", "", err)
		state.statusBar.SetStatus(err.Error(), true)
		state.statusBar.Refresh()
	} else {
		StartSmartboxRunningPage(port, loadedVendors)
	}

	resizeWindow(true)
	state.window.ShowAndRun()
}

func StartSmartboxRunningPage(port string, vendors []string) {
	smartboxPage := widgets.NewSmartboxPage(t("smartbox.running"), t("smartbox.port")+": "+port, vendors)

	fyne.Do(func() {
		state.startPage.Hide()
		state.documentUiMainContainer.RemoveAll()
		state.documentUiMainContainer.Add(smartboxPage)
		state.documentUiMainContainer.Show()
		state.documentUiMainContainer.Refresh()
	})
}

type PathStorage interface {
	String(key string) string
}

func getValidModulePaths(preferences PathStorage) []pkcs11.ModulePath {
	modulePaths := []pkcs11.ModulePath{}

	vendors := []pkcs11.CardVendor{
		pkcs11.CardVendorMup,
		pkcs11.CardVendorPks,
		pkcs11.CardVendorPosta,
		pkcs11.CardVendorHalcom,
		pkcs11.CardVendorEsmart,
	}

	for _, vendor := range vendors {
		preferenceKey := getVendorPreferenceKey(vendor)

		path := strings.TrimSpace(preferences.String(preferenceKey))

		if path != "" {
			modulePaths = append(modulePaths, pkcs11.ModulePath{
				Vendor: vendor,
				Path:   path,
			})
		}
	}

	return modulePaths
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
