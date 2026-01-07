package a2drive

import (
	"github.com/pevans/erc/a2/a2save"
	"github.com/pevans/erc/memory"
)

// Snapshot returns a snapshot of the drive state for serialization.
func (d *Drive) Snapshot() *a2save.DriveState {
	state := &a2save.DriveState{
		MotorOn:      d.motorOn,
		Phase:        d.phase,
		TrackPos:     d.trackPos,
		SectorPos:    d.sectorPos,
		Latch:        d.latch,
		Mode:         d.mode,
		LatchWasRead: d.latchWasRead,
		DiskShifted:  d.diskShifted,
		ImageName:    d.imageName,
		ImageType:    d.imageType,
		WriteProtect: d.writeProtect,
		HasDisk:      d.image != nil,
	}

	if d.image != nil {
		state.ImageData = d.image.Bytes()
	}

	if d.data != nil {
		state.PhysicalData = d.data.Bytes()
	}

	return state
}

// Restore restores the drive state from a snapshot.
func (d *Drive) Restore(state *a2save.DriveState) error {
	d.motorOn = state.MotorOn
	d.phase = state.Phase
	d.trackPos = state.TrackPos
	d.sectorPos = state.SectorPos
	d.latch = state.Latch
	d.mode = state.Mode
	d.latchWasRead = state.LatchWasRead
	d.diskShifted = state.DiskShifted
	d.imageName = state.ImageName
	d.imageType = state.ImageType
	d.writeProtect = state.WriteProtect

	if !state.HasDisk {
		d.image = nil
		d.data = nil
		return nil
	}

	if len(state.ImageData) > 0 {
		d.image = memory.NewSegment(len(state.ImageData))
		if err := d.image.RestoreBytes(state.ImageData); err != nil {
			return err
		}
	}

	if len(state.PhysicalData) > 0 {
		d.data = memory.NewSegment(len(state.PhysicalData))
		if err := d.data.RestoreBytes(state.PhysicalData); err != nil {
			return err
		}
	}

	return nil
}
