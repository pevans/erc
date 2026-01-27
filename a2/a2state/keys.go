package a2state

import "fmt"

const (
	BankDFBlockBank2 = iota
	BankROMSegment
	BankReadAttempts
	BankReadRAM
	BankSysBlockAux
	BankSysBlockSegment
	BankWriteRAM
	Debugger
	DebuggerLookAhead
	DebugImage
	Computer
	DisplayAltChar
	DisplayAuxSegment
	DisplayCol80
	DisplayDoubleHigh
	DisplayHires
	DisplayIou
	DisplayMixed
	DisplayPage2
	DisplayMonochrome
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
	Paused
	PCExpSlot
	PCExpansion
	PCIOSelect
	PCIOStrobe
	PCROMSegment
	PCSlotC3
	PCSlotCX
	SpeakerState
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
	Computer:            "Computer",
	DisplayAltChar:      "DisplayAltChar",
	DisplayAuxSegment:   "DisplayAuxSegment",
	DisplayCol80:        "DisplayCol80",
	DisplayDoubleHigh:   "DisplayDoubleHigh",
	DisplayHires:        "DisplayHires",
	DisplayIou:          "DisplayIou",
	DisplayMixed:        "DisplayMixed",
	DisplayMonochrome:   "DisplayMonochrome",
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
	Paused:              "Paused",
	PCExpSlot:           "PCExpSlot",
	PCExpansion:         "PCExpansion",
	PCIOSelect:          "PCIOSelect",
	PCIOStrobe:          "PCIOStrobe",
	PCROMSegment:        "PCROMSegment",
	PCSlotC3:            "PCSlotC3",
	PCSlotCX:            "PCSlotCX",
	SpeakerState:        "SpeakerState",
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
