package boot

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// A Config is simply a record of things that we allow the user to
// configure for run-time.
type Config struct {
	Log struct {
		File  string
		Level string
	}

	InstructionLog struct {
		File string
	} `toml:"instruction_log"`
}

// NewConfig returns a new configuration object based on the given
// config file. Note that it's ok if that file doesn't exist--in that
// event, we will assume a default configuration.
func NewConfig(file string) (*Config, error) {
	conf := Config{}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	_, err := toml.DecodeFile(file, &conf)
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode config file %s", file)
	}

	return &conf, nil
}

// DefaultConfig simply returns the basic, default configuration that we
// can work with.
func DefaultConfig() *Config {
	return &Config{}
}
