package gui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/internal/gui/widgets"
	"github.com/ubavic/bas-celik/v2/internal/logger"
	"github.com/ubavic/bas-celik/v2/internal/smartbox/pkcs11"
)

const themePreferenceKey = "color-theme"
const languagePreferenceKey = "language"
const pdfScriptPreferenceKey = "pdf-labels-script"
const autoSavePdfKey = "auto-save-pdf"
const autoSaveLocationKey = "auto-save-location"

const runInBackgroundKey = "run-in-background"

const lastUsedDirectoryKey = "last-used-directory"

const smartboxModeKey = "smartbox-mode"
const mupPkcsPathKey = "mup-pkcs-path"
const pksPkcsPathKey = "pks-pkcs-path"
const postaPkcsPathKey = "posta-pkcs-path"
const halcomPkcsPathKey = "halcom-pkcs-path"
const esmartPkcsPathKey = "esmart-pkcs-path"

func showSetupBox() func() {
	return func() {
		preferences := state.app.Preferences()

		colorTheme := preferences.IntWithFallback(themePreferenceKey, 0)
		themeSelect := widget.NewSelect(
			[]string{t("preference.theme.osDetermines"), t("preference.theme.alwaysLight"), t("preference.theme.alwaysDark")},
			func(s string) {})
		themeSelect.SetSelectedIndex(colorTheme)

		language := preferences.IntWithFallback(languagePreferenceKey, 0)
		languageSelect := widget.NewSelect(
			[]string{"Srpski", "Српски", "English", "Русский"},
			func(s string) {},
		)
		languageSelect.SetSelectedIndex(language)

		pdfScript := preferences.IntWithFallback(pdfScriptPreferenceKey, 0)
		pdfScriptSelect := widget.NewSelect(
			[]string{"Latinica", "Ћирилица"},
			func(s string) {},
		)
		pdfScriptSelect.SetSelectedIndex(pdfScript)

		autoSavePdf := preferences.IntWithFallback(autoSavePdfKey, 0)
		autoSavePdfSelect := widget.NewSelect(
			[]string{t("preference.autoSave.no"), t("preference.autoSave.save"), t("preference.autoSave.saveAndOpen")},
			func(s string) {},
		)
		autoSavePdfSelect.SetSelectedIndex(autoSavePdf)

		autoSaveLocation := preferences.String(autoSaveLocationKey)
		autoSaveLocationEntry := widget.NewEntry()
		autoSaveLocationEntry.SetText(autoSaveLocation)
		autoSaveLocationEntry.SetPlaceHolder(t("preference.placeholder.directoryPath"))

		runInBackground := preferences.BoolWithFallback(runInBackgroundKey, false)
		runInBackgroundCheck := widget.NewCheck("", func(b bool) {})
		runInBackgroundCheck.SetChecked(runInBackground)

		mupPkcsEntry := widget.NewEntry()
		mupPkcsEntry.SetText(preferences.String(mupPkcsPathKey))
		mupPkcsEntry.SetPlaceHolder(t("preference.placeholder.modulePath"))
		pksPkcsEntry := widget.NewEntry()
		pksPkcsEntry.SetText(preferences.String(pksPkcsPathKey))
		pksPkcsEntry.SetPlaceHolder(t("preference.placeholder.modulePath"))
		postaPkcsEntry := widget.NewEntry()
		postaPkcsEntry.SetText(preferences.String(postaPkcsPathKey))
		postaPkcsEntry.SetPlaceHolder(t("preference.placeholder.modulePath"))
		halcomPkcsEntry := widget.NewEntry()
		halcomPkcsEntry.SetText(preferences.String(halcomPkcsPathKey))
		halcomPkcsEntry.SetPlaceHolder(t("preference.placeholder.modulePath"))
		esmartPkcsEntry := widget.NewEntry()
		esmartPkcsEntry.SetText(preferences.String(esmartPkcsPathKey))
		esmartPkcsEntry.SetPlaceHolder(t("preference.placeholder.modulePath"))

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

		save := func() {
			preferences.SetInt(themePreferenceKey, themeSelect.SelectedIndex())
			preferences.SetInt(languagePreferenceKey, languageSelect.SelectedIndex())
			preferences.SetInt(pdfScriptPreferenceKey, pdfScriptSelect.SelectedIndex())
			preferences.SetInt(autoSavePdfKey, autoSavePdfSelect.SelectedIndex())
			preferences.SetString(autoSaveLocationKey, autoSaveLocationEntry.Text)
			preferences.SetBool(runInBackgroundKey, runInBackgroundCheck.Checked)
			preferences.SetBool(smartboxModeKey, smartboxModeCheck.Checked)
			preferences.SetString(mupPkcsPathKey, mupPkcsEntry.Text)
			preferences.SetString(pksPkcsPathKey, pksPkcsEntry.Text)
			preferences.SetString(postaPkcsPathKey, postaPkcsEntry.Text)
			preferences.SetString(halcomPkcsPathKey, halcomPkcsEntry.Text)
			preferences.SetString(esmartPkcsPathKey, esmartPkcsEntry.Text)

			dialog.ShowInformation(t("preference.saved"), t("preference.startAgain"), state.window)
		}

		title := canvas.NewText(t("preference.title"), theme.Color(theme.ColorNameForeground))
		title.TextStyle.Bold = true
		title.TextSize = 20

		spacer := widgets.NewSpacer()
		spacer.SetMinWidth(160)

		rows1 := container.New(layout.NewFormLayout(),
			spacer,
			spacer,
			widget.NewLabel(t("preference.runInBackground")),
			runInBackgroundCheck,
			widget.NewLabel(t("preference.theme")),
			themeSelect,
			widget.NewLabel(t("preference.language")),
			languageSelect,
			widget.NewLabel(t("preference.pdfScript")),
			pdfScriptSelect,
			widget.NewLabel(t("preference.autoSave")),
			autoSavePdfSelect,
			widget.NewLabel(t("preference.autoSaveLocation")),
			autoSaveLocationEntry)

		rows2 := container.New(layout.NewFormLayout(),
			spacer,
			spacer,
			widget.NewLabel(t("preference.smartboxMode")),
			smartboxModeCheck,
			widget.NewLabel("MUP"),
			mupPkcsEntry,
			widget.NewLabel("PKS"),
			pksPkcsEntry,
			widget.NewLabel("Pošta"),
			postaPkcsEntry,
			widget.NewLabel("Halcom"),
			halcomPkcsEntry,
			widget.NewLabel("E-Smart"),
			esmartPkcsEntry,
		)

		ghUrl, _ := url.Parse("https://github.com/ubavic/bas-celik")
		authorUrl, _ := ghUrl.Parse("https://ubavic.rs")
		guideUrl, _ := ghUrl.Parse("https://ubavic.rs/e-documents/")
		newVersionInfoLabel := widget.NewLabel("")
		go populateNewVersionInfo(newVersionInfoLabel)

		rows3 := container.New(layout.NewFormLayout(),
			spacer,
			spacer,
			widget.NewLabel(t("about.version")),
			widget.NewLabel(state.version),
			spacer,
			newVersionInfoLabel,
			widget.NewLabel(t("about.moreAboutProgram")),
			widget.NewHyperlink("github.com/ubavic/bas-celik", ghUrl),
			widget.NewLabel(t("about.author")),
			widget.NewHyperlink("Nikola Ubavić", authorUrl),
			widget.NewLabel(t("about.guide")),
			widget.NewHyperlink("ubavic.rs/e-documents", guideUrl),
		)

		afterChangeSmartboxMode = func() {
			rows2.Refresh()
		}

		tabs := container.NewAppTabs(
			container.NewTabItem(t("preference.general"), rows1),
			container.NewTabItem("Smartbox", rows2),
			container.NewTabItem(t("preference.about"), rows3),
		)

		saveButton := widget.NewButton(t("preference.save"), func() {
			save()
			showDocumentUI()
			setTimedStatus(t("preference.saved"))
		})
		saveButton.Importance = widget.HighImportance

		cancelButton := widget.NewButton(t("preference.exit"), func() {
			showDocumentUI()
			setTimedStatus("")
		})

		buttonRow := container.New(layout.NewHBoxLayout(),
			layout.NewSpacer(),
			cancelButton,
			saveButton,
		)

		content := container.New(layout.NewVBoxLayout(),
			title,
			tabs,
			spacer,
			buttonRow)

		state.mainContainer.RemoveAll()
		state.mainContainer.Add(content)
	}
}

func populateNewVersionInfo(versionLabel *widget.Label) {
	newVersion, err := checkForUpdate()
	if err != nil || newVersion == "" {
		logger.Error(err)
		return
	}

	information := ""

	if newVersion != "v"+state.version {
		information = fmt.Sprintf(t("about.newVersionAvailable"), newVersion)
	} else {
		information = t("about.youHaveLatestVersion")
	}

	versionLabel.SetText(information)
}

func checkForUpdate() (string, error) {
	client := http.Client{
		Timeout: time.Second,
	}

	url := `https://api.github.com/repos/ubavic/bas-celik/releases/latest`
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	response := struct {
		TagName string `json:"tag_name"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	return response.TagName, nil
}
