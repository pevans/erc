package obj

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemROM(t *testing.T) {
	rom := SystemROM()
	assert.NotNil(t, rom)
	assert.Greater(t, len(rom), 0, "SystemROM should not be empty")
}

func TestPeripheralROM(t *testing.T) {
	rom := PeripheralROM()
	assert.NotNil(t, rom)
	assert.Greater(t, len(rom), 0, "PeripheralROM should not be empty")
}
