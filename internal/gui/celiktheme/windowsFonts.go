package celiktheme

import "os"

func loadWidowsSystemFonts() (regular []byte, bold []byte, italic []byte, err error) {
	regular, err = os.ReadFile(`C:\Windows\Fonts\segoeui.ttf`)
	if err != nil {
		return
	}

	bold, err = os.ReadFile(`C:\Windows\Fonts\segoeuib.ttf`)
	if err != nil {
		return
	}

	italic, err = os.ReadFile(`C:\Windows\Fonts\segoeuii.ttf`)
	if err != nil {
		return
	}

	return
}
