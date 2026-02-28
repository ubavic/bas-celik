package gui

import (
	"crypto/x509"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
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
	autoSaveMode            AutoSaveMode
	autoSaveLocation        string
	runInBackground         bool
	pdfCyrillicLabels       bool
}

var state State

func StartGui(version string) {
	app := app.New()
	win := app.NewWindow("Baš Čelik")

	preferences := app.Preferences()

	theme := celiktheme.NewTheme(preferences.IntWithFallback(themePreferenceKey, 0))
	app.Settings().SetTheme(theme)

	translation.SetLanguage(preferences.IntWithFallback(languagePreferenceKey, 0))

	statusBar := widgets.NewStatusBar()

	cryptoContainer := container.New(layout.NewVBoxLayout())
	mainPage := container.New(layout.NewVBoxLayout())

	startPage := widgets.NewStartPage()
	startPage.SetStatus("", "", false)

	mainContainer := container.New(layout.NewPaddedLayout())
	win.SetContent(mainContainer)

	runInBackground := preferences.BoolWithFallback(runInBackgroundKey, false)

	state = State{
		app:                     app,
		window:                  win,
		version:                 version,
		mainContainer:           mainContainer,
		cryptoUiContainer:       cryptoContainer,
		documentUiMainContainer: mainPage,
		startPage:               startPage,
		statusBar:               statusBar,
		autoSaveMode:            AutoSaveMode(preferences.Int(autoSavePdfKey)),
		autoSaveLocation:        preferences.String(autoSaveLocationKey),
		runInBackground:         runInBackground,
		pdfCyrillicLabels:       preferences.Int(pdfScriptPreferenceKey) == 1,
	}

	smartboxMode := preferences.BoolWithFallback(smartboxModeKey, false)

	showAboutBox := showAboutBox()
	showSettings := showSetupBox()

	toolbar := widgets.NewToolbar(showAboutBox, showSettings, !smartboxMode)
	state.toolbar = toolbar

	if runInBackground {
		setupTray()
	}

	if smartboxMode {
		startSmartboxUI()
	} else {
		startCardReaderUI()
	}
}

func setupTray() {
	desk, ok := state.app.(desktop.App)
	if !ok {
		return
	}

	m := fyne.NewMenu("Baš Čelik",
		fyne.NewMenuItem(t("tray.show"), func() {
			state.window.Show()
		}))
	desk.SetSystemTrayMenu(m)

	state.window.SetCloseIntercept(func() {
		state.window.Hide()
	})
}
