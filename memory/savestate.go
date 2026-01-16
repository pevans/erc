package memory

import "fmt"

// Bytes returns a copy of the segment's memory contents.
func (s *Segment) Bytes() []uint8 {
	result := make([]uint8, len(s.mem))
	copy(result, s.mem)
	return result
}

// RestoreBytes restores the segment's memory from the given data. Returns an
// error if the data length doesn't match the segment size.
func (s *Segment) RestoreBytes(data []uint8) error {
	if len(data) != len(s.mem) {
		return fmt.Errorf(
			"data length %d doesn't match segment size %d",
			len(data), len(s.mem),
		)
	}

	copy(s.mem, data)
	return nil
}
