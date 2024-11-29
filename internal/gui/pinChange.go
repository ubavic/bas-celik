package gui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/card"
	"github.com/ubavic/bas-celik/internal/gui/widgets"
	"github.com/ubavic/bas-celik/internal/logger"
)

func pinChange(win fyne.Window) func() {
	return func() {
		dialog.ShowConfirm(t("pinChange.title"), t("pinChange.note"), func(changePinContinue bool) {
			if changePinContinue {
				pinForm(win)
			}
		}, win)
	}
}

func pinForm(win fyne.Window) {
	var pinDialog *dialog.CustomDialog

	oldPinEntry := widget.NewPasswordEntry()
	newPinEntry := widget.NewPasswordEntry()
	confirmNewPinEntry := widget.NewPasswordEntry()

	spacer := widgets.NewSpacer()
	spacer.SetMinWidth(200)

	formItems := []*widget.FormItem{
		{Text: t("pinChange.oldPin"), Widget: oldPinEntry},
		{Text: t("pinChange.newPin"), Widget: newPinEntry},
		{Text: t("pinChange.confirmNewPin"), Widget: confirmNewPinEntry},
		{Text: "", Widget: spacer},
	}

	form := &widget.Form{
		Items:      formItems,
		SubmitText: t("pinChange.change"),
		OnSubmit: func() {
			gemaltoCard, ok := state.cardDocument.(*card.Gemalto)
			if !ok {
				pinDialog.Hide()
				return
			}

			if newPinEntry.Text != confirmNewPinEntry.Text {
				err := errors.New(t("pinChange.pinsNotEqual"))
				dialog.ShowError(err, win)
				return
			}

			if !card.ValidatePin(oldPinEntry.Text) {
				err := errors.New(t("pinChange.oldPinFormatError") + " " + t("pinChange.pinFormatExplanation"))
				dialog.ShowError(err, win)
				return
			}

			if !card.ValidatePin(newPinEntry.Text) {
				err := errors.New(t("pinChange.newPinFormatError") + " " + t("pinChange.pinFormatExplanation"))
				dialog.ShowError(err, win)
				return
			}

			if !card.ValidatePin(confirmNewPinEntry.Text) {
				err := errors.New(t("pinChange.confirmNewPinFormatError") + " " + t("pinChange.pinFormatExplanation"))
				dialog.ShowError(err, win)
				return
			}

			if state.cardDocument == nil {
				err := errors.New(t("pinChange.errorNoCard"))
				dialog.ShowError(err, win)
				return
			}

			err := gemaltoCard.ChangePin(newPinEntry.Text, oldPinEntry.Text)
			if err != nil {
				pinDialog.Hide()
				dialog.ShowInformation(t("pinChange.title"), t("pinChange.error"), win)
				logger.Error(err)
				return
			} else {
				pinDialog.Hide()
				dialog.ShowInformation(t("pinChange.title"), t("pinChange.success"), win)
				logger.Info("pin changed")
			}
		},
		CancelText: t("pinChange.cancel"),
		OnCancel: func() {
			pinDialog.Hide()
		},
	}

	pinDialog = dialog.NewCustomWithoutButtons(t("pinChange.title"), form, win)
	pinDialog.Show()
}
