package localization

var cyrillic = []rune{
	'А', 'Б', 'В', 'Г', 'Д', 'Ђ', 'Е', 'Ж', 'З', 'И', 'Ј', 'К', 'Л', 'Љ', 'М', 'Н', 'Њ', 'О', 'П', 'Р', 'С', 'Т', 'Ћ', 'У', 'Ф', 'Х', 'Ц', 'Ч', 'Џ', 'Ш',
	'а', 'б', 'в', 'г', 'д', 'ђ', 'е', 'ж', 'з', 'и', 'ј', 'к', 'л', 'љ', 'м', 'н', 'њ', 'о', 'п', 'р', 'с', 'т', 'ћ', 'у', 'ф', 'х', 'ц', 'ч', 'џ', 'ш',
}

var latin = []string{
	"A", "B", "V", "G", "D", "Đ", "E", "Ž", "Z", "I", "J", "K", "L", "Lj", "M", "N", "Nj", "O", "P", "R", "S", "T", "Ć", "U", "F", "H", "C", "Č", "Dž", "Š",
	"a", "b", "v", "g", "d", "đ", "e", "ž", "z", "i", "j", "k", "l", "lj", "m", "n", "nj", "o", "p", "r", "s", "t", "ć", "u", "f", "h", "c", "č", "dž", "š",
}

func CyrillicToLatin(str string) string {
	out := ""

	for _, r := range str {
		index := indexOfCyrillic(r)
		if index >= 0 {
			out += latin[index]
		} else {
			out += string(r)
		}
	}

	return out
}

func indexOfCyrillic(r rune) int {
	for i, v := range cyrillic {
		if v == r {
			return i
		}
	}
	return -1
}
