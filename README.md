# sauber [![CI workflow status](https://github.com/miguno/sauber/actions/workflows/ci.yml/badge.svg)](https://github.com/miguno/sauber/actions/workflows/ci.yml)

A command line tool that sanitizes filenames on a Synology NAS so the files can
be read and accessed through shared network drives on the NAS.  This solves the
annoying problem that you cannot access files and folders on a shared network
drive of a Synology NAS if their names contain special characters, such as
German umlauts (`Ä`), French accents (`é`), and Polish diacritics (`ł`).

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
user@nas:~$ tree -AF /volume1/music-share
/volume1/music-share/
├── Aehnlich/
├── Protege.mp3
├── Raetsel.mp3
└── Uertuemlich/
    └── Intro.mp3
```

# Installation

Download the [**Latest Release**](https://github.com/miguno/sauber/releases)
for your operating system (see table below) and just run it!

Sauber is a standalone executable, which does not need installation.  It is
recommended, though optional, to rename your downloaded executable to simply
`sauber`.

| Executable Name               | Operating System     | `uname -m` |
|-------------------------------|----------------------|------------|
| `sauber_linux-386`            | Linux x86 32-bit     | `i386`     |
| `sauber_linux-amd64`          | Linux x86 64-bit     | `x86_64`   |
| `sauber_linux-arm`            | Linux ARM 32-bit     | `arm`      |
| `sauber_linux-arm64`          | Linux ARM 64-bit     | `arm64`    |
| `sauber_macos-arm64`          | macOS ARM 64-bit     | `arm64`    |

To find the correct exectuable for your NAS, run `uname -m` in a terminal on
the NAS and match it with the corresponding entry in the table above.

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
Usage:
  sauber [OPTIONS] [<path>]

Application Options:
  -d, --dry-run              Only show what would be done (default mode)
  -f, --force                Make actual changes to filesystem ***modifies your
                             data***
  -n, --max-rename-attempts= Maximum number of rename attempts per file/folder.
                             sauber will terminate when it can not find a
                             sanitized name after this many attempts. (default:
                             100000)
  -s, --silent               Suppress output when sanitizing (ignored when
                             dry-running)
  -t, --truncate=            Max number of characters (actually: bytes) in the
                             sanitized name of a file/folder. Any additional
                             characters are truncated, though file extensions
                             are preserved. Note: Encrypted drives on Synology
                             NAS devices have a limit of 143 characters per
                             file/folder (limit applies to basename, not full
                             path). For details see the Synology DSM Tech Specs
                             or view the summary at
                             https://github.com/miguno/sauber/. (default:
                             999999999)
  -v, --version              Print version information and exit

Help Options:
  -h, --help                 Show this help message

Arguments:
  <path>:                    Path to process, including any sub-folders and
                             files if path is a folder. (Additional positional
                             arguments are ignored.)

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

Step-by-Step:
1. Enable SSH and ssh with Terminal as admin to Synology DiskStation
2. Get root privileges (use admin password) with $ sudo -i
3. Change directory to where the sauber executable is located
4. Execute sauber with $ ./sauber

Pro-Tip, if you want to sanatize a lot of files:
1. Make dry run with sauber and write results to a textfile with $ ./sauber /volume1/Music/ > music.txt
2. Download music.txt and analyse in a text editor (e.g. Visual Studio Code)
Remove all lines, where files are unmodified:
  a) Open search and enable regular expressions, use the following: ^.*\[unmodified\]\n
  b) Replace all with nothing (empty Replace)
3. If the result of renames are OK, actually rename the files with $ ./sauber -f /volume1/Music/


Suggestions? Bugs? Questions? Go to https://github.com/miguno/sauber/
```

# How are names of files and folders sanitized?

Here's a short summary of what sanitization rules you can expect. The exact
rules are defined in [sanitize.go](internal/pkg/sanitize.go), with further
examples in [sanitize_test.go](internal/pkg/sanitize_test.go).

| Original                      | Replacement                      |
|-------------------------------|----------------------------------|
| Ä, Ö, Ü, ä, ö, ü, ß           | Ae, Oe, Ue, ae, oe, ue, ss       |
| !, ?, \|, $                   | _ (underscore)                   |
| – (en dash), — (em dash)      | - (hyphen)                       |
| ạàąâåÅ                        | aaaaaA                           |
| čćçÇČĆ                        | cccCCC                           |
| đĐ                            | dD                               |
| ęéèê                          | eeee                             |
| żźžŻŽ                         | zzzZZ                            |
| Non-printable chars           | - (hyphen)                       |
| Illegal filename chars        | - (hyphen)                       |
| Private Unicode chars         | - (hyphen)                       |
| (and more)                    | (and more)                       |

Also sauber checks for reserved filenames and removes illegal spaces or dots at the end of filenames.

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
* [Synology Help: My file name is garbled. What can I do?](https://kb.synology.com/en-global/DSM/tutorial/garbled_name_smb_FileStation)
  -- did not help in my case (Synology DSM 6.x), as the SMB configuration was
  already set up correctly
* [Synology encryption 143 character limit - does it refer to file name or the entire path?](https://www.reddit.com/r/synology/comments/m93gha/synology_encryption_143_character_limit_does_it/)

# Notes

* Some people who use Synology NAS running DSM 7.x together with macOS clients
  have reported success when using umlauts etc. in file names.  This required
  enabling and configuring the `vfs_fruit` SMB module in the Samba config
  at `/etc/samba/smb.conf`.  Unfortunately, this module is only supported on
  Synology DSM 7.x, not on DSM 6.x (see
  [Reddit discussion](https://www.reddit.com/r/synology/comments/p5bz8t/)).
