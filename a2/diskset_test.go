package a2

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDiskSet(t *testing.T) {
	set := NewDiskSet()

	assert.NotNil(t, set)
	assert.NotNil(t, set.images)
	assert.Equal(t, 0, len(set.images))
	assert.Equal(t, 0, set.current)
}

func TestDiskSetAppend(t *testing.T) {
	tmpDir := t.TempDir()
	validFile := filepath.Join(tmpDir, "valid.dsk")
	err := os.WriteFile(validFile, []byte("test"), 0o644)
	assert.NoError(t, err)

	cases := []struct {
		name    string
		file    string
		wantLen int
		errfn   assert.ErrorAssertionFunc
	}{
		{
			name:    "valid file appends successfully",
			file:    validFile,
			wantLen: 1,
			errfn:   assert.NoError,
		},
		{
			name:    "multiple files append successfully",
			file:    validFile,
			wantLen: 2,
			errfn:   assert.NoError,
		},
		{
			name:    "nonexistent file returns error",
			file:    filepath.Join(tmpDir, "nonexistent.dsk"),
			wantLen: 2,
			errfn:   assert.Error,
		},
	}

	set := NewDiskSet()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := set.Append(c.file)
			c.errfn(t, err)
			assert.Equal(t, c.wantLen, len(set.images))
		})
	}
}

func TestDiskSetDisk(t *testing.T) {
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "disk1.dsk")
	file2 := filepath.Join(tmpDir, "disk2.dsk")

	err := os.WriteFile(file1, []byte("disk1"), 0o644)
	assert.NoError(t, err)
	err = os.WriteFile(file2, []byte("disk2"), 0o644)
	assert.NoError(t, err)

	set := NewDiskSet()
	assert.NoError(t, set.Append(file1))
	assert.NoError(t, set.Append(file2))

	cases := []struct {
		name     string
		index    int
		wantFile string
		readerfn assert.ValueAssertionFunc
		errfn    assert.ErrorAssertionFunc
	}{
		{
			name:     "valid index 0 returns first disk",
			index:    0,
			wantFile: file1,
			readerfn: assert.NotNil,
			errfn:    assert.NoError,
		},
		{
			name:     "valid index 1 returns second disk",
			index:    1,
			wantFile: file2,
			readerfn: assert.NotNil,
			errfn:    assert.NoError,
		},
		{
			name:     "negative index returns error",
			index:    -1,
			wantFile: "",
			readerfn: assert.Nil,
			errfn:    assert.Error,
		},
		{
			name:     "out of bounds index returns error",
			index:    2,
			wantFile: "",
			readerfn: assert.Nil,
			errfn:    assert.Error,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			file, filename, err := set.Disk(c.index)

			c.errfn(t, err)
			c.readerfn(t, file)
			assert.Equal(t, c.wantFile, filename)

			if file != nil {
				assert.NoError(t, file.Close())
			}
		})
	}
}

func TestDiskSetCurrent(t *testing.T) {
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "disk1.dsk")
	file2 := filepath.Join(tmpDir, "disk2.dsk")

	err := os.WriteFile(file1, []byte("disk1"), 0o644)
	assert.NoError(t, err)
	err = os.WriteFile(file2, []byte("disk2"), 0o644)
	assert.NoError(t, err)

	cases := []struct {
		name     string
		files    []string
		current  int
		wantFile string
		errfn    assert.ErrorAssertionFunc
	}{
		{
			name:     "current returns first disk by default",
			files:    []string{file1, file2},
			current:  0,
			wantFile: file1,
			errfn:    assert.NoError,
		},
		{
			name:     "current returns second disk when current is 1",
			files:    []string{file1, file2},
			current:  1,
			wantFile: file2,
			errfn:    assert.NoError,
		},
		{
			name:     "empty diskset returns error",
			files:    []string{},
			current:  0,
			wantFile: "",
			errfn:    assert.Error,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			set := NewDiskSet()
			for _, f := range c.files {
				assert.NoError(t, set.Append(f))
			}
			set.current = c.current

			file, filename, err := set.Current()

			c.errfn(t, err)
			assert.Equal(t, c.wantFile, filename)

			if file != nil {
				assert.NoError(t, file.Close())
			}
		})
	}
}

func TestDiskSetNext(t *testing.T) {
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "disk1.dsk")
	file2 := filepath.Join(tmpDir, "disk2.dsk")
	file3 := filepath.Join(tmpDir, "disk3.dsk")

	err := os.WriteFile(file1, []byte("disk1"), 0o644)
	assert.NoError(t, err)
	err = os.WriteFile(file2, []byte("disk2"), 0o644)
	assert.NoError(t, err)
	err = os.WriteFile(file3, []byte("disk3"), 0o644)
	assert.NoError(t, err)

	set := NewDiskSet()
	assert.NoError(t, set.Append(file1))
	assert.NoError(t, set.Append(file2))
	assert.NoError(t, set.Append(file3))

	cases := []struct {
		name        string
		wantFile    string
		wantCurrent int
	}{
		{
			name:        "first next returns second disk",
			wantFile:    file2,
			wantCurrent: 1,
		},
		{
			name:        "second next returns third disk",
			wantFile:    file3,
			wantCurrent: 2,
		},
		{
			name:        "third next wraps to first disk",
			wantFile:    file1,
			wantCurrent: 0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			file, filename, err := set.Next()

			assert.NoError(t, err)
			assert.NotNil(t, file)
			assert.Equal(t, c.wantFile, filename)
			assert.Equal(t, c.wantCurrent, set.current)

			if file != nil {
				assert.NoError(t, file.Close())
			}
		})
	}
}
