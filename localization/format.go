package localization

import "strings"

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

// Expects a pointer to a date in the format DDMMYYYY.
// Modifies, in place, date to format DD.MM.YYYY.
func FormatDate(in *string) {
	chars := strings.Split(*in, "")
	if len(chars) != 8 {
		return
	}
	if chars[4] == "0" {
		*in = "Nije dostupan"
		return
	}
	*in = chars[0] + chars[1] + "." + chars[2] + chars[3] + "." + chars[4] + chars[5] + chars[6] + chars[7] + "."
}

// Expects a pointer to a date in the format YYYYMMDD.
// Modifies, in place, date to format DD.MM.YYYY.
func FormatDate2(in *string) {
	chars := strings.Split(*in, "")
	if len(chars) != 8 {
		return
	}
	*in = chars[6] + chars[7] + "." + chars[4] + chars[5] + "." + chars[0] + chars[1] + chars[2] + chars[3]
}
