package a2display

import (
	"testing"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/suite"
)

type snapshotSuite struct {
	suite.Suite

	main  *memory.Segment
	aux   *memory.Segment
	state *memory.StateMap
	snap  *Snapshot
}

func (s *snapshotSuite) SetupTest() {
	s.main = memory.NewSegment(0x10000)
	s.aux = memory.NewSegment(0x10000)
	s.state = memory.NewStateMap()
	s.snap = NewSnapshot()

	s.state.SetBool(a2state.DisplayPage2, false)
	s.state.SetBool(a2state.DisplayStore80, false)
	s.state.SetBool(a2state.DisplayHires, false)
}

func TestSnapshotSuite(t *testing.T) {
	suite.Run(t, new(snapshotSuite))
}

func (s *snapshotSuite) TestGetMainReturnsTextBank() {
	s.main.DirectSet(0x400, 0xAB)
	s.aux.DirectSet(0x400, 0xCD)
	s.snap.CopyFromState(s.main, s.aux, s.state)

	s.Equal(uint8(0xAB), s.snap.GetMain(0x400))
}

func (s *snapshotSuite) TestGetAuxReturnsTextBank() {
	s.main.DirectSet(0x400, 0xAB)
	s.aux.DirectSet(0x400, 0xCD)
	s.snap.CopyFromState(s.main, s.aux, s.state)

	s.Equal(uint8(0xCD), s.snap.GetAux(0x400))
}

func (s *snapshotSuite) TestGetMainTextPageBoundaries() {
	s.main.DirectSet(0x400, 0x11)
	s.main.DirectSet(0x7FF, 0x22)
	s.snap.CopyFromState(s.main, s.aux, s.state)

	s.Equal(uint8(0x11), s.snap.GetMain(0x400))
	s.Equal(uint8(0x22), s.snap.GetMain(0x7FF))
	s.Equal(uint8(0), s.snap.GetMain(0x3FF)) // below range
	s.Equal(uint8(0), s.snap.GetMain(0x800)) // above range
}

func (s *snapshotSuite) TestGetAuxTextPageBoundaries() {
	s.aux.DirectSet(0x400, 0x33)
	s.aux.DirectSet(0x7FF, 0x44)
	s.snap.CopyFromState(s.main, s.aux, s.state)

	s.Equal(uint8(0x33), s.snap.GetAux(0x400))
	s.Equal(uint8(0x44), s.snap.GetAux(0x7FF))
	s.Equal(uint8(0), s.snap.GetAux(0x3FF)) // below range
	s.Equal(uint8(0), s.snap.GetAux(0x800)) // above range
}

func (s *snapshotSuite) TestGetMainAndAuxRetainHiresRange() {
	// GetMain and GetAux still work for hires addresses after text additions.
	s.main.DirectSet(0x2000, 0x55)
	s.aux.DirectSet(0x2000, 0x66)
	s.snap.CopyFromState(s.main, s.aux, s.state)

	s.Equal(uint8(0x55), s.snap.GetMain(0x2000))
	s.Equal(uint8(0x66), s.snap.GetAux(0x2000))
}

func (s *snapshotSuite) TestCopyFromStateBothBanks_Default() {
	// 80STORE off, PAGE2 off: captures page 1 ($0400) from both banks.
	s.main.DirectSet(0x400, 0xAB)
	s.aux.DirectSet(0x400, 0xCD)

	s.snap.CopyFromState(s.main, s.aux, s.state)

	s.Equal(uint8(0xAB), s.snap.GetMain(0x400))
	s.Equal(uint8(0xCD), s.snap.GetAux(0x400))
}

func (s *snapshotSuite) TestCopyFromStateBothBanks_Page2WithoutStore80() {
	// PAGE2 on, 80STORE off: captures page 2 ($0800) from both banks.
	s.state.SetBool(a2state.DisplayPage2, true)
	s.state.SetBool(a2state.DisplayStore80, false)

	s.main.DirectSet(0x800, 0xAB)
	s.aux.DirectSet(0x800, 0xCD)

	s.snap.CopyFromState(s.main, s.aux, s.state)

	// Data from $0800 is stored starting at textMain[0], read back via $0400.
	s.Equal(uint8(0xAB), s.snap.GetMain(0x400))
	s.Equal(uint8(0xCD), s.snap.GetAux(0x400))
}

func (s *snapshotSuite) TestCopyFromStateBothBanks_Store80WithPage2() {
	// 80STORE on, PAGE2 on: always captures page 1 ($0400), not page 2.
	s.state.SetBool(a2state.DisplayStore80, true)
	s.state.SetBool(a2state.DisplayPage2, true)

	s.main.DirectSet(0x400, 0xAB)
	s.aux.DirectSet(0x400, 0xCD)
	s.main.DirectSet(0x800, 0xFF) // should not appear in snapshot
	s.aux.DirectSet(0x800, 0xFF)

	s.snap.CopyFromState(s.main, s.aux, s.state)

	s.Equal(uint8(0xAB), s.snap.GetMain(0x400))
	s.Equal(uint8(0xCD), s.snap.GetAux(0x400))
}

func (s *snapshotSuite) TestCopyFromStateMainAndAuxAreSeparate() {
	// GetMain and GetAux return data from their respective banks, not each
	// other's.
	s.main.DirectSet(0x401, 0x11)
	s.aux.DirectSet(0x401, 0x22)

	s.snap.CopyFromState(s.main, s.aux, s.state)

	s.NotEqual(s.snap.GetMain(0x401), s.snap.GetAux(0x401))
}
