package debug

import (
	"strconv"

	"github.com/pkg/errors"
)

func hex(token string, bits int) (int, error) {
	ui64, err := strconv.ParseUint(token, 16, bits)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid hex: \"%v\"", token)
	}

	return int(ui64), nil
}

func integer(token string) (int, error) {
	ui64, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid integer: \"%v\"", token)
	}

	return int(ui64), nil
}
