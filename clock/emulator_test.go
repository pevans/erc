package clock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmulator(t *testing.T) {
	assert.NotNil(t, NewEmulator(1))
}
