package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Spacer struct {
	widget.BaseWidget
	width  float32
	height float32
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
	s.width = max(0, width)
}

func (s *Spacer) SetHeight(height float32) {
	s.height = max(0, height)
}

func (s *Spacer) CreateRenderer() fyne.WidgetRenderer {
	return &SpacerRenderer{
		spacer: s,
	}
}

func (r *SpacerRenderer) Refresh() {}

func (r *SpacerRenderer) Layout(s fyne.Size) {}

func (r *SpacerRenderer) MinSize() fyne.Size {
	width := max(r.spacer.width, theme.Padding())
	height := max(r.spacer.height, theme.Padding())

	return fyne.NewSize(width, height)
}

func (r *SpacerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{}
}

func (r *SpacerRenderer) Destroy() {}
