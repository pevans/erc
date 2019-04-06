package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetLogging(t *testing.T) {
	cases := []struct {
		file, level string
		fn          assert.ErrorAssertionFunc
	}{
		{"", "", assert.NoError},
		{"/dev/null", "", assert.NoError},
		{"/dev/not/here", "", assert.Error},
		{"/dev/null", "warn", assert.NoError},
		{"/dev/null", "fake_level", assert.Error},
	}

	for _, c := range cases {
		c.fn(t, setLogging(c.file, c.level))
	}
}
