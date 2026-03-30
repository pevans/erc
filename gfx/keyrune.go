package gfx

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pevans/erc/input"
)

// keyRunes is a mapping of ebiten Keys to runes. The Apple II is 7-bit ASCII,
// so this works out straightforwardly.
var keyRunes = map[ebiten.Key]rune{
	ebiten.KeyA:            'a',
	ebiten.KeyArrowDown:    rune(0x0A),
	ebiten.KeyArrowLeft:    rune(0x08), // shares a code with Backspace
	ebiten.KeyArrowRight:   rune(0x15),
	ebiten.KeyArrowUp:      rune(0x0B),
	ebiten.KeyB:            'b',
	ebiten.KeyBackquote:    '`',
	ebiten.KeyBackslash:    '\\',
	ebiten.KeyBackspace:    rune(0x08),
	ebiten.KeyBracketLeft:  '[',
	ebiten.KeyBracketRight: ']',
	ebiten.KeyC:            'c',
	ebiten.KeyComma:        ',',
	ebiten.KeyD:            'd',
	ebiten.KeyDelete:       rune(0x7F),
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
	ebiten.KeyEnter:        rune(0x0D),
	ebiten.KeyEqual:        '=',
	ebiten.KeyEscape:       rune(0x1B),
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
	ebiten.KeySpace:        rune(0x20),
	ebiten.KeyT:            't',
	ebiten.KeyTab:          rune(0x09),
	ebiten.KeyU:            'u',
	ebiten.KeyV:            'v',
	ebiten.KeyW:            'w',
	ebiten.KeyX:            'x',
	ebiten.KeyY:            'y',
	ebiten.KeyZ:            'z',
}

// shiftKeyRunes are a map of ebiten Keys to runes, but with the assumption
// that the shift key had been held down. (Shift keys are a modifier rather
// than baked into the ebiten Key itself.)
var shiftKeyRunes = map[ebiten.Key]rune{
	ebiten.KeyA:            'A',
	ebiten.KeyB:            'B',
	ebiten.KeyC:            'C',
	ebiten.KeyD:            'D',
	ebiten.KeyE:            'E',
	ebiten.KeyF:            'F',
	ebiten.KeyG:            'G',
	ebiten.KeyH:            'H',
	ebiten.KeyI:            'I',
	ebiten.KeyJ:            'J',
	ebiten.KeyK:            'K',
	ebiten.KeyL:            'L',
	ebiten.KeyM:            'M',
	ebiten.KeyN:            'N',
	ebiten.KeyO:            'O',
	ebiten.KeyP:            'P',
	ebiten.KeyQ:            'Q',
	ebiten.KeyR:            'R',
	ebiten.KeyS:            'S',
	ebiten.KeyT:            'T',
	ebiten.KeyU:            'U',
	ebiten.KeyV:            'V',
	ebiten.KeyW:            'W',
	ebiten.KeyX:            'X',
	ebiten.KeyY:            'Y',
	ebiten.KeyZ:            'Z',
	ebiten.KeyBackquote:    '~',
	ebiten.KeyBackslash:    '|',
	ebiten.KeyBracketLeft:  '{',
	ebiten.KeyBracketRight: '}',
	ebiten.KeyComma:        '<',
	ebiten.KeyDigit0:       ')',
	ebiten.KeyDigit1:       '!',
	ebiten.KeyDigit2:       '@',
	ebiten.KeyDigit3:       '#',
	ebiten.KeyDigit4:       '$',
	ebiten.KeyDigit5:       '%',
	ebiten.KeyDigit6:       '^',
	ebiten.KeyDigit7:       '&',
	ebiten.KeyDigit8:       '*',
	ebiten.KeyDigit9:       '(',
	ebiten.KeyEqual:        '+',
	ebiten.KeyMinus:        '_',
	ebiten.KeyPeriod:       '>',
	ebiten.KeyQuote:        '"',
	ebiten.KeySemicolon:    ':',
	ebiten.KeySlash:        '?',
}

// KeyToRune returns a rune aligned to a specific ebiten Key and modifier.
func KeyToRune(key ebiten.Key, modifier int) rune {
	if modifier == input.ModShift {
		if r, ok := shiftKeyRunes[key]; ok {
			return r
		}
	}

	r, ok := keyRunes[key]
	if !ok {
		return rune(0)
	}

	return r
}
