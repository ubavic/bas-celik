package widgets

var clipboardHook func(string) bool

func SetClipboard(hook func(string) bool) {
	clipboardHook = hook
}

func copyToClipboard(str string) bool {
	if clipboardHook == nil {
		return false
	}

	return clipboardHook(str)
}
