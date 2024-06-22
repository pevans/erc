package a2sym

import "fmt"

type SwitchMode int

const (
	ModeNone SwitchMode = iota
	ModeR
	ModeR7
	ModeRR
	ModeRW
	ModeW
)

type Switch struct {
	Mode        SwitchMode
	Name        string
	Description string
}

func (mode SwitchMode) String() string {
	switch mode {
	case ModeR:
		return "read"
	case ModeR7:
		return "read, result in high bit"
	case ModeRR:
		return "read, twice consecutively"
	case ModeRW:
		return "read or write"
	case ModeW:
		return "write"
	}

	return "unknown switch mode"
}

func (s Switch) String() string {
	// This is probably a zero-value switch
	if s.Mode == ModeNone {
		return ""
	}

	if len(s.Name) == 0 {
		return fmt.Sprintf("%v (%v)", s.Mode, s.Description)
	}

	if len(s.Description) == 0 {
		return fmt.Sprintf("%v %v", s.Mode, s.Name)
	}

	return fmt.Sprintf("%v %v (%v)", s.Mode, s.Name, s.Description)
}
