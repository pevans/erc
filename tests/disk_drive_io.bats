setup_file() { load disk_drive_io_helper; setup_file; }
setup()      { load disk_drive_io_helper; setup; }
teardown()   { load disk_drive_io_helper; teardown; }

# ---------------------------------------------------------------------------
# Motor Control
# ---------------------------------------------------------------------------

@test "C0E9 turns motor on and enables full-speed mode" {
	disk_run \
		'LDA $C0E8' \
		'LDA $C0E9' \
		'LDA #$01' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ "$(_last_comp "FullSpeed")" == "true" ]]
}

@test "C0E8 turns motors off and disables full-speed mode" {
	disk_run \
		'LDA $C0E9' \
		'LDA $C0E8' \
		'LDA #$01' 'STA $00' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ "$(_last_comp "FullSpeed")" == "false" ]]
}

# ---------------------------------------------------------------------------
# Drive Selection
# ---------------------------------------------------------------------------

@test "C0EB selects drive 2 which reads FF with no disk" {
	disk_run \
		'LDA $C0EB' \
		'LDA $C0E9' \
		'LDA $C0EE' \
		'LDA $C0EC' \
		'STA $00' \
		'LDA $C0E8' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val="$(_last_mem "0000")"
	local dec
	dec=$((16#${val#\$}))
	[[ $dec -eq 255 ]]
}

@test "C0EA switches back to drive 1 after selecting drive 2" {
	disk_run \
		'LDA $C0EB' \
		'LDA $C0EA' \
		'LDA $C0E9' \
		'LDA $C0EE' \
		'LDA $C0EC' \
		'STA $00' \
		'LDA $C0E8' \
		'.halt'
	[[ $status -eq 0 ]]
	local val
	val="$(_last_mem "0000")"
	[[ -n "$val" ]]
}

# ---------------------------------------------------------------------------
# Stepper Motor Phases
# ---------------------------------------------------------------------------

@test "phase switches step the drive head to a new track" {
	# Search for the D5 AA 96 address field prologue on track 0, read the
	# 4-and-4 encoded track number, then step the head forward via phase
	# switches and repeat on the new track. The track numbers must differ.
	DISK_STEPS=10000 disk_run \
		'LDA $C0E9' \
		'LDA $C0EE' \
		'find1: LDA $C0EC' \
		'CMP #$D5' \
		'BNE find1' \
		'LDA $C0EC' \
		'CMP #$AA' \
		'BNE find1' \
		'LDA $C0EC' \
		'CMP #$96' \
		'BNE find1' \
		'LDA $C0EC' \
		'LDA $C0EC' \
		'LDA $C0EC' \
		'STA $00' \
		'LDA $C0EC' \
		'STA $01' \
		'LDA $C0E1' \
		'LDA $C0E3' \
		'LDA $C0E5' \
		'LDA $C0E7' \
		'LDA $C0E1' \
		'LDA $C0E3' \
		'find2: LDA $C0EC' \
		'CMP #$D5' \
		'BNE find2' \
		'LDA $C0EC' \
		'CMP #$AA' \
		'BNE find2' \
		'LDA $C0EC' \
		'CMP #$96' \
		'BNE find2' \
		'LDA $C0EC' \
		'LDA $C0EC' \
		'LDA $C0EC' \
		'STA $02' \
		'LDA $C0EC' \
		'STA $03' \
		'LDA $C0E8' \
		'.halt'
	[[ $status -eq 0 ]]
	# Decode 4-and-4 track numbers from each address field
	local t0 t1
	t0=$(decode_4and4 "$(_last_mem "0000")" "$(_last_mem "0001")")
	t1=$(decode_4and4 "$(_last_mem "0002")" "$(_last_mem "0003")")
	# First search should find track 0
	[[ $t0 -eq 0 ]]
	# After phase switches, should be on a different track
	[[ $t1 -ne 0 ]]
}

@test "reverse phase switches step the drive head backward" {
	# Step forward past track 0, record the track number, then step
	# backward using the reverse phase sequence and verify the head moved
	# to a lower-numbered track.
	DISK_STEPS=10000 disk_run \
		'LDA $C0E9' \
		'LDA $C0EE' \
		'LDA $C0E1' \
		'LDA $C0E3' \
		'LDA $C0E5' \
		'LDA $C0E7' \
		'LDA $C0E1' \
		'LDA $C0E3' \
		'fwd: LDA $C0EC' \
		'CMP #$D5' \
		'BNE fwd' \
		'LDA $C0EC' \
		'CMP #$AA' \
		'BNE fwd' \
		'LDA $C0EC' \
		'CMP #$96' \
		'BNE fwd' \
		'LDA $C0EC' \
		'LDA $C0EC' \
		'LDA $C0EC' \
		'STA $00' \
		'LDA $C0EC' \
		'STA $01' \
		'LDA $C0E1' \
		'LDA $C0E7' \
		'LDA $C0E5' \
		'LDA $C0E3' \
		'bwd: LDA $C0EC' \
		'CMP #$D5' \
		'BNE bwd' \
		'LDA $C0EC' \
		'CMP #$AA' \
		'BNE bwd' \
		'LDA $C0EC' \
		'CMP #$96' \
		'BNE bwd' \
		'LDA $C0EC' \
		'LDA $C0EC' \
		'LDA $C0EC' \
		'STA $02' \
		'LDA $C0EC' \
		'STA $03' \
		'LDA $C0E8' \
		'.halt'
	[[ $status -eq 0 ]]
	local fwd bwd
	fwd=$(decode_4and4 "$(_last_mem "0000")" "$(_last_mem "0001")")
	bwd=$(decode_4and4 "$(_last_mem "0002")" "$(_last_mem "0003")")
	# Forward track should be past track 0
	[[ $fwd -gt 0 ]]
	# Backward track should be less than forward track
	[[ $bwd -lt $fwd ]]
}

# ---------------------------------------------------------------------------
# Read/Write Mode Switching
# ---------------------------------------------------------------------------

@test "C0EF sets write mode and C0EE restores read mode" {
	disk_run \
		'LDA $C0E9' \
		'LDA $C0EE' \
		'LDA $C0EC' \
		'STA $00' \
		'LDA $C0EF' \
		'LDA $C0EC' \
		'STA $01' \
		'LDA $C0EE' \
		'LDA $C0EC' \
		'STA $02' \
		'LDA $C0E8' \
		'.halt'
	[[ $status -eq 0 ]]
	local raw0 raw1 raw2
	raw0="$(_last_mem "0000")"
	raw1="$(_last_mem "0001")"
	raw2="$(_last_mem "0002")"
	local v0 v1 v2
	v0=$((16#${raw0#\$}))
	v1=$((16#${raw1#\$}))
	v2=$((16#${raw2#\$}))
	# In read mode, all physical disk bytes are GCR-encoded (>= $96)
	[[ $v0 -ge $((16#96)) ]]
	# In write mode, C0EC returns 0 (no read occurs)
	[[ $v1 -eq 0 ]]
	# Back in read mode, disk data returns again
	[[ $v2 -ge $((16#96)) ]]
}

# ---------------------------------------------------------------------------
# Write Protection
# ---------------------------------------------------------------------------

@test "C0EE bit 7 is set when write protection is enabled" {
	# Toggle write protection on via ctrl-a,w shortcut at steps 0 and 1,
	# then read C0EE to observe the write-protect status in bit 7.
	DISK_KEYS="0:ctrl-a,1:w" disk_run \
		'NOP' \
		'NOP' \
		'LDA $C0E9' \
		'LDA $C0EE' \
		'STA $00' \
		'LDA $C0E8' \
		'.halt'
	[[ $status -eq 0 ]]
	[[ "$(_last_comp "WriteProtect")" == "true" ]]
	local val
	val="$(_last_mem "0000")"
	[[ -n "$val" ]]
	local dec
	dec=$((16#${val#\$}))
	[[ $((dec & 128)) -ne 0 ]]
}
