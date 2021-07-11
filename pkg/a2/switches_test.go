package a2

import "github.com/pevans/erc/pkg/data"

func (s *a2Suite) TestMapSoftSwitches() {
	var (
		addr int
		ok   bool
	)

	// Testing for BankAuxiliary switches
	for addr = 0; addr < 0x200; addr++ {
		_, ok = s.comp.RMap[addr]
		s.Equal(true, ok)

		_, ok = s.comp.WMap[addr]
		s.Equal(true, ok)
	}

	for addr = 0x400; addr < 0x800; addr++ {
		_, ok = s.comp.RMap[addr]
		s.Equal(true, ok)

		_, ok = s.comp.WMap[addr]
		s.Equal(true, ok)
	}

	for addr = 0x2000; addr < 0x4000; addr++ {
		_, ok = s.comp.RMap[addr]
		s.Equal(true, ok)

		_, ok = s.comp.WMap[addr]
		s.Equal(true, ok)
	}

	// Testing all cases where ROM or bank-addressable RAM could be
	// found
	for addr = 0xD000; addr < 0x10000; addr++ {
		_, ok = s.comp.RMap[addr]
		s.Equal(true, ok)

		_, ok = s.comp.WMap[addr]
		s.Equal(true, ok)
	}

	rmapModifiers := []data.DByte{
		kbDataAndStrobe,
		kbAnyKeyDown,
		0xC013,
		0xC014,
		0xC018,
		0xC01C,
		0xC01D,
		0xC054,
		0xC055,
		0xC056,
		0xC057,
		0xC080,
		0xC081,
		0xC082,
		0xC083,
		0xC088,
		0xC089,
		0xC08A,
		0xC08B,
		0xC011,
		0xC012,
		0xC016,
		0xC01A,
		0xC01B,
		0xC01E,
		0xC01F,
		0xC050,
		0xC051,
		0xC052,
		0xC053,
		0xC05E,
		0xC05F,
		0xC07E,
		0xC07F,
	}

	wmapModifiers := []data.DByte{
		0xC000,
		0xC001,
		0xC002,
		0xC003,
		0xC004,
		0xC005,
		0xC054,
		0xC055,
		0xC056,
		0xC057,
		0xC008,
		0xC009,
		0xC00C,
		0xC00D,
		0xC00E,
		0xC00F,
		0xC050,
		0xC051,
		0xC052,
		0xC053,
		0xC05E,
		0xC05F,
		0xC07E,
		0xC07F,
	}

	// Here we simply test all of the possible RMap and WMap switches
	// which modify something
	for _, addr := range rmapModifiers {
		_, ok = s.comp.RMap[addr.Int()]
		s.Truef(ok, "addr=%x", addr)
	}

	for _, addr := range wmapModifiers {
		_, ok = s.comp.WMap[addr.Int()]
		s.Equal(true, ok)
	}
}
