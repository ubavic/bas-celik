package gui

import (
	"crypto/x509"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/card"
	"github.com/ubavic/bas-celik/v2/internal/gui/celiktheme"
	"github.com/ubavic/bas-celik/v2/internal/gui/translation"
	"github.com/ubavic/bas-celik/v2/internal/gui/widgets"
)

type State struct {
	version                 string
	mu                      sync.Mutex
	app                     fyne.App
	window                  fyne.Window
	mainContainer           *fyne.Container
	documentUi              *fyne.Container
	cryptoUiContainer       *fyne.Container
	cryptoUi                *fyne.Container
	startPage               *widgets.StartPage
	documentUiMainContainer *fyne.Container
	toolbar                 *widgets.Toolbar
	statusBar               *widgets.StatusBar
	cardDocument            card.CardDocument
	selectedCert            int
	certs                   []x509.Certificate
	certsSelectorButtons    []*widget.Button
}

var state State

func StartGui(version string) {
	app := app.New()
	win := app.NewWindow("Baš Čelik")

	theme := celiktheme.NewTheme(app.Preferences().IntWithFallback(themePreferenceKey, 1))
	app.Settings().SetTheme(theme)

	translation.SetLanguage(app.Preferences().IntWithFallback(languagePreferenceKey, 1))

	statusBar := widgets.NewStatusBar()

	cryptoContainer := container.New(layout.NewVBoxLayout())
	mainPage := container.New(layout.NewVBoxLayout())

	startPage := widgets.NewStartPage()
	startPage.SetStatus("", "", false)

	mainContainer := container.New(layout.NewPaddedLayout())
	win.SetContent(mainContainer)

	state = State{
		app:                     app,
		window:                  win,
		version:                 version,
		mainContainer:           mainContainer,
		cryptoUiContainer:       cryptoContainer,
		documentUiMainContainer: mainPage,
		startPage:               startPage,
		statusBar:               statusBar,
	}

	smartboxMode := app.Preferences().BoolWithFallback(smartboxModeKey, false)

	showAboutBox := showAboutBox()
	showSettings := showSetupBox()

	toolbar := widgets.NewToolbar(showAboutBox, showSettings, !smartboxMode)
	state.toolbar = toolbar

	if smartboxMode {
		startSmartboxUI()
	} else {
		startCardReaderUI()
	}
}
