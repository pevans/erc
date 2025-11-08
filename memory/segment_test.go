package memory

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSegment(t *testing.T) {
	assert.NotEqual(t, nil, NewSegment(1))
}

func TestSet(t *testing.T) {
	s := NewSegment(100)
	addr := 1
	val := uint8(123)

	s.Set(addr, val)
	assert.Equal(t, val, s.Mem[addr])

	assert.Panics(t, func() {
		s.Set(-1, val)
	})

	assert.Panics(t, func() {
		s.Set(cap(s.Mem)+1, val)
	})
}

func TestGet(t *testing.T) {
	s := NewSegment(100)
	addr := 1
	val := uint8(123)

	s.Mem[addr] = val
	assert.Equal(t, val, s.Get(addr))

	assert.Panics(t, func() {
		_ = s.Get(-1)
	})

	assert.Panics(t, func() {
		_ = s.Get(cap(s.Mem) + 1)
	})
}

func TestSize(t *testing.T) {
	s := NewSegment(100)
	assert.Equal(t, 100, s.Size())
}

func TestCopySlice(t *testing.T) {
	type test struct {
		written   int
		size      int
		start     int
		byteSlice []uint8
		errfn     assert.ErrorAssertionFunc
	}

	cases := map[string]test{
		"negative start": {
			written:   0,
			size:      1,
			start:     -1,
			byteSlice: []uint8{},
			errfn:     assert.Error,
		},
		"no size": {
			written:   0,
			size:      0,
			start:     1,
			byteSlice: []uint8{1, 2},
			errfn:     assert.Error,
		},
		"normal": {
			written:   3,
			size:      5,
			start:     0,
			byteSlice: []uint8{1, 2, 3},
			errfn:     assert.NoError,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			s := NewSegment(c.size)

			writ, err := s.CopySlice(c.start, c.byteSlice)
			assert.Equal(t, c.written, writ)
			c.errfn(t, err)
		})
	}
}

func TestWriteFile(t *testing.T) {
	const (
		size  = 1
		value = 0x33
		file  = "/tmp/segment.writefile"
	)

	s := NewSegment(size)

	s.Set(0, value)

	// We should not be able to write a bad file
	assert.Error(t, s.WriteFile(""))

	// But should be able to write a good file
	assert.NoError(t, s.WriteFile(file))

	// We should be able to see the file, if we looked...
	ns := NewSegment(size)
	assert.NoError(t, ns.ReadFile(file))
	assert.Equal(t, uint8(value), ns.Get(0))
	assert.NoError(t, os.Remove(file))
}

func TestReadFile(t *testing.T) {
	s := NewSegment(256)

	// See if we return an error for a bad file of some kind
	assert.Error(t, s.ReadFile(""))

	// See that we don't return an error for a real file
	assert.NoError(t, s.ReadFile("../data/logical.sector"))

	// Make sure we have some real data
	assert.NotEqual(t, uint8(0), s.Get(0))
}

func TestExtractFrom(t *testing.T) {
	type test struct {
		destSize int
		srcSize  int
		start    int
		end      int
		written  int
		errfn    assert.ErrorAssertionFunc
	}

	cases := map[string]test{
		"negative start": {
			destSize: 10,
			srcSize:  10,
			start:    -1,
			end:      5,
			written:  0,
			errfn:    assert.Error,
		},
		"end beyond source": {
			destSize: 10,
			srcSize:  10,
			start:    5,
			end:      15,
			written:  0,
			errfn:    assert.Error,
		},
		"normal extraction": {
			destSize: 10,
			srcSize:  20,
			start:    5,
			end:      10,
			written:  5,
			errfn:    assert.NoError,
		},
		"full segment": {
			destSize: 10,
			srcSize:  10,
			start:    0,
			end:      10,
			written:  10,
			errfn:    assert.NoError,
		},
		"empty range": {
			destSize: 10,
			srcSize:  10,
			start:    5,
			end:      5,
			written:  0,
			errfn:    assert.NoError,
		},
		"destination too small": {
			destSize: 3,
			srcSize:  10,
			start:    0,
			end:      5,
			written:  0,
			errfn:    assert.Error,
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			src := NewSegment(c.srcSize)
			dest := NewSegment(c.destSize)

			// Fill source with test data
			for i := range c.srcSize {
				src.Set(i, uint8(i+100))
			}

			writ, err := dest.ExtractFrom(src, c.start, c.end)
			assert.Equal(t, c.written, writ)
			c.errfn(t, err)

			// Verify data was copied correctly on success
			if err == nil && c.written > 0 {
				for i := range c.written {
					assert.Equal(t, uint8(c.start+i+100), dest.Get(i),
						"byte at position %d should match source", i)
				}
			}
		})
	}
}

func TestGet16BigEndian(t *testing.T) {
	s := NewSegment(100)

	t.Run("basic 16-bit number", func(t *testing.T) {
		s.Set(10, 0x12) // MSB at addr
		s.Set(11, 0x34) // LSB at addr+1
		assert.Equal(t, uint16(0x1234), s.Get16BigEndian(10))
	})

	t.Run("maximum values", func(t *testing.T) {
		s.Set(20, 0xFF) // MSB
		s.Set(21, 0xFF) // LSB
		assert.Equal(t, uint16(0xFFFF), s.Get16BigEndian(20))
	})

	t.Run("zero values", func(t *testing.T) {
		s.Set(30, 0x00) // MSB
		s.Set(31, 0x00) // LSB
		assert.Equal(t, uint16(0x0000), s.Get16BigEndian(30))
	})

	t.Run("out of bound access", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = s.Get16BigEndian(-1)
		})

		assert.Panics(t, func() {
			_ = s.Get16BigEndian(99)
		})
	})
}

