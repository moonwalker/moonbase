package version

import (
	"fmt"
	"runtime/debug"
)

var (
	commit = "dev"
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

func ShortRev() string {
	return fmt.Sprintf("%.7s", commit)
}
