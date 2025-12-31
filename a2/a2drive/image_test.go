package a2drive

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageType(t *testing.T) {
	cases := []struct {
		name  string
		fname string
		want  int
		errfn assert.ErrorAssertionFunc
	}{
		{"do file", "something.do", a2enc.DOS33, assert.NoError},
		{"dsk file", "something.dsk", a2enc.DOS33, assert.NoError},
		{"nib file", "something.nib", a2enc.Nibble, assert.NoError},
		{"po file", "something.po", a2enc.ProDOS, assert.NoError},
		{"bad file", "bad", -1, assert.Error},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			typ, err := ImageType(c.fname)
			assert.Equal(t, c.want, typ)
			c.errfn(t, err)
		})
	}
}

func TestDriveLoad(t *testing.T) {
	d := NewDrive()

	data, err := os.Open("../../data/logical.disk")
	require.NoError(t, err)
	defer data.Close()

	require.NoError(t, d.Load(data, "something.dsk"))

	assert.Equal(t, a2enc.DOS33, d.imageType)
	assert.NotNil(t, d.image)
	assert.NotNil(t, d.data)
}

func TestDriveSave(t *testing.T) {
	logSeg := memory.NewSegment(a2enc.DosSize)
	for i := range logSeg.Size() {
		logSeg.Set(i, uint8(i%256))
	}

	encodedDOS33, err := a2enc.Encode(a2enc.DOS33, logSeg)
	require.NoError(t, err)

	encodedProDOS, err := a2enc.Encode(a2enc.ProDOS, logSeg)
	require.NoError(t, err)

	// Create a modified logical segment to test that changes to encoded data
	// are properly saved
	modifiedLogSeg := memory.NewSegment(a2enc.DosSize)
	for i := range modifiedLogSeg.Size() {
		modifiedLogSeg.Set(i, uint8(i%256))
	}

	modifiedLogSeg.Set(0, 0xAA)
	modifiedLogSeg.Set(100, 0xBB)
	modifiedLogSeg.Set(1000, 0xCC)

	encodedModified, err := a2enc.Encode(a2enc.DOS33, modifiedLogSeg)
	require.NoError(t, err)

	cases := []struct {
		name          string
		imageName     string
		imageType     int
		data          *memory.Segment
		expectedLog   *memory.Segment
		checkFileStat bool
		errfn         assert.ErrorAssertionFunc
	}{
		{
			name:          "empty ImageName does nothing",
			imageName:     "",
			imageType:     a2enc.DOS33,
			data:          encodedDOS33,
			expectedLog:   logSeg,
			checkFileStat: false,
			errfn:         assert.NoError,
		},
		{
			name:          "empty Data does nothing",
			imageName:     "something!",
			imageType:     a2enc.DOS33,
			data:          nil,
			expectedLog:   nil,
			checkFileStat: false,
			errfn:         assert.NoError,
		},
		{
			name:          "DOS33 saves successfully",
			imageName:     filepath.Join(t.TempDir(), "test_dos33.dsk"),
			imageType:     a2enc.DOS33,
			data:          encodedDOS33,
			expectedLog:   logSeg,
			checkFileStat: true,
			errfn:         assert.NoError,
		},
		{
			name:          "ProDOS saves successfully",
			imageName:     filepath.Join(t.TempDir(), "test_prodos.po"),
			imageType:     a2enc.ProDOS,
			data:          encodedProDOS,
			expectedLog:   logSeg,
			checkFileStat: true,
			errfn:         assert.NoError,
		},
		{
			name:          "modified data is saved correctly",
			imageName:     filepath.Join(t.TempDir(), "test_modified.dsk"),
			imageType:     a2enc.DOS33,
			data:          encodedModified,
			expectedLog:   modifiedLogSeg,
			checkFileStat: true,
			errfn:         assert.NoError,
		},
		{
			name:          "invalid image type returns error",
			imageName:     filepath.Join(t.TempDir(), "test_invalid.dsk"),
			imageType:     99,
			data:          encodedDOS33,
			expectedLog:   nil,
			checkFileStat: false,
			errfn:         assert.Error,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := NewDrive()

			d.imageName = c.imageName
			d.imageType = c.imageType
			d.image = logSeg
			d.data = c.data

			err := d.Save()
			c.errfn(t, err)

			if c.checkFileStat {
				_, statErr := os.Stat(c.imageName)
				assert.NoError(t, statErr)

				// Read the saved file back and verify contents
				savedBytes, err := os.ReadFile(c.imageName)
				require.NoError(t, err)
				assert.Equal(t, c.expectedLog.Size(), len(savedBytes))

				// Verify every byte matches the expected logical segment
				for i := range c.expectedLog.Size() {
					assert.Equal(t,
						c.expectedLog.Get(i), savedBytes[i],
						"byte mismatch at offset %v", i,
					)
				}
			}
		})
	}
}
