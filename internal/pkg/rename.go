package internal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func Rename(isActualRun bool, node *FsNode, config Config) error {
	if node == nil {
		return errors.New("node must not be nil")
	}
	renameAttemptsThusFar := 0
	for renameAttemptsThusFar < config.MaxRenameAttemptsPerPath {
		candidateName, err := sanitizeWithCounter(*node, renameAttemptsThusFar, config)
		if err != nil {
			return err
		}
		// Rename only when needed
		if candidateName != node.name {
			// Safe to rename (unless the filesystem changed out-of-band in the
			// meantime, unbeknownst to us)
			if !node.HasSiblingOfName(candidateName) {
				node.name = candidateName
				if isActualRun {
					err := os.Rename(node.RenamePath(), node.Path())
					if err != nil {
						return err
					}
				} else {
					if !config.SilentMode {
						fmt.Println(
							color.RedString(node.originalPath),
							"=>", color.GreenString(node.Path()))
					}
				}
				break
			}
		} else {
			if !isActualRun && !config.SilentMode {
				if node.originalPath == node.Path() {
					fmt.Println(node.originalPath, "[unmodified]")
				} else {
					// path changed because at least one parent directory has
					// been renamed
					fmt.Println(
						color.RedString(node.originalPath),
						"=>", color.GreenString(node.Path()))
				}
			}
			break
		}
		renameAttemptsThusFar++
	}

	if renameAttemptsThusFar >= config.MaxRenameAttemptsPerPath {
		return fmt.Errorf("failed to rename '%s' (no rename attempts left)", node.originalPath)
	} else {
		for _, child := range node.children {
			err := Rename(isActualRun, child, config)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func sanitizeWithCounter(node FsNode, renameAttemptsThusFar int, config Config) (string, error) {
	if !(renameAttemptsThusFar >= 0) {
		log.Fatalf("renameAttemptsThusFar must be >= 0, you provided %d", renameAttemptsThusFar)
	}
	if !(config.MaxBasenameLength > 0) {
		log.Fatalf("maxRenameAttempts must be > 0, you provided %d", config.MaxRenameAttemptsPerPath)
	}
	candidate := Sanitize(node.name)
	candidate, err := truncateName(candidate, node.isDir, config.MaxBasenameLength)
	if err != nil {
		return "", err
	}
	if renameAttemptsThusFar > 0 {
		digits := numDigits(config.MaxRenameAttemptsPerPath - 1)
		formatString := fmt.Sprintf("%%s_%%0%dd", digits)
		candidate = fmt.Sprintf(formatString, candidate, renameAttemptsThusFar)
	}
	return candidate, nil
}

func numDigits(n int) int {
	if n == 0 {
		return 1
	}
	count := 0
	for n != 0 {
		n /= 10
		count += 1
	}
	return count
}

func truncateName(name string, isDir bool, maxBasenameLength int) (string, error) {
	if !(maxBasenameLength >= 1) {
		return "", fmt.Errorf("maxBasenameLength must be >= 1, you provided %d", maxBasenameLength)
	}
	if len(name) <= maxBasenameLength {
		return name, nil
	} else {
		if isDir {
			return name[:maxBasenameLength], nil
		} else {
			extension := filepath.Ext(name)
			if extension != "" {
				if len(extension) > maxBasenameLength {
					return "",
						fmt.Errorf("could not truncate name '%s' to %d characters while preserving file extension '%s'",
							name,
							maxBasenameLength,
							extension)
				} else {
					nameWithoutExtension := name[:maxBasenameLength-len(extension)]
					return nameWithoutExtension + extension, nil
				}
			} else {
				return name[:maxBasenameLength], nil
			}
		}
	}
}
