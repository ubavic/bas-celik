package gui

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/card"
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/internal/gui/celiktheme"
	"github.com/ubavic/bas-celik/internal/gui/reader"
	"github.com/ubavic/bas-celik/internal/gui/translation"
	"github.com/ubavic/bas-celik/internal/gui/widgets"
)

type State struct {
	mu            sync.Mutex
	verbose       bool
	window        *fyne.Window
	startPage     *widgets.StartPage
	toolbar       *widgets.Toolbar
	spacer        *widgets.Spacer
	mainPage      *fyne.Container
	mainContainer *fyne.Container
	statusBar     *widgets.StatusBar
	cardDocument  card.CardDocument
}

var state State

func StartGui(verbose_ bool, version string) {
	app := app.New()
	win := app.NewWindow("Baš Čelik")

	theme := celiktheme.NewTheme(app.Preferences().IntWithFallback(themePreferenceKey, 1))
	app.Settings().SetTheme(theme)

	translation.SetLanguage(app.Preferences().IntWithFallback(languagePreferenceKey, 1))

	showAboutBox := showAboutBox(win, version)
	showSettings := showSetupBox(win, app)
	changePin := pinChange(win)

	widgets.SetClipboard(CopyToClipboard)

	readerSelection := make(chan string)

	statusBar := widgets.NewStatusBar()
	toolbar := widgets.NewToolbar(showAboutBox, showSettings, changePin)
	spacer := widgets.NewSpacer()

	reader.NewPoller(toolbar, readerSelection)

	startPage := widgets.NewStartPage()
	startPage.SetStatus("", "", false)

	mainPage := container.New(layout.NewVBoxLayout())
	rows := container.New(layout.NewVBoxLayout(), toolbar, spacer, startPage, mainPage)
	columns := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), rows, layout.NewSpacer())
	mainContainer := container.New(layout.NewPaddedLayout(), columns)

	state = State{
		verbose:       verbose_,
		toolbar:       toolbar,
		startPage:     startPage,
		window:        &win,
		spacer:        spacer,
		statusBar:     statusBar,
		mainPage:      mainPage,
		mainContainer: mainContainer,
	}

	win.SetContent(mainContainer)

	go cardLoop(readerSelection)

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

	(*state.window).Resize(state.mainContainer.MinSize())
}

func setStartPage(status, explanation string, err error) {
	state.mu.Lock()
	defer state.mu.Unlock()

	isError := false
	if err != nil {
		isError = true
	}

	if state.verbose && isError {
		fmt.Println(err)
	}

	state.startPage.SetStatus(status, explanation, isError)
	state.startPage.Refresh()

	state.mainPage.RemoveAll()

	state.mainPage.Hide()
	state.startPage.Show()

	(*state.window).Resize(state.mainContainer.MinSize())
}

func setStatus(status string, err error) {
	isError := false
	if err != nil {
		isError = true
	}

	if state.verbose && isError {
		fmt.Println(err)
	}

	state.statusBar.SetStatus(status, isError)
	state.statusBar.Refresh()
}

func updateMedicalDocHandler(doc *document.MedicalDocument) func() {
	return func() {
		err := doc.UpdateValidUntilDateFromRfzo()
		if err != nil {
			dialog.ShowInformation(t("error.error"), t("error.dataUpdate"), *state.window)
			return
		}

		setStatus(t("ui.updateSuccessful"), nil)
		setUI(doc)
	}
}

func CopyToClipboard(str string) bool {
	if state.window == nil {
		return false
	}

	win := *state.window
	clipboard := win.Clipboard()
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

func t(id string) string {
	return translation.Translate(id)
}
