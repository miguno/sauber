package internal

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Sanitize(filename string) string {
	s := replacer.Replace(filename)
	s = removeAccents(s)
	s = replaceNonPrintableChars(s)
	s = replacePrivateUnicodeChars(s)
	s = removeTrailingSpaceOrDot(s)
	s = changeReservedFilenames(s)
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
	// illegal file name characters for Windows / Mac / Linux
	"<", "-",
	">", "-",
	":", "-",
	"\"", "-",
	"/", "-",
	"\\", "-",
	"|", "-",
	"?", "-",
	"*", "-",
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

// This function uses a regular expression to match any control characters,
// which are the non-printable characters in ASCII and Unicode.
// The [:cntrl:] character class matches all control characters.
// The ReplaceAllString method of the regexp package is used to replace
// all matches with an empty string, effectively removing them from the string.
func replaceNonPrintableChars(filename string) string {
	reg := regexp.MustCompile("[[:cntrl:]]")
	return reg.ReplaceAllString(filename, "-")
}

// This function uses a regular expression to match any character that
// is a private use character (\\p{Co}) or a invisible formatting indicator
// (\\p{Cf}) or any code point to which no character has been assigned (\\p{Cn}).
func replacePrivateUnicodeChars(str string) string {
	reg := regexp.MustCompile("\\p{Co}|\\p{Cf}")
	return reg.ReplaceAllString(str, "-")
}

// If the given filename is a reserved filename, the function adds "_changed" to
// the original filename and returns the new filename.
// If the given filename is not a reserved filename,
// the function returns the original filename.
func changeReservedFilenames(filename string) string {
	if isReservedFilename(strings.ToUpper(filename)) {
		return filename + "_changed"
	}
	return filename
}

// This function first checks if the given filename is one of the reserved
// filenames by comparing it to a list of reserved filenames.
// The strings.EqualFold method is used to perform a case-insensitive
// comparison between the reserved filenames and the given filename.
func isReservedFilename(filename string) bool {
	reservedFilenames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "LPT1", "LPT2", "LPT3", "LOCK$"}
	for _, reservedFilename := range reservedFilenames {
		if strings.EqualFold(reservedFilename, filename) {
			return true
		}
	}
	return false
}

// This function uses a loop to repeatedly remove spaces and dots
// from the end of the filename until there are no more spaces or dots left.
// The strings.TrimSuffix method is called twice in each iteration to
// remove both spaces and dots.
// If the given filename does not end with a space or a dot,
// the function simply returns the original filename.
func removeTrailingSpaceOrDot(filename string) string {
	for strings.HasSuffix(filename, " ") || strings.HasSuffix(filename, ".") {
		filename = strings.TrimSuffix(filename, " ")
		filename = strings.TrimSuffix(filename, ".")
	}
	return filename
}
