package main

import (
	"github.com/pevans/erc/cmd"
)

// ConfigFile is the default (relative) location of our configuration file.
const ConfigFile = `.erc/config.toml`

func main() {
	cmd.Execute()
}
