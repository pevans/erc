package elog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Let's not be prescriptive about how the line looks -- just look for things
// in the line
func TestInstructionString(t *testing.T) {
	t.Run("with instruction and operand", func(t *testing.T) {
		ln := Instruction{
			Instruction:     "ABC",
			PreparedOperand: "$123",
		}

		str := ln.String()

		assert.Contains(t, str, ln.Instruction)
		assert.Contains(t, str, ln.PreparedOperand)
	})

	t.Run("with comment", func(t *testing.T) {
		ln := Instruction{
			Instruction:     "ABC",
			PreparedOperand: "$123",
			Comment:         "comment is free",
		}

		str := ln.String()

		assert.Contains(t, str, ln.Instruction)
		assert.Contains(t, str, ln.PreparedOperand)

		// Test for the semicolon since our assembly "notation" uses that as
		// the indicator
		assert.Contains(t, str, ln.Comment)
	})
}
