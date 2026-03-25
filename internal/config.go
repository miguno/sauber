package internal

type Config struct {
	SkipDirectories          map[string]bool
	MaxRenameAttemptsPerPath int
	MaxBasenameLength        int
	SilentMode               bool
}

var DefaultSkipDirectories = map[string]bool{
	"@eaDir": true, // special directory on Synology NAS
}
