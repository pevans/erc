setup_file() { load disk_images_helper; setup_file; }
setup()      { load disk_images_helper; setup; }
teardown()   { load disk_images_helper; teardown; }

# ---------------------------------------------------------------------------
# Section 3: Image Formats
# ---------------------------------------------------------------------------

@test "encode accepts .dsk input" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
}

@test "encode accepts .do input" {
	asm 'NOP' '.halt'
	cp "$TMP/test.dsk" "$TMP/test.do"
	encode "$TMP/test.do" "$TMP/out.enc"
	[[ $status -eq 0 ]]
}

@test "encode accepts .po input" {
	asm 'NOP' '.halt'
	cp "$TMP/test.dsk" "$TMP/test.po"
	encode "$TMP/test.po" "$TMP/out.enc"
	[[ $status -eq 0 ]]
}

@test "encode rejects .nib input" {
	make_zeros "$TMP/test.nib" 232960
	encode "$TMP/test.nib" "$TMP/out.enc"
	[[ $status -ne 0 ]]
	[[ "$output" == *"already in nibble format"* ]]
}

@test "encode rejects wrong-size input" {
	make_zeros "$TMP/test.dsk" 1024
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -ne 0 ]]
	[[ "$output" == *"unexpected size"* ]]
}

# ---------------------------------------------------------------------------
# Section 6: Physical Track Layout
# ---------------------------------------------------------------------------

@test "encoded output is 223440 bytes" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	local size
	size=$(wc -c <"$TMP/out.enc" | tr -d ' ')
	[[ "$size" -eq 223440 ]]
}

@test "gap1 is 48 FF bytes at start of track 0" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	for i in 0 1 23 46 47; do
		[[ "$(byte_at "$TMP/out.enc" "$i")" == "ff" ]]
	done
}

@test "track 1 gap1 starts at offset 6384" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# Track 1 starts at PhysTrackLen = 6384 (0x18F0), with gap1 bytes
	[[ "$(byte_at "$TMP/out.enc" 6384)" == "ff" ]]
	[[ "$(byte_at "$TMP/out.enc" 6431)" == "ff" ]]
	# Byte at 6432 (6384+48) should be address prologue D5
	[[ "$(byte_at "$TMP/out.enc" 6432)" == "d5" ]]
}

@test "second sector starts 396 bytes after first" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# First sector at offset 48 (after gap1), second at 48+396=444
	[[ "$(byte_at "$TMP/out.enc" 444)" == "d5" ]]
	[[ "$(byte_at "$TMP/out.enc" 445)" == "aa" ]]
	[[ "$(byte_at "$TMP/out.enc" 446)" == "96" ]]
}

# ---------------------------------------------------------------------------
# Section 7: Address Field
# ---------------------------------------------------------------------------

@test "address field prologue is D5 AA 96" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# First address field starts at offset 48 (after gap1)
	[[ "$(byte_at "$TMP/out.enc" 48)" == "d5" ]]
	[[ "$(byte_at "$TMP/out.enc" 49)" == "aa" ]]
	[[ "$(byte_at "$TMP/out.enc" 50)" == "96" ]]
}

@test "volume marker is 4-and-4 encoded FE" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# 4-and-4 of 0xFE: first=((0xFE>>1)&0x55)|0xAA=0xFF, second=(0xFE&0x55)|0xAA=0xFE
	[[ "$(byte_at "$TMP/out.enc" 51)" == "ff" ]]
	[[ "$(byte_at "$TMP/out.enc" 52)" == "fe" ]]
}

@test "track 0 sector 0 encoded in address field" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# 4-and-4 of 0x00: first=0xAA, second=0xAA
	# Track at bytes 53-54
	[[ "$(byte_at "$TMP/out.enc" 53)" == "aa" ]]
	[[ "$(byte_at "$TMP/out.enc" 54)" == "aa" ]]
	# Sector at bytes 55-56
	[[ "$(byte_at "$TMP/out.enc" 55)" == "aa" ]]
	[[ "$(byte_at "$TMP/out.enc" 56)" == "aa" ]]
}

