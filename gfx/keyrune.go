package gfx

import "github.com/hajimehoshi/ebiten/v2"

var keyRunes = map[ebiten.Key]rune{
	ebiten.KeyA:            'a',
	ebiten.KeyB:            'b',
	ebiten.KeyBackquote:    '`',
	ebiten.KeyBackslash:    '\\',
	ebiten.KeyBackspace:    '\b',
	ebiten.KeyBracketLeft:  '[',
	ebiten.KeyBracketRight: ']',
	ebiten.KeyC:            'c',
	ebiten.KeyComma:        ',',
	ebiten.KeyD:            'd',
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
	ebiten.KeyE:            'e',
	ebiten.KeyEnter:        '\r',
	ebiten.KeyEqual:        '=',
	ebiten.KeyEscape:       rune(0x1b),
	ebiten.KeyF:            'f',
	ebiten.KeyG:            'g',
	ebiten.KeyH:            'h',
	ebiten.KeyI:            'i',
	ebiten.KeyJ:            'j',
	ebiten.KeyK:            'k',
	ebiten.KeyL:            'l',
	ebiten.KeyM:            'm',
	ebiten.KeyMinus:        '-',
	ebiten.KeyN:            'n',
	ebiten.KeyO:            'o',
	ebiten.KeyP:            'p',
	ebiten.KeyPeriod:       '.',
	ebiten.KeyQ:            'q',
	ebiten.KeyQuote:        '\'',
	ebiten.KeyR:            'r',
	ebiten.KeyS:            's',
	ebiten.KeySemicolon:    ';',
	ebiten.KeySlash:        '/',
	ebiten.KeySpace:        ' ',
	ebiten.KeyT:            't',
	ebiten.KeyTab:          '\t',
	ebiten.KeyU:            'u',
	ebiten.KeyV:            'v',
	ebiten.KeyW:            'w',
	ebiten.KeyX:            'x',
	ebiten.KeyY:            'y',
	ebiten.KeyZ:            'z',
}

func KeyToRune(key ebiten.Key) rune {
	if r, ok := keyRunes[key]; ok {
		return r
	}

	return rune(0)
}
