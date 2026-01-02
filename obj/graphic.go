package obj

import _ "embed"

//go:embed png/debug.png
var debugPNG []byte

func DebugPNG() []byte {
	return debugPNG
}

//go:embed png/pause.png
var pausePNG []byte

func PausePNG() []byte {
	return pausePNG
}

//go:embed png/resume.png
var resumePNG []byte

func ResumePNG() []byte {
	return resumePNG
}

//go:embed png/writeable.png
var writeablePNG []byte

func WriteablePNG() []byte {
	return writeablePNG
}

//go:embed png/write-protected.png
var writeProtectedPNG []byte

func WriteProtectedPNG() []byte {
	return writeProtectedPNG
}
