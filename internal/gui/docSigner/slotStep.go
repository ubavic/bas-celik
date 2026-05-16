package docSigning

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type slotStep struct {
	slotSelect    *widget.Select
	refreshButton *widget.Button
	active        bool
}

func (s *slotStep) activate() *fyne.Container {
	if s.slotSelect == nil {
		s.slotSelect = widget.NewSelect(nil, nil)
		s.slotSelect.PlaceHolder = gWizard.t("xmlSign.step.selectSlot")
		s.slotSelect.OnChanged = func(_ string) {
			gWizard.selectedSlotIdx = s.slotSelect.SelectedIndex()
			gWizard.nextBtn.Enable()
		}
	}
	s.slotSelect.Disable()

	if s.refreshButton == nil {
		s.refreshButton = widget.NewButton(gWizard.t("xmlSign.refresh"), func() {
			s.loadSlots()
		})
	}
	s.refreshButton.DisableableWidget.Enable()

	go s.loadSlots()

	return container.New(layout.NewHBoxLayout(), s.slotSelect, s.refreshButton)
}

func (s *slotStep) loadSlots() {
	if gWizard.session == nil {
		return
	}

	fyne.Do(func() {
		s.slotSelect.Disable()
		gWizard.nextBtn.Disable()
	})

	slotIds, slotNames, err := gWizard.session.ListSlots()

	fyne.Do(func() {
		if err != nil || len(slotIds) == 0 {
			msg := gWizard.t("error.noReader")
			if err != nil {
				msg += ": " + err.Error()
			}
			return
		}

		gWizard.slotIds = slotIds
		gWizard.slotNames = slotNames

		s.slotSelect.Options = slotNames
		s.slotSelect.Enable()

		if len(slotNames) > 0 {
			s.slotSelect.SetSelectedIndex(0)
		}

		s.slotSelect.Refresh()
	})
}

func (s *slotStep) complete() *stepError {
	if gWizard.selectedSlotIdx < 0 {
		return &stepError{"", nil}
	}

	s.slotSelect.DisableableWidget.Disable()
	s.refreshButton.DisableableWidget.Disable()

	return nil
}

func (s *slotStep) revert() {
	gWizard.selectedSlotIdx = -1
	gWizard.slotIds = nil
	gWizard.slotNames = nil
}
