package obj

// SystemROM returns the embedded system ROM data.
func SystemROM() []uint8 {
	return systemROM
}

// PeripheralROM returns the embedded peripheral ROM data.
func PeripheralROM() []uint8 {
	return peripheralROM
}
