package gui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ubavic/bas-celik/internal/gui/widgets"
)

func showAboutBox(win fyne.Window, version string) func() {
	verLabel := widget.NewLabelWithStyle(t("about.version")+": "+version, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	moreLabel := widget.NewLabel(t("about.moreAboutProgram"))
	url, _ := url.Parse("https://github.com/ubavic/bas-celik")
	linkLabel := widget.NewHyperlink("github.com/ubavic/bas-celik", url)
	spacer := widgets.NewSpacer()
	hBox := container.NewHBox(moreLabel, linkLabel)
	vBox := container.NewVBox(verLabel, hBox, spacer)

	return func() {
		dialog.ShowCustom(
			t("about.title"),
			t("about.close"),
			vBox,
			win,
		)
	}
}
