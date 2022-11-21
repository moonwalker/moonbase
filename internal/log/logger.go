package log

import (
	"os"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/moonwalker/moonbase/internal/version"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zlog.Logger = zlog.With().
		Str("name", version.Name).
		Str("version", version.ShortRev()).
		Logger()

	if version.IsDev() {
		zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func Info() *zerolog.Event {
	return zlog.Info()
}