@test "address checksum is volume XOR track XOR sector" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# For track 0, sector 0: checksum = 0xFE ^ 0x00 ^ 0x00 = 0xFE
	# 4-and-4 of 0xFE = FF FE
	[[ "$(byte_at "$TMP/out.enc" 57)" == "ff" ]]
	[[ "$(byte_at "$TMP/out.enc" 58)" == "fe" ]]
}

@test "address field epilogue is DE AA EB" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	[[ "$(byte_at "$TMP/out.enc" 59)" == "de" ]]
	[[ "$(byte_at "$TMP/out.enc" 60)" == "aa" ]]
	[[ "$(byte_at "$TMP/out.enc" 61)" == "eb" ]]
}

@test "track 1 address field has track number 1" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# Track 1 address field at 6384+48=6432: D5 AA 96 [vol] [track] [sect] [cksum] DE AA EB
	# 4-and-4 of track 1: first=((1>>1)&0x55)|0xAA=0xAA, second=(1&0x55)|0xAA=0xAB
	[[ "$(byte_at "$TMP/out.enc" 6437)" == "aa" ]]
	[[ "$(byte_at "$TMP/out.enc" 6438)" == "ab" ]]
}

@test "physical sector 5 has sector 5 in address field" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# Physical sector 5 on track 0: offset = 48 + (5 * 396) = 2028
	# Prologue D5 AA 96
	[[ "$(byte_at "$TMP/out.enc" 2028)" == "d5" ]]
	[[ "$(byte_at "$TMP/out.enc" 2029)" == "aa" ]]
	[[ "$(byte_at "$TMP/out.enc" 2030)" == "96" ]]
	# 4-and-4 of sector 5: first=((5>>1)&0x55)|0xAA=0xAA, second=(5&0x55)|0xAA=0xAF
	# Sector bytes at prologue(3) + volume(2) + track(2) = offset+7
	[[ "$(byte_at "$TMP/out.enc" 2035)" == "aa" ]]
	[[ "$(byte_at "$TMP/out.enc" 2036)" == "af" ]]
}

@test "track 34 has valid address field at expected offset" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# Track 34 starts at 34 * 6384 = 217056; address field at 217056 + 48 = 217104
	[[ "$(byte_at "$TMP/out.enc" 217104)" == "d5" ]]
	[[ "$(byte_at "$TMP/out.enc" 217105)" == "aa" ]]
	[[ "$(byte_at "$TMP/out.enc" 217106)" == "96" ]]
	# Volume 0xFE: 4-and-4 = ff fe; track 34 (0x22): 4-and-4 = bb aa
	[[ "$(byte_at "$TMP/out.enc" 217107)" == "ff" ]]
	[[ "$(byte_at "$TMP/out.enc" 217108)" == "fe" ]]
	[[ "$(byte_at "$TMP/out.enc" 217109)" == "bb" ]]
	[[ "$(byte_at "$TMP/out.enc" 217110)" == "aa" ]]
}

# ---------------------------------------------------------------------------
# Section 8: Data Field
# ---------------------------------------------------------------------------

@test "gap2 is 6 FF bytes between address and data fields" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# Address field ends at byte 62, gap2 is bytes 62-67
	for i in 62 63 64 65 66 67; do
		[[ "$(byte_at "$TMP/out.enc" "$i")" == "ff" ]]
	done
}

@test "data field prologue is D5 AA AD" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# Data field starts at byte 68 (48 gap1 + 14 address + 6 gap2)
	[[ "$(byte_at "$TMP/out.enc" 68)" == "d5" ]]
	[[ "$(byte_at "$TMP/out.enc" 69)" == "aa" ]]
	[[ "$(byte_at "$TMP/out.enc" 70)" == "ad" ]]
}

