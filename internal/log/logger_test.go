package log

import (
	"fmt"
	"testing"

	zlog "github.com/rs/zerolog/log"
)

func TestLoggerMethods(t *testing.T) {
	err := fmt.Errorf("expected: %s actual: %s", "a", "b")
	msg := "error happened"
	Error(err).Msg(msg)

	zlog.Logger = liveLogger
	Error(err).Msg(msg)

	Error(nil).Msg("")
}
