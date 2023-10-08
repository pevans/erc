package clog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugf(t *testing.T) {
	b := new(strings.Builder)
	stdlog = NewChannel(b)

	go stdlog.WriteLoop(shutdown)

	Debug("test")
	shutdown <- true

	assert.Contains(t, b.String(), "test")
}

func TestInfof(t *testing.T) {
	b := new(strings.Builder)
	stdlog = NewChannel(b)
	Level = LevelDebug

	t.Run("info doesn't log anything at debug level", func(t *testing.T) {
		go stdlog.WriteLoop(shutdown)

		Info("test")
		shutdown <- true

		assert.NotContains(t, b.String(), "test")
	})

	Level = LevelInfo
	t.Run("info does log anything at info level", func(t *testing.T) {
		go stdlog.WriteLoop(shutdown)

		Info("test")
		shutdown <- true

		assert.Contains(t, b.String(), "test")
	})
}

func TestErrorf(t *testing.T) {
	b := new(strings.Builder)
	stdlog = NewChannel(b)

	// Error will always print something out, so we don't need to test
	// the log level.
	go stdlog.WriteLoop(shutdown)

	// Error will also send a shutdown signal to stdlog, so we don't
	// need to send a message here.
	assert.Panics(t, func() {
		Error("test")
	})

	assert.Contains(t, b.String(), "test")
}

func TestShutdown(t *testing.T) {
	sl := len(shutdown)

	go func() {
		Shutdown()
	}()

	<-shutdown
	assert.Equal(t, len(shutdown), sl)
}
