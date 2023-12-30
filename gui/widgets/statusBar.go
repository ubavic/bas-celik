package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type StatusBar struct {
	widget.BaseWidget
	status    string
	apiStatus string
	err       bool
}

type StatusBarRenderer struct {
	bar        *StatusBar
	statusText *canvas.Text
	apiText    *canvas.Text
}

func NewStatusBar() *StatusBar {
	statusBar := &StatusBar{
		status:    "",
		apiStatus: "",
		err:       true,
	}
	statusBar.ExtendBaseWidget(statusBar)
	return statusBar
}

func (sb *StatusBar) SetStatus(status string, err bool) {
	sb.status = status
	sb.err = err
}

func (sb *StatusBar) SetApiStatus(status string) {
	sb.apiStatus = status
}

func (sb *StatusBar) CreateRenderer() fyne.WidgetRenderer {
	statusText := canvas.NewText(sb.status, theme.ForegroundColor())
	statusText.TextSize = 11
	statusText.Color = theme.ErrorColor()

	apiText := canvas.NewText(sb.apiStatus, theme.ForegroundColor())
	apiText.TextSize = 11

	return &StatusBarRenderer{
		bar:        sb,
		statusText: statusText,
		apiText:    apiText,
	}
}

func (r *StatusBarRenderer) Refresh() {
	r.statusText.Text = r.bar.status
	r.apiText.Text = r.bar.apiStatus

	if r.bar.err {
		r.statusText.Color = theme.ErrorColor()
	} else {
		r.statusText.Color = theme.ForegroundColor()
	}

	r.statusText.Refresh()
	r.apiText.Refresh()
}

func (r *StatusBarRenderer) Layout(s fyne.Size) {
	ts1 := fyne.MeasureText(r.statusText.Text, r.statusText.TextSize, r.statusText.TextStyle)
	r.statusText.Move(fyne.Position{X: theme.Padding(), Y: 0})
	r.apiText.Move(fyne.Position{X: theme.Padding(), Y: theme.Padding() + ts1.Height})
}

func (r *StatusBarRenderer) MinSize() fyne.Size {
	ts1 := fyne.MeasureText(r.statusText.Text, r.statusText.TextSize, r.statusText.TextStyle)
	ts2 := fyne.MeasureText(r.apiText.Text, r.apiText.TextSize, r.apiText.TextStyle)
	return fyne.NewSize(fyne.Max(ts1.Width, ts2.Width)+theme.Padding(), ts1.Height+ts2.Height+theme.Padding())
}

func (r *StatusBarRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.apiText, r.statusText}
}

func (r *StatusBarRenderer) Destroy() {}
