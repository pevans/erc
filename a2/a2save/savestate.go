// Package a2save provides types for serializing and deserializing emulator
// state for save-state functionality.
package a2save

// SaveStateVersion is the current version of the save state format. Increment
// this when making breaking changes to the format.
const SaveStateVersion = 1

// SaveState is the top-level container for all emulator state.
type SaveState struct {
	Version       int
	CPU           CPUState
	Main          []uint8
	Aux           []uint8
	StateFlags    StateFlags
	Drive1        DriveState
	Drive2        DriveState
	SelectedDrive int // 1 or 2
	DiskSet       DiskSetState
	Speed         int
}

// CPUState captures all CPU register and internal state.
type CPUState struct {
	PC           uint16
	LastPC       uint16
	A            uint8
	X            uint8
	Y            uint8
	P            uint8
	S            uint8
	CycleCounter uint64
	Opcode       uint8
	Operand      uint16
	EffAddr      uint16
	EffVal       uint8
	AddrMode     int
	ReadOp       bool
}

// StateFlags captures boolean and integer state from the StateMap.
// Only includes values needed for restoration (excludes Segment pointers).
type StateFlags struct {
	// Bank state
	BankDFBlockBank2 bool
	BankReadAttempts int
	BankReadRAM      bool
	BankWriteRAM     bool
	BankSysBlockAux  bool

	// Display state
	DisplayAltChar    bool
	DisplayCol80      bool
	DisplayDoubleHigh bool
	DisplayHires      bool
	DisplayIou        bool
	DisplayMixed      bool
	DisplayPage2      bool
	DisplayStore80    bool
	DisplayText       bool

	// Keyboard state
	KBKeyDown uint8
	KBLastKey uint8
	KBStrobe  uint8

	// Memory state
	MemReadAux  bool
	MemWriteAux bool

	// PC/Slot state
	PCExpSlot   int
	PCExpansion bool
	PCIOSelect  bool
	PCIOStrobe  bool
	PCSlotC3    bool
	PCSlotCX    bool
}

// DriveState captures all floppy drive state.
type DriveState struct {
	MotorOn      bool
	Phase        int
	TrackPos     int
	SectorPos    int
	Latch        uint8
	Mode         int
	LatchWasRead bool
	DiskShifted  bool
	ImageName    string
	ImageType    int
	WriteProtect bool
	HasDisk      bool
	ImageData    []uint8
	PhysicalData []uint8
}

// DiskSetState captures the disk set configuration.
type DiskSetState struct {
	Images  []string
	Current int
}
