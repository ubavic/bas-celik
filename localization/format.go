package localization

func FormatYesNo(a bool, script Script) string {
	if script == Latin {
		if a {
			return "Da"
		} else {
			return "Ne"
		}
	} else {
		if a {
			return "Да"
		} else {
			return "Не"
		}
	}
}
