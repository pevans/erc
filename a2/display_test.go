package a2

import (
	"github.com/pevans/erc/a2/a2display"
	"github.com/pevans/erc/a2/a2state"
)

func (s *a2Suite) TestIsVerticalBlank() {
	s.Run("returns consistent result based on cycle position", func() {
		_ = s.comp.Boot()

		cycles := s.comp.CPU.CycleCounter() % a2display.ScanCycleCount
		result := s.comp.IsVerticalBlank()

		if cycles < 12480 {
			s.False(result, "cycles %d < 12480, should not be in vblank", cycles)
		} else {
			s.True(result, "cycles %d >= 12480, should be in vblank", cycles)
		}
	})

	s.Run("changes state after executing instructions", func() {
		_ = s.comp.Boot()

		seenFalse := false
		seenTrue := false

		for range 9000 {
			s.comp.Main.Set(int(s.comp.CPU.PC), 0xEA)
			_ = s.comp.CPU.Execute()

			cycles := s.comp.CPU.CycleCounter() % a2display.ScanCycleCount
			result := s.comp.IsVerticalBlank()

			expectedInVBlank := cycles >= 12480
			s.Equal(expectedInVBlank, result,
				"at %d cycles (mod %d), expected vblank=%v but got %v",
				s.comp.CPU.CycleCounter(), cycles, expectedInVBlank, result)

			if result {
				seenTrue = true
			} else {
				seenFalse = true
			}

			if seenTrue && seenFalse {
				return
			}
		}

		s.True(seenTrue && seenFalse, "should have seen both true and false states")
	})
}

func (s *a2Suite) TestRenderFlash() {
	s.Run("updates lastFlashOn and triggers redraw when flash state differs", func() {
		// At cycle 0, flashOn=true. Pre-seed lastFlashOn=false so they
		// differ.
		s.comp.lastFlashOn = false
		s.comp.State.SetBool(a2state.DisplayRedraw, false)

		s.comp.Render()

		s.True(s.comp.lastFlashOn)
	})

	s.Run("does not trigger redraw when flash state is unchanged", func() {
		// Seed lastFlashOn to match what Render will compute (true at cycle
		// 0).
		s.comp.lastFlashOn = true
		s.comp.State.SetBool(a2state.DisplayRedraw, false)

		s.comp.Render()

		s.False(s.comp.State.Bool(a2state.DisplayRedraw))
	})
}
