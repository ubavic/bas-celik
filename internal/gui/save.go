package gui

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/ubavic/bas-celik/document"
	"github.com/ubavic/bas-celik/internal/logger"
)

func savePdf(doc document.Document) func() {
	return func() {
		if doc == nil {
			return
		}

		pdf, fileName, err := doc.BuildPdf()

		if err != nil {
			setStatus("error.generatingPdf", fmt.Errorf("generating PDF: %w", err))
			return
		}

		dialog := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
			if w == nil || err != nil {
				setStatus("error.writingPdf", fmt.Errorf("writing PDF: %w", err))
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

		dialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
		dialog.SetFileName(fileName)

		lastUsedDirectoryURI := getLastUsedDirectory()
		if lastUsedDirectoryURI != nil {
			dialog.SetLocation(lastUsedDirectoryURI)
		}

		dialog.Show()
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

		dialog := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
			if w == nil || err != nil {
				setStatus("error.writingXlsx", fmt.Errorf("writing xlsx: %w", err))
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

		dialog.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		dialog.SetFileName(fileName)

		lastUsedDirectoryURI := getLastUsedDirectory()
		if lastUsedDirectoryURI != nil {
			dialog.SetLocation(lastUsedDirectoryURI)
		}

		dialog.Show()
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
