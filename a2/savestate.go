package a2

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/pevans/erc/a2/a2save"
	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/gfx"
	"github.com/pevans/erc/obj"
)

// pauseForStateOp pauses the emulator and waits for the ProcessLoop to
// actually stop before returning. Returns whether the emulator was already
// paused (so caller knows whether to unpause).
func (c *Computer) pauseForStateOp() bool {
	wasPaused := c.State.Bool(a2state.Paused)
	if !wasPaused {
		c.State.SetBool(a2state.Paused, true)

		// Give ProcessLoop time to see the pause and stop
		time.Sleep(150 * time.Millisecond)
	}

	return wasPaused
}

// resumeAfterStateOp unpauses the emulator if it wasn't paused before.
func (c *Computer) resumeAfterStateOp(wasPaused bool) {
	if !wasPaused {
		c.State.SetBool(a2state.Paused, false)
	}
}

// SaveState writes the current emulator state to the specified file.
func (c *Computer) SaveState(filename string) error {
	wasPaused := c.pauseForStateOp()
	defer c.resumeAfterStateOp(wasPaused)

	state := &a2save.SaveState{
		Version:    a2save.SaveStateVersion,
		CPU:        *c.CPU.Snapshot(),
		Main:       c.Main.Bytes(),
		Aux:        c.Aux.Bytes(),
		StateFlags: *c.snapshotStateFlags(),
		Drive1:     *c.Drive1.Snapshot(),
		Drive2:     *c.Drive2.Snapshot(),
		Speed:      c.speed,
		DiskSet:    *c.Disks.Snapshot(),
	}

	if c.SelectedDrive == c.Drive1 {
		state.SelectedDrive = 1
	} else {
		state.SelectedDrive = 2
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create save state file: %w", err)
	}
	defer file.Close() //nolint:errcheck

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(state); err != nil {
		return fmt.Errorf("could not encode save state: %w", err)
	}

	return nil
}

// LoadState restores emulator state from the specified file.
func (c *Computer) LoadState(filename string) error {
	wasPaused := c.pauseForStateOp()
	defer c.resumeAfterStateOp(wasPaused)

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open save state file: %w", err)
	}
	defer file.Close() //nolint:errcheck

	var state a2save.SaveState
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&state); err != nil {
		return fmt.Errorf("could not decode save state: %w", err)
	}

	if state.Version > a2save.SaveStateVersion {
		return fmt.Errorf(
			"save state version %d not supported (max supported: %d)",
			state.Version, a2save.SaveStateVersion,
		)
	}

	// Restore CPU
	c.CPU.Restore(&state.CPU)

	// Restore memory
	if err := c.Main.RestoreBytes(state.Main); err != nil {
		return fmt.Errorf("could not restore main memory: %w", err)
	}

	if err := c.Aux.RestoreBytes(state.Aux); err != nil {
		return fmt.Errorf("could not restore aux memory: %w", err)
	}

	// Restore state flags
	c.restoreStateFlags(&state.StateFlags)

	// Restore drives
	if err := c.Drive1.Restore(&state.Drive1); err != nil {
		return fmt.Errorf("could not restore drive 1: %w", err)
	}

	if err := c.Drive2.Restore(&state.Drive2); err != nil {
		return fmt.Errorf("could not restore drive 2: %w", err)
	}

	// Restore selected drive
	if state.SelectedDrive == 2 {
		c.SelectedDrive = c.Drive2
	} else {
		c.SelectedDrive = c.Drive1
	}

	// Restore disk set
	c.Disks.Restore(&state.DiskSet)

	// Restore speed
	c.SetSpeed(state.Speed)

	// Rebuild segment references in StateMap
	c.rebuildSegmentReferences()

	// Clear keyboard state to avoid stuck keys from the loaded state
	c.ClearKeys()

	// Force display redraw
	c.State.SetBool(a2state.DisplayRedraw, true)

	return nil
}

