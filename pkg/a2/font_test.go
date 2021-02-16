package a2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemFont(t *testing.T) {
	assert.NotNil(t, SystemFont())
}
