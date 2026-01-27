package debug

import (
	"fmt"
	"strconv"
)

func hex(token string, bits int) (int, error) {
	ui64, err := strconv.ParseUint(token, 16, bits)
	if err != nil {
		return 0, fmt.Errorf("invalid hex: \"%v\": %w", token, err)
	}

	return int(ui64), nil
}

func integer(token string) (int, error) {
	ui64, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: \"%v\": %w", token, err)
	}

	return int(ui64), nil
}
