package mos

import "github.com/pevans/erc/a2/a2save"

// Snapshot returns a snapshot of the CPU state for serialization.
func (c *CPU) Snapshot() *a2save.CPUState {
	return &a2save.CPUState{
		PC:           c.PC,
		LastPC:       c.LastPC,
		A:            c.A,
		X:            c.X,
		Y:            c.Y,
		P:            c.P,
		S:            c.S,
		CycleCounter: c.cycleCounter,
		Opcode:       c.opcode,
		Operand:      c.Operand,
		EffAddr:      c.EffAddr,
		EffVal:       c.EffVal,
		AddrMode:     c.AddrMode,
		ReadOp:       c.ReadOp,
	}
}

// Restore restores the CPU state from a snapshot.
func (c *CPU) Restore(state *a2save.CPUState) {
	c.PC = state.PC
	c.LastPC = state.LastPC
	c.A = state.A
	c.X = state.X
	c.Y = state.Y
	c.P = state.P
	c.S = state.S
	c.cycleCounter = state.CycleCounter
	c.opcode = state.Opcode
	c.Operand = state.Operand
	c.EffAddr = state.EffAddr
	c.EffVal = state.EffVal
	c.AddrMode = state.AddrMode
	c.ReadOp = state.ReadOp
}
