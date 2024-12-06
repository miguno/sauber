package internal

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Sanitize(filename string) string {
	s := replacer.Replace(filename)
	s = removeAccents(s)
	s = replaceInvisibleCharsWithHyphen(s)
	s = replacePrivateUseCharsWithHyphen(s)
	return s
}

var replacer = strings.NewReplacer(
	"ß", "ss",
	"Ä", "Ae",
	"Ö", "Oe",
	"Ü", "Ue",
	"ä", "ae",
	"ö", "oe",
	"ü", "ue",
	// U+0308 Combining Diaeresis (https://en.wikipedia.org/wiki/Diaeresis_(diacritic)
	// This character is a Non-spacing Mark and inherits its script property
	// from the preceding character. The character is also known as
	// 'double dot above', 'umlaut', 'Greek dialytika', and 'double derivative'.
	//
	// Note how the 'Ä' below is actually two chars: an 'A' followed by U+0308,
	// i.e., the two dots to be added on top of the 'A'. Try it yourself: put
	// your cursor to the left of the 'Ä', then use the arrow keys on your
	// keyboard to move the cursor to the right. You will notice that you need
	// two key presses to get across 'Ä'. (This may not work in all text
	// editors, such as neovim. It does work in GoLand IDE, for example.)
	"Ä", "Ae", // A<0308> (see above)
	"Ö", "Oe", // O<0308> (see above)
	"Ü", "Ue", // U<0308> (see above)
	"ä", "ae", // a<0308> (see above)
	"ö", "oe", // o<0308> (see above)
	"ü", "ue", // u<0308> (see above)
	// a few hardcoded rules to make repeated `.` and `,` more pleasant
	",,,,,", ",",
	",,,,", ",",
	",,,", ",",
	",,", ",",
	".....", ".",
	"....", ".",
	"…", "...", // horizontal ellipsis
)

func removeAccents(s string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		runes.Map(func(r rune) rune {
			switch r {
			case 'ą':
				return 'a'
			case 'ć':
				return 'c'
			case 'đ':
				return 'd'
			case 'Đ':
				return 'D'
			case 'ę':
				return 'e'
			case 'ł':
				return 'l'
			case 'ń':
				return 'n'
			case 'ó':
				return 'o'
			case 'ś':
				return 's'
			case 'ż':
				return 'z'
			case 'ź':
				return 'z'
			case '%':
				return '_'
			case '?':
				return '_'
			case '!':
				return '_'
			case '|':
				return '_'
			case '$':
				return '_'
			case '–': // en dash
				return '-' // hyphen
			case '—': // em dash
				return '-' // hyphen
			}
			return r
		}),
		norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func replaceInvisibleCharsWithHyphen(str string) string {
	// `\p{Cf}`: invisible formatting indicator
	reg := regexp.MustCompile(`\p{Cf}`)
	return reg.ReplaceAllString(str, "-")
}

func replacePrivateUseCharsWithHyphen(str string) string {
	// `\p{Co}`: any code point reserved for private use
	reg := regexp.MustCompile(`\p{Co}`)
	return reg.ReplaceAllString(str, "-")
}
