package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbsoluteRoot(t *testing.T) {
	root := FsNode{
		name:         "/",
		originalPath: "/",
		isDir:        true,
	}
	root.AddNestedChild("/angband.exe", false)
	root.AddNestedChild("/csgo", true)
	root.AddNestedChild("/csgo/g2", true)
	root.AddNestedChild("/csgo/g2/academy", true)
	root.AddNestedChild("/csgo/navi", true)
	root.AddNestedChild("/csgo/navi/b1t", false)
	root.AddNestedChild("/csgo/navi/s1mple", false)
	expected := []string{
		"/[d]",
		"/angband.exe",
		"/csgo[d]",
		"/csgo/g2[d]",
		"/csgo/g2/academy[d]",
		"/csgo/navi[d]",
		"/csgo/navi/b1t",
		"/csgo/navi/s1mple",
	}
	assert.Equal(t, expected, root.PathsDecorated())
}

func TestRelativeRoot(t *testing.T) {
	root := FsNode{
		name:         "csgo",
		originalPath: "csgo",
		isDir:        true,
	}
	root.AddNestedChild("csgo/g2", true)
	root.AddNestedChild("csgo/navi", true)
	root.AddNestedChild("csgo/navi/academy", true)
	root.AddNestedChild("csgo/navi/academy/UnnamedPlayer", false)
	root.AddNestedChild("csgo/navi/perfecto", false)
	root.AddNestedChild("csgo/zeus", false)

	expected := []string{
		"csgo[d]",
		"csgo/g2[d]",
		"csgo/navi[d]",
		"csgo/navi/academy[d]",
		"csgo/navi/academy/UnnamedPlayer",
		"csgo/navi/perfecto",
		"csgo/zeus",
	}
	assert.Equal(t, expected, root.PathsDecorated())
}

func TestRelChildShouldBePlacedUnderAbsRoot(t *testing.T) {
	root := FsNode{
		name:         "csgo",
		originalPath: "csgo",
		isDir:        true,
	}
	root.AddNestedChild("/absolutePath", false)
	expected := []string{"csgo[d]", "csgo/absolutePath"}
	assert.Equal(t, expected, root.PathsDecorated())
}

func TestAbsChildShouldBePlacedUnderAbsRoot(t *testing.T) {
	root := FsNode{
		name:         "/csgo",
		originalPath: "/csgo",
		isDir:        true,
	}
	root.AddNestedChild("/absolutePath", false)
	expected := []string{"/csgo[d]", "/csgo/absolutePath"}
	assert.Equal(t, expected, root.PathsDecorated())
}

func TestIsRoot(t *testing.T) {
	root := FsNode{
		name:         "/csgo",
		originalPath: "/csgo",
		isDir:        true,
	}
	assert.True(t, root.IsRoot())
	assert.False(t, root.HasParents())
}

func TestHasParents(t *testing.T) {
	root := FsNode{
		name:         "/csgo",
		originalPath: "/csgo",
		isDir:        true,
	}
	root.AddNestedChild("/absolutePath", false)
	child := root.children[0]
	assert.False(t, child.IsRoot())
	assert.True(t, child.HasParents())
}
