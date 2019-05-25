package a2

import "github.com/pevans/erc/pkg/data"

// SwitchCheck is a collection of functions which return and set a
// mode in a computer, and is used to produce setters, unsetters,
// issetters, etc. for various soft switches.
type SwitchCheck struct {
	mode    func(c *Computer) int
	setMode func(c *Computer, mode int)
}

// IsSetter returns a read map function that checks if a given flag is
// set in our mode.
func (s *SwitchCheck) IsSetter(flag int) ReadMapFn {
	return func(c *Computer, addr data.Addressor) data.Byte {
		if s.mode(c)&flag > 0 {
			return data.Byte(0x80)
		}

		return data.Byte(0x0)
	}
}

// SetterR returns a read map function which sets a flag directly to our
// mode and returns a positive byte signal.
func (s *SwitchCheck) SetterR(flag int) ReadMapFn {
	return func(c *Computer, addr data.Addressor) data.Byte {
		s.setMode(c, flag)
		return data.Byte(0x80)
	}
}

// ReSetterR returns a read map function which sets a flag alongside
// existing flags (using bitwise-OR logic) and returns a positive byte
// signal.
func (s *SwitchCheck) ReSetterR(flag int) ReadMapFn {
	return func(c *Computer, addr data.Addressor) data.Byte {
		s.setMode(c, s.mode(c)|flag)
		return data.Byte(0x80)
	}
}

// UnSetterR returns a read map function which unsets a flag and
// preserves existing flags and returns a negative byte signal.
func (s *SwitchCheck) UnSetterR(flag int) ReadMapFn {
	return func(c *Computer, addr data.Addressor) data.Byte {
		s.setMode(c, s.mode(c) & ^flag)
		return data.Byte(0x0)
	}
}

// SetterW returns a write map function which sets a flag directly.
func (s *SwitchCheck) SetterW(flag int) WriteMapFn {
	return func(c *Computer, addr data.Addressor, val data.Byte) {
		s.setMode(c, flag)
	}
}

// ReSetterW returns a write map function which sets a flag and
// preserves existing flags alongside it.
func (s *SwitchCheck) ReSetterW(flag int) WriteMapFn {
	return func(c *Computer, addr data.Addressor, val data.Byte) {
		s.setMode(c, s.mode(c)|flag)
	}
}

// UnSetterW returns a write map function which unsets a flag but
// preserves existing flags alongside it.
func (s *SwitchCheck) UnSetterW(flag int) WriteMapFn {
	return func(c *Computer, addr data.Addressor, val data.Byte) {
		s.setMode(c, s.mode(c) & ^flag)
	}
}
