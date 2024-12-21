package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/internal/gui/widgets"
)

const themePreferenceKey = "color-theme"
const languagePreferenceKey = "language"
const lastUsedDirectoryKey = "last-used-directory"

func showSetupBox(win fyne.Window, app fyne.App) func() {
	preferences := app.Preferences()

	return func() {
		colorTheme := preferences.IntWithFallback(themePreferenceKey, 0)
		themeSelect := widget.NewSelect(
			[]string{t("preference.theme.osDetermines"), t("preference.theme.alwaysLight"), t("preference.theme.alwaysDark")},
			func(s string) {})
		themeSelect.SetSelectedIndex(colorTheme)

		language := preferences.IntWithFallback(languagePreferenceKey, 0)
		languageSelect := widget.NewSelect(
			[]string{"Srpski", "Српски", "English"},
			func(s string) {},
		)
		languageSelect.SetSelectedIndex(language)

		formItems := []*widget.FormItem{
			{Text: t("preference.theme"), Widget: themeSelect},
			{Text: t("preference.language"), Widget: languageSelect},
			{Text: "", Widget: &widgets.Spacer{}},
		}

		onExit := func(save bool) {
			if !save {
				return
			}

			preferences.SetInt(themePreferenceKey, themeSelect.SelectedIndex())
			preferences.SetInt(languagePreferenceKey, languageSelect.SelectedIndex())

			dialog.ShowInformation(t("preference.saved"), t("preference.startAgain"), win)
		}

		dialog.ShowForm(
			t("preference.title"),
			t("preference.save"),
			t("preference.exit"),
			formItems,
			onExit,
			win)
	}
}