// SetStateSlot sets the state slot for the computer and shows a message
// on-screen to indicate the slot.
func (c *Computer) SetStateSlot(n int) {
	c.stateSlot = n
	c.ShowText(fmt.Sprintf("use state slot: %v", n))
}

// SaveStateSlot saves state in a standardized filename based on the first
// disk in the diskset and our current state slot
func (c *Computer) SaveStateSlot() error {
	err := c.SaveState(fmt.Sprintf("%v.%v.state", c.Disks.Name(), c.stateSlot))
	if err == nil {
		_ = gfx.ShowStatus(obj.StateSavePNG())
		c.ShowText(fmt.Sprintf("saved state: %v", c.stateSlot))
	}

	return err
}

// LoadStateSlot loads state from a standardized filename based on the first
// disk in the diskset and our current state slot
func (c *Computer) LoadStateSlot() error {
	err := c.LoadState(fmt.Sprintf("%v.%v.state", c.Disks.Name(), c.stateSlot))
	if err == nil {
		_ = gfx.ShowStatus(obj.StateLoadPNG())
		c.ShowText(fmt.Sprintf("loaded state: %v", c.stateSlot))
	}

	return err
}

// snapshotStateFlags captures the boolean and integer state flags from the StateMap.
func (c *Computer) snapshotStateFlags() *a2save.StateFlags {
	return &a2save.StateFlags{
		// Bank state
		BankDFBlockBank2: c.State.Bool(a2state.BankDFBlockBank2),
		BankReadAttempts: c.State.Int(a2state.BankReadAttempts),
		BankReadRAM:      c.State.Bool(a2state.BankReadRAM),
		BankWriteRAM:     c.State.Bool(a2state.BankWriteRAM),
		BankSysBlockAux:  c.State.Bool(a2state.BankSysBlockAux),

		// Display state
		DisplayAltChar:    c.State.Bool(a2state.DisplayAltChar),
		DisplayCol80:      c.State.Bool(a2state.DisplayCol80),
		DisplayDoubleHigh: c.State.Bool(a2state.DisplayDoubleHigh),
		DisplayHires:      c.State.Bool(a2state.DisplayHires),
		DisplayIou:        c.State.Bool(a2state.DisplayIou),
		DisplayMixed:      c.State.Bool(a2state.DisplayMixed),
		DisplayPage2:      c.State.Bool(a2state.DisplayPage2),
		DisplayStore80:    c.State.Bool(a2state.DisplayStore80),
		DisplayText:       c.State.Bool(a2state.DisplayText),

		// Keyboard state
		KBKeyDown: c.State.Uint8(a2state.KBKeyDown),
		KBLastKey: c.State.Uint8(a2state.KBLastKey),
		KBStrobe:  c.State.Uint8(a2state.KBStrobe),

		// Memory state
		MemReadAux:  c.State.Bool(a2state.MemReadAux),
		MemWriteAux: c.State.Bool(a2state.MemWriteAux),

		// PC/Slot state
		PCExpSlot:   c.State.Int(a2state.PCExpSlot),
		PCExpansion: c.State.Bool(a2state.PCExpansion),
		PCIOSelect:  c.State.Bool(a2state.PCIOSelect),
		PCIOStrobe:  c.State.Bool(a2state.PCIOStrobe),
		PCSlotC3:    c.State.Bool(a2state.PCSlotC3),
		PCSlotCX:    c.State.Bool(a2state.PCSlotCX),
	}
}

