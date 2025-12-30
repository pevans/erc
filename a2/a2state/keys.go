package a2state

import "fmt"

const (
	noop = iota

	BankDFBlockBank2
	BankROMSegment
	BankReadAttempts
	BankReadRAM
	BankSysBlockAux
	BankSysBlockSegment
	BankWriteRAM
	Debugger
	DebuggerLookAhead
	DebugImage
	DiskComputer
	DisplayAltChar
	DisplayAuxSegment
	DisplayCol80
	DisplayDoubleHigh
	DisplayHires
	DisplayIou
	DisplayMixed
	DisplayPage2
	DisplayRedraw
	DisplayStore80
	DisplayText
	InstructionReadOp
	KBKeyDown
	KBLastKey
	KBStrobe
	MemAuxSegment
	MemMainSegment
	MemReadAux
	MemReadSegment
	MemWriteAux
	MemWriteSegment
	PCExpSlot
	PCExpansion
	PCIOSelect
	PCIOStrobe
	PCROMSegment
	PCSlotC3
	PCSlotCX
)

var keyStringMap = map[int]string{
	BankDFBlockBank2:    "BankDFBlockBank2",
	BankROMSegment:      "BankROMSegment",
	BankReadAttempts:    "BankReadAttempts",
	BankReadRAM:         "BankReadRAM",
	BankSysBlockAux:     "BankSysBlockAux",
	BankSysBlockSegment: "BankSysBlockSegment",
	BankWriteRAM:        "BankWriteRAM",
	Debugger:            "Debugger",
	DebuggerLookAhead:   "DebuggerLookAhead",
	DebugImage:          "DebugImage",
	DiskComputer:        "DiskComputer",
	DisplayAltChar:      "DisplayAltChar",
	DisplayAuxSegment:   "DisplayAuxSegment",
	DisplayCol80:        "DisplayCol80",
	DisplayDoubleHigh:   "DisplayDoubleHigh",
	DisplayHires:        "DisplayHires",
	DisplayIou:          "DisplayIou",
	DisplayMixed:        "DisplayMixed",
	DisplayPage2:        "DisplayPage2",
	DisplayRedraw:       "DisplayRedraw",
	DisplayStore80:      "DisplayStore80",
	DisplayText:         "DisplayText",
	InstructionReadOp:   "InstructionReadOp",
	KBKeyDown:           "KBKeyDown",
	KBLastKey:           "KBLastkey",
	KBStrobe:            "KBStrobe",
	MemAuxSegment:       "MemAuxSegment",
	MemMainSegment:      "MemMainSegment",
	MemReadAux:          "MemReadAux",
	MemReadSegment:      "MemReadSegment",
	MemWriteAux:         "MemWriteAux",
	MemWriteSegment:     "MemWriteSegment",
	PCExpSlot:           "PCExpSlot",
	PCExpansion:         "PCExpansion",
	PCIOSelect:          "PCIOSelect",
	PCIOStrobe:          "PCIOStrobe",
	PCROMSegment:        "PCROMSegment",
	PCSlotC3:            "PCSlotC3",
	PCSlotCX:            "PCSlotCX",
}

func KeyToString(key any) string {
	intKey, ok := key.(int)
	if !ok {
		return fmt.Sprintf("unknown non-int key (%v)", key)
	}

	if name, ok := keyStringMap[intKey]; ok {
		return name
	}

	return fmt.Sprintf("unknown (%v)", key)
}
