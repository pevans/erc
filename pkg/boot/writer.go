package boot

import "os"

// OpenFile will attempt to open the given filename for writing log
// data out. If the filename is empty, it will not return an error, but
// will return a nil file value.
func OpenFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return nil, nil
	}

	return os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
}