// restoreStateFlags restores the boolean and integer state flags to the StateMap.
func (c *Computer) restoreStateFlags(flags *a2save.StateFlags) {
	// Bank state
	c.State.SetBool(a2state.BankDFBlockBank2, flags.BankDFBlockBank2)
	c.State.SetInt(a2state.BankReadAttempts, flags.BankReadAttempts)
	c.State.SetBool(a2state.BankReadRAM, flags.BankReadRAM)
	c.State.SetBool(a2state.BankWriteRAM, flags.BankWriteRAM)
	c.State.SetBool(a2state.BankSysBlockAux, flags.BankSysBlockAux)

	// Display state
	c.State.SetBool(a2state.DisplayAltChar, flags.DisplayAltChar)
	c.State.SetBool(a2state.DisplayCol80, flags.DisplayCol80)
	c.State.SetBool(a2state.DisplayDoubleHigh, flags.DisplayDoubleHigh)
	c.State.SetBool(a2state.DisplayHires, flags.DisplayHires)
	c.State.SetBool(a2state.DisplayIou, flags.DisplayIou)
	c.State.SetBool(a2state.DisplayMixed, flags.DisplayMixed)
	c.State.SetBool(a2state.DisplayPage2, flags.DisplayPage2)
	c.State.SetBool(a2state.DisplayStore80, flags.DisplayStore80)
	c.State.SetBool(a2state.DisplayText, flags.DisplayText)

	// Keyboard state
	c.State.SetUint8(a2state.KBKeyDown, flags.KBKeyDown)
	c.State.SetUint8(a2state.KBLastKey, flags.KBLastKey)
	c.State.SetUint8(a2state.KBStrobe, flags.KBStrobe)

	// Memory state
	c.State.SetBool(a2state.MemReadAux, flags.MemReadAux)
	c.State.SetBool(a2state.MemWriteAux, flags.MemWriteAux)

	// PC/Slot state
	c.State.SetInt(a2state.PCExpSlot, flags.PCExpSlot)
	c.State.SetBool(a2state.PCExpansion, flags.PCExpansion)
	c.State.SetBool(a2state.PCIOSelect, flags.PCIOSelect)
	c.State.SetBool(a2state.PCIOStrobe, flags.PCIOStrobe)
	c.State.SetBool(a2state.PCSlotC3, flags.PCSlotC3)
	c.State.SetBool(a2state.PCSlotCX, flags.PCSlotCX)
}

// rebuildSegmentReferences restores the Segment pointers in StateMap
// based on the current boolean state.
func (c *Computer) rebuildSegmentReferences() {
	// Core segment references
	c.State.SetSegment(a2state.MemMainSegment, c.Main)
	c.State.SetSegment(a2state.MemAuxSegment, c.Aux)
	c.State.SetSegment(a2state.DisplayAuxSegment, c.Aux)
	c.State.SetSegment(a2state.BankROMSegment, c.ROM)
	c.State.SetSegment(a2state.PCROMSegment, c.ROM)

	// Rebuild read/write segment pointers based on flags
	if c.State.Bool(a2state.MemReadAux) {
		c.State.SetSegment(a2state.MemReadSegment, c.Aux)
	} else {
		c.State.SetSegment(a2state.MemReadSegment, c.Main)
	}

	if c.State.Bool(a2state.MemWriteAux) {
		c.State.SetSegment(a2state.MemWriteSegment, c.Aux)
	} else {
		c.State.SetSegment(a2state.MemWriteSegment, c.Main)
	}

	if c.State.Bool(a2state.BankSysBlockAux) {
		c.State.SetSegment(a2state.BankSysBlockSegment, c.Aux)
	} else {
		c.State.SetSegment(a2state.BankSysBlockSegment, c.Main)
	}

	// Restore DiskComputer reference
	c.State.SetAny(a2state.DiskComputer, c)
}

// Snapshot returns a snapshot of the DiskSet for serialization.
func (set *DiskSet) Snapshot() *a2save.DiskSetState {
	images := make([]string, len(set.images))
	copy(images, set.images)

	return &a2save.DiskSetState{
		Images:  images,
		Current: set.current,
	}
}

// Restore restores the DiskSet from a snapshot.
func (set *DiskSet) Restore(state *a2save.DiskSetState) {
	set.images = make([]string, len(state.Images))
	copy(set.images, state.Images)
	set.current = state.Current
}
