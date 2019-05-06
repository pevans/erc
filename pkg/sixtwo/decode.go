package sixtwo

import "github.com/pevans/erc/pkg/data"

type decoder struct {
	logSeg    *data.Segment
	physSeg   *data.Segment
	imageType int
	logOff    int
	physOff   int
}

//  00    01    02    03    04    05    06    07    08    09    0A    0B    0C    0D    0E    0F
var conv6bit = []data.Byte{
	0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, // 00
	0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0X00, 0X04, 0XFF, 0XFF, 0X08, 0X0C, 0XFF, 0X10, 0X14, 0X18, // 10
	0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0X1C, 0X20, 0XFF, 0XFF, 0XFF, 0X24, 0X28, 0X2C, 0X30, 0X34, // 20
	0XFF, 0XFF, 0X38, 0X3C, 0X40, 0X44, 0X48, 0X4C, 0XFF, 0X50, 0X54, 0X58, 0X5C, 0X60, 0X64, 0X68, // 30
	0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0X6C, 0XFF, 0X70, 0X74, 0X78, // 40
	0XFF, 0XFF, 0XFF, 0X7C, 0XFF, 0XFF, 0X80, 0X84, 0XFF, 0X88, 0X8C, 0X90, 0X94, 0X98, 0X9C, 0XA0, // 50
	0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XA4, 0XA8, 0XAC, 0XFF, 0XB0, 0XB4, 0XB8, 0XBC, 0XC0, 0XC4, 0XC8, // 60
	0XFF, 0XFF, 0XCC, 0XD0, 0XD4, 0XD8, 0XDC, 0XE0, 0XFF, 0XE4, 0XE8, 0XEC, 0XF0, 0XF4, 0XF8, 0XFC, // 70
}

func Decode(imageType int, src *data.Segment) (*data.Segment, error) {
	dec := &decoder{
		physSeg:   src,
		logSeg:    data.NewSegment(DosSize),
		imageType: imageType,
	}

	for track := 0; track < NumTracks; track++ {
		dec.writeTrack(track)
	}

	return dec.logSeg, nil
}

func (d *decoder) writeTrack(track int) {
	d.physOff = (track * PhysTrackLen) + PhysTrackHeader

	for sect := 0; sect < NumSectors; sect++ {
		d.logOff = (track * LogTrackLen) +
			(logicalSector(d.imageType, sect) * LogSectorLen)

		d.writeSector(track, sect)
	}
}

func HeaderOK(seg *data.Segment, off int) bool {
	addr := data.Int(off)

	return seg.Get(addr) == data.Byte(0xD5) &&
		seg.Get(data.Plus(addr, 1)) == data.Byte(0xAA) &&
		seg.Get(data.Plus(addr, 2)) == data.Byte(0xAD)
}

func (d *decoder) writeSector(track, sect int) {
	// Skip header and the data marker
	d.physOff += PhysSectorHeader + 3

	var (
		conv = make([]data.Byte, 0x157)
		xor  = make([]data.Byte, 0x156)
	)

	for i := 0; i < 0x157; i++ {
		conv[i] = conv6bit[d.physSeg.Get(data.DByte(d.physOff+i))&0x7F]
	}

	for i, lval := 0, data.Byte(0); i < 0x156; i++ {
		xor[i] = lval ^ conv[i]
		lval = xor[i]
	}

	for i := data.Byte(0); i < 0x56; i++ {
		var (
			offac = i + 0xAC
			off56 = i + 0x56

			vac = (xor[int(offac)+0x56] & 0xFC) | ((xor[i] & 0x80) >> 7) | ((xor[i] & 0x40) >> 5)
			v56 = (xor[int(off56)+0x56] & 0xFC) | ((xor[i] & 0x20) >> 5) | ((xor[i] & 0x10) >> 3)
			v00 = (xor[i+0x56] & 0xFC) | ((xor[i] & 0x08) >> 3) | ((xor[i] & 0x04) >> 1)
		)

		if offac >= 0xAC {
			d.logSeg.Set(data.DByte(d.physOff+int(offac)), vac)
		}

		d.logSeg.Set(data.DByte(d.physOff+int(off56)), v56)
		d.logSeg.Set(data.DByte(d.physOff+int(i)), v00)
	}

	d.logOff += LogSectorLen
}
