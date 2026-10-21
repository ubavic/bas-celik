package gui

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/ubavic/bas-celik/v2/document"
	"github.com/ubavic/bas-celik/v2/internal/logger"
    "os"
    "os/exec"
)

func savePdf(doc document.Document) func() {
	return func() {
		if doc == nil {
			return
		}

		switch doc := doc.(type) {
		case *document.IdDocument:
			if state.pdfCyrillicLabels {
				doc.SetCyrillicLabels()
			}
		}

		pdf, fileName, err := doc.BuildPdf()

		if err != nil {
			setStatus("error.generatingPdf", fmt.Errorf("generating PDF: %w", err))
			return
		}

		fDialog := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
			if err != nil {
				setStatus("error.writingPdf", fmt.Errorf("writing PDF: %w", err))
				return
			}

			if w == nil {
				return
			}

			saveLastUsedDirectory(w.URI())

			_, err = w.Write(pdf)
			if err != nil {
				setStatus("error.writingPdf", fmt.Errorf("writing PDF: %w", err))
				return
			}

			err = w.Close()
			if err != nil {
				setStatus("error.writingPdf", fmt.Errorf("writing PDF: %w", err))
				return
			}

			setStatus("ui.pdfSaved", nil)
		}, state.window)

		fDialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
		fDialog.SetFileName(fileName)
		fDialog.SetView(dialog.ListView)

		lastUsedDirectoryURI := getLastUsedDirectory()
		if lastUsedDirectoryURI != nil {
			fDialog.SetLocation(lastUsedDirectoryURI)
		}

		fDialog.Show()
	}
}

func saveXlsx(doc document.Document) func() {
	return func() {
		if doc == nil {
			return
		}

		excel, fileName, err := doc.BuildExcel()

		if err != nil {
			setStatus("error.generatingXlsx", fmt.Errorf("generating Xlsx: %w", err))
			return
		}

		fDialog := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
			if err != nil {
				setStatus("error.writingXlsx", fmt.Errorf("writing xlsx: %w", err))
				return
			}

			if w == nil {
				return
			}

			saveLastUsedDirectory(w.URI())

			_, err = w.Write(excel)
			if err != nil {
				setStatus("error.writingXlsx", fmt.Errorf("writing xlsx: %w", err))
				return
			}

			err = w.Close()
			if err != nil {
				setStatus("error.writingXlsx", fmt.Errorf("writing xlsx: %w", err))
				return
			}

			setStatus("ui.xlsxSaved", nil)
		}, state.window)

		fDialog.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		fDialog.SetFileName(fileName)
		fDialog.SetView(dialog.ListView)

		lastUsedDirectoryURI := getLastUsedDirectory()
		if lastUsedDirectoryURI != nil {
			fDialog.SetLocation(lastUsedDirectoryURI)
		}

		fDialog.Show()
	}
}

func saveLastUsedDirectory(uri fyne.URI) {
	directoryPath := filepath.Dir(uri.Path())

	if directoryPath == "." || directoryPath == "" {
		return
	}

	preferences := state.app.Preferences()

	preferences.SetString(lastUsedDirectoryKey, directoryPath)
}

func getLastUsedDirectory() fyne.ListableURI {
	preferences := state.app.Preferences()
	lastUsedDirectory := preferences.String(lastUsedDirectoryKey)

	if lastUsedDirectory == "" {
		return nil
	}

	fileURI := storage.NewFileURI(lastUsedDirectory)
	listableURI, err := storage.ListerForURI(fileURI)
	if err != nil {
		logger.Error(err)
		return nil
	}

	return listableURI
}

    //Dodajemo print funkciju

func printPdf(doc document.Document) func() {
	return func() {
		if doc == nil {
			return
		}

		switch doc := doc.(type) {
		case *document.IdDocument:
			if state.pdfCyrillicLabels {
				doc.SetCyrillicLabels()
			}
		}

		pdf, _, err := doc.BuildPdf()
		if err != nil {
			setStatus("error.generatingPdf", fmt.Errorf("generating PDF: %w", err))
			return
		}

		// Pravimo privremeni PDF fajl
		tempFile, err := os.CreateTemp("", "bas-celik-print-*.pdf")
		if err != nil {
			logger.Error(fmt.Errorf("creating temp print file: %w", err))
			return
		}
		tempPath := tempFile.Name()

		if _, err := tempFile.Write(pdf); err != nil {
			tempFile.Close()
			os.Remove(tempPath)
			logger.Error(fmt.Errorf("writing temp print file: %w", err))
			return
		}
		tempFile.Close()

		// Šaljemo na podrazumevani štampač preko CUPS komande 'lpr'
		cmd := exec.Command("lpr", tempPath)
		err = cmd.Run()

		// Brišemo privremeni fajl
		_ = os.Remove(tempPath)

		if err != nil {
			logger.Error(fmt.Errorf("running lpr: %w", err))
			return
		}

		// Obaveštenje u statusnoj traci
		setTimedStatus("Dokument poslat na štampu")
	}
}
