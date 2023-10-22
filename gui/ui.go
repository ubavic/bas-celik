package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"github.com/ebfe/scard"
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/widgets"
)

var statusBar *widgets.StatusBar
var window *fyne.Window
var verbose bool
var startPageOn bool // possible data races
var startPage *widgets.StartPage

func StartGui(ctx *scard.Context, verbose_ bool) {
	startPageOn = true
	verbose = verbose_

	app := app.New()
	win := app.NewWindow("Baš Čelik")
	window = &win
	app.Settings().SetTheme(MyTheme{})

	statusBar = widgets.NewStatusBar()

	startPage = widgets.NewStartPage()
	startPage.SetStatus("", "", false)

	go pooler(ctx)

	win.SetContent(container.New(layout.NewPaddedLayout(), startPage))
	win.ShowAndRun()
}

func setUI(doc document.Document) {
	pdfHandler := savePdf(window, doc)
	ui := doc.BuildUI(pdfHandler, statusBar)
	columns := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), ui, layout.NewSpacer())
	container := container.New(layout.NewPaddedLayout(), columns)
	(*window).SetContent(container)

	(*window).Resize(container.MinSize())
	startPageOn = false
}

func setStartPage(status, explanation string, err error) {
	isError := false
	if err != nil {
		isError = true
	}

	if verbose && isError {
		fmt.Println(err)
	}

	startPage.SetStatus(status, explanation, isError)
	startPage.Refresh()

	if !startPageOn {
		(*window).SetContent(container.New(layout.NewPaddedLayout(), startPage))
		startPageOn = true
	}
}

func setStatus(status string, err error) {
	isError := false
	if err != nil {
		isError = true
	}

	if verbose && isError {
		fmt.Println(err)
	}

	statusBar.SetStatus(status, isError)
	statusBar.Refresh()
}

func savePdf(win *fyne.Window, doc document.Document) func() {
	return func() {
		pdf, fileName, err := doc.BuildPdf()

		if err != nil {
			setStatus("Greška pri generisanju PDF-a", fmt.Errorf("generating PDF: %w", err))
			return
		}

		dialog := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
			if w == nil || err != nil {
				return
			}

			_, err = w.Write(pdf)
			if err != nil {
				setStatus("Greška pri zapisivanju PDF-a", fmt.Errorf("writing PDF: %w", err))
				return
			}

			err = w.Close()
			if err != nil {
				setStatus("Greška pri zapisivanju PDF-a", fmt.Errorf("writing PDF: %w", err))
				return
			}

			setStatus("PDF sačuvan", nil)
		}, *win)

		dialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
		dialog.SetFileName(fileName)

		dialog.Show()
	}
}
