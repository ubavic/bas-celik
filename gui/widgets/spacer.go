package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Spacer struct {
	widget.BaseWidget
}

type SpacerRenderer struct {
	spacer *Spacer
}

func NewSpacer() *Spacer {
	spacer := &Spacer{}
	spacer.ExtendBaseWidget(spacer)
	return spacer
}

func (s *Spacer) CreateRenderer() fyne.WidgetRenderer {
	return &SpacerRenderer{
		spacer: s,
	}
}

func (r *SpacerRenderer) Refresh() {}

func (r *SpacerRenderer) Layout(s fyne.Size) {}

func (r *SpacerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(theme.Padding(), 5*theme.Padding())
}

func (r *SpacerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{}
}

func (r *SpacerRenderer) Destroy() {}
