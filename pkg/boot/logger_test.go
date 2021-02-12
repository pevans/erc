package boot

import (
	"log"
	"os"
	"strings"
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
				assert.Equal(t, os.Stdout, w.logger.Writer())
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

func TestLogUseOutput(t *testing.T) {
	l, _ := DefaultConfig().NewLogger()

	l.UseOutput()
	assert.Equal(t, os.Stdout, log.Writer())
}

func testErrLevel(
	t *testing.T,
	level LogLevel,
	s string,
	fn func(*Logger, string),
	compar assert.ComparisonAssertionFunc,
) {
	// Test with a valid level
	b := new(strings.Builder)
	l := &Logger{
		logger: log.New(b, "", log.LstdFlags),
		Level:  level,
	}

	fn(l, s)
	compar(t, b.String(), s)
}

func TestLogError(t *testing.T) {
	fn := func(l *Logger, s string) {
		l.Error(s)
	}

	testErrLevel(t, LogError, "heyf", fn, assert.Contains)
	testErrLevel(t, LogNothing, "heyg", fn, assert.NotContains)
}

func TestLogErrorf(t *testing.T) {
	fn := func(l *Logger, s string) {
		l.Errorf("%s", s)
	}

	testErrLevel(t, LogError, "heyf", fn, assert.Contains)
	testErrLevel(t, LogNothing, "heyg", fn, assert.NotContains)
}

func TestLogDebug(t *testing.T) {
	fn := func(l *Logger, s string) {
		l.Debug(s)
	}

	testErrLevel(t, LogDebug, "heyf", fn, assert.Contains)
	testErrLevel(t, LogError, "heyg", fn, assert.NotContains)
}

func TestLogDebugf(t *testing.T) {
	fn := func(l *Logger, s string) {
		l.Debugf("%s", s)
	}

	testErrLevel(t, LogDebug, "heyf", fn, assert.Contains)
	testErrLevel(t, LogError, "heyg", fn, assert.NotContains)
}
