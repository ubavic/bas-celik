package docSigning

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/ubavic/bas-celik/v2/internal/gui/widgets"
	"github.com/ubavic/bas-celik/v2/internal/pkcs11"
)

var gWizard *SignWizard

type SignWizard struct {
	Panel      *fyne.Container
	stepsBox   *fyne.Container
	backBtn    *widget.Button
	nextBtn    *widget.Button
	steps      []wizardStep
	containers []*fyne.Container

	loadedVendors     []pkcs11.CardVendor
	selectedVendorIdx int

	slotIds         []uint
	slotNames       []string
	selectedSlotIdx int

	session         *pkcs11.PkcsModuleSession
	certs           []pkcs11.NamedCert
	selectedCertIdx int
	pinEntry        *widget.Entry
	loginErr        *widget.Label

	window fyne.Window
	t      func(string, ...any) string

	onStatusMsg func(string)
}

func InitSignWizard(win fyne.Window, translate func(string, ...any) string, loadedVendors []pkcs11.CardVendor, onStatusMsg func(string)) *SignWizard {
	if gWizard != nil {
		return gWizard
	}

	signWizard := SignWizard{
		loadedVendors:     loadedVendors,
		window:            win,
		t:                 translate,
		selectedVendorIdx: -1,
		selectedSlotIdx:   -1,
		selectedCertIdx:   -1,
		onStatusMsg:       onStatusMsg,
	}

	heading := widget.NewLabel(translate("xmlSign.title"))
	heading.TextStyle.Bold = true

	stepsBox := container.New(layout.NewVBoxLayout(), heading)

	backBtn := widget.NewButton(translate("xmlSign.back"), func() { signWizard.goBack() })
	nextBtn := widget.NewButton(translate("xmlSign.next"), func() { signWizard.goNext() })
	nextBtn.Importance = widget.HighImportance
	backBtn.Disable()
	nextBtn.Disable()

	signWizard.stepsBox = stepsBox
	signWizard.backBtn = backBtn
	signWizard.nextBtn = nextBtn

	scroll := container.NewVScroll(stepsBox)
	scroll.SetMinSize(fyne.NewSize(600, 400))
	buttonRow := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), backBtn, nextBtn)
	signWizard.Panel = container.New(layout.NewVBoxLayout(), scroll, buttonRow)

	gWizard = &signWizard

	gWizard.appendStep(&moduleStep{})

	return gWizard
}

func (w *SignWizard) Deinit() {
	if w.session != nil {
		w.session.CloseSession()
		w.session = nil
	}
	pkcs11.Deinit()
	gWizard = nil
}

func (w *SignWizard) currentStepIdx() int {
	return len(w.steps) - 1
}

func (w *SignWizard) removeLastStep() {
	n := len(w.steps)
	if n == 0 {
		return
	}

	if n <= len(w.containers) {
		w.stepsBox.Remove(w.containers[n-1])
		w.containers = w.containers[:n-1]
	}

	w.steps = w.steps[:n-1]
	w.stepsBox.Refresh()
}

func (w *SignWizard) appendStep(step wizardStep) {
	c := step.activate()
	spacer := widgets.NewSpacer()
	spacer.SetMinWidth(30)
	wrapper := container.New(layout.NewVBoxLayout(), c, spacer)
	w.steps = append(w.steps, step)
	w.containers = append(w.containers, wrapper)
	w.stepsBox.Add(wrapper)
	w.stepsBox.Refresh()
}

func (w *SignWizard) goBack() {
	current := w.currentStepIdx()
	if current <= 0 {
		return
	}

	if current >= 4 {
		// Leaving signing mode: go back to cert selection, keep session open.
		for w.currentStepIdx() >= 4 {
			w.removeLastStep()
		}

		// Reactivate certStep (index 3) — re-enables the select in place.
		w.steps[3].activate()
		w.nextBtn.SetText(w.t("xmlSign.next"))
		if w.selectedCertIdx >= 0 {
			w.nextBtn.Enable()
		} else {
			w.nextBtn.Disable()
		}
		w.backBtn.Enable()
		w.stepsBox.Refresh()
		return
	}

	// Steps 1, 2, or 3: remove current, reactivate previous.
	w.removeLastStep()
	n := len(w.steps)
	if n > 0 {
		w.steps[n-1].activate()
	}

	switch w.currentStepIdx() {
	case 0:
		w.session = nil
		w.nextBtn.SetText(w.t("xmlSign.next"))
		if w.selectedVendorIdx >= 0 {
			w.nextBtn.Enable()
		} else {
			w.nextBtn.Disable()
		}
		w.backBtn.Disable()
	case 1:
		w.selectedSlotIdx = -1
		w.nextBtn.SetText(w.t("xmlSign.next"))
		w.nextBtn.Disable()
		w.backBtn.Enable()
	case 2:
		// Close the PKCS#11 session but keep the module context alive so the
		// user can re-enter the PIN and log in again without restarting.
		if w.session != nil {
			w.session.CloseSession()
		}
		w.certs = nil
		w.selectedCertIdx = -1
		w.nextBtn.SetText(w.t("xmlSign.next"))
		if w.pinEntry != nil && w.pinEntry.Text != "" {
			w.nextBtn.Enable()
		} else {
			w.nextBtn.Disable()
		}
		w.backBtn.Enable()
	}

	w.stepsBox.Refresh()
}

func (w *SignWizard) goNext() {
	step := w.steps[w.currentStepIdx()]

	err := step.complete()
	if err != nil && err.msg != "" {
		if w.onStatusMsg != nil {
			w.onStatusMsg(w.t(err.msg))
		}
		return
	}

	switch w.currentStepIdx() {
	case 0:
		w.nextBtn.Disable()
		w.backBtn.Enable()
		w.appendStep(&slotStep{})
	case 1:
		w.nextBtn.Disable()
		w.backBtn.Enable()
		w.appendStep(&loginStep{})
	case 2:
		// loginStep.complete() starts an async goroutine that appends certStep
		// on success; nothing else to do here.
	case 3:
		// certStep completed: add the sign step.
		w.nextBtn.Disable()
		w.backBtn.Enable()
		w.appendStep(&signStep{})
	}
}
