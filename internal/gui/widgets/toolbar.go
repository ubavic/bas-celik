package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/v2/internal/gui/translation"
)

type Toolbar struct {
	widget.BaseWidget
	readers           []string
	onOpenAbout       func()
	onOpenPreferences func()
	onReaderChange    func(string)
	selectedReader    string
	showReaders       bool
}

type ToolbarRenderer struct {
	toolbar           *Toolbar
	aboutButton       *widget.Button
	preferencesButton *widget.Button
	container         *fyne.Container
	readersLabel      *widget.Label
	readersSelect     *widget.Select
}

func NewToolbar(onOpenAbout, onOpenPreferences func(), showReaders bool) *Toolbar {
	toolbar := &Toolbar{
		readers:           nil,
		onOpenAbout:       onOpenAbout,
		onOpenPreferences: onOpenPreferences,
		showReaders:       showReaders,
	}

	toolbar.ExtendBaseWidget(toolbar)
	return toolbar
}

func (t *Toolbar) HookReaderChange(hook func(string)) {
	t.onReaderChange = hook
}

func (t *Toolbar) CreateRenderer() fyne.WidgetRenderer {
	label := widget.NewLabel(translation.Translate("ui.reader"))

	onChange := func(reader string) {
		if t.onReaderChange != nil {
			t.onReaderChange(reader)
		}
	}

	readersSelect := widget.NewSelect(t.readers, onChange)

	preferencesButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), t.onOpenPreferences)
	preferencesButton.Importance = widget.LowImportance

	aboutButton := widget.NewButtonWithIcon("", theme.InfoIcon(), t.onOpenAbout)
	aboutButton.Importance = widget.LowImportance

	var horizontalContainer *fyne.Container
	if t.showReaders {
		horizontalContainer = container.New(layout.NewHBoxLayout(), label, readersSelect, layout.NewSpacer(), preferencesButton, aboutButton)
	} else {
		horizontalContainer = container.New(layout.NewHBoxLayout(), layout.NewSpacer(), preferencesButton, aboutButton)
	}

	return &ToolbarRenderer{
		toolbar:           t,
		aboutButton:       aboutButton,
		preferencesButton: preferencesButton,
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
	}

	r.aboutButton.Refresh()
}

func (r *ToolbarRenderer) Layout(s fyne.Size) {
	availableWidth := s.Width
	availableWidth -= r.aboutButton.Size().Width
	availableWidth -= r.preferencesButton.MinSize().Width
	if r.toolbar.showReaders {
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
		objects = append(objects, r.readersSelect)
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
