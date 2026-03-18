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
	Computer
	DebugImage
	Debugger
	DebuggerLookAhead
	DiskIndex
	DisplayAltChar
	DisplayAuxSegment
	DisplayCol80
	DisplayDoubleHigh
	DisplayHires
	DisplayIou
	DisplayMixed
	DisplayMonochrome
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
	Paused
	SpeakerState
	Speed
	StateSlot
	VolumeLevel
	VolumeMuted
	WriteProtect
)

var keyStringMap = map[int]string{
	BankDFBlockBank2:    "BankDFBlockBank2",
	BankROMSegment:      "BankROMSegment",
	BankReadAttempts:    "BankReadAttempts",
	BankReadRAM:         "BankReadRAM",
	BankSysBlockAux:     "BankSysBlockAux",
	BankSysBlockSegment: "BankSysBlockSegment",
	BankWriteRAM:        "BankWriteRAM",
	Computer:            "Computer",
	DebugImage:          "DebugImage",
	Debugger:            "Debugger",
	DebuggerLookAhead:   "DebuggerLookAhead",
	DiskIndex:           "DiskIndex",
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
	PCExpSlot:           "PCExpSlot",
	PCExpansion:         "PCExpansion",
	PCIOSelect:          "PCIOSelect",
	PCIOStrobe:          "PCIOStrobe",
	PCROMSegment:        "PCROMSegment",
	PCSlotC3:            "PCSlotC3",
	PCSlotCX:            "PCSlotCX",
	Paused:              "Paused",
	SpeakerState:        "SpeakerState",
	Speed:               "Speed",
	StateSlot:           "StateSlot",
	VolumeLevel:         "VolumeLevel",
	VolumeMuted:         "VolumeMuted",
	WriteProtect:        "WriteProtect",
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
