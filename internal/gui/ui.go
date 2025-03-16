package gui

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/ubavic/bas-celik/v2/card"
	"github.com/ubavic/bas-celik/v2/internal/gui/celiktheme"
	"github.com/ubavic/bas-celik/v2/internal/gui/translation"
	"github.com/ubavic/bas-celik/v2/internal/gui/widgets"
)

type State struct {
	mu            sync.Mutex
	app           fyne.App
	window        fyne.Window
	startPage     *widgets.StartPage
	toolbar       *widgets.Toolbar
	spacer        *widgets.Spacer
	mainPage      *fyne.Container
	mainContainer *fyne.Container
	statusBar     *widgets.StatusBar
	cardDocument  card.CardDocument
}

var state State
var version string

func StartGui(_version string) {
	version = _version

	app := app.New()
	win := app.NewWindow("Baš Čelik")

	theme := celiktheme.NewTheme(app.Preferences().IntWithFallback(themePreferenceKey, 1))
	app.Settings().SetTheme(theme)

	translation.SetLanguage(app.Preferences().IntWithFallback(languagePreferenceKey, 1))

	startCardReaderUI(app, win)
}
