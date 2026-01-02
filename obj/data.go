package obj

import (
	_ "embed"
)

//go:embed apple2e.rom
var systemROM []uint8

// SystemROM returns the embedded system ROM. This would be mapped over
// $D000-$FFFF.
func SystemROM() []uint8 {
	return systemROM
}

//go:embed peripheral.rom
var peripheralROM []uint8

// PeripheralROM returns the embedded peripheral ROM. This would be mapped
// over $C000-$CFFF; the first page is zero bytes because that would normally
// be used for soft switches.
func PeripheralROM() []uint8 {
	return peripheralROM
}
