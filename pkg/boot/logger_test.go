package boot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLogLevel(t *testing.T) {
	type test struct {
		level string
		want  int
	}

	cases := map[string]test{
		"has debug": {
			level: "debug",
			want:  LogDebug,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			conf := Config{
				Log: configLog{
					Level: c.level,
				},
			}

			assert.Equal(t, c.want, conf.LogLevel())
		})
	}
}

func TestConfigNewLogger(t *testing.T) {
	type test struct {
		fileName string
		errfn    assert.ErrorAssertionFunc
	}

	cases := map[string]test{
		"has a good file": {
			fileName: "/dev/null",
			errfn:    assert.NoError,
		},

		"has a bad file": {
			fileName: "/",
			errfn:    assert.Error,
		},

		"has no file": {
			fileName: "",
			errfn:    assert.NoError,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			conf := Config{
				Log: configLog{
					File: c.fileName,
				},
			}

			_, err := conf.NewLogger()
			c.errfn(t, err)
		})
	}

}
