package version

import (
	"fmt"
	"runtime/debug"
)

const (
	Name      = "moonbase"
	defCommit = "dev"
)

var (
	commit = defCommit
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			commit = kv.Value
		}
	}
}

func IsDev() bool {
	return ShortRev() == defCommit
}

func ShortRev() string {
	return fmt.Sprintf("%.7s", commit)
}
