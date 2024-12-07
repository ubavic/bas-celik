package gui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/internal/gui/widgets"
	"github.com/ubavic/bas-celik/internal/logger"
)

func showAboutBox(win fyne.Window, version string) func() {
	verLabel := widget.NewLabelWithStyle(t("about.version")+": "+version, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	moreLabel := widget.NewLabel(t("about.moreAboutProgram"))
	url, _ := url.Parse("https://github.com/ubavic/bas-celik")
	linkLabel := widget.NewHyperlink("github.com/ubavic/bas-celik", url)
	spacer := widgets.NewSpacer()
	verButton := widget.NewButton(t("about.checkVersion"), func() {
		newVersion, err := checkForUpdate()
		if err != nil || newVersion == "" {
			logger.Error(err)
			return
		}

		if newVersion != "v"+version {
			dialog.ShowInformation(t("about.version"), fmt.Sprintf(t("about.newVersionAvailable"), newVersion), win)
		} else {
			dialog.ShowInformation(t("about.version"), t("about.youHaveLatestVersion"), win)
		}
	})

	hBox0 := container.NewHBox(verLabel, verButton)
	hBox1 := container.NewHBox(moreLabel, linkLabel)
	vBox := container.NewVBox(hBox0, hBox1, spacer)

	return func() {
		dialog.ShowCustom(
			t("about.title"),
			t("about.close"),
			vBox,
			win,
		)
	}
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
