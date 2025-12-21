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

func startCardReaderUI() {
	widgets.SetClipboard(copyToClipboard)

	spacer := widgets.NewSpacer()

	poller, pollerErr := reader.NewPoller(state.toolbar, connectToCard)

	rows := container.New(layout.NewVBoxLayout(), state.toolbar, spacer, state.startPage, state.documentUiMainContainer, state.cryptoUiContainer)
	columns := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), rows, layout.NewSpacer())

	state.documentUi = columns

	state.mainContainer.Add(state.documentUi)

	state.cryptoUiContainer.Hide()

	if pollerErr == nil {
		poller.StartPoller()
	} else {
		setStartPage("error.contextFail", "", pollerErr)
	}

	state.window.ShowAndRun()
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

	fyne.DoAndWait(func() {
		state.documentUiMainContainer.RemoveAll()
		state.documentUiMainContainer.Add(page)
		state.documentUiMainContainer.Add(buttonBar)

		state.startPage.Hide()
		state.documentUiMainContainer.Show()
	})

	resizeWindow(false)
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

	state.documentUiMainContainer.RemoveAll()

	state.documentUiMainContainer.Hide()
	state.startPage.Show()

	resizeWindow(true)
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
	fyne.Do(func() {
		state.statusBar.SetStatus(status, isError)
		state.statusBar.Refresh()
	})
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
	setTimedStatus(label)

	return true
}

func setTimedStatus(label string) {
	state.statusBar.SetStatus(label, false)
	state.statusBar.Refresh()
	go func() {
		time.Sleep(2 * time.Second)
		if state.statusBar.GetStatus() == label {
			fyne.Do(func() {
				state.statusBar.SetStatus("", false)
				state.statusBar.Refresh()
			})
		}
	}()
}

func setTimedStatusError(labelKey string, err error) {
	if err == nil {
		return
	}

	logger.Error(err)

	label := t(labelKey)

	state.statusBar.SetStatus(label, false)
	state.statusBar.Refresh()
	go func() {
		time.Sleep(2 * time.Second)
		if state.statusBar.GetStatus() == label {
			fyne.Do(func() {
				state.statusBar.SetStatus("", false)
				state.statusBar.Refresh()
			})
		}
	}()
}

func showDocumentUI() {
	fyne.DoAndWait(func() {
		state.mainContainer.RemoveAll()
		state.mainContainer.Add(state.documentUi)
	})
}

func resizeWindow(keepCurrentSize bool) {
	minSize := state.documentUi.MinSize()

	if keepCurrentSize {
		currentSize := state.mainContainer.Size()

		if currentSize.Height > minSize.Height {
			minSize.Height = currentSize.Height
		}

		if currentSize.Width > minSize.Width {
			minSize.Width = currentSize.Width
		}
	}

	fyne.Do(func() {
		state.window.Resize(minSize)
	})
}
