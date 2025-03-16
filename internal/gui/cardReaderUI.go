package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/document"
	"github.com/ubavic/bas-celik/v2/internal/gui/reader"
	"github.com/ubavic/bas-celik/v2/internal/gui/translation"
	"github.com/ubavic/bas-celik/v2/internal/gui/widgets"
	"github.com/ubavic/bas-celik/v2/internal/logger"
)

func startCardReaderUI(app fyne.App, win fyne.Window) {
	showAboutBox := showAboutBox(win, version)
	showSettings := showSetupBox(win, app)
	changePin := pinChange(win)

	widgets.SetClipboard(copyToClipboard)

	statusBar := widgets.NewStatusBar()
	toolbar := widgets.NewToolbar(showAboutBox, showSettings, changePin)
	spacer := widgets.NewSpacer()

	poller, pollerErr := reader.NewPoller(toolbar, connectToCard)

	startPage := widgets.NewStartPage()
	startPage.SetStatus("", "", false)

	mainPage := container.New(layout.NewVBoxLayout())
	rows := container.New(layout.NewVBoxLayout(), toolbar, spacer, startPage, mainPage)
	columns := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), rows, layout.NewSpacer())
	mainContainer := container.New(layout.NewPaddedLayout(), columns)

	state = State{
		app:           app,
		window:        win,
		toolbar:       toolbar,
		startPage:     startPage,
		spacer:        spacer,
		statusBar:     statusBar,
		mainPage:      mainPage,
		mainContainer: mainContainer,
	}

	win.SetContent(mainContainer)

	if pollerErr == nil {
		poller.StartPoller()
	} else {
		setStartPage("error.contextFail", "", pollerErr)
	}

	win.ShowAndRun()
}

func setUI(doc document.Document) {
	state.mu.Lock()
	defer state.mu.Unlock()

	var page *fyne.Container
	buttonBarObjects := []fyne.CanvasObject{state.statusBar, layout.NewSpacer()}

	switch doc := doc.(type) {
	case *document.IdDocument:
		page = pageID(doc)
	case *document.MedicalDocument:
		updateButton := widget.NewButton(t("ui.update"), updateMedicalDocHandler(doc))
		buttonBarObjects = append(buttonBarObjects, updateButton)
		page = pageMedical(doc)
	case *document.VehicleDocument:
		page = pageVehicle(doc)
	}

	savePdfButton := widget.NewButton(t("ui.savePdf"), savePdf(doc))
	saveXlsxButton := widget.NewButton(t("ui.saveXlsx"), saveXlsx(doc))
	buttonBarObjects = append(buttonBarObjects, saveXlsxButton, savePdfButton)

	buttonBar := container.New(layout.NewHBoxLayout(), buttonBarObjects...)

	state.mainPage.RemoveAll()
	state.mainPage.Add(page)
	state.mainPage.Add(buttonBar)

	state.startPage.Hide()
	state.mainPage.Show()

	state.window.Resize(state.mainContainer.MinSize())
}

func setStartPage(statusId, explanationId string, err error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	status := t(statusId)
	explanation := t(explanationId)

	isError := false
	if err != nil {
		isError = true
	}

	if isError {
		logger.Error(err)
	} else {
		logger.Info(translation.EnglishTranslation(statusId) + " " + translation.EnglishTranslation(explanationId))
	}

	state.startPage.SetStatus(status, explanation, isError)
	state.startPage.Refresh()

	state.mainPage.RemoveAll()

	state.mainPage.Hide()
	state.startPage.Show()

	state.window.Resize(state.mainContainer.MinSize())
}

func setStatus(statusId string, err error) {
	isError := false
	if err != nil {
		isError = true
	}

	if isError {
		logger.Error(err)
	} else {
		logger.Info(translation.EnglishTranslation(statusId))
	}

	status := t(statusId)
	state.statusBar.SetStatus(status, isError)
	state.statusBar.Refresh()
}

func updateMedicalDocHandler(doc *document.MedicalDocument) func() {
	return func() {
		err := doc.UpdateValidUntilDateFromRfzo()
		if err != nil {
			logger.Error(fmt.Errorf("updating medical information: %w", err))
			dialog.ShowInformation(t("error.error"), t("error.dataUpdate"), state.window)
			return
		}

		setStatus("ui.updateSuccessful", nil)
		setUI(doc)
	}
}

func copyToClipboard(str string) bool {
	if state.window == nil {
		return false
	}

	clipboard := state.window.Clipboard()
	if clipboard == nil {
		return false
	}

	label := t("ui.contentCopied")

	clipboard.SetContent(str)
	state.statusBar.SetStatus(label, false)
	state.statusBar.Refresh()
	go func() {
		time.Sleep(2 * time.Second)
		if state.statusBar.GetStatus() == label {
			state.statusBar.SetStatus("", false)
			state.statusBar.Refresh()
		}
	}()

	return true
}
