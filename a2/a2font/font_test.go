package a2font_test

import (
	"testing"

	"github.com/pevans/erc/a2/a2font"
	"github.com/stretchr/testify/assert"
)

func TestSystemFont40(t *testing.T) {
	assert.NotNil(t, a2font.SystemFont40())
}
