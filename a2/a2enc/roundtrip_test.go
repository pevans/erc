package a2enc_test

import (
	"testing"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
)

func TestRoundtrip(t *testing.T) {
	original := memory.NewSegment(a2enc.DosSize)
	for i := range original.Size() {
		sector := i / 256
		original.Set(i, uint8(sector))
	}

	encoded, err := a2enc.Encode(a2enc.DOS33, original)
	assert.NoError(t, err)
	assert.NotNil(t, encoded)

	decoded, err := a2enc.Decode(a2enc.DOS33, encoded)
	assert.NoError(t, err)
	assert.NotNil(t, decoded)

	for i := range original.Size() {
		origByte := original.Get(i)
		decodedByte := decoded.Get(i)
		if origByte != decodedByte {
			sector := i / 256
			offsetInSector := i % 256
			t.Errorf(
				"mismatch at offset $%04X (sector %d, byte %d): expected $%02X, got $%02X",
				i, sector, offsetInSector, origByte, decodedByte,
			)
		}
	}
}
