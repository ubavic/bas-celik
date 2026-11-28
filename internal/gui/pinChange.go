package gui

import (
	"errors"
	"strconv"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/card"
	"github.com/ubavic/bas-celik/v2/internal/gui/widgets"
	"github.com/ubavic/bas-celik/v2/internal/logger"
)

func pinChange() func() {
	return func() {
		state.mu.Lock()
		gemaltoCard, ok := state.cardDocument.(*card.Gemalto)
		if !ok {
			state.mu.Unlock()
			return
		}

		triesLeft, err := gemaltoCard.PinTriesLeft()
		state.mu.Unlock()
		if err != nil {
			logger.Error(err)
			return
		}

		pinForm(gemaltoCard, triesLeft)
	}
}

func pinForm(gemaltoCard *card.Gemalto, triesLeft int) {
	var pinDialog *dialog.CustomDialog

	triesLeftLabel := widget.NewLabel(strconv.Itoa(triesLeft))
	oldPinEntry := widget.NewPasswordEntry()
	newPinEntry := widget.NewPasswordEntry()
	confirmNewPinEntry := widget.NewPasswordEntry()

	spacer := widgets.NewSpacer()
	spacer.SetMinWidth(200)

	formItems := []*widget.FormItem{
		{Text: t("pinChange.triesLeft"), Widget: triesLeftLabel},
		{Text: t("pinChange.oldPin"), Widget: oldPinEntry},
		{Text: t("pinChange.newPin"), Widget: newPinEntry},
		{Text: t("pinChange.confirmNewPin"), Widget: confirmNewPinEntry},
		{Text: "", Widget: spacer},
	}

	form := &widget.Form{
		Items:      formItems,
		SubmitText: t("pinChange.change"),
		OnSubmit: func() {
			if newPinEntry.Text != confirmNewPinEntry.Text {
				err := errors.New(t("pinChange.pinsNotEqual"))
				dialog.ShowError(err, state.window)
				return
			}

			if !card.ValidatePin(oldPinEntry.Text) {
				err := errors.New(t("pinChange.oldPinFormatError") + " " + t("pinChange.pinFormatExplanation"))
				dialog.ShowError(err, state.window)
				return
			}

			if !card.ValidatePin(newPinEntry.Text) {
				err := errors.New(t("pinChange.newPinFormatError") + " " + t("pinChange.pinFormatExplanation"))
				dialog.ShowError(err, state.window)
				return
			}

			if !card.ValidatePin(confirmNewPinEntry.Text) {
				err := errors.New(t("pinChange.confirmNewPinFormatError") + " " + t("pinChange.pinFormatExplanation"))
				dialog.ShowError(err, state.window)
				return
			}

			state.mu.Lock()
			if state.cardDocument != gemaltoCard {
				state.mu.Unlock()
				pinDialog.Hide()
				err := errors.New(t("pinChange.errorNoCard"))
				dialog.ShowError(err, state.window)
				return
			}

			triesLeft, err := gemaltoCard.ChangePin(newPinEntry.Text, oldPinEntry.Text)
			state.mu.Unlock()
			if err != nil {
				pinDialog.Hide()
				message := t("pinChange.error")
				if triesLeft > -1 {
					message += "\n" + t("pinChange.triesLeft") + ":" + strconv.Itoa(triesLeft)
				}
				dialog.ShowInformation(t("pinChange.title"), message, state.window)
				logger.Error(err)
				return
			} else {
				pinDialog.Hide()
				dialog.ShowInformation(t("pinChange.title"), t("pinChange.success"), state.window)
				logger.Info("pin changed")
			}
		},
		CancelText: t("pinChange.cancel"),
		OnCancel: func() {
			pinDialog.Hide()
		},
	}

	pinDialog = dialog.NewCustomWithoutButtons(t("pinChange.title"), form, state.window)
	pinDialog.Show()
}
