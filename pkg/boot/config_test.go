package boot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	type test struct {
		file string
		efn  assert.ErrorAssertionFunc
		vfn  assert.ValueAssertionFunc
	}

	cases := map[string]test{
		"no file, default bevhavior": {
			file: "",
			efn:  assert.NoError,
			vfn:  assert.NotNil,
		},
		"not a file": {
			file: "/tmp",
			efn:  assert.Error,
			vfn:  assert.Nil,
		},
		"invalid file": {
			file: "/does/not/exist",
			efn:  assert.Error,
			vfn:  assert.Nil,
		},
		"valid file": {
			file: "../../data/config.toml",
			efn:  assert.NoError,
			vfn:  assert.NotNil,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(tt *testing.T) {
			conf, err := NewConfig(c.file)
			c.efn(tt, err)
			c.vfn(tt, conf)
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	assert.NotNil(t, DefaultConfig())
}
