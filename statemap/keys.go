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

func KeyToString(key any) string {
	intKey, ok := key.(int)
	if !ok {
		return fmt.Sprintf("unknown non-int key (%v)", key)
	}

	switch intKey {
	case BankRead:
		return "BankRead"
	case BankWrite:
		return "BankWrite"
	case BankDFBlock:
		return "BankDFBlock"
	case BankSysBlock:
		return "BankSysBlock"
	case BankReadAttempts:
		return "BankReadAttempts"
	case BankSysBlockSegment:
		return "BankSysBlockSegment"
	case BankROMSegment:
		return "BankROMSegment"
	case Debugger:
		return "Debugger"
	case DebuggerLookAhead:
		return "DebuggerLookAhead"
	case DiskComputer:
		return "DiskComputer"
	case DisplayAltChar:
		return "DisplayAltChar"
	case DisplayCol80:
		return "DisplayCol80"
	case DisplayStore80:
		return "DisplayStore80"
	case DisplayPage2:
		return "DisplayPage2"
	case DisplayText:
		return "DisplayText"
	case DisplayMixed:
		return "DisplayMixed"
	case DisplayHires:
		return "DisplayHires"
	case DisplayIou:
		return "DisplayIou"
	case DisplayDoubleHigh:
		return "DisplayDoubleHigh"
	case DisplayRedraw:
		return "DisplayRedraw"
	case DisplayAuxSegment:
		return "DisplayAuxSegment"
	case InstructionReadOp:
		return "InstructionReadOp"
	case KBLastKey:
		return "KBLastkey"
	case KBStrobe:
		return "KBStrobe"
	case KBKeyDown:
		return "KBKeyDown"
	case MemRead:
		return "MemRead"
	case MemWrite:
		return "MemWrite"
	case MemReadSegment:
		return "MemReadSegment"
	case MemWriteSegment:
		return "MemWriteSegment"
	case MemAuxSegment:
		return "MemAuxSegment"
	case MemMainSegment:
		return "MemMainSegment"
	case PCExpansion:
		return "PCExpansion"
	case PCExpSlot:
		return "PCExpSlot"
	case PCIOSelect:
		return "PCIOSelect"
	case PCIOStrobe:
		return "PCIOStrobe"
	case PCROMSegment:
		return "PCROMSegment"
	case PCSlotC3:
		return "PCSlotC3"
	case PCSlotCX:
		return "PCSlotCX"
	}

	return fmt.Sprintf("unknown (%v)", key)
}
