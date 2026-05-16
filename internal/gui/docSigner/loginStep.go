package docSigning

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/internal/logger"
)

type loginStep struct {
	pinEntry *widget.Entry
	errLabel *widget.Label
}

func (l *loginStep) activate() *fyne.Container {
	if l.pinEntry == nil {
		l.pinEntry = widget.NewPasswordEntry()
		l.pinEntry.PlaceHolder = gWizard.t("xmlSign.step.login")
		l.pinEntry.OnChanged = func(s string) {
			if s != "" {
				gWizard.nextBtn.Enable()
			} else {
				gWizard.nextBtn.Disable()
			}
		}
	}

	if l.errLabel == nil {
		l.errLabel = widget.NewLabel("")
	}

	gWizard.pinEntry = l.pinEntry
	l.pinEntry.SetText("")
	l.pinEntry.Enable()

	l.errLabel.SetText("")
	l.errLabel.Hide()
	gWizard.loginErr = l.errLabel

	gWizard.nextBtn.Disable()

	return container.New(layout.NewVBoxLayout(), l.pinEntry, l.errLabel)
}

func (l *loginStep) complete() *stepError {
	if l.pinEntry == nil || l.pinEntry.Text == "" {
		return &stepError{"", nil}
	}

	pin := l.pinEntry.Text
	slotId := int(gWizard.slotIds[gWizard.selectedSlotIdx])

	l.pinEntry.SetText("XXXXXXXX")
	l.pinEntry.DisableableWidget.Disable()

	gWizard.nextBtn.Disable()
	gWizard.backBtn.Disable()
	if l.errLabel != nil {
		l.errLabel.Hide()
	}

	go func() {
		err := gWizard.session.OpenSessionAndLogin(pin, slotId)
		fyne.Do(func() {
			if err != nil {
				l.showError(gWizard.t("xmlSign.loginError"))
				logger.Error(err)
				gWizard.nextBtn.Enable()
				gWizard.backBtn.Enable()
				l.pinEntry.DisableableWidget.Enable()
				l.pinEntry.SetText("")
				return
			}

			certs, err := gWizard.session.GetCertificates()
			if err != nil || len(certs) == 0 {
				msg := gWizard.t("crypto.noCertFound")
				if err != nil {
					msg += ": " + err.Error()
				}
				l.showError(msg)
				gWizard.nextBtn.Enable()
				gWizard.backBtn.Enable()
				l.pinEntry.SetText("")
				return
			}

			gWizard.certs = certs
			if l.errLabel != nil {
				l.errLabel.Hide()
			}

			gWizard.nextBtn.SetText(gWizard.t("xmlSign.next"))
			gWizard.nextBtn.Disable()
			gWizard.backBtn.Enable()

			gWizard.appendStep(&certStep{})
		})
	}()

	return nil
}

func (l *loginStep) showError(msg string) {
	if l.errLabel == nil {
		return
	}
	l.errLabel.SetText(msg)
	l.errLabel.Show()
	l.errLabel.Refresh()
}

func (l *loginStep) revert() {
	gWizard.loginErr = nil
	gWizard.pinEntry = nil
}
