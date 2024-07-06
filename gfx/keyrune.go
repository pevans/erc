package gfx

import "github.com/hajimehoshi/ebiten/v2"

var keyRunes = map[ebiten.Key]rune{
	ebiten.KeyA:            'A',
	ebiten.KeyB:            'B',
	ebiten.KeyBackquote:    '`',
	ebiten.KeyBackslash:    '\\',
	ebiten.KeyBackspace:    '\b',
	ebiten.KeyBracketLeft:  '[',
	ebiten.KeyBracketRight: ']',
	ebiten.KeyC:            'C',
	ebiten.KeyComma:        ',',
	ebiten.KeyD:            'D',
	ebiten.KeyDelete:       rune(127),
	ebiten.KeyDigit0:       '0',
	ebiten.KeyDigit1:       '1',
	ebiten.KeyDigit2:       '2',
	ebiten.KeyDigit3:       '3',
	ebiten.KeyDigit4:       '4',
	ebiten.KeyDigit5:       '5',
	ebiten.KeyDigit6:       '6',
	ebiten.KeyDigit7:       '7',
	ebiten.KeyDigit8:       '8',
	ebiten.KeyDigit9:       '9',
	ebiten.KeyE:            'E',
	ebiten.KeyEnter:        '\r',
	ebiten.KeyEqual:        '=',
	ebiten.KeyEscape:       rune(0x1b),
	ebiten.KeyF:            'F',
	ebiten.KeyG:            'G',
	ebiten.KeyH:            'H',
	ebiten.KeyI:            'I',
	ebiten.KeyJ:            'J',
	ebiten.KeyK:            'K',
	ebiten.KeyL:            'L',
	ebiten.KeyM:            'M',
	ebiten.KeyMinus:        '-',
	ebiten.KeyN:            'N',
	ebiten.KeyO:            'O',
	ebiten.KeyP:            'P',
	ebiten.KeyPeriod:       '.',
	ebiten.KeyQ:            'Q',
	ebiten.KeyQuote:        '\'',
	ebiten.KeyR:            'R',
	ebiten.KeyS:            'S',
	ebiten.KeySemicolon:    ';',
	ebiten.KeySlash:        '/',
	ebiten.KeySpace:        ' ',
	ebiten.KeyT:            'T',
	ebiten.KeyTab:          '\t',
	ebiten.KeyU:            'U',
	ebiten.KeyV:            'V',
	ebiten.KeyW:            'W',
	ebiten.KeyX:            'X',
	ebiten.KeyY:            'Y',
	ebiten.KeyZ:            'Z',
}

func KeyToRune(key ebiten.Key) rune {
	if r, ok := keyRunes[key]; ok {
		return r
	}

	return rune(0)
}
