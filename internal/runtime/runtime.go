package runtime

import (
	"fmt"
	"runtime/debug"
)

const (
	Name = "moonbase"
)

var (
	dev    = true
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
			dev = false
			commit = kv.Value
		}
	}
}

func IsDev() bool {
	return dev
}

func ShortRev() string {
	return fmt.Sprintf("%.7s", commit)
}
