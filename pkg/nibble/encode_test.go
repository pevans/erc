package nibble

import (
	"testing"
)

func TestEncode(t *testing.T) {
	// A funny thing about the encode procedure is it's identical to the
	// decode procedure, so the decode test should suffice.
	TestDecode(t)
}
