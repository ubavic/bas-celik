package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UnknownCardPage struct {
	widget.BaseWidget
	status    string
	atr       string
	copyLabel string
}

type UnknownCardPageRenderer struct {
	page         *UnknownCardPage
	statusText   *canvas.Text
	atrText      *canvas.Text
	copyAtrLabel *copyAtrLabel
	container    *fyne.Container
}

func NewUnknownCardPage() *UnknownCardPage {
	page := &UnknownCardPage{
		status:    "",
		atr:       "",
		copyLabel: "",
	}
	page.ExtendBaseWidget(page)
	return page
}

func (p *UnknownCardPage) SetContent(status, atr, copyLabel string) {
	p.status = status
	p.atr = atr
	p.copyLabel = copyLabel
}

func (p *UnknownCardPage) CreateRenderer() fyne.WidgetRenderer {
	statusText := canvas.NewText(p.status, theme.Color(theme.ColorNameError))
	statusText.TextSize = 16
	statusText.Color = theme.Color(theme.ColorNameError)

	atrText := canvas.NewText(atrLine(p.atr), theme.Color(theme.ColorNameForeground))
	atrText.TextSize = 12
	atrText.Color = theme.Color(theme.ColorNameForeground)

	copyAtr := newCopyAtrLabel(p.copyLabel, func() {
		copyToClipboard(p.atr)
	})

	box := container.New(layout.NewVBoxLayout(), statusText, atrText, copyAtr)
	c := container.New(layout.NewCenterLayout(), box)

	return &UnknownCardPageRenderer{
		page:         p,
		statusText:   statusText,
		atrText:      atrText,
		copyAtrLabel: copyAtr,
		container:    c,
	}
}

func (r *UnknownCardPageRenderer) Refresh() {
	r.statusText.Text = r.page.status
	r.atrText.Text = atrLine(r.page.atr)
	r.copyAtrLabel.SetText(r.page.copyLabel)

	r.statusText.Color = theme.Color(theme.ColorNameError)

	fyne.Do(func() {
		r.statusText.Refresh()
		r.atrText.Refresh()
		r.copyAtrLabel.Refresh()
	})
}

func (r *UnknownCardPageRenderer) Layout(s fyne.Size) {
	r.container.Resize(s)
}

func (r *UnknownCardPageRenderer) MinSize() fyne.Size {
	return fyne.NewSize(500, 300)
}

func (r *UnknownCardPageRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}

func (r *UnknownCardPageRenderer) Destroy() {}

func atrLine(atr string) string {
	return "ATR: " + atr
}

type copyAtrLabel struct {
	widget.BaseWidget
	text     string
	onTapped func()
}

type copyAtrLabelRenderer struct {
	label *copyAtrLabel
	text  *canvas.Text
}

func newCopyAtrLabel(text string, onTapped func()) *copyAtrLabel {
	label := &copyAtrLabel{
		text:     text,
		onTapped: onTapped,
	}
	label.ExtendBaseWidget(label)
	return label
}

func (l *copyAtrLabel) SetText(text string) {
	l.text = text
}

func (l *copyAtrLabel) Tapped(*fyne.PointEvent) {
	if l.onTapped != nil {
		l.onTapped()
	}
}

func (l *copyAtrLabel) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (l *copyAtrLabel) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(l.text, theme.Color(theme.ColorNameHyperlink))
	text.TextSize = 12
	return &copyAtrLabelRenderer{
		label: l,
		text:  text,
	}
}

func (r *copyAtrLabelRenderer) Refresh() {
	r.text.Text = r.label.text
	r.text.Color = theme.Color(theme.ColorNameHyperlink)
	r.text.Refresh()
}

func (r *copyAtrLabelRenderer) Layout(s fyne.Size) {
	r.text.Resize(s)
}

func (r *copyAtrLabelRenderer) MinSize() fyne.Size {
	return r.text.MinSize()
}

func (r *copyAtrLabelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.text}
}

func (r *copyAtrLabelRenderer) Destroy() {}
