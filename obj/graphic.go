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

//go:embed png/disk1.png
var disk1PNG []byte

func Disk1PNG() []byte {
	return disk1PNG
}

//go:embed png/disk2.png
var disk2PNG []byte

func Disk2PNG() []byte {
	return disk2PNG
}

//go:embed png/disk3.png
var disk3PNG []byte

func Disk3PNG() []byte {
	return disk3PNG
}

//go:embed png/disk4.png
var disk4PNG []byte

func Disk4PNG() []byte {
	return disk4PNG
}

//go:embed png/disk5.png
var disk5PNG []byte

func Disk5PNG() []byte {
	return disk5PNG
}

//go:embed png/disk6.png
var disk6PNG []byte

func Disk6PNG() []byte {
	return disk6PNG
}

//go:embed png/disk7.png
var disk7PNG []byte

func Disk7PNG() []byte {
	return disk7PNG
}

//go:embed png/disk8.png
var disk8PNG []byte

func Disk8PNG() []byte {
	return disk8PNG
}

//go:embed png/disk9.png
var disk9PNG []byte

func Disk9PNG() []byte {
	return disk9PNG
}

//go:embed png/disk10.png
var disk10PNG []byte

func Disk10PNG() []byte {
	return disk10PNG
}
