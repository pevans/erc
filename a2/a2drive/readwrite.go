package a2drive

// LoadLatch reads the byte at the current drive position with respect to
// track and position. No change to the latch will occur if we cannot detect
// that the disk position has shifted since the last LoadLatch was called.
func (d *Drive) LoadLatch() {
	if d.data == nil {
		return
	}

	if d.diskShifted {
		d.latch = d.data.DirectGet(d.dataPosition())
		d.diskShifted = false
		d.latchWasRead = false
	}
}

// PeekLatch returns the value of the drive latch. Unlike ReadLatch, it does
// not record a latch read; it simply returns whatever is in the latch.
func (d *Drive) PeekLatch() uint8 {
	return d.latch
}

// ReadLatch returns the byte that is currently loaded in the drive latch.
// If this data has not been read before, it is returned unmodified. If it has
// been read before, then it will be returned with the high bit set to zero.
func (d *Drive) ReadLatch() uint8 {
	if d.data == nil {
		return 0xFF
	}

	if d.latchWasRead {
		return d.latch & 0x7F
	}

	d.latchWasRead = true

	return d.latch
}

// ReadMode returns true if the drive is in read mode (is able to read data
// from the disk)
func (d *Drive) ReadMode() bool {
	return d.mode == readMode
}

// SetLatch sets the value of the drive latch to val.
func (d *Drive) SetLatch(val uint8) {
	d.latch = val
	d.latchWasRead = false
}

// SetReadMode sets the drive to read mode.
func (d *Drive) SetReadMode() {
	d.mode = readMode
}

// SetWriteMode sets the drive to write mode.
func (d *Drive) SetWriteMode() {
	d.mode = writeMode
}

// SetWriteProtect will change the writeProtect status of a drive to the given
// status.
func (d *Drive) SetWriteProtect(status bool) {
	d.writeProtect = status
}

// ToggleWriteProtect flips the status of write protection for the disk in a
// drive. If it was true, it becomes false, and vice-versa.
func (d *Drive) ToggleWriteProtect() {
	d.writeProtect = !d.writeProtect
}

// WriteLatch writes the byte in the drive's latch to the disk loaded in the
// drive. The drive must be in WriteMode for this operation to succeed, and
// the motor must be on. WriteLatch will not write any data if the latch byte
// does not have its high bit set to 1. The byte in the latch will be written
// to the current position of the drive head on the disk (with respect to
// track and sector).
func (d *Drive) WriteLatch() {
	if d.data == nil {
		return
	}

	if d.WriteMode() && d.MotorOn() && d.latch&0x80 > 0 {
		d.data.DirectSet(d.dataPosition(), d.latch)
	}
}

// WriteMode returns true if the drive is in write mode (is able to write data
// to the disk). This does not take write protection into account; it's
// possible for a drive to be in write mode but still be unable to write to a
// disk that is write-protected.
func (d *Drive) WriteMode() bool {
	return d.mode == writeMode
}

// WriteProtected returns true if the disk in the drive is write-protected
// (can't be written to).
func (d *Drive) WriteProtected() bool {
	return d.writeProtect
}
