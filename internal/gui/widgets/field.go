package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Field struct {
	widget.BaseWidget
	name, value string
	minWidth    float32
	hovered     bool
	copied      bool
}

type FieldRenderer struct {
	field      *Field
	background *canvas.Rectangle
	nameText   *canvas.Text
	valueLabel *widget.Label
}

func NewField(name, value string, minWidth float32) *Field {
	field := &Field{
		name:     name,
		value:    value,
		minWidth: minWidth,
	}
	field.ExtendBaseWidget(field)
	return field
}

func (f *Field) CreateRenderer() fyne.WidgetRenderer {
	nameText := canvas.NewText(f.name, theme.Color(theme.ColorNameForeground))
	nameText.TextSize = 11

	valueText := widget.NewLabelWithStyle(f.value, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	valueText.Wrapping = fyne.TextWrapWord
	valueText.Resize(fyne.NewSize(f.minWidth, valueText.MinSize().Height))

	background := canvas.NewRectangle(color.Transparent)
	background.CornerRadius = theme.InputRadiusSize()

	return &FieldRenderer{
		field:      f,
		background: background,
		nameText:   nameText,
		valueLabel: valueText,
	}
}

func (f *Field) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (f *Field) MouseIn(*desktop.MouseEvent) {
	f.copied = false
	f.hovered = true
	f.Refresh()
}

func (f *Field) MouseMoved(*desktop.MouseEvent) {
}

func (f *Field) MouseOut() {
	f.hovered = false
	f.Refresh()
}

func (f *Field) Tapped(*fyne.PointEvent) {
	if copyToClipboard(f.value) {
		f.copied = true
	}
	f.Refresh()
}

func (r *FieldRenderer) Refresh() {
	if r.field.hovered && !r.field.copied {
		r.background.FillColor = theme.Color(theme.ColorNameButton)
	} else {
		r.background.FillColor = color.Transparent
	}
	r.valueLabel.Refresh()
	r.background.Refresh()
}

func (r *FieldRenderer) Layout(s fyne.Size) {
	r.nameText.Move(fyne.Position{X: theme.Padding(), Y: 0})
	r.valueLabel.Resize(s.SubtractWidthHeight(0, 2*theme.Padding()))
	r.valueLabel.Move(fyne.Position{X: -theme.Padding(), Y: theme.Padding()})
	r.background.Resize(s)
}

func (r *FieldRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.field.minWidth+2*theme.Padding(), r.valueLabel.MinSize().Height-theme.Padding())
}

func (r *FieldRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.valueLabel, r.nameText}
}

func (r *FieldRenderer) Destroy() {}
