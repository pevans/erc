package asmrec

import "io"

// A Recorder is an interface which allows you to record assembly
// instructions that have been executed.
type Recorder interface {
	Record(w io.Writer) error
}
