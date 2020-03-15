package boot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	cases := []struct {
		file string
		fn   assert.ErrorAssertionFunc
	}{
		{"", assert.NoError},
		{"/tmp", assert.Error},
		{"./test.toml", assert.NoError},
	}

	for _, c := range cases {
		_, err := NewConfig(c.file)
		c.fn(t, err, c.file)
	}
}

func TestDefaultConfig(t *testing.T) {
	assert.NotNil(t, DefaultConfig())
}
