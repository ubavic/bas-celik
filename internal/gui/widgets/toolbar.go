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
	selectedReader    string
	readerChanged     bool
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

	onChange := func(s string) {
		if s != t.selectedReader {
			t.selectedReader = s
			t.readerChanged = true
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

	if len(r.toolbar.readers) == 0 {
		r.readersSelect.Selected = ""
		r.readersSelect.PlaceHolder = "Nema"
		r.readersSelect.Disable()
	} else if len(r.toolbar.readers) == 1 {
		r.readersSelect.Selected = r.toolbar.readers[0]
		r.readersSelect.Disable()
	} else {
		r.readersSelect.Enable()
	}

	if r.toolbar.pinChangeEnabled {
		r.pinChangeButton.Enable()
	} else {
		r.pinChangeButton.Disable()
	}

	if r.readersSelect.Selected == "" && len(r.toolbar.readers) > 0 {
		r.toolbar.selectedReader = r.toolbar.readers[0]
		r.readersSelect.Selected = r.toolbar.readers[0]
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

func (r *Toolbar) SetReaders(readers []string) {
	selectFirstReader := false

	if len(r.readers) == 0 && len(readers) > 0 {
		selectFirstReader = true
	}

	r.readers = make([]string, len(readers))
	copy(r.readers, readers)

	if selectFirstReader {
		r.selectedReader = r.readers[0]
		r.readerChanged = true
	}

	r.Refresh()
}

func (r *Toolbar) ReaderChanged() bool {
	if r.readerChanged {
		r.readerChanged = false
		return true
	}

	return false
}

func (r *Toolbar) GetReaderName() string {
	return r.selectedReader
}
