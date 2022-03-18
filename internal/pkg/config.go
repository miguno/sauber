package internal

type Config struct {
	SkipDirectories          map[string]bool
	MaxRenameAttemptsPerPath int
	MaxBasenameLength        int
	SilentMode               bool
}

func DefaultConfig() Config {
	var skipDirectories = map[string]bool{
		"@eaDir": true, // special directory on Synology NAS
	}
	return Config{
		SkipDirectories:          skipDirectories,
		MaxRenameAttemptsPerPath: 100000,
		MaxBasenameLength:        999999999,
		SilentMode:               false,
	}
}
