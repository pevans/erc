package a2

import "github.com/stretchr/testify/assert"

func (s *a2Suite) TestDefineSoftSwitches() {
	var addr int

	// Testing for BankAuxiliary switches
	for addr = 0; addr < 0x200; addr++ {
		assert.Contains(s.T(), s.comp.RMap, addr)
		assert.Contains(s.T(), s.comp.WMap, addr)
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
	}

	// Here we simply test all of the possible RMap and WMap switches
	// which modify something
	for _, addr = range rmapModifiers {
		assert.Contains(s.T(), s.comp.RMap, addr)
	}

	for _, addr = range wmapModifiers {
		assert.Contains(s.T(), s.comp.WMap, addr)
	}
}
