package disasm

import (
	"os"

	"github.com/pevans/erc/pkg/asmrec"
)

type SourceMap struct {
	smap map[int]asmrec.Recorder
}

func NewSourceMap() *SourceMap {
	sm := new(SourceMap)
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
	target := os.Getenv("DISASM_LOG")
	if target == "" {
		return nil
	}

	w, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0755)
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
