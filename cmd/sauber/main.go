package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"

	internal "github.com/miguno/sauber/internal/pkg"
)

// Version is used to inject version information during the project build process (see `justfile`).
var Version = "development"

// TODO: Support multiple positional args as input locations, e.g. `sauber *.mp3`
// TODO: Increase test coverage
func main() {
	type OptionsArgs struct {
		Folder string `description:"Path to process, including any sub-folders and files if path is a folder. (Additional positional arguments are ignored.)" positional-arg-name:"<path>"`
	}
	var Options struct {
		DryRun            bool `short:"d" long:"dry-run" description:"Only show what would be done (default mode)"`
		ActualRun         bool `short:"f" long:"force" description:"Make actual changes to filesystem ***modifies your data***"`
		MaxRenameAttempts int  `short:"n" long:"max-rename-attempts" default:"100000" description:"Maximum number of rename attempts per file/folder. sauber will terminate when it can not find a sanitized name after this many attempts."`
		Silent            bool `short:"s" long:"silent" description:"Suppress output when sanitizing (ignored when dry-running)"`
		Truncate          int  `short:"t" long:"truncate" default:"999999999" description:"Max number of characters (actually: bytes) in the sanitized name of a file/folder. Any additional characters are truncated, though file extensions are preserved. Note: Encrypted drives on Synology NAS devices have a limit of 143 characters per file/folder (limit applies to basename, not full path). For details see the Synology DSM Tech Specs or view the summary at https://github.com/miguno/sauber/."`
		Version           bool `short:"v" long:"version" description:"Print version information and exit"`
		//Folder            string `required:"1" positional-args:"yes" positional-arg-name:"folder" value-name:"foo"`
		Args OptionsArgs `positional-args:"yes"`
	}
	parser := flags.NewParser(&Options, flags.Default)
	// The `args` return value is ignored because the only positional argument
	// we accept is already parsed into `OptionsArgs.Folder` automatically
	// (think: `OptionsArgs.Folder = popd(args)`, thus reducing the args count
	// by 1).
	_, err := parser.Parse()

	if Options.Version {
		_, _ = fmt.Fprintf(os.Stderr, "sauber version: %s\n", Version)
		os.Exit(0)
	}
	if err == flags.ErrHelp || (Options.Args.Folder == "") {
		parser.WriteHelp(os.Stderr)
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
	if !(Options.MaxRenameAttempts >= 1) {
		log.Fatalf("max number of rename attempts must be >= 1, you provided %d", Options.MaxRenameAttempts)
	}
	if !(Options.Truncate >= 1) {
		log.Fatalf("max number of characters in the name of a file/dir must be >= 1, you provided %d",
			Options.Truncate)
	}

	config := internal.Config{
		SkipDirectories:          internal.DefaultSkipDirectories,
		MaxRenameAttemptsPerPath: Options.MaxRenameAttempts,
		MaxBasenameLength:        Options.Truncate,
		SilentMode:               Options.Silent,
	}

	if Options.Args.Folder != "" {
		rootPath := Options.Args.Folder
		root, err := internal.Find(rootPath, config.SkipDirectories)
		if err != nil {
			log.Fatalf("failed to access or list contents of '%s', because %s",
				rootPath, err.Error())
		}
		isActualRun := Options.ActualRun && !Options.DryRun
		process(isActualRun, root, config)
	}
}

func process(isActualRun bool, node *internal.FsNode, config internal.Config) {
	err := internal.Rename(isActualRun, node, config)
	if err != nil {
		log.Fatal(err.Error())
	}
}
