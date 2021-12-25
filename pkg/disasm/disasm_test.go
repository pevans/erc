package disasm

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(os.Stdout)
	shutdown <- true

	assert.NotNil(t, dislog)
}

func TestAvailable(t *testing.T) {
	Init(os.Stdout)

	assert.True(t, Available())
	shutdown <- true
	dislog = nil
	assert.False(t, Available())
}

func TestShutdown(t *testing.T) {
	b := new(strings.Builder)

	Init(b)
	dislog.Printf("test")
	Shutdown()
	assert.Nil(t, dislog)
	assert.Contains(t, b.String(), "test")

	b.Reset()
	Shutdown()
	assert.NotContains(t, b.String(), "test")
}

func TestMap(t *testing.T) {
	b := new(strings.Builder)

	Init(b)
	Map(123, "test123")
	Shutdown()

	// Note that Map() does NOT guarantee the address is written to the
	// channel log. That's up to the disassembler logic that calls
	// Map().
	assert.Contains(t, b.String(), "test123")
}
