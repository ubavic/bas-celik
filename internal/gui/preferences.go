package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/internal/gui/widgets"
)

const themePreferenceKey = "color-theme"

func showSetupBox(win fyne.Window, app fyne.App) func() {
	preferences := app.Preferences()

	return func() {
		colorTheme := preferences.IntWithFallback(themePreferenceKey, 0)
		themeSelect := widget.NewSelect([]string{"Operativni sistem određuje", "Uvek svetla", "Uvek tamna"}, func(s string) {})
		themeSelect.SetSelectedIndex(colorTheme)

		formItems := []*widget.FormItem{
			{Text: "Tema aplikacije", Widget: themeSelect},
			{Text: "", Widget: &widgets.Spacer{}},
		}

		onExit := func(save bool) {
			if !save {
				return
			}

			preferences.SetInt(themePreferenceKey, themeSelect.SelectedIndex())

			dialog.ShowInformation("Podešavanja sačuvana", "Ponovo pokrenite aplikaciju da biste videli promene.", win)
		}

		dialog.ShowForm("Podešavanja", "Sačuvaj", "Izađi", formItems, onExit, win)
	}
}
