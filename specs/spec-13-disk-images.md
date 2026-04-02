---
Specification: 13
Category: Storage
Drafted At: 2026-03-26
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes how Apple II disk images are represented, encoded,
decoded, and used within erc. It covers the three supported image formats
(DOS 3.3, ProDOS, and nibble), the encoding schemes that translate between
logical and physical data, the drive emulation that reads and writes encoded
data, and the disk controller soft switches that connect the drive to the CPU.

# 2. Disk Geometry

An Apple II 5.25-inch floppy disk contains 35 concentric tracks, numbered 0
through 34. Each track contains 16 sectors, numbered 0 through 15. Each
logical sector holds 256 bytes. A full disk therefore holds:

    35 tracks x 16 sectors x 256 bytes = 143,360 bytes (140 KB)

The drive head is positioned using half-track steps, giving 70 possible
positions (0-69). The actual track number is the half-track position divided
by 2 (integer division).

# 3. Image Formats

Erc recognizes three disk image formats, determined by file extension.

## 3.1. DOS 3.3 (.dsk, .do)

A 143,360-byte file containing logical sector data. Sectors are stored in the
file in logical order, but the logical-to-physical mapping uses the DOS 3.3
interleave table (see section 5). This is the most common format for Apple II
disk images.

## 3.2. ProDOS (.po)

A 143,360-byte file, identical in size to DOS 3.3, but using the ProDOS
sector interleave table (see section 5). The file layout is the same -- 35
tracks of 16 sectors of 256 bytes each -- but which logical sector
corresponds to which physical sector differs.

## 3.3. Nibble (.nib)

A 232,960-byte file (35 tracks x 6,656 bytes per track) containing
already-encoded physical data. Nibble images bypass the encoding and decoding
process entirely. They exist because some software uses tricks that store data
in areas that the standard encoding scheme treats as padding or overhead.

# 4. Logical vs. Physical Data

The terms "logical" and "physical" refer to two representations of the same
data:

- **Logical data** is the raw 256 bytes per sector that software actually
  uses -- program code, file contents, and so on. Logical data is what you
  find in a .dsk or .po file.

- **Physical data** is what the disk drive hardware reads and writes. It
  includes the logical data, but that data has been encoded to be safe for
  the disk medium, and wrapped in metadata fields, checksums, and sync gaps.

The Apple II disk drive cannot distinguish an intentional zero bit from a read
error. The encoding scheme ensures there is never more than two consecutive
zero bits in the data stream. This is the fundamental reason physical encoding
exists.

Why not have all data in a physical format? Certainly some software is shared
in that form. The benefit of the logical format is that the data is available
for outside editing. With it, you can view and modify the program code or some
data within the program with a hex editor. Such a task would be much more
difficult with data in the physical format.

# 5. Sector Interleaving

When physically encoding a track, sectors are not written in logical order.
Instead, a sector interleave table maps physical sector positions to logical
sector numbers. This allows the operating system to read sequential logical
sectors without waiting for a full disk rotation between each one.

## 5.1. DOS 3.3 Interleave

    Physical:  0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15
    Logical:   0  7 14  6 13  5 12  4 11  3 10  2  9  1  8 15

## 5.2. ProDOS Interleave

    Physical:  0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15
    Logical:   0  8  1  9  2 10  3 11  4 12  5 13  6 14  7 15

## 5.3. Direction of Mapping

The interleave table answers the question: "when encoding physical sector N,
which logical sector should I read from?" During encoding, for each physical
sector position, the encoder looks up the corresponding logical sector and
copies that sector's 256 bytes into the physical stream. During decoding, the
reverse mapping is applied.

# 6. Physical Track Layout

A physically encoded track consists of a leading gap followed by 16 sectors.
Each sector contains an address field, a gap, a data field, and a trailing
gap.

## 6.1. Track Structure

    [ gap1 (48 bytes) ]
    [ sector 0 ][ sector 1 ] ... [ sector 15 ]

## 6.2. Sector Structure

    [ address field ][ gap2 (6 bytes) ][ data field ][ gap3 (27 bytes) ]

## 6.3. Gap Bytes

All gap bytes are 0xFF (self-sync bytes). They give the drive hardware time
to synchronize and give software a window to prepare for the next field.

- **Gap 1**: 48 bytes at the start of each track.
- **Gap 2**: 6 bytes between the address and data fields of a sector.
- **Gap 3**: 27 bytes after the data field of each sector.

## 6.4. Physical Track and Sector Sizes

Each physical sector occupies 396 bytes (0x18C). A full physical track
occupies 6,384 bytes (0x18F0), which is the 48 gap1 bytes plus 16 sectors of
396 bytes each.

For nibble images, the track length is 6,656 bytes (0x1A00). The difference
exists because nibble images preserve the full raw track data as it appears on
the physical medium, including additional padding and sync bytes that the
standard encoding scheme does not generate.

