package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pborman/getopt/v2"

	"github.com/miguno/sauber/internal/pkg"
)

// Version is used to inject version information during the project build process (see `justfile`).
var Version = "development"

// TODO: Increase test coverage
func main() {
	config := internal.DefaultConfig()

	getopt.StringLong("help", 'h', "", "print this usage information and exit")
	optionHelpFlag := getopt.Lookup("help").SetFlag()

	getopt.StringLong("version", 'v', "", "print version information and exit")
	optionVersionFlag := getopt.Lookup("version").SetFlag()

	getopt.StringLong("dry-run", 'd', "", "only show what would be done (default mode)")
	optionDryRunFlag := getopt.Lookup("dry-run").SetFlag()

	getopt.StringLong("force", 'f', "",
		"make actual changes to filesystem ***modifies your data***")
	optionActualRunFlag := getopt.Lookup("force").SetFlag()

	maxRenameAttemptsFlag := getopt.IntLong("max-rename-attempts", 'n', config.MaxRenameAttemptsPerPath,
		`
maximum number of rename attempts per file/directory;
sauber will terminate when it can not find a sanitized
name after this many attempts
`)

	getopt.StringLong("silent", 's', "",
		"suppress output when sanitizing (ignored when dry-running)")
	optionSilentFlag := getopt.Lookup("silent").SetFlag()

	truncateFlag := getopt.IntLong("truncate", 't', config.MaxBasenameLength,
		`
max number of characters (actually: bytes) in the sanitized name
of a file/dir; any additional characters are truncated, though
file extensions are preserved;
Note: Encrypted drives on Synology NAS devices have a limit
of 143 characters per file/dir (limit applies to basename,
not full path). For details see the Synology DSM Tech Specs
or view the summary at https://github.com/miguno/sauber/.
`)

	getopt.Parse()

	if optionVersionFlag.Seen() {
		_, _ = fmt.Fprintf(os.Stderr, "sauber version: %s\n", Version)
		os.Exit(0)
	}
	if optionHelpFlag.Seen() || getopt.NArgs() == 0 {
		getopt.Usage()
		s := `
sauber sanitizes the names of files and directories by replacing umlauts,
accents, and similar diacritics.  By default, it performs a dry run to
let you verify any changes it would make.

Examples:
  # Perform a dry run of sanitizing /volume1/music (including) and all
  # its sub-directories and files.  This operation does not modify
  # any data, it only reports what sauber *would* change.
  $ sauber /volume1/music

  # Sanitize /volume1/music (including) and all its sub-directories and files.
  # *** WARNING: This command modifies your data! Always do a dry run first! ***
  $ sauber --force /volume1/music

Suggestions? Bugs? Questions? Go to https://github.com/miguno/sauber/`
		_, _ = fmt.Fprintln(os.Stderr, s)
		os.Exit(1)
	}
	if !(*maxRenameAttemptsFlag >= 1) {
		log.Fatalf("max number of rename attempts must be >= 1, you provided %d", *maxRenameAttemptsFlag)
	}
	config.MaxRenameAttemptsPerPath = *maxRenameAttemptsFlag

	if !(*truncateFlag >= 1) {
		log.Fatalf("max number of characters in the name of a file/dir must be >= 1, you provided %d",
			*truncateFlag)
	}
	config.MaxBasenameLength = *truncateFlag

	if optionSilentFlag.Seen() {
		config.SilentMode = true
	}

	if getopt.NArgs() > 0 {
		rootPath := getopt.Arg(0)
		root, err := internal.Find(rootPath, config.SkipDirectories)
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to access or list contents of '%s', because %s",
				rootPath, err.Error()))
		}
		isActualRun := optionActualRunFlag.Seen() && !optionDryRunFlag.Seen()
		process(isActualRun, root, config)
	}
}

func process(isActualRun bool, node *internal.FsNode, config internal.Config) {
	err := internal.Rename(isActualRun, node, config)
	if err != nil {
		log.Fatal(err.Error())
	}
}
