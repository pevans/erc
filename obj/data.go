package obj

import (
	_ "embed"
)

//go:embed apple2e.rom
var systemROM []uint8

//go:embed peripheral.rom
var peripheralROM []uint8