# 7. Address Field

The address field identifies which track and sector the following data field
belongs to. It uses 4-and-4 encoding for its metadata.

## 7.1. Address Field Layout

    D5 AA 96          prologue (3 bytes)
    [volume]           4-and-4 encoded (2 bytes)
    [track]            4-and-4 encoded (2 bytes)
    [sector]           4-and-4 encoded (2 bytes)
    [checksum]         4-and-4 encoded (2 bytes)
    DE AA EB           epilogue (3 bytes)

The checksum is the XOR of the volume, track, and sector values (before
4-and-4 encoding).

## 7.2. Volume Marker

The volume byte is hard-coded to 0xFE. Original Apple II disks used this to
identify which disk was inserted, but in practice it is always the same value
in erc.

## 7.3. 4-and-4 Encoding

4-and-4 encoding splits one byte into two bytes, each with enough high bits
set to be safe for the disk medium:

    first  = ((value >> 1) & 0x55) | 0xAA
    second = (value & 0x55) | 0xAA

To decode:

    value = ((first & 0x55) << 1) | (second & 0x55)

# 8. Data Field

The data field contains the actual sector data, encoded using 6-and-2
encoding. It transforms 256 logical bytes into 343 encoded data bytes, which
are then translated through a GCR (group coded recording) lookup table.

## 8.1. Data Field Layout

    D5 AA AD           prologue (3 bytes)
    [encoded data]     343 GCR-translated bytes
    [checksum]         1 GCR-translated byte
    DE AA EB           epilogue (3 bytes)

## 8.2. Prologue and Epilogue

The prologue bytes D5 AA AD mark the start of the data field. The epilogue
bytes DE AA EB mark the end. These byte sequences are chosen so they cannot
appear in encoded data, making them unambiguous markers.

Note: the address field epilogue (DE AA EB) is not strictly validated during
decoding, because some Apple II software writes partial or nonstandard
epilogues. The data field epilogue is validated.

## 8.3. 6-and-2 Encoding

The 256 bytes of a logical sector are split into two buffers:

- **Six-block**: 256 bytes, each containing the upper 6 bits of a source
  byte (shifted right by 2).
- **Two-block**: 86 bytes, each containing the lower 2 bits from three
  source bytes, packed together.

The packing works as follows. For each source byte at index `i`:

    six[i] = source[i] >> 2
    rev    = bit-reverse of the low 2 bits: ((source[i] & 2) >> 1) | ((source[i] & 1) << 1)
    two[i % 86] |= rev << ((i / 86) * 2)

The two-block has 86 entries because 256 / 3 rounds up to 86 (the last entry
only packs bits from one byte rather than three).

## 8.4. XOR Chaining

Before being written to the physical stream, the six-block and two-block
bytes are XOR-chained. This serves as an error-detection mechanism:

