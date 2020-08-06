package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntAddr(t *testing.T) {
	assert.Equal(t, 123, Int(123).Addr())
}

func TestByteAddr(t *testing.T) {
	type test struct {
		i    int
		want int
	}

	cases := map[string]test{
		"zero":     {i: 0, want: 0},
		"123":      {i: 123, want: 123},
		"overflow": {i: 258, want: 2},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			assert.Equal(t, c.want, Byte(c.i).Addr())
		})
	}
}

func TestDByteAddr(t *testing.T) {
	type test struct {
		i    int
		want int
	}

	cases := map[string]test{
		"zero":     {i: 0, want: 0},
		"123":      {i: 123, want: 123},
		"256":      {i: 256, want: 256},
		"overflow": {i: 65537, want: 1},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			assert.Equal(t, c.want, DByte(c.i).Addr())
		})
	}
}

func TestPlus(t *testing.T) {
	var (
		n257   = 257
		n65537 = 65537
	)

	assert.Equal(t, 246, Plus(Int(123), 123).Addr())

	// The interesting thing here is I want to test the behavior around
	// overflows. But Go is pretty smart -- it knows Byte is uint8, and
	// if you try Byte(257), it will catch you at compilation time. So
	// the way around this is to assign 257 to some integer variable,
	// and do a constructor of that -- Byte(something).
	assert.Equal(t, 124, Plus(Byte(n257), 123).Addr())
	assert.Equal(t, 124, Plus(DByte(n65537), 123).Addr())
}
