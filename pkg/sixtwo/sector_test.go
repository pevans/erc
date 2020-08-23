package sixtwo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogicalSector(t *testing.T) {
	// Test that weird sectors return as sector zero
	assert.Equal(t, 0, logicalSector(0, -1))
	assert.Equal(t, 0, logicalSector(0, 16))

	// Test that normal sectors in known image types return as expected
	sector := 7
	assert.Equal(t, dosSectorTable[sector], logicalSector(DOS33, sector))
	assert.Equal(t, proSectorTable[sector], logicalSector(ProDOS, sector))

	// Test that unknown image types return the sector as given
	assert.Equal(t, sector, logicalSector(-1, sector))
}
