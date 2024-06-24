package statemap

import "fmt"

const (
	noop = iota

	BankDFBlock
	BankROMSegment
	BankRead
	BankReadAttempts
	BankSysBlock
	BankSysBlockSegment
	BankWrite
	Debugger
	DebuggerLookAhead
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
	MemRead
	MemReadSegment
	MemWrite
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
	BankRead:            "BankRead",
	BankWrite:           "BankWrite",
	BankDFBlock:         "BankDFBlock",
	BankSysBlock:        "BankSysBlock",
	BankReadAttempts:    "BankReadAttempts",
	BankSysBlockSegment: "BankSysBlockSegment",
	BankROMSegment:      "BankROMSegment",
	Debugger:            "Debugger",
	DebuggerLookAhead:   "DebuggerLookAhead",
	DiskComputer:        "DiskComputer",
	DisplayAltChar:      "DisplayAltChar",
	DisplayCol80:        "DisplayCol80",
	DisplayStore80:      "DisplayStore80",
	DisplayPage2:        "DisplayPage2",
	DisplayText:         "DisplayText",
	DisplayMixed:        "DisplayMixed",
	DisplayHires:        "DisplayHires",
	DisplayIou:          "DisplayIou",
	DisplayDoubleHigh:   "DisplayDoubleHigh",
	DisplayRedraw:       "DisplayRedraw",
	DisplayAuxSegment:   "DisplayAuxSegment",
	InstructionReadOp:   "InstructionReadOp",
	KBLastKey:           "KBLastkey",
	KBStrobe:            "KBStrobe",
	KBKeyDown:           "KBKeyDown",
	MemRead:             "MemRead",
	MemWrite:            "MemWrite",
	MemReadSegment:      "MemReadSegment",
	MemWriteSegment:     "MemWriteSegment",
	MemAuxSegment:       "MemAuxSegment",
	MemMainSegment:      "MemMainSegment",
	PCExpansion:         "PCExpansion",
	PCExpSlot:           "PCExpSlot",
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
