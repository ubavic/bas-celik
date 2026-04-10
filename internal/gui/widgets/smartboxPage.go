package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const boxW = 300
const boxH = 50

type SmartboxPage struct {
	widget.BaseWidget
	message     string
	portMessage string
	vendors     []string
}

func NewSmartboxPage(message, portMessage string, vendors []string) *SmartboxPage {
	smartboxPage := &SmartboxPage{
		message:     message,
		portMessage: portMessage,
		vendors:     vendors,
	}

	smartboxPage.ExtendBaseWidget(smartboxPage)

	return smartboxPage
}

func (sp *SmartboxPage) CreateRenderer() fyne.WidgetRenderer {
	messageText := canvas.NewText(sp.message, theme.Color(theme.ColorNameForeground))
	messageText.TextStyle.Bold = true
	messageText.Alignment = fyne.TextAlignCenter

	portText := canvas.NewText(sp.portMessage, theme.Color(theme.ColorNameForeground))
	portText.Alignment = fyne.TextAlignCenter

	texts := make([]*canvas.Text, 0, len(sp.vendors))
	boxes := make([]*canvas.Rectangle, 0, len(sp.vendors))

	for _, v := range sp.vendors {
		rec := canvas.NewRectangle(theme.Color(theme.ColorNameButton))
		rec.Resize(fyne.NewSize(boxW, boxH))
		rec.Move(fyne.NewPos(20, 20))
		rec.CornerRadius = theme.InputRadiusSize()
		boxes = append(boxes, rec)

		text := canvas.NewText(v, theme.Color(theme.ColorNameForeground))
		text.Alignment = fyne.TextAlignCenter
		texts = append(texts, text)
	}

	return &SmartboxPageRenderer{
		page:        sp,
		portText:    portText,
		messageText: messageText,
		texts:       texts,
		boxes:       boxes,
	}
}

type SmartboxPageRenderer struct {
	page        *SmartboxPage
	messageText *canvas.Text
	portText    *canvas.Text
	boxes       []*canvas.Rectangle
	texts       []*canvas.Text
}

func (r *SmartboxPageRenderer) Refresh() {
	r.messageText.Refresh()
	r.portText.Refresh()

	for _, b := range r.boxes {
		b.Refresh()
	}

	for _, t := range r.texts {
		t.Refresh()
	}
}

func (r *SmartboxPageRenderer) Layout(s fyne.Size) {
	pad := theme.Padding()
	x, y := pad, 3*pad
	width := s.Width - 2*pad

	r.messageText.Resize(fyne.NewSize(width, r.messageText.MinSize().Height))
	r.messageText.Move(fyne.NewPos(x, y))
	y += r.messageText.MinSize().Height + pad

	r.portText.Resize(fyne.NewSize(width, r.messageText.MinSize().Height))
	r.portText.Move(fyne.NewPos(x, y))
	y += r.portText.MinSize().Height + pad*2

	vendorY := y + 6*pad
	boxX := (width - boxW) / 2
	for i := range r.boxes {
		b := r.boxes[i]
		t := r.texts[i]
		boxSize := fyne.NewSize(boxW, boxH)
		b.Resize(boxSize)
		b.Move(fyne.NewPos(boxX, vendorY))
		t.Resize(boxSize)
		t.Move(fyne.NewPos(boxX, vendorY))
		vendorY += boxH + pad
	}
}

func (r *SmartboxPageRenderer) MinSize() fyne.Size {
	return fyne.NewSize(500, 300)
}

func (r *SmartboxPageRenderer) Objects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0, 2*len(r.boxes)+2)
	objects = append(objects, r.messageText, r.portText)

	for _, b := range r.boxes {
		objects = append(objects, b)
	}

	for _, t := range r.texts {
		objects = append(objects, t)
	}

	return objects
}

func (r *SmartboxPageRenderer) Destroy() {}
