package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	assert.Equal(t,
		"",
		Sanitize(""))
	assert.Equal(t,
		"foo bar",
		Sanitize("foo bar"))
	assert.Equal(t,
		".DS_Store",
		Sanitize(".DS_Store"))
	assert.Equal(t,
		"@eaDir",
		Sanitize("@eaDir"))
	assert.Equal(t,
		"--",
		Sanitize("–—"),
		"replace en-dash and em-dash with hyphens")
	assert.Equal(t,
		"... New.txt",
		Sanitize("… New.text"),
		"replace horizontal ellipsis")
	assert.Equal(t,
		"Your book -The Art of the Last of Us...-..eml",
		Sanitize("Your book The Art of the Last of Us.....eml"))
	assert.Equal(t,
		"The Rise and Fall of Command & Conquer - Documentary.mp4",
		Sanitize("The Rise and Fall of Command & Conquer  Documentary.mp4"))
	assert.Equal(t,
		"LPT3_changed",
		Sanitize("LPT3"))
	assert.Equal(t,
		"Text.doc",
		Sanitize("Text.doc . "))

	replacedSpecials := "!?%|$"
	for _, c := range replacedSpecials {
		assert.Equal(t, "_", Sanitize(string(c)))
	}

	keptSpecials := "@#,.-_()[]{}"
	assert.Equal(t, keptSpecials, Sanitize(keptSpecials), "keep some special characters")

	// https://en.wikipedia.org/wiki/Diacritic
	assert.Equal(t,
		"aaaaaAcccCCCdDeeeeiilnoooossSSuuuzzzZZ",
		Sanitize("ạàąâåÅčćçÇČĆđĐęéèêîìłńóôộớšśŚŠùûůżźžŻŽ"),
		"remove diacritics from a-z and A-Z")

	// Example song filenames in the wild
	assert.Equal(t,
		"1-03 Urtuemlicher Titan .mp3",
		Sanitize("1-03 Urtümlicher Titan ....mp3"))
	assert.Equal(t,
		"4-10 Ein Portraet ueber Torquemada.mp3",
		Sanitize("4-10 Ein Porträt über Torquemada.mp3"))
	// This test string uses U+0308 Combining Diaeresis for the Ü and ä, even
	// though you can't "see" this easily. Try moving your cursor through the
	// string from left to right, you'll notice that it takes two key presses
	assert.Equal(t,
		"Lovecraft Ueber _Ein Portraet Torquemadas_.mp3",
		Sanitize("Lovecraft Über _Ein Porträt Torquemadas_.mp3"))
}
