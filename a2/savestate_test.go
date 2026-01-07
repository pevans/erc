package a2

import (
	"os"
	"path/filepath"

	"github.com/pevans/erc/a2/a2state"
)

func (s *a2Suite) TestSaveStateCreatesFile() {
	tmpDir := s.T().TempDir()
	filename := filepath.Join(tmpDir, "test.state")

	err := s.comp.SaveState(filename)
	s.NoError(err)

	_, err = os.Stat(filename)
	s.NoError(err)
}

func (s *a2Suite) TestSaveStateLoadStateRoundTrip() {
	tmpDir := s.T().TempDir()
	filename := filepath.Join(tmpDir, "test.state")

	// Set some distinctive CPU state
	s.comp.CPU.PC = 0x1234
	s.comp.CPU.A = 0xAB
	s.comp.CPU.X = 0xCD
	s.comp.CPU.Y = 0xEF
	s.comp.CPU.S = 0x42
	s.comp.CPU.P = 0x37

	// Set some memory (use DirectSet for Aux to bypass soft switches)
	s.comp.Main.Set(0x1000, 0xDE)
	s.comp.Main.Set(0x1001, 0xAD)
	s.comp.Aux.DirectSet(0x2000, 0xBE)
	s.comp.Aux.DirectSet(0x2001, 0xEF)

	// Set some state flags
	s.comp.State.SetBool(a2state.DisplayHires, true)
	s.comp.State.SetBool(a2state.MemReadAux, true)
	s.comp.State.SetUint8(a2state.KBLastKey, 0x41)

	// Set speed
	s.comp.SetSpeed(3)

	// Save state
	err := s.comp.SaveState(filename)
	s.NoError(err)

	// Create a fresh computer
	newComp := NewComputer(1)
	_ = newComp.Boot()

	// Verify the new computer has different state
	s.NotEqual(uint16(0x1234), newComp.CPU.PC)
	s.NotEqual(uint8(0xAB), newComp.CPU.A)

	// Load state into new computer
	err = newComp.LoadState(filename)
	s.NoError(err)

	// Verify CPU state was restored
	s.Equal(uint16(0x1234), newComp.CPU.PC)
	s.Equal(uint8(0xAB), newComp.CPU.A)
	s.Equal(uint8(0xCD), newComp.CPU.X)
	s.Equal(uint8(0xEF), newComp.CPU.Y)
	s.Equal(uint8(0x42), newComp.CPU.S)
	s.Equal(uint8(0x37), newComp.CPU.P)

	// Verify memory was restored (use DirectGet for Aux to bypass soft switches)
	s.Equal(uint8(0xDE), newComp.Main.Get(0x1000))
	s.Equal(uint8(0xAD), newComp.Main.Get(0x1001))
	s.Equal(uint8(0xBE), newComp.Aux.DirectGet(0x2000))
	s.Equal(uint8(0xEF), newComp.Aux.DirectGet(0x2001))

	// Verify state flags were restored
	s.True(newComp.State.Bool(a2state.DisplayHires))
	s.True(newComp.State.Bool(a2state.MemReadAux))
	s.Equal(uint8(0x41), newComp.State.Uint8(a2state.KBLastKey))

	// Verify speed was restored
	s.Equal(3, newComp.speed)
}

func (s *a2Suite) TestSaveStateSelectedDrive() {
	tmpDir := s.T().TempDir()
	filename := filepath.Join(tmpDir, "test.state")

	// Select drive 2
	s.comp.SelectedDrive = s.comp.Drive2

	err := s.comp.SaveState(filename)
	s.NoError(err)

	newComp := NewComputer(1)
	_ = newComp.Boot()

	// Should default to drive 1
	s.Equal(newComp.Drive1, newComp.SelectedDrive)

	err = newComp.LoadState(filename)
	s.NoError(err)

	// Should be drive 2 after load
	s.Equal(newComp.Drive2, newComp.SelectedDrive)
}

func (s *a2Suite) TestLoadStateFileNotFound() {
	err := s.comp.LoadState("/nonexistent/path/to/file.state")
	s.Error(err)
}

func (s *a2Suite) TestLoadStateInvalidFile() {
	tmpDir := s.T().TempDir()
	filename := filepath.Join(tmpDir, "invalid.state")

	// Write some garbage data
	err := os.WriteFile(filename, []byte("not a valid gob file"), 0o644)
	s.NoError(err)

	err = s.comp.LoadState(filename)
	s.Error(err)
}

func (s *a2Suite) TestDiskSetSnapshotRestore() {
	set := NewDiskSet()
	set.images = []string{"disk1.dsk", "disk2.dsk", "disk3.dsk"}
	set.current = 2

	snapshot := set.Snapshot()
	s.Equal([]string{"disk1.dsk", "disk2.dsk", "disk3.dsk"}, snapshot.Images)
	s.Equal(2, snapshot.Current)

	newSet := NewDiskSet()
	newSet.Restore(snapshot)

	s.Equal([]string{"disk1.dsk", "disk2.dsk", "disk3.dsk"}, newSet.images)
	s.Equal(2, newSet.current)
}

func (s *a2Suite) TestSegmentReferencesRebuilt() {
	tmpDir := s.T().TempDir()
	filename := filepath.Join(tmpDir, "test.state")

	// Set MemReadAux to true
	s.comp.State.SetBool(a2state.MemReadAux, true)
	s.comp.State.SetSegment(a2state.MemReadSegment, s.comp.Aux)

	err := s.comp.SaveState(filename)
	s.NoError(err)

	newComp := NewComputer(1)
	_ = newComp.Boot()

	err = newComp.LoadState(filename)
	s.NoError(err)

	// The segment reference should be rebuilt based on the flag
	s.True(newComp.State.Bool(a2state.MemReadAux))
	s.Equal(newComp.Aux, newComp.State.Segment(a2state.MemReadSegment))
}
