package main

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Config struct {
	Log struct {
		File  string
		Level string
	}
}

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

func DefaultConfig() *Config {
	return &Config{}
}
