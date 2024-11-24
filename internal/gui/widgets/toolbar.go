package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/internal/gui/translation"
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

func NewToolbar(onOpenAbout, onOpenPreferences, onPinChange func()) *Toolbar {
	toolbar := &Toolbar{
		readers:           nil,
		onOpenAbout:       onOpenAbout,
		onOpenPreferences: onOpenPreferences,
		onPinChange:       onPinChange,
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

	pinChangeButton := widget.NewButtonWithIcon("", theme.VisibilityOffIcon(), t.onPinChange)
	pinChangeButton.Importance = widget.LowImportance
	pinChangeButton.Disable()

	preferencesButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), t.onOpenPreferences)
	preferencesButton.Importance = widget.LowImportance

	aboutButton := widget.NewButtonWithIcon("", theme.InfoIcon(), t.onOpenAbout)
	aboutButton.Importance = widget.LowImportance

	container := container.New(layout.NewHBoxLayout(), label, readersSelect, layout.NewSpacer(), pinChangeButton, preferencesButton, aboutButton)

	return &ToolbarRenderer{
		toolbar:           t,
		aboutButton:       aboutButton,
		preferencesButton: preferencesButton,
		pinChangeButton:   pinChangeButton,
		container:         container,
		readersLabel:      label,
		readersSelect:     readersSelect,
	}
}

func (r *ToolbarRenderer) Refresh() {
	r.readersSelect.SetOptions(r.toolbar.readers)
	r.readersSelect.Selected = r.toolbar.selectedReader

	if len(r.toolbar.readers) <= 1 {
		r.readersSelect.Disable()
	} else {
		r.readersSelect.Enable()
	}

	if r.toolbar.pinChangeEnabled {
		r.pinChangeButton.Enable()
	} else {
		r.pinChangeButton.Disable()
	}

	r.readersSelect.Refresh()
	r.aboutButton.Refresh()
}

func (r *ToolbarRenderer) Layout(s fyne.Size) {
	availableWidth := s.Width
	availableWidth -= r.aboutButton.Size().Width
	availableWidth -= r.preferencesButton.MinSize().Width
	availableWidth -= r.pinChangeButton.MinSize().Width
	availableWidth -= r.readersLabel.MinSize().Width
	availableWidth -= 2 * theme.InnerPadding()
	r.container.Resize(s)
	r.readersSelect.Resize(fyne.Size{Width: availableWidth, Height: s.Height})
}

func (r *ToolbarRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *ToolbarRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.readersLabel, r.readersSelect, r.aboutButton, r.container}
}

func (r *ToolbarRenderer) Destroy() {}

func (r *Toolbar) SetReaders(readers []string, selectedReader string) {
	r.readers = make([]string, len(readers))
	copy(r.readers, readers)

	r.selectedReader = selectedReader

	r.Refresh()
}
