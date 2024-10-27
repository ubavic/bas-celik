package localization

import "strings"

func FormatYesNo(a bool, script Language) string {
	if script == SrLatin {
		if a {
			return "Da"
		} else {
			return "Ne"
		}
	} else if script == SrCyrillic {
		if a {
			return "Да"
		} else {
			return "Не"
		}
	} else {
		if a {
			return "Yes"
		} else {
			return "No"
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
func FormatDateYMD(in *string) {
	chars := strings.Split(*in, "")
	if len(chars) != 8 {
		return
	}
	*in = chars[6] + chars[7] + "." + chars[4] + chars[5] + "." + chars[0] + chars[1] + chars[2] + chars[3]
}

// Joins list of strings into a single string
// separating them with a comma and a space.
// Empty strings are skipped.
func JoinWithComma(strs ...string) string {
	var nonemptyStrings []string
	for _, str := range strs {
		if str != "" {
			nonemptyStrings = append(nonemptyStrings, str)
		}
	}

	return strings.Join(nonemptyStrings, ", ")
}
