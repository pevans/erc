package a2drive

import (
	"testing"

	"github.com/pevans/erc/a2/a2enc"
	"github.com/stretchr/testify/assert"
)

func TestNewDrive(t *testing.T) {
	d := NewDrive()

	assert.NotNil(t, d)
	assert.True(t, d.ReadMode())
	assert.Equal(t, a2enc.DOS33, d.imageType)
}

func TestMotorOn(t *testing.T) {
	d := NewDrive()

	assert.NotNil(t, d)
	assert.False(t, d.MotorOn())

	d.StartMotor()
	assert.True(t, d.MotorOn())

	d.StopMotor()
	assert.False(t, d.MotorOn())
}
