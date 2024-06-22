package asm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Let's not be prescriptive about how the line looks -- just look for
// things in the line
func TestLineString(t *testing.T) {
	t.Run("with instruction and operand", func(t *testing.T) {
		ln := Line{
			Instruction: "ABC",
			Operand:     "$123",
		}

		str := ln.String()

		assert.Contains(t, str, ln.Instruction)
		assert.Contains(t, str, ln.Operand)
	})

	t.Run("with comment", func(t *testing.T) {
		ln := Line{
			Instruction: "ABC",
			Operand:     "$123",
			Comment:     "comment is free",
		}

		str := ln.String()

		assert.Contains(t, str, ln.Instruction)
		assert.Contains(t, str, ln.Operand)

		// Test for the semicolon since our assembly "notation" uses
		// that as the indicator
		assert.Contains(t, str, "; "+ln.Comment)
	})
}
