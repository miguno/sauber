# sauber [![CI workflow status](https://github.com/miguno/sauber/actions/workflows/ci.yml/badge.svg)](https://github.com/miguno/sauber/actions/workflows/ci.yml)

A command line tool that sanitizes filenames on a Synology NAS so the files can
be read and accessed through shared network drives on the NAS.

```sh
# Before
user@nas:~$ tree -AF /volume1/music-share
/volume1/music-share/
├── Ähnlich/
├── Protégé.mp3
├── Rätsel.mp3
└── Ürtümlich/
    └── Intro.mp3

# After
user@nas:~$ sauber --force /volume1/music-share
user@nas:~$ ls /volume1/music-share
/volume1/music-share/
├── Aehnlich/
├── Protege.mp3
├── Raetsel.mp3
└── Uertuemlich/
    └── Intro.mp3
```

# Download

See [Releases](https://github.com/miguno/sauber/releases).

The download is a standalone executable, which does not need installation.
Just run it! It is recommended, though optional, to rename your downloaded
binary to simply `sauber`.

| Executable Name               | Operating System     |
|-------------------------------|----------------------|
| `sauber_linux-386`            | Linux x86 32-bit     |
| `sauber_linux-amd64`          | Linux x86 64-bit     |
| `sauber_macos-arm64`          | macOS ARM 64-bit     |

# Usage

sauber sanitizes the names of files and directories by replacing umlauts,
accents, and similar diacritics.  By default, it performs a dry run to
let you verify any changes it would make.

> You must run `sauber` directly on the Synology NAS!
>
> Due to the nature of the problem it solves, `sauber` must operate directly on
> the filesystem of the NAS attached storage.  You do not need to run `sauber`
> as user `root`, however.

```
$ sauber -h
Usage: sauber [-dfhsv] [-n value] [-t value] [parameters ...]
 -d, --dry-run  only show what would be done (default mode)
 -f, --force    make actual changes to filesystem ***modifies your data***
 -h, --help     print this usage information and exit
 -n, --max-rename-attempts=value
                maximum number of rename attempts per file/directory;
                sauber will terminate when it can not find a sanitized
                name after this many attempts [100000]
 -s, --silent   suppress output when sanitizing (ignored when dry-running)
 -t, --truncate=value
                max number of characters (actually: bytes) in the sanitized name
                of a file/dir; any additional characters are truncated, though
                file extensions are preserved;
                Note: Encrypted drives on Synology NAS devices have a limit
                of 143 characters per file/dir (limit applies to basename,
                not full path). For details see the Synology DSM Tech Specs
                or view the summary at https://github.com/miguno/sauber/. [999999999]
 -v, --version  print version information and exit

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

Suggestions? Bugs? Questions? Go to https://github.com/miguno/sauber/
```

# Why do I need sauber?

If you are reading this, you are likely a fellow Synology NAS user.
Many users run into errors on Synology NAS devices when the names of files
and directories contain non-ASCII characters, such as German umlauts
(`Ä`), French accents (`é`), and Polish diacritics (`ł`).

Files with non-ASCII characters work perfectly on the NAS itself,
but accessing the files remotely via SMB shares does not work: the NAS even
reports misleading error messages like "The file is missing" or "The file is
corrupt".  Any such files on the NAS can no longer be read, copied, renamed,
etc. from other connected devices such as computers and mobile phones.

```sh
# Example error when trying to read a file with non-ASCII characters from
# a shared network drive on a Synology NAS.
alice@laptop $ ls -l /mnt/music-share/Rätsel.mp3
[/mnt/nas-share/Rätsel.mp3: No such file or directory (os error 2)]
```

The easiest remedy is to rename such files by replacing all non-ASCII
characters with ASCII characters directly on the NAS device itself.
This is the purpose of `sauber`.

# References

* [Synology DSM Technical Specifications](https://www.synology.com/en-global/dsm/7.0/software_spec/dsm):
  see "Storage Manager" > "Specifications" > "General" to view the
  "Maximum file name length" and "Maximum path name length" limitations
  for the supported filesystems.
    * For ext4 (see note below in case of encrypted shared folders):
        * Maximum file name length: 255 bytes
        * Maximum path name length: 4,096 bytes
    * For btrfs (see note below in case of encrypted shared folders):
        * Maximum file name length: 255 bytes
        * Maximum path name length: 4,096 bytes
    * Note that different character encodings may contain different data sizes
      (e.g., a character with UTF-8 encoding may contain 1 to 4 bytes).
    * For **encrypted** shared folders, the length of file/folder name should
      be within 143 characters (up to about 47 characters for non-Latin
      languages), and the length of the file path should be within 2,048
      characters.
* [Synology encryption 143 character limit - does it refer to file name or the entire path?](https://www.reddit.com/r/synology/comments/m93gha/synology_encryption_143_character_limit_does_it/)

