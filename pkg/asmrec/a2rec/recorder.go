package a2rec

import (
	"fmt"
	"io"
	"strings"

	"github.com/pevans/erc/pkg/data"
)

type Recorder struct {
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
		rel := int8(r.Operand)
		return fmt.Sprintf("%d", rel)
	case "ZPG":
		return fmt.Sprintf("$%02X", r.Operand)
	case "ZPX":
		return fmt.Sprintf("$%02X,X", r.Operand)
	case "ZPY":
		return fmt.Sprintf("$%02X,Y", r.Operand)
	}

	return ""
}

func (r *Recorder) Record(w io.Writer) error {
	str := fmt.Sprintf(
		`%04X %02X`, r.PC, r.Opcode,
	)

	operand := r.FormatOperand()

	// If it's greater than 255, then we have two-byte operand, so print
	// the MSB now.
	if r.Operand > 0xFF {
		str += fmt.Sprintf(` %02X`, r.Operand>>8)
	}

	// Even if operand is zero, we may still need to print something out
	// for it. So the next check actually depends on whether we have a
	// printed operand or not.
	if operand != "" {
		str += fmt.Sprintf(` %02X`, r.Operand&0xFF)
	}

	if len(str) < 13 {
		str += strings.Repeat(" ", 13-len(str))
	}

	str += fmt.Sprintf(` %3s %7s`, r.Inst, operand)
	str += fmt.Sprintf(
		" ; A=%02X X=%02X Y=%02X S=%02X P=%02X EffAddr=%04X EffVal=%02X\n",
		r.A, r.X, r.Y, r.S, r.P, r.EffAddr, r.EffVal,
	)

	_, err := fmt.Fprint(w, str)
	return err
}
