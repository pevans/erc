package disasm

import (
	"os"

	"github.com/pevans/erc/pkg/asmrec"
)

type SourceMap struct {
	file string
	smap map[int]asmrec.Recorder
}

func NewSourceMap(file string) *SourceMap {
	sm := new(SourceMap)
	sm.file = file
	sm.smap = make(map[int]asmrec.Recorder)

	return sm
}

func (sm *SourceMap) Map(addr int, rec asmrec.Recorder) bool {
	if _, ok := sm.smap[addr]; ok {
		return false
	}

	sm.smap[addr] = rec
	return true
}

func (sm *SourceMap) WriteLog() error {
	if sm.file == "" {
		return nil
	}

	w, err := os.OpenFile(sm.file, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	for _, rec := range sm.smap {
		if err := rec.Record(w); err != nil {
			return err
		}
	}

	return nil
}
