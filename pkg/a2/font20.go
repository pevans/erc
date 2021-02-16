package a2

import "github.com/pevans/erc/pkg/gfx"

func font20(f *gfx.Font) {
	f.DefineGlyph(0x20, []byte{ // SP
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x21, []byte{ // !
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
	})

	f.DefineGlyph(0x22, []byte{ // "
		0, 0, 1, 0, 1, 0, 0,
		0, 0, 1, 0, 1, 0, 0,
		0, 0, 1, 0, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x23, []byte{ // #
		0, 0, 1, 0, 1, 0, 0,
		0, 0, 1, 0, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 0, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 1, 0, 1, 0, 0,
		0, 0, 1, 0, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x24, []byte{ // $
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 1, 0, 1, 0, 0, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 1, 0, 1, 0,
		0, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x25, []byte{ // %
		0, 1, 1, 0, 0, 0, 0,
		0, 1, 1, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 1, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x26, []byte{ // &
		0, 0, 1, 0, 0, 0, 0,
		0, 1, 0, 1, 0, 0, 0,
		0, 1, 0, 1, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 1, 0, 1, 0, 1, 0,
		0, 1, 0, 0, 1, 0, 0,
		0, 0, 1, 1, 0, 1, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x27, []byte{ // '
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x28, []byte{ // (
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x29, []byte{ // )
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x2a, []byte{ // *
		0, 0, 0, 1, 0, 0, 0,
		0, 1, 0, 1, 0, 1, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 1, 0, 1, 0, 1, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x2b, []byte{ // +
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x2c, []byte{ // ,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x2d, []byte{ // -
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x2e, []byte{ // .
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x2f, []byte{ // /
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x30, []byte{ // 0
		0, 0, 1, 1, 1, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 1, 0, 0, 1, 1, 0,
		0, 1, 0, 1, 0, 1, 0,
		0, 1, 1, 0, 0, 1, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x31, []byte{ // 1
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x32, []byte{ // 2
		0, 0, 1, 1, 1, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 1, 1, 0, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x33, []byte{ // 3
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x34, []byte{ // 4
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 1, 1, 0, 0,
		0, 0, 1, 0, 1, 0, 0,
		0, 1, 0, 0, 1, 0, 0,
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x35, []byte{ // 5
		0, 1, 1, 1, 1, 1, 0,
		0, 1, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x36, []byte{ // 6
		0, 0, 0, 1, 1, 1, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x37, []byte{ // 7
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x38, []byte{ // 8
		0, 0, 1, 1, 1, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 0, 1, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x39, []byte{ // 9
		0, 0, 1, 1, 1, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 0, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 1, 1, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x3a, []byte{ // :
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x3b, []byte{ // ;
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x3c, []byte{ // <
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x3d, []byte{ // =
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x3e, []byte{ // >
		0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})

	f.DefineGlyph(0x3f, []byte{ // ?
		0, 0, 1, 1, 1, 0, 0,
		0, 1, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	})
}