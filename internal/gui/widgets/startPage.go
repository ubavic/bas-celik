package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type StartPage struct {
	widget.BaseWidget
	status      string
	explanation string
	err         bool
}

type StartPageRenderer struct {
	page            *StartPage
	statusText      *canvas.Text
	explanationText *canvas.Text
	container       *fyne.Container
}

func NewStartPage() *StartPage {
	statusBar := &StartPage{
		status:      "",
		explanation: "",
		err:         false,
	}
	statusBar.ExtendBaseWidget(statusBar)
	return statusBar
}

func (sb *StartPage) SetStatus(status, explanation string, err bool) {
	sb.status = status
	sb.explanation = explanation
	sb.err = err
}

func (sb *StartPage) CreateRenderer() fyne.WidgetRenderer {
	statusText := canvas.NewText(sb.status, theme.Color(theme.ColorNameForeground))
	statusText.TextSize = 16
	statusText.Color = theme.Color(theme.ColorNameError)

	explanationText := canvas.NewText(sb.status, theme.Color(theme.ColorNameForeground))
	explanationText.TextSize = 11
	explanationText.Color = theme.Color(theme.ColorNameForeground)

	box := container.New(layout.NewVBoxLayout(), statusText, explanationText)
	container := container.New(layout.NewCenterLayout(), box)

	return &StartPageRenderer{
		page:            sb,
		statusText:      statusText,
		explanationText: explanationText,
		container:       container,
	}
}

func (r *StartPageRenderer) Refresh() {
	r.statusText.Text = r.page.status
	r.explanationText.Text = r.page.explanation

	if r.page.err {
		r.statusText.Color = theme.Color(theme.ColorNameError)
	} else {
		r.statusText.Color = theme.Color(theme.ColorNameForeground)
	}

	r.statusText.Refresh()
	r.explanationText.Refresh()
}

func (r *StartPageRenderer) Layout(s fyne.Size) {
	r.container.Resize(s)
}

func (r *StartPageRenderer) MinSize() fyne.Size {
	return fyne.NewSize(500, 300)
}

func (r *StartPageRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

func (r *StartPageRenderer) Destroy() {}
