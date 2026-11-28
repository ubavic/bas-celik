package gui

import (
	"fmt"
	"runtime"
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

	rows := container.New(layout.NewVBoxLayout(), state.toolbar, spacer, state.startPage, state.unknownCardPage, state.documentUiMainContainer, state.cryptoUiContainer)
	columns := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), rows, layout.NewSpacer())

	state.documentUi = columns

	state.mainContainer.Add(state.documentUi)

	state.cryptoUiContainer.Hide()

	if pollerErr == nil {
		poller.StartPoller()
	} else {
		explanation := ""
		if runtime.GOOS == "linux" {
			explanation = t("error.contextFail.linux")
		}
		setStartPage("error.contextFail", explanation, pollerErr)
	}

	state.window.ShowAndRun()
}

func setUI(doc document.Document) {
	fyne.Do(func() {
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

		state.documentUiMainContainer.RemoveAll()
		state.documentUiMainContainer.Add(page)
		state.documentUiMainContainer.Add(buttonBar)

		state.startPage.Hide()
		state.unknownCardPage.Hide()
		state.documentUiMainContainer.Show()

		resizeWindowNow(false)
	})
}

// safe to call from any goroutine.
func setStartPage(statusId, explanation string, err error) {
	status := t(statusId)

	isError := err != nil

	if isError {
		logger.Error(err)
	} else {
		logger.Info(translation.EnglishTranslation(statusId))
	}

	fyne.Do(func() {
		state.startPage.SetStatus(status, explanation, isError)
		state.startPage.Refresh()

		state.documentUiMainContainer.RemoveAll()

		state.documentUiMainContainer.Hide()
		state.unknownCardPage.Hide()
		state.startPage.Show()

		resizeWindowNow(true)
	})
}

// Safe to call from any goroutine.
func setUnknownCardPage(atr string, err error) {
	if err != nil {
		logger.Error(err)
	}

	status := t("error.unknownCard")
	copyLabel := t("error.unknownCard.copyAtr")

	fyne.Do(func() {
		state.unknownCardPage.SetContent(status, atr, copyLabel)
		state.unknownCardPage.Refresh()

		state.documentUiMainContainer.RemoveAll()

		state.documentUiMainContainer.Hide()
		state.startPage.Hide()
		state.unknownCardPage.Show()

		resizeWindowNow(true)
	})
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
		go func() {
			err := doc.UpdateValidUntilDateFromRfzo()
			if err != nil {
				logger.Error(fmt.Errorf("updating medical information: %w", err))
				fyne.Do(func() {
					dialog.ShowInformation(t("error.error"), t("error.dataUpdate"), state.window)
				})
				return
			}

			setStatus("ui.updateSuccessful", nil)
			setUI(doc)
		}()
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

// Safe to call from any goroutine.
func setTimedStatus(label string) {
	showTimedStatus(label)
}

// Safe to call from any goroutine.
func setTimedStatusError(labelKey string, err error) {
	if err == nil {
		return
	}

	logger.Error(err)

	showTimedStatus(t(labelKey))
}

func showTimedStatus(label string) {
	fyne.Do(func() {
		state.statusBar.SetStatus(label, false)
		state.statusBar.Refresh()
	})

	time.AfterFunc(2*time.Second, func() {
		// Compare and clear on the Fyne goroutine so the check and the update
		// are atomic with respect to other status changes.
		fyne.Do(func() {
			if state.statusBar.GetStatus() == label {
				state.statusBar.SetStatus("", false)
				state.statusBar.Refresh()
			}
		})
	})
}

// showDocumentUI must be called from the Fyne goroutine
func showDocumentUI() {
	state.mainContainer.RemoveAll()
	state.mainContainer.Add(state.documentUi)
}

// Safe to call from any goroutine.
func resizeWindow(keepCurrentSize bool) {
	fyne.Do(func() {
		resizeWindowNow(keepCurrentSize)
	})
}

// resizeWindowNow must be called from the Fyne goroutine.
func resizeWindowNow(keepCurrentSize bool) {
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

	state.window.Resize(minSize)
}
