package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Spacer struct {
	widget.BaseWidget
	minWidth float32
}

type SpacerRenderer struct {
	spacer *Spacer
}

func NewSpacer() *Spacer {
	spacer := &Spacer{}
	spacer.ExtendBaseWidget(spacer)
	return spacer
}

func (s *Spacer) SetMinWidth(width float32) {
	if width < 0 {
		width = 0
	}

	s.minWidth = width
}

func (s *Spacer) CreateRenderer() fyne.WidgetRenderer {
	return &SpacerRenderer{
		spacer: s,
	}
}

func (r *SpacerRenderer) Refresh() {}

func (r *SpacerRenderer) Layout(s fyne.Size) {}

func (r *SpacerRenderer) MinSize() fyne.Size {
	width := r.spacer.minWidth

	if r.spacer.minWidth == 0 {
		width = theme.Padding()
	}

	return fyne.NewSize(width, 5*theme.Padding())
}

func (r *SpacerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{}
}

func (r *SpacerRenderer) Destroy() {}
