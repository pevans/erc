package a2

import "github.com/stretchr/testify/assert"

func (s *a2Suite) TestDefineSoftSwitches() {
	var (
		addr int
		ok   bool
	)

	// Testing for BankAuxiliary switches
	for addr = 0; addr < 0x200; addr++ {
		_, ok = s.comp.RMap[addr]
		assert.Equal(s.T(), true, ok)

		_, ok = s.comp.WMap[addr]
		assert.Equal(s.T(), true, ok)
	}

	// Testing all cases where ROM or bank-addressable RAM could be
	// found
	for addr = 0xD000; addr < 0x10000; addr++ {
		_, ok = s.comp.RMap[addr]
		assert.Equal(s.T(), true, ok)

		_, ok = s.comp.WMap[addr]
		assert.Equal(s.T(), true, ok)
	}

	rmapModifiers := []int{
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
	}

	wmapModifiers := []int{
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
	}

	// Here we simply test all of the possible RMap and WMap switches
	// which modify something
	for _, addr = range rmapModifiers {
		_, ok = s.comp.RMap[addr]
		assert.Equal(s.T(), true, ok)
	}

	for _, addr = range wmapModifiers {
		_, ok = s.comp.WMap[addr]
		assert.Equal(s.T(), true, ok)
	}
}
