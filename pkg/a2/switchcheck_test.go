package a2

import (
	"github.com/pevans/erc/pkg/data"
)

const (
	mockNoMode   = 0
	mockMode     = 1
	mockModeB    = 2
	mockStubAddr = data.Int(0)
	mockStubVal  = data.Byte(0)
)

func mockSwitchCheck(modeState *int) *SwitchCheck {
	return &SwitchCheck{
		mode:    func(c *Computer) int { return mockMode },
		setMode: func(c *Computer, mode int) { *modeState = mode },
	}
}

func (s *a2Suite) TestIsSetter() {
	msc := mockSwitchCheck(nil)

	rfn := msc.IsSetter(mockMode)
	s.Equal(data.Byte(0x80), rfn(s.comp, mockStubAddr))

	rfn = msc.IsSetter(mockNoMode)
	s.Equal(data.Byte(0x0), rfn(s.comp, mockStubAddr))
}

func (s *a2Suite) TestIsOpSetter() {
	msc := mockSwitchCheck(nil)

	rfn := msc.IsOpSetter(mockMode)
	s.Equal(data.Byte(0x0), rfn(s.comp, mockStubAddr))

	rfn = msc.IsOpSetter(mockNoMode)
	s.Equal(data.Byte(0x80), rfn(s.comp, mockStubAddr))
}

func (s *a2Suite) TestSetterR() {
	mst := mockNoMode
	msc := mockSwitchCheck(&mst)

	rfn := msc.SetterR(mockMode)
	s.Equal(data.Byte(0x80), rfn(s.comp, mockStubAddr))
	s.Equal(mockMode, mst)
}

func (s *a2Suite) TestReSetterR() {
	mst := mockNoMode
	msc := mockSwitchCheck(&mst)

	rfn := msc.ReSetterR(mockModeB)
	s.Equal(data.Byte(0x80), rfn(s.comp, mockStubAddr))
	s.Equal(mockMode|mockModeB, mst)
}

func (s *a2Suite) TestUnSetterR() {
	mst := mockNoMode
	msc := mockSwitchCheck(&mst)

	rfn := msc.UnSetterR(mockModeB)
	s.Equal(data.Byte(0x0), rfn(s.comp, mockStubAddr))
	s.Equal(mockMode, mst)
}

func (s *a2Suite) TestSetterW() {
	mst := mockNoMode
	msc := mockSwitchCheck(&mst)

	wfn := msc.SetterW(mockMode)
	wfn(s.comp, mockStubAddr, mockStubVal)
	s.Equal(mockMode, mst)
}

func (s *a2Suite) TestReSetterW() {
	mst := mockNoMode
	msc := mockSwitchCheck(&mst)

	wfn := msc.ReSetterW(mockModeB)
	wfn(s.comp, mockStubAddr, mockStubVal)
	s.Equal(mockMode|mockModeB, mst)
}

func (s *a2Suite) TestUnSetterW() {
	mst := mockNoMode
	msc := mockSwitchCheck(&mst)

	wfn := msc.UnSetterW(mockModeB)
	wfn(s.comp, mockStubAddr, mockStubVal)
	s.Equal(mockMode, mst)
}
