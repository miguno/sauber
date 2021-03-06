package internal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Find(rootPath string, skipSet map[string]bool) (*FsNode, error) {
	if _, err := os.Stat(rootPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("'%s' does not exist", rootPath)
	}
	// Important to get rid of ".." and "." pollution in paths
	rootPath = filepath.Clean(rootPath)
	var rootNode *FsNode = nil
	err := filepath.WalkDir(rootPath,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !skipPath(path, skipSet) {
				if rootNode != nil {
					rootNode.AddNestedChild(path, info.IsDir())
				} else {
					rootNode = &FsNode{
						name:         filepath.Base(path),
						originalPath: path,
						isDir:        info.IsDir(),
					}
				}
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return rootNode, nil
}

func skipPath(path string, skipSet map[string]bool) bool {
	subPaths := strings.Split(path, string(os.PathSeparator))
	for _, sp := range subPaths {
		if skipSet[sp] {
			return true
		}
	}
	return false
}