func TestGet16LittleEndian(t *testing.T) {
	s := NewSegment(100)

	t.Run("basic 16-bit number", func(t *testing.T) {
		s.Set(10, 0x34) // LSB at addr
		s.Set(11, 0x12) // MSB at addr+1
		assert.Equal(t, uint16(0x1234), s.Get16LittleEndian(10))
	})

	t.Run("maximum values", func(t *testing.T) {
		s.Set(20, 0xFF) // LSB
		s.Set(21, 0xFF) // MSB
		assert.Equal(t, uint16(0xFFFF), s.Get16LittleEndian(20))
	})

	t.Run("zero values", func(t *testing.T) {
		s.Set(30, 0x00) // LSB
		s.Set(31, 0x00) // MSB
		assert.Equal(t, uint16(0x0000), s.Get16LittleEndian(30))
	})

	t.Run("out of bound access", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = s.Get16LittleEndian(-1)
		})

		assert.Panics(t, func() {
			_ = s.Get16LittleEndian(99)
		})
	})
}

func TestGet16(t *testing.T) {
	s := NewSegment(100)

	t.Run("little endian", func(t *testing.T) {
		s.Endianness = LittleEndian
		s.Set(10, 0x34) // LSB at addr
		s.Set(11, 0x12) // MSB at addr+1
		assert.Equal(t, uint16(0x1234), s.Get16(10))
	})

	t.Run("big endian", func(t *testing.T) {
		s.Endianness = BigEndian
		s.Set(20, 0x12) // MSB at addr
		s.Set(21, 0x34) // LSB at addr+1
		assert.Equal(t, uint16(0x1234), s.Get16(20))
	})
}

func TestSet16BigEndian(t *testing.T) {
	s := NewSegment(100)

	t.Run("basic 16-bit number", func(t *testing.T) {
		s.Set16BigEndian(10, 0x1234)
		assert.Equal(t, uint8(0x12), s.Get(10)) // MSB at addr
		assert.Equal(t, uint8(0x34), s.Get(11)) // LSB at addr+1
	})

	t.Run("maximum values", func(t *testing.T) {
		s.Set16BigEndian(20, 0xFFFF)
		assert.Equal(t, uint8(0xFF), s.Get(20)) // MSB
		assert.Equal(t, uint8(0xFF), s.Get(21)) // LSB
	})

	t.Run("zero values", func(t *testing.T) {
		s.Set16BigEndian(30, 0x0000)
		assert.Equal(t, uint8(0x00), s.Get(30)) // MSB
		assert.Equal(t, uint8(0x00), s.Get(31)) // LSB
	})

	t.Run("out of bound access", func(t *testing.T) {
		assert.Panics(t, func() {
			s.Set16BigEndian(-1, 0x1234)
		})

		assert.Panics(t, func() {
			s.Set16BigEndian(99, 0x1234)
		})
	})
}

func TestSet16LittleEndian(t *testing.T) {
	s := NewSegment(100)

	t.Run("basic 16-bit number", func(t *testing.T) {
		s.Set16LittleEndian(10, 0x1234)
		assert.Equal(t, uint8(0x34), s.Get(10)) // LSB at addr
		assert.Equal(t, uint8(0x12), s.Get(11)) // MSB at addr+1
	})

	t.Run("maximum values", func(t *testing.T) {
		s.Set16LittleEndian(20, 0xFFFF)
		assert.Equal(t, uint8(0xFF), s.Get(20)) // LSB
		assert.Equal(t, uint8(0xFF), s.Get(21)) // MSB
	})

	t.Run("zero values", func(t *testing.T) {
		s.Set16LittleEndian(30, 0x0000)
		assert.Equal(t, uint8(0x00), s.Get(30)) // LSB
		assert.Equal(t, uint8(0x00), s.Get(31)) // MSB
	})

	t.Run("out of bound access", func(t *testing.T) {
		assert.Panics(t, func() {
			s.Set16LittleEndian(-1, 0x1234)
		})

		assert.Panics(t, func() {
			s.Set16LittleEndian(99, 0x1234)
		})
	})
}

func TestSet16(t *testing.T) {
	s := NewSegment(100)

	t.Run("big endian", func(t *testing.T) {
		s.Endianness = BigEndian
		s.Set16(20, 0x1234)
		assert.Equal(t, uint8(0x12), s.Get(20)) // MSB at addr
		assert.Equal(t, uint8(0x34), s.Get(21)) // LSB at addr+1
	})

	t.Run("little endian", func(t *testing.T) {
		s.Endianness = LittleEndian
		s.Set16(30, 0x7856)
		assert.Equal(t, uint8(0x56), s.Get(30)) // LSB at addr
		assert.Equal(t, uint8(0x78), s.Get(31)) // MSB at addr+1
	})
}
