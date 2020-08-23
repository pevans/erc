package a2rec

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatOperand(t *testing.T) {
	type test struct {
		mode string
		vfn  assert.ValueAssertionFunc
	}

	cases := []test{
		{"", assert.Empty},
		{"not a real mode", assert.Empty},
		{"ACC", assert.Empty},
		{"IMP", assert.Empty},
		{"BY2", assert.Empty},
		{"BY3", assert.Empty},
		{"ABS", assert.NotEmpty},
		{"ABX", assert.NotEmpty},
		{"ABY", assert.NotEmpty},
		{"IDX", assert.NotEmpty},
		{"IDY", assert.NotEmpty},
		{"IND", assert.NotEmpty},
		{"IMM", assert.NotEmpty},
		{"REL", assert.NotEmpty},
		{"ABS", assert.NotEmpty},
		{"ZPG", assert.NotEmpty},
		{"ZPX", assert.NotEmpty},
		{"ZPY", assert.NotEmpty},
	}

	for _, c := range cases {
		t.Run(c.mode, func(tt *testing.T) {
			r := Recorder{Mode: c.mode}
			c.vfn(tt, r.FormatOperand())
		})
	}
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

	assert.NoError(t, r.Record(w))

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

	r.Operand = 0x114
	w.Reset()
	assert.NoError(t, r.Record(w))
	output = w.String()
	assert.Contains(t, output, "14 01")
}
