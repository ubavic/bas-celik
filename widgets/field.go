package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Field struct {
	widget.BaseWidget
	name, value string
	minWidth    float32
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
	nameText := canvas.NewText(f.name, theme.ForegroundColor())
	nameText.TextSize = 11

	valueText := widget.NewLabelWithStyle(f.value, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	valueText.Wrapping = fyne.TextWrapWord
	valueText.Resize(fyne.NewSize(f.minWidth, valueText.MinSize().Height))

	return &FieldRenderer{
		field:      f,
		background: canvas.NewRectangle(color.RGBA{A: 0x00}),
		nameText:   nameText,
		valueLabel: valueText,
	}
}

func (r *FieldRenderer) Refresh() {
	r.background.Refresh()
	r.valueLabel.Refresh()
}

func (r *FieldRenderer) Layout(s fyne.Size) {
	r.nameText.Move(fyne.Position{X: theme.Padding(), Y: theme.Padding()})
	r.valueLabel.Move(fyne.Position{X: -theme.Padding(), Y: 2 * theme.Padding()})
	r.background.Resize(s)
}

func (r *FieldRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.field.minWidth+2*theme.Padding(), r.valueLabel.MinSize().Height)
}

func (r *FieldRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.valueLabel, r.nameText}
}

func (r *FieldRenderer) Destroy() {}
