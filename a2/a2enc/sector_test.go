package a2enc_test

import (
	"testing"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/stretchr/testify/assert"
)

func TestLogicalSector(t *testing.T) {
	// Test that weird sectors return as sector zero
	assert.Equal(t, 0, a2enc.LogicalSector(0, -1))
	assert.Equal(t, 0, a2enc.LogicalSector(0, 16))

	// Test that normal sectors in known image types return as expected
	sector := 7
	assert.Equal(t, 0x4, a2enc.LogicalSector(a2enc.DOS33, sector))
	assert.Equal(t, 0xb, a2enc.LogicalSector(a2enc.ProDOS, sector))

	// Test that unknown image types return the sector as given
	assert.Equal(t, sector, a2enc.LogicalSector(-1, sector))
}
