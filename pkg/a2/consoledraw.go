package a2

import (
	"github.com/gdamore/tcell"
	"github.com/pevans/erc/pkg/data"
	"github.com/pevans/erc/pkg/emu"
)

type ConsoleScreen struct {
	screen tcell.Screen
}

func (s *ConsoleScreen) Draw(bytes data.Getter, ctx emu.DrawContext) error {
	return nil
}

func (s *ConsoleScreen) Init() error {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	scr, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	if err := scr.Init(); err != nil {
		return err
	}

	style := tcell.StyleDefault.
		Foreground(tcell.ColorGray).
		Background(tcell.ColorBlack)

	scr.SetStyle(style)
	scr.Clear()

	s.screen = scr

	return nil
}
