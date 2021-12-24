package clog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugf(t *testing.T) {
	b := new(strings.Builder)
	stdlog = NewChannel(b)

	go stdlog.WriteLoop(Shutdown)

	Debug("test")
	Shutdown <- true

	assert.Contains(t, b.String(), "test")
}

func TestInfof(t *testing.T) {
	b := new(strings.Builder)
	stdlog = NewChannel(b)
	Level = LevelDebug

	t.Run("info doesn't log anything at debug level", func(t *testing.T) {
		go stdlog.WriteLoop(Shutdown)

		Info("test")
		Shutdown <- true

		assert.NotContains(t, b.String(), "test")
	})

	Level = LevelInfo
	t.Run("info does log anything at info level", func(t *testing.T) {
		go stdlog.WriteLoop(Shutdown)

		Info("test")
		Shutdown <- true

		assert.Contains(t, b.String(), "test")
	})
}

func TestErrorf(t *testing.T) {
	b := new(strings.Builder)
	stdlog = NewChannel(b)

	// Error will always print something out, so we don't need to test
	// the log level.
	go stdlog.WriteLoop(Shutdown)

	// Error will also send a Shutdown signal to stdlog, so we don't
	// need to send a message here.
	assert.Panics(t, func() {
		Error("test")
	})

	assert.Contains(t, b.String(), "test")
}
