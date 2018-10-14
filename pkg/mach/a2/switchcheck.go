package a2

import "github.com/pevans/erc/pkg/mach"

type SwitchCheck struct {
	mode    func(c *Computer) int
	setMode func(c *Computer, mode int)
}

func (s *SwitchCheck) IsSetter(flag int) ReadMapFn {
	return func(c *Computer, addr mach.Addressor) mach.Byte {
		if s.mode(c)&flag > 0 {
			return mach.Byte(0x80)
		}

		return mach.Byte(0x0)
	}
}

func (s *SwitchCheck) SetterR(flag int) ReadMapFn {
	return func(c *Computer, addr mach.Addressor) mach.Byte {
		s.setMode(c, flag)
		return mach.Byte(0x80)
	}
}

func (s *SwitchCheck) ReSetterR(flag int) ReadMapFn {
	return func(c *Computer, addr mach.Addressor) mach.Byte {
		s.setMode(c, s.mode(c)|flag)
		return mach.Byte(0x80)
	}
}

func (s *SwitchCheck) UnSetterR(flag int) ReadMapFn {
	return func(c *Computer, addr mach.Addressor) mach.Byte {
		s.setMode(c, s.mode(c) & ^flag)
		return mach.Byte(0x0)
	}
}

func (s *SwitchCheck) SetterW(flag int) WriteMapFn {
	return func(c *Computer, addr mach.Addressor, val mach.Byte) {
		s.setMode(c, flag)
	}
}

func (s *SwitchCheck) ReSetterW(flag int) WriteMapFn {
	return func(c *Computer, addr mach.Addressor, val mach.Byte) {
		s.setMode(c, s.mode(c)|flag)
	}
}

func (s *SwitchCheck) UnSetterW(flag int) WriteMapFn {
	return func(c *Computer, addr mach.Addressor, val mach.Byte) {
		s.setMode(c, s.mode(c) & ^flag)
	}
}
