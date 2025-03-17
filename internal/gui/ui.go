package gui

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/ubavic/bas-celik/v2/card"
	"github.com/ubavic/bas-celik/v2/internal/gui/celiktheme"
	"github.com/ubavic/bas-celik/v2/internal/gui/translation"
	"github.com/ubavic/bas-celik/v2/internal/gui/widgets"
)

type State struct {
	mu            sync.Mutex
	app           fyne.App
	window        fyne.Window
	version       string
	startPage     *widgets.StartPage
	toolbar       *widgets.Toolbar
	spacer        *widgets.Spacer
	mainPage      *fyne.Container
	mainContainer *fyne.Container
	statusBar     *widgets.StatusBar
	cardDocument  card.CardDocument
}

var state State

func StartGui(version string) {
	app := app.New()
	win := app.NewWindow("Baš Čelik")

	theme := celiktheme.NewTheme(app.Preferences().IntWithFallback(themePreferenceKey, 1))
	app.Settings().SetTheme(theme)

	translation.SetLanguage(app.Preferences().IntWithFallback(languagePreferenceKey, 1))

	statusBar := widgets.NewStatusBar()

	mainPage := container.New(layout.NewVBoxLayout())

	startPage := widgets.NewStartPage()
	startPage.SetStatus("", "", false)

	mainContainer := container.New(layout.NewPaddedLayout())
	win.SetContent(mainContainer)

	state = State{
		app:           app,
		window:        win,
		version:       version,
		mainContainer: mainContainer,
		mainPage:      mainPage,
		startPage:     startPage,
		statusBar:     statusBar,
	}

	if app.Preferences().BoolWithFallback(smartboxModeKey, false) {
		startSmartboxUI()
	} else {
		startCardReaderUI()
	}
}
