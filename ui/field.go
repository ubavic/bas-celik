package ui

import (
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
	field               *Field
	background          *canvas.Rectangle
	nameText, valueText *canvas.Text
}

func newField(name string, minWidth float32) *Field {
	field := &Field{
		name:     name,
		value:    "",
		minWidth: minWidth,
	}
	field.ExtendBaseWidget(field)
	return field
}

func (f *Field) CreateRenderer() fyne.WidgetRenderer {
	nameText := canvas.NewText(f.name, theme.ForegroundColor())
	nameText.TextSize = 11

	valueText := canvas.NewText(f.value, theme.ForegroundColor())
	valueText.TextStyle = fyne.TextStyle{Bold: true}

	return &FieldRenderer{
		field:      f,
		background: canvas.NewRectangle(theme.BackgroundColor()),
		nameText:   nameText,
		valueText:  valueText,
	}
}

func (r *FieldRenderer) Refresh() {
	r.valueText.Text = r.field.value
	r.background.Refresh()
	r.valueText.Refresh()
}

func (r *FieldRenderer) Layout(s fyne.Size) {
	r.nameText.Move(fyne.Position{X: theme.Padding(), Y: theme.Padding()})
	r.valueText.Move(fyne.Position{X: theme.Padding(), Y: 15 + theme.Padding()})
	r.background.Resize(s)
}

func (r *FieldRenderer) MinSize() fyne.Size {
	ts1 := fyne.MeasureText(r.nameText.Text, r.nameText.TextSize, r.nameText.TextStyle)
	return fyne.NewSize(r.field.minWidth+2*theme.Padding(), ts1.Height+ts1.Height+2*theme.Padding())
}

func (r *FieldRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.valueText, r.nameText}
}

func (r *FieldRenderer) Destroy() {}
