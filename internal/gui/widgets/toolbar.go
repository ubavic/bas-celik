package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/internal/gui/icon"
	"github.com/ubavic/bas-celik/v2/internal/gui/translation"
)

type Toolbar struct {
	widget.BaseWidget
	readers           []string
	onOpenAbout       func()
	onOpenPreferences func()
	onPinChange       func()
	onReaderChange    func(string)
	selectedReader    string
	pinChangeEnabled  bool
	showReaders       bool
}

type ToolbarRenderer struct {
	toolbar           *Toolbar
	aboutButton       *widget.Button
	preferencesButton *widget.Button
	pinChangeButton   *widget.Button
	container         *fyne.Container
	readersLabel      *widget.Label
	readersSelect     *widget.Select
}

func NewToolbar(onOpenAbout, onOpenPreferences, onPinChange func(), showReaders bool) *Toolbar {
	toolbar := &Toolbar{
		readers:           nil,
		onOpenAbout:       onOpenAbout,
		onOpenPreferences: onOpenPreferences,
		onPinChange:       onPinChange,
		showReaders:       showReaders,
	}

	toolbar.ExtendBaseWidget(toolbar)
	return toolbar
}

func (t *Toolbar) HookReaderChange(hook func(string)) {
	t.onReaderChange = hook
}

func (t *Toolbar) EnablePinChange() {
	t.pinChangeEnabled = true
	t.Refresh()
}

func (t *Toolbar) DisablePinChange() {
	t.pinChangeEnabled = false
	t.Refresh()
}

func (t *Toolbar) CreateRenderer() fyne.WidgetRenderer {
	label := widget.NewLabel(translation.Translate("ui.reader"))

	onChange := func(reader string) {
		if t.onReaderChange != nil {
			t.onReaderChange(reader)
		}
	}

	readersSelect := widget.NewSelect(t.readers, onChange)

	pinChangeButton := widget.NewButtonWithIcon("", icon.PinThemedResource, t.onPinChange)
	pinChangeButton.Importance = widget.LowImportance
	pinChangeButton.Disable()

	preferencesButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), t.onOpenPreferences)
	preferencesButton.Importance = widget.LowImportance

	aboutButton := widget.NewButtonWithIcon("", theme.InfoIcon(), t.onOpenAbout)
	aboutButton.Importance = widget.LowImportance

	var horizontalContainer *fyne.Container
	if t.showReaders {
		horizontalContainer = container.New(layout.NewHBoxLayout(), label, readersSelect, layout.NewSpacer(), pinChangeButton, preferencesButton, aboutButton)
	} else {
		horizontalContainer = container.New(layout.NewHBoxLayout(), layout.NewSpacer(), preferencesButton, aboutButton)
	}

	return &ToolbarRenderer{
		toolbar:           t,
		aboutButton:       aboutButton,
		preferencesButton: preferencesButton,
		pinChangeButton:   pinChangeButton,
		container:         horizontalContainer,
		readersLabel:      label,
		readersSelect:     readersSelect,
	}
}

func (r *ToolbarRenderer) Refresh() {
	if r.toolbar.showReaders {
		r.readersSelect.SetOptions(r.toolbar.readers)
		r.readersSelect.Selected = r.toolbar.selectedReader

		if len(r.toolbar.readers) <= 1 {
			r.readersSelect.Disable()
		} else {
			r.readersSelect.Enable()
		}

		r.readersSelect.Refresh()

		if r.toolbar.pinChangeEnabled {
			r.pinChangeButton.Enable()
		} else {
			r.pinChangeButton.Disable()
		}
	}

	r.aboutButton.Refresh()
}

func (r *ToolbarRenderer) Layout(s fyne.Size) {
	availableWidth := s.Width
	availableWidth -= r.aboutButton.Size().Width
	availableWidth -= r.preferencesButton.MinSize().Width
	if r.toolbar.showReaders {
		availableWidth -= r.pinChangeButton.MinSize().Width
		availableWidth -= r.readersLabel.MinSize().Width
	}
	availableWidth -= 2 * theme.InnerPadding()
	r.container.Resize(s)
	r.readersSelect.Resize(fyne.Size{Width: availableWidth, Height: s.Height})
}

func (r *ToolbarRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *ToolbarRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.aboutButton, r.preferencesButton, r.container}

	if r.toolbar.showReaders {
		objects = append(objects, r.readersSelect, r.pinChangeButton)
	}

	return objects
}

func (r *ToolbarRenderer) Destroy() {}

func (r *Toolbar) SetReaders(readers []string, selectedReader string) {
	r.readers = make([]string, len(readers))
	copy(r.readers, readers)

	r.selectedReader = selectedReader

	r.Refresh()
}
