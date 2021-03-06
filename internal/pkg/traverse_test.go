package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldErrorWhenRootDoesNotExist(t *testing.T) {
	rootNode, err := Find("/thisPathShouldNeverExistOnTheDeviceThatRunsTheTests",
		DefaultSkipDirectories)
	assert.Empty(t, rootNode)
	assert.Error(t, err)
}

func TestShouldFindSingleFile(t *testing.T) {
	rootNode, _ := Find("../../test/traverse/root-single-file",
		DefaultSkipDirectories)
	expected := []string{"../../test/traverse/root-single-file"}
	assert.Equal(t, expected, (*rootNode).PathsDecorated())
}

func TestBasicFind(t *testing.T) {
	rootNode, _ := Find("../../test/traverse/root-basic/.",
		DefaultSkipDirectories)
	expected := []string{
		"../../test/traverse/root-basic[d]",
		"../../test/traverse/root-basic/Foo!Bar?Lorem#[d]",
		"../../test/traverse/root-basic/Foo!Bar?Lorem#/intro.mp3",
		"../../test/traverse/root-basic/Größe.mp3",
		"../../test/traverse/root-basic/Urtümlich[d]",
		"../../test/traverse/root-basic/Urtümlich/Ähnliche",
		"../../test/traverse/root-basic/foo[d]",
		"../../test/traverse/root-basic/foo/README.md",
	}
	assert.Equal(t, expected, (*rootNode).PathsDecorated())
}

func TestSkipPath(t *testing.T) {
	assert.Equal(t, false, skipPath("/bar", DefaultSkipDirectories))
	assert.Equal(t, false, skipPath("foo/bar", DefaultSkipDirectories))
	assert.Equal(t, false, skipPath("../foo/bar", DefaultSkipDirectories))

	assert.Equal(t, true, skipPath("foo/@eaDir", DefaultSkipDirectories))
	assert.Equal(t, true, skipPath("foo/@eaDir/", DefaultSkipDirectories))
	assert.Equal(t, true, skipPath("foo/@eaDir/synology", DefaultSkipDirectories))
}
