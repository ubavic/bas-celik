package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Group struct {
	widget.BaseWidget
	name    string
	objects []fyne.CanvasObject
}

type GroupRenderer struct {
	group      *Group
	nameText   *canvas.Text
	column     *fyne.Container
	baseHeight float32
}

func NewGroup(name string, objects ...fyne.CanvasObject) *Group {
	group := &Group{
		name:    name,
		objects: objects,
	}
	group.ExtendBaseWidget(group)
	return group
}

func (g *Group) CreateRenderer() fyne.WidgetRenderer {
	nameText := canvas.NewText(g.name, theme.Color(theme.ColorNameForeground))
	nameText.TextStyle.Bold = true
	nameText.TextSize = 14

	baseHeight := 3 * theme.Padding()
	if len(g.objects) > 0 {
		baseHeight = g.objects[0].MinSize().Height
	}

	column := container.New(layout.NewVBoxLayout(), g.objects...)

	return &GroupRenderer{
		group:      g,
		nameText:   nameText,
		column:     column,
		baseHeight: baseHeight,
	}
}

func (r *GroupRenderer) Refresh() {
	r.column.Refresh()
	r.nameText.Refresh()
}

func (r *GroupRenderer) Layout(s fyne.Size) {
	textHeight := r.nameText.MinSize().Height
	r.nameText.Move(fyne.NewPos(2*theme.Padding(), r.baseHeight-textHeight))
	r.column.Move(fyne.Position{X: theme.Padding(), Y: r.baseHeight + theme.Padding()})
	r.column.Layout.Layout(r.group.objects, s.SubtractWidthHeight(2*theme.Padding(), 0))
}

func (r *GroupRenderer) MinSize() fyne.Size {
	width := r.column.MinSize().Width + 2*theme.Padding()
	height := r.column.MinSize().Height + r.baseHeight + theme.Padding()
	return fyne.NewSize(width, height)
}

func (r *GroupRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.column, r.nameText}
}

func (r *GroupRenderer) Destroy() {}
