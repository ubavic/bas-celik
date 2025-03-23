package gui

import (
	"runtime"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/internal/gui/widgets"
	"github.com/ubavic/bas-celik/v2/internal/smartbox/pkcs11"
)

const themePreferenceKey = "color-theme"
const languagePreferenceKey = "language"
const smartboxModeKey = "smartbox-mode"
const lastUsedDirectoryKey = "last-used-directory"

const mupPkcsPathKey = "mup-pkcs-path"
const pksPkcsPathKey = "pks-pkcs-path"
const postaPkcsPathKey = "posta-pkcs-path"
const halcomPkcsPathKey = "halcom-pkcs-path"
const esmartPkcsPathKey = "esmart-pkcs-path"

func showSetupBox() func() {
	preferences := state.app.Preferences()

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

		mupPkcsEntry := widget.NewEntry()
		mupPkcsEntry.SetText(preferences.String(mupPkcsPathKey))
		pksPkcsEntry := widget.NewEntry()
		pksPkcsEntry.SetText(preferences.String(pksPkcsPathKey))
		postaPkcsEntry := widget.NewEntry()
		postaPkcsEntry.SetText(preferences.String(postaPkcsPathKey))
		halcomPkcsEntry := widget.NewEntry()
		halcomPkcsEntry.SetText(preferences.String(halcomPkcsPathKey))
		esmartPkcsEntry := widget.NewEntry()
		esmartPkcsEntry.SetText(preferences.String(esmartPkcsPathKey))

		var afterChangeSmartboxMode func()

		changeSmartboxMode := func(mode bool) {
			if mode {
				mupPkcsEntry.Enable()
				pksPkcsEntry.Enable()
				postaPkcsEntry.Enable()
				halcomPkcsEntry.Enable()
				esmartPkcsEntry.Enable()

				if mupPkcsEntry.Text+pksPkcsEntry.Text+postaPkcsEntry.Text+halcomPkcsEntry.Text+esmartPkcsEntry.Text == "" {
					mupPkcsEntry.SetText(pkcs11.GetDefaultPath(pkcs11.CardVendorMup, runtime.GOOS))
					pksPkcsEntry.SetText(pkcs11.GetDefaultPath(pkcs11.CardVendorPks, runtime.GOOS))
					postaPkcsEntry.SetText(pkcs11.GetDefaultPath(pkcs11.CardVendorPosta, runtime.GOOS))
					halcomPkcsEntry.SetText(pkcs11.GetDefaultPath(pkcs11.CardVendorHalcom, runtime.GOOS))
					esmartPkcsEntry.SetText(pkcs11.GetDefaultPath(pkcs11.CardVendorEsmart, runtime.GOOS))
				}
			} else {
				mupPkcsEntry.Disable()
				pksPkcsEntry.Disable()
				postaPkcsEntry.Disable()
				halcomPkcsEntry.Disable()
				esmartPkcsEntry.Disable()
			}

			if afterChangeSmartboxMode != nil {
				afterChangeSmartboxMode()
			}
		}

		smartboxMode := preferences.BoolWithFallback(smartboxModeKey, false)
		smartboxModeCheck := widget.NewCheck("", changeSmartboxMode)
		smartboxModeCheck.SetChecked(smartboxMode)

		if !smartboxMode {
			changeSmartboxMode(false)
		}

		formItems := []*widget.FormItem{
			{Text: t("preference.theme"), Widget: themeSelect},
			{Text: t("preference.language"), Widget: languageSelect},
			{Text: t("preference.smartboxMode"), Widget: smartboxModeCheck},
			{Text: "MUP", Widget: mupPkcsEntry},
			{Text: "PKS", Widget: pksPkcsEntry},
			{Text: "Pošta", Widget: postaPkcsEntry},
			{Text: "Halcom", Widget: halcomPkcsEntry},
			{Text: "E-Smart", Widget: esmartPkcsEntry},
			{Text: "", Widget: &widgets.Spacer{}},
		}

		save := func() {
			preferences.SetInt(themePreferenceKey, themeSelect.SelectedIndex())
			preferences.SetInt(languagePreferenceKey, languageSelect.SelectedIndex())
			preferences.SetBool(smartboxModeKey, smartboxModeCheck.Checked)
			preferences.SetString(mupPkcsPathKey, mupPkcsEntry.Text)
			preferences.SetString(pksPkcsPathKey, pksPkcsEntry.Text)
			preferences.SetString(postaPkcsPathKey, postaPkcsEntry.Text)
			preferences.SetString(halcomPkcsPathKey, halcomPkcsEntry.Text)
			preferences.SetString(esmartPkcsPathKey, esmartPkcsEntry.Text)

			dialog.ShowInformation(t("preference.saved"), t("preference.startAgain"), state.window)
		}

		form := widget.Form{
			Items:      formItems,
			SubmitText: t("preference.save"),
			CancelText: t("preference.exit"),
			OnSubmit: func() {
				save()
				showDocumentUI()
				setTimedStatus(t("preference.saved"))
			},
			OnCancel: func() {
				showDocumentUI()
				setTimedStatus("")
			},
		}

		afterChangeSmartboxMode = func() {
			form.Refresh()
		}

		title := canvas.NewText(t("preference.title"), theme.Color(theme.ColorNameForeground))
		title.TextStyle.Bold = true
		title.TextSize = 20

		rows := container.New(layout.NewVBoxLayout(), title, &form)

		state.mainContainer.RemoveAll()
		state.mainContainer.Add(rows)
	}
}
