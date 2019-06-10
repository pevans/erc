package asmrec

import "io"

type Recorder interface {
	Record(w io.Writer) error
}
