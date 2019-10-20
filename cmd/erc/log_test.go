package main

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	type test struct {
		level    string
		expected log.Level
	}

	cases := map[string]test{
		"blank":   {"", log.ErrorLevel},
		"invalid": {"...", log.ErrorLevel},
		"valid":   {"warn", log.WarnLevel},
	}

	for desc, c := range cases {
		t.Run(desc, func(tt *testing.T) {
			assert.Equal(tt, c.expected, logLevel(c.level))
		})
	}
}

func TestOpenLogFile(t *testing.T) {
	type test struct {
		file  string
		nilFn assert.ValueAssertionFunc
		errFn assert.ErrorAssertionFunc
	}

	cases := map[string]test{
		"no file":     {"", assert.Nil, assert.NoError},
		"a good file": {"/dev/null", assert.NotNil, assert.NoError},
		"a bad file":  {"/dev/bad_file", assert.Nil, assert.Error},
	}

	for desc, c := range cases {
		t.Run(desc, func(tt *testing.T) {
			file, err := openLogFile(c.file)
			c.nilFn(tt, file)
			c.errFn(tt, err)
		})
	}
}

func TestSetLogging(t *testing.T) {
	assert.NoError(t, setLogging("", ""))
	assert.Error(t, setLogging("/dev/bleh", ""))
}
