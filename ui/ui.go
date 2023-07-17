package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/document"
)

var nameF, birthDateF, sexF, personalNumberF, birthPlaceF, addressF, addressDateF *Field
var issuedByF, documentNumberF, issueDateF, expiryDateF *Field
var statusBar *StatusBar
var imgWidget *canvas.Image
var ui *fyne.Container
var doc *document.Document

var verbose bool

func SetUI(w *fyne.Window, d *document.Document, ver bool) *fyne.Container {
	doc = d
	verbose = ver

	nameF = newField("Ime, ime roditelja, prezime", 350)
	birthDateF = newField("Datum rođenja", 100)
	sexF = newField("Pol", 80)
	personalNumberF = newField("JMBG", 200)
	birthRow := container.New(layout.NewHBoxLayout(), sexF, birthDateF, personalNumberF)
	birthPlaceF = newField("Mesto rođenja, opština i država", 350)
	addressF = newField("Prebivalište i adresa stana", 350)
	addressDateF = newField("Datum promene adrese", 10)
	issuedByF = newField("Dokument izdaje", 10)
	documentNumberF = newField("Broj dokumenta", 100)
	issueDateF = newField("Datum izdavanja", 200)
	expiryDateF = newField("Važi do", 200)
	dateRow := container.New(layout.NewHBoxLayout(), issueDateF, expiryDateF)
	colRight := container.New(layout.NewVBoxLayout(), nameF, birthRow, birthPlaceF, addressF, addressDateF, issuedByF, documentNumberF, dateRow)

	imgWidget = canvas.NewImageFromImage(d.DefaultPhoto)
	imgWidget.SetMinSize(fyne.Size{Width: 200, Height: 250})
	imgWidget.FillMode = canvas.ImageFillContain
	colLeft := container.New(layout.NewVBoxLayout(), imgWidget)
	cols := container.New(layout.NewHBoxLayout(), colLeft, colRight)

	statusBar = newStatusBar()
	saveButton := widget.NewButton("Sačuvaj PDF", savePdf(w))
	buttonBar := container.New(layout.NewHBoxLayout(), statusBar, layout.NewSpacer(), saveButton)

	ui = container.New(layout.NewVBoxLayout(), cols, buttonBar)

	return ui
}

func UpdateUI() {
	nameF.value = doc.GivenName + ", " + doc.ParentName + ", " + doc.Surname
	birthDateF.value = doc.DateOfBirth
	sexF.value = doc.Sex
	personalNumberF.value = doc.PersonalNumber
	birthPlaceF.value = doc.FormatPlaceOfBirth()
	addressF.value = doc.FormatAddress()
	addressDateF.value = doc.AddressDate

	issuedByF.value = doc.IssuingAuthority
	issueDateF.value = doc.IssuingDate
	expiryDateF.value = doc.ExpiryDate
	documentNumberF.value = doc.DocumentNumber

	imgWidget.Image = doc.Photo

	ui.Refresh()
}

func SetStatus(status string, err error) {
	isError := false
	if err != nil {
		isError = true
	}

	if verbose && isError {
		fmt.Println(err)
	}

	statusBar.SetStatus(status, isError)
	statusBar.Refresh()
}

func savePdf(win *fyne.Window) func() {
	return func() {
		if doc.Loaded {
			pdf, fileName, err := doc.Pdf()

			if err != nil {
				SetStatus(
					"Greška pri generisanju PDF-a",
					fmt.Errorf("generating PDF: %w", err))
				return
			}

			dialog := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
				if w == nil || err != nil {
					return
				}

				_, err = w.Write(pdf)
				if err != nil {
					SetStatus(
						"Greška pri zapisivanju PDF-a",
						fmt.Errorf("writing PDF: %w", err))
					return
				}

				err = w.Close()
				if err != nil {
					SetStatus(
						"Greška pri zapisivanju PDF-a",
						fmt.Errorf("writing PDF: %w", err))
					return
				}

				SetStatus("PDF sačuvan", nil)
			}, *win)

			dialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
			dialog.SetFileName(fileName)

			dialog.Show()
		}
	}
}
