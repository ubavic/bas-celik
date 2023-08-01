package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type StatusBar struct {
	widget.BaseWidget
	status string
	err    bool
}

type StatusBarRenderer struct {
	bar        *StatusBar
	statusText *canvas.Text
}

func NewStatusBar() *StatusBar {
	statusBar := &StatusBar{
		status: "",
		err:    true,
	}
	statusBar.ExtendBaseWidget(statusBar)
	return statusBar
}

func (sb *StatusBar) SetStatus(status string, err bool) {
	sb.status = status
	sb.err = err
}

func (sb *StatusBar) CreateRenderer() fyne.WidgetRenderer {
	statusText := canvas.NewText(sb.status, theme.ForegroundColor())
	statusText.TextSize = 11
	statusText.Color = theme.ErrorColor()

	return &StatusBarRenderer{
		bar:        sb,
		statusText: statusText,
	}
}

func (r *StatusBarRenderer) Refresh() {
	r.statusText.Text = r.bar.status

	if r.bar.err {
		r.statusText.Color = theme.ErrorColor()
	} else {
		r.statusText.Color = theme.ForegroundColor()
	}

	r.statusText.Refresh()
}

func (r *StatusBarRenderer) Layout(s fyne.Size) {
	r.statusText.Move(fyne.Position{X: theme.Padding(), Y: 3 * theme.Padding()})
}

func (r *StatusBarRenderer) MinSize() fyne.Size {
	ts1 := fyne.MeasureText(r.statusText.Text, r.statusText.TextSize, r.statusText.TextStyle)
	return fyne.NewSize(ts1.Width+2*theme.Padding(), ts1.Height)
}

func (r *StatusBarRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.statusText}
}

func (r *StatusBarRenderer) Destroy() {}
