package gui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/internal/logger"
)

func checkVersionOnStartup() {
	time.Sleep(time.Second)

	preferences := state.app.Preferences()
	if !preferences.BoolWithFallback(checkVersionOnStartupKey, true) {
		return
	}

	go func() {
		newVersion, err := checkForUpdate()
		if err != nil {
			logger.Error(fmt.Errorf("checking version: %w", err))
			return
		}
		if newVersion == "" || newVersion == "v"+state.version {
			return
		}

		state.app.SendNotification(fyne.NewNotification("Baš Čelik", fmt.Sprintf(t("about.newVersionAvailable"), newVersion)))
	}()
}

func populateNewVersionInfo(versionLabel *widget.Label) {
	newVersion, err := checkForUpdate()
	if err != nil {
		logger.Error(fmt.Errorf("populating version info: %w", err))
		return
	}
	if newVersion == "" {
		return
	}

	information := ""
	if newVersion != "v"+state.version {
		information = fmt.Sprintf(t("about.newVersionAvailable"), newVersion)
	} else {
		information = t("about.youHaveLatestVersion")
	}

	fyne.Do(func() {
		versionLabel.SetText(information)
	})
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
