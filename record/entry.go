package record

import "fmt"

// Tag constants for the three categories of observable state.
const (
	TagMem  = "mem"
	TagReg  = "reg"
	TagComp = "comp"
)

// An Entry is a record of a single value changing from one state to another
// at a given execution step.
type Entry struct {
	Step int
	Tag  string
	Name string
	Old  any
	New  any
}

// String renders the entry in the canonical format:
//
// step N: tag name old -> new
func (e Entry) String() string {
	return fmt.Sprintf("step %d: %s %s %s -> %s",
		e.Step, e.Tag, e.Name, formatValue(e.Old), formatValue(e.New))
}

func formatValue(v any) string {
	switch val := v.(type) {
	case uint8:
		return fmt.Sprintf("$%02X", val)
	case uint16:
		return fmt.Sprintf("$%04X", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
