package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FsNode struct {
	// The basename of the node. For example, if the absolute path is
	// "/csgo/faze/karrigan", then name is "karrigan".  The value is
	// initialized from the original path, but will be changed by sauber
	// during the rename operation.
	name string
	// The original, pre-sauber path of the node on the actual filesystem.
	// Can be an absolute or relative path.  Needed for detecting renaming
	// collisions and for the eventual rename operation on the actual
	// filesystem.
	// Examples: "/foo/Ähnlich", "hello/dolly"
	originalPath string
	isDir        bool
	parent       *FsNode
	children     []*FsNode
}

func (node FsNode) OriginalPath() string {
	return node.originalPath
}

func (node FsNode) OriginalName() string {
	return filepath.Base(node.originalPath)
}

func (node FsNode) Apply(f func(n FsNode)) {
	f(node)
	for _, child := range node.children {
		(*child).Apply(f)
	}
}

func (node FsNode) Print() {
	node.Apply(func(n FsNode) {
		fmt.Println(n.Path())
	})
}

// AddNestedChild Updates the node's tree with the given path, creating any
// missing child nodes in the hierarchy as needed (think: sub-directories).
//
// For example, if node path is "/csgo", then adding path "faze/ropz" changes
// the node's tree to have paths "/csgo", "/csgo/faze", "csgo/faze/ropz".
//
// Note to maintainers:
// This function must operate on `*FsNode` to be able to mutate the tree.
func (node *FsNode) AddNestedChild(originalPath string, isDir bool) {
	relPath := strings.TrimPrefix(originalPath, node.originalPath)
	subPaths := strings.Split(relPath, string(os.PathSeparator))
	if len(subPaths) > 0 && subPaths[0] == "" {
		// Needed for when root node has a relative path
		subPaths = subPaths[1:]
	}
	if len(subPaths) > 0 {
		var child *FsNode = nil
		firstLevelPath := subPaths[0]
		for _, n := range node.children {
			if n.name == firstLevelPath {
				child = n
			}
		}
		if child == nil {
			var isDirectory bool
			if len(subPaths) > 1 {
				// This child has children itself, which means it MUST be a directory.
				isDirectory = true
			} else {
				// Leaf node, i.e., we're at the end of the directory tree.
				// The passed parameter then determines whether this node is
				// a directory or not (likely from an `os.Stat()` call).
				isDirectory = isDir
			}
			child = &FsNode{
				name:         firstLevelPath,
				originalPath: filepath.Clean(node.originalPath + string(os.PathSeparator) + firstLevelPath),
				isDir:        isDirectory,
				parent:       node,
				children:     nil,
			}
			node.children = append(node.children, child)
		}
		if len(subPaths) > 1 {
			remainingPath := filepath.Join(subPaths[1:]...)
			child.AddNestedChild(remainingPath, isDir)
		}
	}
}

func (node FsNode) Path() string {
	var path string
	if node.parent != nil {
		path = (*node.parent).Path()
	}
	path = filepath.Join(path, node.name)
	if node.IsRoot() {
		originalParentPath := filepath.Dir(node.OriginalPath())
		path = filepath.Join(originalParentPath, path)
	}
	return filepath.Clean(path)
}

func (node FsNode) PathDecorated() string {
	path := node.Path()
	if node.isDir {
		path += "[d]"
	}
	return path
}

// RenamePath returns the node's current path during a running rename
// operation. The path can differ from the original path because
// one of the node's parents might have been renamed in the meantime.
//
// Used as the source path when renaming the node on the actual
// filesystem.
func (node FsNode) RenamePath() string {
	var path string
	if node.parent != nil {
		path = (*node.parent).Path()
	}
	path = filepath.Join(path, node.OriginalName())
	if node.IsRoot() {
		originalParentPath := filepath.Dir(node.OriginalPath())
		path = filepath.Join(originalParentPath, path)
	}
	return filepath.Clean(path)
}

func (node FsNode) Paths() []string {
	return node.paths(false)

}

func (node FsNode) PathsDecorated() []string {
	return node.paths(true)
}

func (node FsNode) paths(isDecorate bool) []string {
	var paths []string
	if isDecorate {
		paths = append(paths, node.PathDecorated())
	} else {
		paths = append(paths, node.Path())
	}
	for _, child := range node.children {
		var childPaths []string
		if isDecorate {
			childPaths = (*child).PathsDecorated()
		} else {
			childPaths = (*child).Paths()
		}
		for _, cp := range childPaths {
			paths = append(paths, cp)
		}
	}
	return paths
}

func (node FsNode) IsRoot() bool {
	return node.parent == nil
}

func (node FsNode) Siblings() []FsNode {
	if node.parent != nil {
		var siblings []FsNode
		for _, n := range node.parent.children {
			if n.name != node.name {
				siblings = append(siblings, *n)
			}
		}
		return siblings
	} else {
		return nil
	}
}

// HasSiblingOfName returns true if any of the node's siblings is already
// named the same as the provided name, false otherwise.
//
// Used to detect name collisions during rename, and thus to determine
// whether we can safely rename the node to the given name.
func (node FsNode) HasSiblingOfName(name string) bool {
	for _, sibling := range node.Siblings() {
		if sibling.name == name {
			return true
		}
	}
	return false
}
