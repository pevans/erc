package a2rec

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatOperand(t *testing.T) {
	var r Recorder

	// If we don't know what the heck this recorder holds (whether empty
	// or weird), then we should demonstrate that we return an empty
	// string
	assert.Empty(t, r.FormatOperand())
	r.Mode = "not a real thing"
	assert.Empty(t, r.FormatOperand())

	// However, if the mode is a real address mode that an Apple II
	// understands, then we should be able to get something that looks
	// nonempty
	r.Mode = "ABS"
	assert.NotEmpty(t, r.FormatOperand())
}

func TestRecord(t *testing.T) {
	var (
		oldCounter = counter
		w          = new(strings.Builder)
		r          = Recorder{
			PC:      0x12,
			Opcode:  0x13,
			Operand: 0x14,
			Mode:    "ABS",
			Inst:    "ADC",
		}
	)

	r.Record(w)

	output := w.String()

	// There's a few things we want to see in the output based on the
	// input. There's other stuff, but we don't want to be super
	// exhaustive. We just want to see what we passed in in some form.
	assert.Contains(t, output, "ADC")  // instruction
	assert.Contains(t, output, "0012") // program counter value
	assert.Contains(t, output, "13")   // opcode
	assert.Contains(t, output, "14")   // operand

	// We should see the global (yuck) counter incremented by 1
	assert.Equal(t, counter, oldCounter+1)
}
