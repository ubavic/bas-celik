package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Toolbar struct {
	widget.BaseWidget
	readers           []string
	onOpenAbout       func()
	onOpenPreferences func()
	selectedReader    string
	readerChanged     bool
}

type ToolbarRenderer struct {
	toolbar       *Toolbar
	aboutButton   *widget.Button
	container     *fyne.Container
	readersLabel  *widget.Label
	readersSelect *widget.Select
}

func NewToolbar(onOpenAbout, onOpenPreferences func()) *Toolbar {
	toolbar := &Toolbar{
		readers:           nil,
		onOpenAbout:       onOpenAbout,
		onOpenPreferences: onOpenPreferences,
	}
	toolbar.ExtendBaseWidget(toolbar)
	return toolbar
}

func (t *Toolbar) CreateRenderer() fyne.WidgetRenderer {
	label := widget.NewLabel("Čitač:")

	onChange := func(s string) {
		if s != t.selectedReader {
			t.selectedReader = s
			t.readerChanged = true
		}
	}
	readersSelect := widget.NewSelect(t.readers, onChange)

	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), t.onOpenPreferences)
	settingsButton.Importance = widget.LowImportance

	aboutButton := widget.NewButtonWithIcon("", theme.InfoIcon(), t.onOpenAbout)
	aboutButton.Importance = widget.LowImportance

	container := container.New(layout.NewHBoxLayout(), label, readersSelect, layout.NewSpacer(), settingsButton, aboutButton)

	return &ToolbarRenderer{
		toolbar:       t,
		aboutButton:   aboutButton,
		container:     container,
		readersLabel:  label,
		readersSelect: readersSelect,
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

	if r.readersSelect.Selected == "" && len(r.toolbar.readers) > 0 {
		r.toolbar.selectedReader = r.toolbar.readers[0]
		r.readersSelect.Selected = r.toolbar.readers[0]
	}

	r.readersSelect.Refresh()
	r.aboutButton.Refresh()
}

func (r *ToolbarRenderer) Layout(s fyne.Size) {
	r.container.Resize(s)
}

func (r *ToolbarRenderer) MinSize() fyne.Size {
	return r.container.MinSize()
}

func (r *ToolbarRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.readersLabel, r.readersSelect, r.aboutButton, r.container}
}

func (r *ToolbarRenderer) Destroy() {}

func (r *Toolbar) SetReaders(readers []string) {
	r.readers = make([]string, len(readers))
	copy(r.readers, readers)
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
