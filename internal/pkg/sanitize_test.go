package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	assert.Equal(t, "", Sanitize(""))
	assert.Equal(t, "foo bar", Sanitize("foo bar"))
	assert.Equal(t, ".DS_Store", Sanitize(".DS_Store"))
	assert.Equal(t, "@eaDir", Sanitize("@eaDir"))
	assert.Equal(t, "--", Sanitize("–—"), "replace en-dash and em-dash with hyphens")
	assert.Equal(t, "...", Sanitize("…"), "replace horizontal ellipsis")
	assert.Equal(t, "x-...- x", Sanitize("x... x"), "replace private use characters with hyphens")
	// \u200D : zero width joiner (ZWJ)
	// \u200E : left-to-right mark (LRM)
	// \u200F : right-to-left mark (RLM)
	// \u2060 : word joiner
	assert.Equal(t, "----", Sanitize("\u200D\u200E\u200F\u2060"), "replace invisible characters with hyphens")

	replacedSpecials := "!?%|$"
	for _, c := range replacedSpecials {
		assert.Equal(t, "_", Sanitize(string(c)))
	}

	keptSpecials := "@#,.-_()[]{}"
	assert.Equal(t, keptSpecials, Sanitize(keptSpecials), "keep some special characters")

	// https://en.wikipedia.org/wiki/Diacritic
	assert.Equal(t,
		"aaaaaAcccCCCdDeeeeiilnoooossSSuuuzzzZZ",
		Sanitize("ạàąâåÅčćçÇČĆđĐęéèêîìłńóôộớšśŚŠùûůżźžŻŽ"), "remove diacritics from a-z and A-Z")

	// Example song filenames in the wild
	assert.Equal(t, "1-03 Urtuemlicher Titan .mp3", Sanitize("1-03 Urtümlicher Titan ....mp3"))
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
