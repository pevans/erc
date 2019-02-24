package mach

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntAddr(t *testing.T) {
	i := Int(123)

	assert.Equal(t, 123, i.Addr())
}

func TestByteAddr(t *testing.T) {
	b := Byte(123)

	assert.Equal(t, 123, b.Addr())
}

func TestDByteAddr(t *testing.T) {
	db := DByte(123)

	assert.Equal(t, 123, db.Addr())
}

func TestPlus(t *testing.T) {
	a := Int(123)
	b := Plus(a, 123)

	assert.Equal(t, 246, b.Addr())
}
