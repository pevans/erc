package mach

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteAddr(t *testing.T) {
	b := Byte(123)

	assert.Equal(t, 123, b.Addr())
}

func TestDByteAddr(t *testing.T) {
	db := DByte(123)

	assert.Equal(t, 123, db.Addr())
}
