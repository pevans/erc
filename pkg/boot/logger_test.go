package boot

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLogLevel(t *testing.T) {
	type test struct {
		level string
		want  LogLevel
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
		fileName     string
		expectStdout bool
		errfn        assert.ErrorAssertionFunc
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
			fileName:     "",
			expectStdout: true,
			errfn:        assert.NoError,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			conf := Config{
				Log: configLog{
					File: c.fileName,
				},
			}

			w, err := conf.NewLogger()
			c.errfn(t, err)

			if c.expectStdout {
				assert.Equal(t, os.Stdout, w.log.Writer())
			}
		})
	}

}

func TestLogCanLog(t *testing.T) {
	e := Logger{Level: LogError}
	d := Logger{Level: LogDebug}

	assert.True(t, e.CanLog(LogError))
	assert.False(t, e.CanLog(LogDebug))

	assert.True(t, d.CanLog(LogError))
	assert.True(t, d.CanLog(LogDebug))
}
