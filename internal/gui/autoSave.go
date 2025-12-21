package gui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ubavic/bas-celik/v2/document"
)

type AutoSaveMode int

const (
	AutoSaveNone = AutoSaveMode(iota)
	AutoSaveSave
	AutoSaveSaveAndOpen
)

func autoSave(doc document.Document) {
	if doc == nil {
		return
	}

	if state.autoSaveMode <= 0 || state.autoSaveMode > 2 {
		return
	}

	if state.autoSaveLocation == "" {
		return
	}

	fileInfo, err := os.Stat(state.autoSaveLocation)
	if err != nil {
		setTimedStatusError("autoSave.error.general", fmt.Errorf("autosave: stating directory: %w", err))
		return
	}

	if !fileInfo.IsDir() {
		setTimedStatusError("autoSave.error.nonDirectory", fmt.Errorf("location '%s' is not directory", state.autoSaveLocation))
		return
	}

	pdf, name, err := doc.BuildPdf()
	if err != nil {
		setTimedStatusError("autoSave.error.general", fmt.Errorf("autosave: building pdf: %w", err))
		return
	}

	pdfPath := filepath.Join(state.autoSaveLocation, time.Now().Format("2006-01-02_15-04")+"_"+name)

	err = os.WriteFile(pdfPath, pdf, 0600)
	if err != nil {
		setTimedStatusError("autoSave.error.general", fmt.Errorf("autosave: saving data: %w", err))
		return
	}

	if state.autoSaveMode == 2 {
		err = openFile(pdfPath)
		if err != nil {
			setTimedStatusError("autoSave.error.fileOpen", fmt.Errorf("autosave: opening saved file: %w", err))
			return
		}
	}

	setTimedStatus(t("autoSave.done"))
}

func openFile(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	default:
		return fmt.Errorf("unsupported OS")
	}

	return cmd.Start()
}
