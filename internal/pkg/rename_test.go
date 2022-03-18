package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateName(t *testing.T) {
	var s string

	// directories
	s, _ = truncateName("1234567890", true, 5)
	assert.Equal(t, "12345", s)

	s, _ = truncateName("1234567890", true, 1)
	assert.Equal(t, "1", s)

	_, err := truncateName("1234567890", true, 0)
	assert.Error(t, err)

	// files
	s, _ = truncateName("1234567890", false, 5)
	assert.Equal(t, "12345", s)
	s, _ = truncateName("1234567890.txt", false, 5)
	assert.Equal(t, "1.txt", s)
	s, _ = truncateName("1234567890.abcdefghij", false, 12)
	assert.Equal(t, "1.abcdefghij", s)
	s, _ = truncateName("1234567890.abcdefghij", false, 11)
	assert.Equal(t, ".abcdefghij", s)

	_, err2 := truncateName("1234567890", false, 0)
	assert.Error(t, err2)
	_, err3 := truncateName("1234567890.abcdefghij", false, 10)
	assert.Error(t, err3)
}
