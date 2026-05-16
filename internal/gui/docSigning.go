package gui

import (
	"errors"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	docSigning "github.com/ubavic/bas-celik/v2/internal/gui/docSigner"
	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
)

func startXmlSigningUI() {
	preferences := state.app.Preferences()
	modulePaths := getValidModulePaths(preferences)

	rows := container.New(layout.NewVBoxLayout(),
		state.toolbar,
		state.startPage,
		state.documentUiMainContainer,
		state.statusBar,
	)
	state.documentUi = rows
	state.mainContainer.Add(state.documentUi)

	loadedVendors, err := pkcs11.LoadModules(modulePaths)

	if len(loadedVendors) == 0 {
		if err == nil {
			err = errors.New("no module path is set")
		}
		setStartPage("xmlSign.noModules", t("xmlSign.noModulesExplanation"), err)

		state.window.SetOnClosed(func() {
			pkcs11.Deinit()
		})
		resizeWindow(true)
		state.window.ShowAndRun()
		return
	}

	w := docSigning.InitSignWizard(state.window, t, loadedVendors, setTimedStatus)
	state.window.SetOnClosed(w.Deinit)

	state.startPage.Hide()
	state.documentUiMainContainer.RemoveAll()
	state.documentUiMainContainer.Add(w.Panel)
	state.documentUiMainContainer.Show()

	resizeWindow(false)
	state.window.ShowAndRun()
}
