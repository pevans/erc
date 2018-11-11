package a2

import (
	"github.com/pevans/erc/pkg/mach"
)

type Decoder struct {
	src       *mach.Segment
	dst       *mach.Segment
	imageType int
}

//  00    01    02    03    04    05    06    07    08    09    0A    0B    0C    0D    0E    0F
var conv6bit = []mach.Byte{
	0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, // 00
	0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0X00, 0X04, 0XFF, 0XFF, 0X08, 0X0C, 0XFF, 0X10, 0X14, 0X18, // 10
	0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0X1C, 0X20, 0XFF, 0XFF, 0XFF, 0X24, 0X28, 0X2C, 0X30, 0X34, // 20
	0XFF, 0XFF, 0X38, 0X3C, 0X40, 0X44, 0X48, 0X4C, 0XFF, 0X50, 0X54, 0X58, 0X5C, 0X60, 0X64, 0X68, // 30
	0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0X6C, 0XFF, 0X70, 0X74, 0X78, // 40
	0XFF, 0XFF, 0XFF, 0X7C, 0XFF, 0XFF, 0X80, 0X84, 0XFF, 0X88, 0X8C, 0X90, 0X94, 0X98, 0X9C, 0XA0, // 50
	0XFF, 0XFF, 0XFF, 0XFF, 0XFF, 0XA4, 0XA8, 0XAC, 0XFF, 0XB0, 0XB4, 0XB8, 0XBC, 0XC0, 0XC4, 0XC8, // 60
	0XFF, 0XFF, 0XCC, 0XD0, 0XD4, 0XD8, 0XDC, 0XE0, 0XFF, 0XE4, 0XE8, 0XEC, 0XF0, 0XF4, 0XF8, 0XFC, // 70
}

func NewDecoder(imgType int, src *mach.Segment) *Decoder {
	return &Decoder{
		src:       src,
		imageType: imgType,
	}
}

func (d *Decoder) DecodeNIB() (*mach.Segment, error) {
	dst := mach.NewSegment(d.src.Size())
	_, err := dst.CopySlice(0, d.src.Mem)

	if err != nil {
		return nil, err
	}

	return dst, nil
}

func (d *Decoder) DecodeDOS() (*mach.Segment, error) {
	d.dst = mach.NewSegment(DD140K)
	doff := 0

	for track := 0; track < DDNumTracks; track++ {
		doff += d.DecodeTrack(track, doff)
	}

	return d.dst, nil
}

func (d *Decoder) DecodeTrack(track, doff int) int {
	soff := (track * PhysTrackLen) + PhysTrackHeader

	for sect := 0; sect < DDNumSectors; sect++ {
		doff := (track * LogTrackLen) + (LogicalSector(d.imageType, sect) * LogSectorLen)
		_ = d.DecodeSector(track, sect, doff, soff)
		soff += PhysSectorLen
	}

	return LogTrackLen
}

func (d *Decoder) DecodeSector(track, sect, doff, soff int) int {
	// Skip header and the data marker
	soff += PhysSectorHeader + 3

	conv := make([]mach.Byte, 0x157)
	for i := 0; i < 0x157; i++ {
		conv[i] = conv6bit[d.src.Get(mach.DByte(soff+i))&0x7F]
	}

	xor := make([]mach.Byte, 0x156)
	for i, lval := 0, mach.Byte(0); i < 0x156; i++ {
		xor[i] = lval ^ conv[i]
		lval = xor[i]
	}

	for i := mach.Byte(0); i < 0x56; i++ {
		var (
			offac, off56  mach.Byte
			vac, v56, v00 mach.Byte
		)

		offac = i + 0xAC
		off56 = i + 0x56

		vac = (xor[int(offac)+0x56] & 0xFC) | ((xor[i] & 0x80) >> 7) | ((xor[i] & 0x40) >> 5)
		v56 = (xor[int(off56)+0x56] & 0xFC) | ((xor[i] & 0x20) >> 5) | ((xor[i] & 0x10) >> 3)
		v00 = (xor[i+0x56] & 0xFC) | ((xor[i] & 0x08) >> 3) | ((xor[i] & 0x04) >> 1)

		if offac >= 0xAC {
			d.dst.Set(mach.DByte(doff+int(offac)), vac)
		}

		d.dst.Set(mach.DByte(doff+int(off56)), v56)
		d.dst.Set(mach.DByte(doff+int(i)), v00)
	}

	return LogSectorLen
}
