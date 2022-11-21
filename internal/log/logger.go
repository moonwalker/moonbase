package log

import (
	"os"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/moonwalker/moonbase/internal/runtime"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if runtime.IsDev() {
		zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		zlog.Logger = zlog.With().
			Str("name", runtime.Name).
			Str("version", runtime.ShortRev()).
			Logger()
	}
}

func Info() *zerolog.Event {
	return zlog.Info()
}

func Error(err error) *zerolog.Event {
	return zlog.Error().Err(err)
}