1. The first two-block byte is written as-is.
2. Each subsequent two-block byte is XOR'd with the previous two-block byte.
3. The first six-block byte is XOR'd with the last two-block byte.
4. Each subsequent six-block byte is XOR'd with the previous six-block byte.
5. A final checksum byte (the last six-block value, un-XOR'd) is appended.

During decoding, the XOR chain is reversed to reconstruct the original
values. If the final checksum does not resolve to zero, the data is corrupt.

## 8.5. GCR Translation

Every 6-bit value (0x00-0x3F) is translated through a 64-entry GCR lookup
table before being written to the physical stream. The table maps each 6-bit
value to an 8-bit byte that guarantees no more than one consecutive zero bit.

    0x00-0x0F: 96 97 9A 9B 9D 9E 9F A6 A7 AB AC AD AE AF B2 B3
    0x10-0x1F: B4 B5 B6 B7 B9 BA BB BC BD BE BF CB CD CE CF D3
    0x20-0x2F: D6 D7 D9 DA DB DC DD DE DF E5 E6 E7 E9 EA EB EC
    0x30-0x3F: ED EE EF F2 F3 F4 F5 F6 F7 F9 FA FB FC FD FE FF

During decoding, a reverse map is constructed from the same table to
translate physical bytes back into 6-bit values.

# 9. Encoding and Decoding

## 9.1. Encoding (Logical to Physical)

Given a logical disk image (143,360 bytes), encoding produces a physical
segment. For DOS 3.3 and ProDOS images, this segment is 223,440 bytes
(35 tracks x 6,384 bytes per track). This is smaller than the 232,960 bytes
of a nibble image because the encoder produces only the exact gaps, fields,
and metadata needed -- no extra padding. Nibble images are larger because they
preserve the full raw track contents as captured from a real disk, which
includes additional sync bytes and other artifacts. The process:

1. For each of the 35 tracks:
   a. Write gap1 (48 x 0xFF).
   b. For each physical sector 0-15:
      - Look up the corresponding logical sector from the interleave table.
      - Read 256 bytes from the logical segment at the appropriate offset.
      - Write the address field (4-and-4 encoded metadata).
      - Write gap2 (6 x 0xFF).
      - Write the data field (6-and-2 encoded, XOR-chained, GCR-translated).
      - Write gap3 (27 x 0xFF).

For nibble images, encoding is a no-op -- the segment is returned unchanged.

## 9.2. Decoding (Physical to Logical)

Given a physical segment, decoding produces a logical segment (143,360
bytes). The process:

1. For each of the 35 tracks:
   a. For each physical sector 0-15:
      - Scan for the address field prologue (D5 AA 96).
      - Decode the 4-and-4 metadata and validate the checksum.
      - Verify track and sector numbers match expectations.
      - Scan for the data field prologue (D5 AA AD).
      - Read and GCR-decode 343 bytes of encoded data plus 1 checksum byte.
      - Reverse the XOR chain to recover the six-block and two-block.
      - Reconstruct the 256 logical bytes and validate the checksum.
      - Scan for the data field epilogue (DE AA EB).
      - Write the 256 bytes to the correct position in the logical segment.

If any step fails -- a prologue is not found, a checksum does not match, or
a track/sector number is unexpected -- the decoder returns an error
identifying the track and sector where the failure occurred. Decoding does not
attempt to skip or recover from corrupt sectors.

For nibble images, decoding is a no-op.

# 10. Drive Emulation

The drive emulation models a Disk II floppy drive. Each drive maintains two
memory segments: the original logical image, and the working physical
(encoded) data.

## 10.1. Loading

When a disk image is loaded:

1. The file extension determines the image type (DOS 3.3, ProDOS, or nibble).
2. The file size is validated against the expected size for that type.
3. The raw bytes are stored in the image segment.
4. The image segment is encoded to produce the physical data segment.
5. The sector position is reset to 0; the track position is left unchanged.

## 10.2. Saving

When a disk image is saved:

1. The physical data segment is decoded back to logical form.
2. The logical bytes are written to the original file.

This round-trip means that any writes the emulated software makes to the
physical data are correctly translated back to logical form.

## 10.3. Write Protection

Write protection is managed independently of the disk image or drive state.
It is toggled by the user (for example, via a keyboard shortcut or command
line flag) rather than being an inherent property of the image file. When
write protection is enabled, all write operations to the disk are silently
ignored. The write-protect status is reported to the CPU via the $C0EE soft
switch (bit 7).

## 10.4. Multiple Drives

The Apple II supports two Disk II drives. Each drive maintains its own
independent state: motor, track position, sector position, latch, mode, and
write protection. Only one drive can be selected at a time, via the $C0EA and
$C0EB soft switches. All disk controller operations apply to the currently
selected drive.

## 10.5. The Latch

All data transfer between the disk and the CPU passes through a single-byte
buffer called the latch. The latch acts as an airlock:

- **LoadLatch** copies the byte at the current disk position into the latch.
  This only happens if the disk has shifted since the last load.
- **ReadLatch** returns the latch value. On the first read, the value is
  returned unmodified. On subsequent reads (without a new LoadLatch), the
  high bit is cleared (masked with 0x7F).
- **SetLatch** sets the latch to a value provided by the CPU (for writes).
- **WriteLatch** writes the latch value to the disk at the current position.
  This only succeeds if the drive is in write mode, the motor is on, the disk
  is not write-protected, and the latch value has its high bit set.

## 10.6. Position

The drive tracks two positions:

- **Track position** (0-69 in half-tracks): controlled by the stepper motor
  phases. The actual track number is `trackPos / 2`. The physical Disk II
  drive has mechanical stops that prevent the head from moving below track 0
  or above track 34. Erc enforces the same limits by clamping the half-track
  position to the range 0-69.
- **Sector position** (offset within the current track): advanced by 1 byte
  after each read or write operation. Wraps around the track length, since
  the disk is circular.

Erc does not emulate precise disk rotation. A real drive spins continuously,
and the sector position advances with the physical rotation of the disk. Erc
instead advances the sector position by exactly one byte per shift operation
($C0EC access). This is a common simplification used by other Apple II
emulators and is sufficient for all known software.

The absolute position in the physical data segment is:

    (track * trackLen) + sectorPos

where `trackLen` is 6,384 for DOS 3.3/ProDOS or 6,656 for nibble images.

The current sector number can be derived from the sector position:

    sector = sectorPos / physSectorLen

where `physSectorLen` is 396 bytes (0x18C).

## 10.7. Stepper Motor Phases

The drive head moves via a stepper motor with 4 phases. Phase transitions
determine whether the head steps inward (+1 half-track), outward (-1), or
stays put (0). A phase table maps each (current phase, new phase)
combination to a step direction:

    Phase transitions (rows = current, columns = new):
          1   2   3   4
    1:    0   1   0  -1
    2:   -1   0   1   0
    3:    0  -1   0   1
    4:    1   0  -1   0

The new phase is derived from the soft switch address: odd addresses 0x1,
0x3, 0x5, 0x7 correspond to phases 1-4 respectively. Even addresses 0x0,
0x2, 0x4, 0x6 turn the corresponding phase off. Turning a phase off does not
move the head; only turning a phase on triggers a step calculation.

# 11. Disk Controller Soft Switches

The Disk II controller occupies 16 soft switch addresses whose base depends
on which slot the controller card is installed in. The base address is
$C080 + (slot * $10). In erc, the disk controller is in slot 6, giving a
range of $C0E0-$C0EF. A controller in slot 7 would use $C0F0-$C0FF instead.

Most Apple II software accesses these switches using absolute indexed
addressing with the X register set to the slot number times 16 (e.g.,
`LDA $C08C,X` with X=$60 for slot 6). This allows the same code to work with
a controller in any slot.

The addresses below assume slot 6. Both reads and writes to these addresses
trigger the same behavior.

## 11.1. Switch Table

    $C0E0-$C0E7   Phase control (stepper motor)         (base+$0 - base+$7)
    $C0E8         Motor off (both drives)               (base+$8)
    $C0E9         Motor on (selected drive)             (base+$9)
    $C0EA         Select drive 1                        (base+$A)
    $C0EB         Select drive 2                        (base+$B)
    $C0EC         Shift (read or write one byte)        (base+$C)
    $C0ED         Load/peek latch                       (base+$D)
    $C0EE         Set read mode (write-protect bit 7)   (base+$E)
    $C0EF         Set write mode                        (base+$F)

## 11.2. The Shift Operation ($C0EC)

This is the primary data transfer switch. Its behavior depends on the drive
mode:

- **Read mode** (or write-protected): load a byte from the disk into the
  latch, return the latch value, and advance the sector position by 1.
- **Write mode**: write the latch value to the disk and advance the sector
  position by 1.

The shift operation works regardless of whether the motor is on or off. In
practice, software only accesses $C0EC while the motor is running, but the
emulator does not gate reads on the motor state. Write operations do require
the motor to be on (see section 10.5).

## 11.3. The Latch Switch ($C0ED)

- **On write** (when in write mode and motor is on): sets the latch to the
  written value.
- **On read** (when in read mode): returns the current latch value without
  loading new data from the disk.

## 11.4. Full-Speed Mode

When the drive motor turns on, erc switches to full-speed emulation (skipping
cycle timing). This is because the Apple II's RWTS (Read/Write Track/Sector)
routines contain tight timing loops calibrated to the disk's rotation speed.
Emulating these loops at the correct clock speed would make disk operations
painfully slow. Full-speed mode is disabled when the motor turns off, or when
the speaker is actively producing sound.

Real Disk II drives keep the motor spinning briefly after the off switch is
hit (roughly one second). Erc does not emulate this spin-down delay -- the
motor stops immediately when the off switch is accessed.

# 12. Boot Sequence

The Disk II controller card includes a small boot ROM (sometimes called the
P6 ROM) that is mapped into the slot's ROM space ($C600 for slot 6). When the
Apple II powers on or resets, the CPU begins executing the system ROM, which
eventually transfers control to the controller's boot ROM. The boot ROM loads
one page (256 bytes) from track 0, sector 0 of the disk in drive 1 into
memory at $0800 and jumps to that address. This boot sector typically contains
a short loader that reads the rest of the operating system from disk using the
RWTS routines.

The provision of the boot ROM itself is outside the scope of this spec. The
details of what the boot sector code does vary by operating system (DOS 3.3,
ProDOS, etc.) and are also not part of this spec -- they are determined by the
contents of the disk image itself.

# 13. Data Flow Summary

## 13.1. Loading and Reading

    .dsk file (logical, 143,360 bytes)
      --> a2enc.Encode --> physical segment (223,440 bytes)
      --> stored in drive.data
      --> CPU reads via $C0EC soft switch
      --> LoadLatch, ReadLatch, Shift

## 13.2. Writing and Saving

    CPU writes via $C0ED (SetLatch) then $C0EC (WriteLatch + Shift)
      --> bytes written to drive.data (physical segment)
      --> a2enc.Decode --> logical segment (143,360 bytes)
      --> written to .dsk file

## 13.3. Nibble Images

    .nib file (physical, 232,960 bytes)
      --> no encoding needed, used directly as drive.data
      --> CPU reads/writes as normal
      --> no decoding needed for save, written directly
