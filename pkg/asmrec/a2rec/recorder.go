package a2rec

import (
	"fmt"
	"io"
	"strings"

	"github.com/pevans/erc/pkg/data"
)

// A Recorder is an object which satisfies the asmrec.Recorder interface
// and allows us to record assembly instructions for the Apple 2.
type Recorder struct {
	PrintState bool

	PC      data.DByte
	Opcode  data.Byte
	Operand data.DByte
	Mode    string
	Inst    string
	A       data.Byte
	X       data.Byte
	Y       data.Byte
	S       data.Byte
	P       data.Byte
	EffAddr data.DByte
	EffVal  data.Byte
}

var counter = 0

// FormatOperand will produce a formatted string that represents an
// operand to an instruction.
func (r *Recorder) FormatOperand() string {
	switch r.Mode {
	case "ACC", "IMP", "BY2", "BY3":
		return ""
	case "ABS":
		return fmt.Sprintf("$%04X", r.Operand)
	case "ABX":
		return fmt.Sprintf("$%04X,X", r.Operand)
	case "ABY":
		return fmt.Sprintf("$%04X,Y", r.Operand)
	case "IDX":
		return fmt.Sprintf("($%02X,X)", r.Operand)
	case "IDY":
		return fmt.Sprintf("($%02X),Y", r.Operand)
	case "IND":
		return fmt.Sprintf("($%04X)", r.Operand)
	case "IMM":
		return fmt.Sprintf("#$%02X", r.Operand)
	case "REL":
		newAddr := r.PC + r.Operand + 2

		// It's signed, so the effect of the operand should be negative w/r/t
		// two's complement.
		if r.Operand >= 0x80 {
			newAddr -= 256
		}

		return fmt.Sprintf("$%04X", newAddr)
	case "ZPG":
		return fmt.Sprintf("$%02X", r.Operand)
	case "ZPX":
		return fmt.Sprintf("$%02X,X", r.Operand)
	case "ZPY":
		return fmt.Sprintf("$%02X,Y", r.Operand)
	}

	return ""
}

// Record will print out the idealized form of an opcode-operand
// sequence as an assembly instruction.
func (r *Recorder) Record(w io.Writer) error {
	str := fmt.Sprintf(
		`%04X %02X`, r.PC, r.Opcode,
	)

	pstatus := []rune("NVUBDIZC")
	operand := r.FormatOperand()

	// If it's greater than 255, then we have two-byte operand, so print
	// the MSB now.
	if r.Operand > 0xFF {
		str += fmt.Sprintf(` %02X %02X`, r.Operand&0xFF, r.Operand>>8)
	} else if r.Operand > 0x00 {
		str += fmt.Sprintf(` %02X`, r.Operand&0xFF)
	}

	if len(str) < 13 {
		str += strings.Repeat(" ", 13-len(str))
	}

	for i := 7; i >= 0; i-- {
		bit := (r.P >> uint(i)) & 1
		if bit == 0 {
			pstatus[7-i] = '.'
		}
	}

	str += fmt.Sprintf(` %3s %7s`, r.Inst, operand)

	if r.PrintState {
		str += fmt.Sprintf(
			" ; A=%02X X=%02X Y=%02X P=%02X S=%02X (%s) EA=%04X EV=%02X +%d",
			r.A, r.X, r.Y, r.P, r.S, string(pstatus),
			r.EffAddr, r.EffVal, counter,
		)
	}

	str += "\n"

	counter++

	_, err := fmt.Fprint(w, str)
	return err
}