@test "data field epilogue is DE AA EB" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# Data field: 3 prologue + 343 encoded = 346, epilogue at 68+346=414
	[[ "$(byte_at "$TMP/out.enc" 414)" == "de" ]]
	[[ "$(byte_at "$TMP/out.enc" 415)" == "aa" ]]
	[[ "$(byte_at "$TMP/out.enc" 416)" == "eb" ]]
}

@test "gap3 is 27 FF bytes after data field" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# Data field epilogue ends at byte 417, gap3 is bytes 417-443
	[[ "$(byte_at "$TMP/out.enc" 417)" == "ff" ]]
	[[ "$(byte_at "$TMP/out.enc" 430)" == "ff" ]]
	[[ "$(byte_at "$TMP/out.enc" 443)" == "ff" ]]
}

@test "encoded data bytes are valid GCR values" {
	asm 'NOP' '.halt'
	encode "$TMP/test.dsk" "$TMP/out.enc"
	[[ $status -eq 0 ]]
	# GCR values are >= 0x96. Check first few data bytes after prologue.
	for offset in 71 72 73 74 75; do
		local val
		val="$(byte_at "$TMP/out.enc" "$offset")"
		# Convert to decimal and check >= 0x96 (150)
		[[ $((16#$val)) -ge $((16#96)) ]]
	done
}

# ---------------------------------------------------------------------------
# Section 9: Encoding
# ---------------------------------------------------------------------------

@test "encoding different data produces different output" {
	# Encode an all-zero disk
	make_zeros "$TMP/zeros.dsk" 143360
	encode "$TMP/zeros.dsk" "$TMP/zeros.enc"
	[[ $status -eq 0 ]]

	# Encode a disk with a program
	asm 'LDA #$FF' 'STA $00' '.halt'
	encode "$TMP/test.dsk" "$TMP/prog.enc"
	[[ $status -eq 0 ]]

	# The outputs should differ (data fields differ, address fields same)
	! cmp -s "$TMP/zeros.enc" "$TMP/prog.enc"
}

@test ".dsk and .po produce different output for non-uniform data" {
	make_patterned "$TMP/patterned.dsk"
	encode "$TMP/patterned.dsk" "$TMP/dos.enc"
	[[ $status -eq 0 ]]

	cp "$TMP/patterned.dsk" "$TMP/patterned.po"
	encode "$TMP/patterned.po" "$TMP/pro.enc"
	[[ $status -eq 0 ]]

	# Different interleave tables must produce different physical output
	! cmp -s "$TMP/dos.enc" "$TMP/pro.enc"
}

@test ".dsk and .po encode identically for uniform data" {
	# When all sectors are identical (zeros), interleave doesn't matter.
	make_zeros "$TMP/uniform.dsk" 143360
	encode "$TMP/uniform.dsk" "$TMP/dos.enc"
	[[ $status -eq 0 ]]

	cp "$TMP/uniform.dsk" "$TMP/uniform.po"
	encode "$TMP/uniform.po" "$TMP/pro.enc"
	[[ $status -eq 0 ]]

	cmp -s "$TMP/dos.enc" "$TMP/pro.enc"
}

# ---------------------------------------------------------------------------
# Section 9.2: Decoding (round-trip)
# ---------------------------------------------------------------------------

@test "encode then decode round-trips to identical .dsk" {
	make_patterned "$TMP/original.dsk"
	encode "$TMP/original.dsk" "$TMP/encoded.enc"
	[[ $status -eq 0 ]]

	decode "$TMP/encoded.enc" "$TMP/decoded.dsk"
	[[ $status -eq 0 ]]

	cmp -s "$TMP/original.dsk" "$TMP/decoded.dsk"
}

@test "encode then decode round-trips to identical .po" {
	make_patterned "$TMP/original.dsk"
	cp "$TMP/original.dsk" "$TMP/original.po"
	encode "$TMP/original.po" "$TMP/encoded.enc"
	[[ $status -eq 0 ]]

	decode "$TMP/encoded.enc" "$TMP/decoded.po"
	[[ $status -eq 0 ]]

	cmp -s "$TMP/original.po" "$TMP/decoded.po"
}

# ---------------------------------------------------------------------------
# Section 10: Drive Emulation
# ---------------------------------------------------------------------------

@test "disk boots and reaches program" {
	disk_run 'STA $00' '.halt'
	[[ $status -eq 0 ]]
}

@test "write protection defaults to off" {
	# Write to $00 so state.log is guaranteed to exist, then check that
	# WriteProtect never appeared as "true".
	disk_run 'LDA #$01' 'STA $00' '.halt'
	[[ $status -eq 0 ]]
	! grep -q 'comp WriteProtect.*true' "$OUT/state.log"
}

# ---------------------------------------------------------------------------
# Section 11: Soft Switches
# ---------------------------------------------------------------------------

@test "C0EC reads a byte from disk" {
	disk_run \
		'LDA $C0E9' \
		'LDA $C0EE' \
		'LDA $C0EC' \
		'STA $00' \
		'LDA $C0E8' \
		'.halt'
	[[ $status -eq 0 ]]
	# The byte read from disk should be stored at $00
	local val
	val="$(_last_mem "0000")"
	# We can't predict the exact value, but it should have been written
	[[ -n "$val" ]]
}

@test "consecutive C0EC reads advance sector position" {
	disk_run \
		'LDA $C0E9' \
		'LDA $C0EE' \
		'LDA $C0EC' \
		'STA $00' \
		'LDA $C0EC' \
		'STA $01' \
		'LDA $C0EC' \
		'STA $02' \
		'LDA $C0E8' \
		'.halt'
	[[ $status -eq 0 ]]
	local v0 v1 v2
	v0="$(_last_mem "0000")"
	v1="$(_last_mem "0001")"
	v2="$(_last_mem "0002")"
	[[ -n "$v0" ]]
	[[ -n "$v1" ]]
	[[ -n "$v2" ]]
	# At least one pair must differ, proving the position advances
	[[ "$v0" != "$v1" || "$v1" != "$v2" ]]
}

@test "C0ED returns latch without advancing position" {
	disk_run \
		'LDA $C0E9' \
		'LDA $C0EE' \
		'LDA $C0EC' \
		'STA $00' \
		'LDA $C0ED' \
		'STA $01' \
		'LDA $C0ED' \
		'STA $02' \
		'LDA $C0EC' \
		'STA $03' \
		'LDA $C0E8' \
		'.halt'
	[[ $status -eq 0 ]]
	local v0 v1 v2 v3
	v0="$(_last_mem "0000")"
	v1="$(_last_mem "0001")"
	v2="$(_last_mem "0002")"
	v3="$(_last_mem "0003")"
	[[ -n "$v1" ]]
	[[ -n "$v2" ]]
	# Two C0ED reads without intervening C0EC must return the same value
	[[ "$v1" == "$v2" ]]
}

@test "C0EE in read mode reports write-protect status" {
	disk_run \
		'LDA $C0E9' \
		'LDA $C0EE' \
		'STA $00' \
		'LDA $C0E8' \
		'.halt'
	[[ $status -eq 0 ]]
	# Disk is not write-protected, so bit 7 of $C0EE should be 0.
	# The value stored should have bit 7 clear.
	local val
	val="$(_last_mem "0000")"
	[[ -n "$val" ]]
	local dec
	dec=$((16#${val#\$}))
	[[ $((dec & 128)) -eq 0 ]]
}

# ---------------------------------------------------------------------------
# Section 12: Boot Sequence
# ---------------------------------------------------------------------------

@test "boot loads program from track 0 into memory" {
	disk_run \
		'LDA #$AB' \
		'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val="$(_last_mem "0000")"
	local hex
	hex="$(printf '%02X' "0x${val#\$}")"
	[[ "$hex" == "AB" ]]
}
