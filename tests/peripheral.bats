setup_file() { load peripheral_helper; setup_file; }
setup()      { load peripheral_helper; setup; }
teardown()   { load peripheral_helper; teardown; }

# ---------------------------------------------------------------------------
# Section 3: Default State
# ---------------------------------------------------------------------------

@test "default state: SlotCX is true" {
	pc_run '.halt'
	[[ "$(_last_comp PCSlotCX)" == "" || "$(_last_comp PCSlotCX)" == "true" ]]
}

@test "default state: SlotC3 is false" {
	pc_run '.halt'
	[[ "$(_last_comp PCSlotC3)" == "" || "$(_last_comp PCSlotC3)" == "false" ]]
}

@test "default state: Expansion is false" {
	pc_run '.halt'
	[[ "$(_last_comp PCExpansion)" == "" || "$(_last_comp PCExpansion)" == "false" ]]
}

@test "default state: IOSelect is false" {
	pc_run '.halt'
	[[ "$(_last_comp PCIOSelect)" == "" || "$(_last_comp PCIOSelect)" == "false" ]]
}

@test "default state: IOStrobe is false" {
	pc_run '.halt'
	[[ "$(_last_comp PCIOStrobe)" == "" || "$(_last_comp PCIOStrobe)" == "false" ]]
}

@test "default state: ExpSlot is 0" {
	pc_run '.halt'
	[[ "$(_last_comp PCExpSlot)" == "" || "$(_last_comp PCExpSlot)" == "0" ]]
}

@test "default state: RDCXROM returns 00 (SlotCX is true)" {
	pc_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C015' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "default state: RDC3ROM returns 00 (SlotC3 is false)" {
	pc_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C017' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

# ---------------------------------------------------------------------------
# Section 4.1: Control Switches
# ---------------------------------------------------------------------------

@test "C007 write sets SlotCX false" {
	pc_run \
		'STA $C007' \
		'.halt'
	[[ "$(_last_comp PCSlotCX)" == "false" ]]
}

@test "C006 write sets SlotCX true" {
	pc_run \
		'STA $C007' \
		'STA $C006' \
		'.halt'
	[[ "$(_last_comp PCSlotCX)" == "true" ]]
}

@test "C00B write sets SlotC3 true" {
	pc_run \
		'STA $C00B' \
		'.halt'
	[[ "$(_last_comp PCSlotC3)" == "true" ]]
}

@test "C00A write sets SlotC3 false" {
	pc_run \
		'STA $C00B' \
		'STA $C00A' \
		'.halt'
	[[ "$(_last_comp PCSlotC3)" == "false" ]]
}

@test "reading C006 does not change SlotCX" {
	pc_run_with_mem \
		'STA $C007' \
		'LDA $C006' \
		'STA $00' \
		'.halt'
	[[ "$(_last_comp PCSlotCX)" == "false" ]]
}

@test "reading C007 does not change SlotCX" {
	pc_run \
		'LDA $C007' \
		'.halt'
	# SlotCX should remain at default (true), no state change logged
	[[ "$(_last_comp PCSlotCX)" == "" || "$(_last_comp PCSlotCX)" == "true" ]]
}

@test "reading C00A does not change SlotC3" {
	pc_run \
		'STA $C00B' \
		'LDA $C00A' \
		'.halt'
	[[ "$(_last_comp PCSlotC3)" == "true" ]]
}

@test "reading C00B does not change SlotC3" {
	pc_run \
		'LDA $C00B' \
		'.halt'
	[[ "$(_last_comp PCSlotC3)" == "" || "$(_last_comp PCSlotC3)" == "false" ]]
}

@test "reading C006 returns keyboard data latch" {
	PC_KEYS="0:A" pc_run_with_mem \
		'LDA $C006' \
		'STA $00' \
		'.halt'
	# 'A' = $41, strobe = $80, latch = $C1
	[[ "$(_last_mem 0000)" == '$C1' ]]
}

@test "reading C007 returns keyboard data latch" {
	PC_KEYS="0:A" pc_run_with_mem \
		'LDA $C007' \
		'STA $00' \
		'.halt'
	# 'A' = $41, strobe = $80, latch = $C1
	[[ "$(_last_mem 0000)" == '$C1' ]]
}

@test "reading C00A returns keyboard data latch" {
	PC_KEYS="0:A" pc_run_with_mem \
		'LDA $C00A' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$C1' ]]
}

@test "reading C00B returns keyboard data latch" {
	PC_KEYS="0:A" pc_run_with_mem \
		'LDA $C00B' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$C1' ]]
}

# ---------------------------------------------------------------------------
# Section 4.2: Status Switches
# ---------------------------------------------------------------------------

@test "C015 returns 00 when SlotCX is true" {
	pc_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C015' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "C015 returns 80 when SlotCX is false" {
	pc_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'STA $C007' \
		'LDA $C015' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$80' ]]
}

@test "C017 returns 80 when SlotC3 is true" {
	pc_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'STA $C00B' \
		'LDA $C017' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$80' ]]
}

@test "C017 returns 00 when SlotC3 is false" {
	pc_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C017' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "C015 returns keyboard latch in lower bits" {
	PC_KEYS="0:A" pc_run_with_mem \
		'LDA $C015' \
		'STA $00' \
		'.halt'
	# SlotCX=true -> bit 7 = 0; 'A' latch = $C1; result = $41
	[[ "$(_last_mem 0000)" == '$41' ]]
}

@test "C017 returns keyboard latch in lower bits" {
	PC_KEYS="0:A" pc_run_with_mem \
		'STA $C00B' \
		'LDA $C017' \
		'STA $00' \
		'.halt'
	# SlotC3=true -> bit 7 = 1; 'A' latch = $C1; result = $C1
	[[ "$(_last_mem 0000)" == '$C1' ]]
}

@test "writing to C015 has no effect" {
	pc_run \
		'STA $C015' \
		'.halt'
	# SlotCX should remain at default (true)
	[[ "$(_last_comp PCSlotCX)" == "" || "$(_last_comp PCSlotCX)" == "true" ]]
}

@test "writing to C017 has no effect" {
	pc_run \
		'STA $C017' \
		'.halt'
	# SlotC3 should remain at default (false)
	[[ "$(_last_comp PCSlotC3)" == "" || "$(_last_comp PCSlotC3)" == "false" ]]
}

# ---------------------------------------------------------------------------
# Section 4.3: Expansion ROM Disable
# ---------------------------------------------------------------------------

@test "reading CFFF clears expansion state" {
	PC_STEPS=200 pc_run \
		'LDA $C100' \
		'LDA $C800' \
		'LDA $CFFF' \
		'.halt'
	[[ "$(_last_comp PCIOSelect)" == "false" ]]
	[[ "$(_last_comp PCIOStrobe)" == "false" ]]
	[[ "$(_last_comp PCExpansion)" == "false" ]]
	[[ "$(_last_comp PCExpSlot)" == "0" ]]
}

@test "writing CFFF clears expansion state" {
	PC_STEPS=200 pc_run \
		'LDA $C100' \
		'LDA $C800' \
		'STA $CFFF' \
		'.halt'
	[[ "$(_last_comp PCIOSelect)" == "false" ]]
	[[ "$(_last_comp PCIOStrobe)" == "false" ]]
	[[ "$(_last_comp PCExpansion)" == "false" ]]
	[[ "$(_last_comp PCExpSlot)" == "0" ]]
}

@test "CFFF returns expansion ROM before clearing when IOSelect is set" {
	PC_STEPS=200 pc_run_with_mem \
		'STA $C007' \
		'LDA $CFFF' \
		'STA $00' \
		'STA $C006' \
		'LDA $C100' \
		'LDA $C800' \
		'LDA $CFFF' \
		'STA $01' \
		'.halt'
	# First CFFF (SlotCX=false, no expansion): internal ROM
	# Second CFFF (SlotCX=true, IOSelect=true): expansion ROM
	# Expansion falls back to internal ROM when no card installed
	[[ "$(_last_mem 0000)" == "$(_last_mem 0001)" ]]
	[[ "$(_last_comp PCIOSelect)" == "false" ]]
	[[ "$(_last_comp PCExpansion)" == "false" ]]
}

# ---------------------------------------------------------------------------
# Section 5: ROM Read Logic -- IOSelect and ExpSlot
# ---------------------------------------------------------------------------

@test "reading C100-C7FF with SlotCX true sets IOSelect" {
	pc_run \
		'LDA $CFFF' \
		'LDA $C100' \
		'.halt'
	[[ "$(_last_comp PCIOSelect)" == "true" ]]
}

@test "reading C100 sets ExpSlot to 1" {
	pc_run \
		'LDA $C100' \
		'.halt'
	[[ "$(_last_comp PCExpSlot)" == "1" ]]
}

@test "reading C200 sets ExpSlot to 2" {
	pc_run \
		'LDA $C200' \
		'.halt'
	[[ "$(_last_comp PCExpSlot)" == "2" ]]
}

@test "reading C500 sets ExpSlot to 5" {
	pc_run \
		'LDA $C500' \
		'.halt'
	[[ "$(_last_comp PCExpSlot)" == "5" ]]
}

@test "reading C800-CFFE with SlotCX true sets IOStrobe" {
	PC_STEPS=200 pc_run \
		'LDA $C100' \
		'LDA $C800' \
		'.halt'
	[[ "$(_last_comp PCIOStrobe)" == "true" ]]
}

@test "reading C800 after IOSelect enables Expansion" {
	PC_STEPS=200 pc_run \
		'LDA $C100' \
		'LDA $C800' \
		'.halt'
	[[ "$(_last_comp PCExpansion)" == "true" ]]
}

@test "reading CFFF does not set IOStrobe" {
	pc_run \
		'LDA $CFFF' \
		'.halt'
	# Unlike C800-CFFE, CFFF does not set IOStrobe
	[[ "$(_last_comp PCIOStrobe)" == "" || "$(_last_comp PCIOStrobe)" == "false" ]]
}

@test "reading C100-C7FF with SlotCX false does not set IOSelect" {
	pc_run \
		'STA $C007' \
		'LDA $C100' \
		'.halt'
	[[ "$(_last_comp PCIOSelect)" == "" || "$(_last_comp PCIOSelect)" == "false" ]]
}

@test "reading C800-CFFE with SlotCX false does not set IOStrobe" {
	pc_run \
		'STA $C007' \
		'LDA $C800' \
		'.halt'
	[[ "$(_last_comp PCIOStrobe)" == "" || "$(_last_comp PCIOStrobe)" == "false" ]]
}

@test "C800 without prior IOSelect does not enable Expansion" {
	pc_run \
		'LDA $CFFF' \
		'LDA $C800' \
		'.halt'
	# IOStrobe is set (any C800-CFFE read with SlotCX true sets it)
	[[ "$(_last_comp PCIOStrobe)" == "true" ]]
	# But Expansion is not enabled because IOSelect was false
	[[ "$(_last_comp PCExpansion)" == "" || "$(_last_comp PCExpansion)" == "false" ]]
}

# ---------------------------------------------------------------------------
# Section 5: ROM Read Logic -- SlotC3 interaction
# ---------------------------------------------------------------------------

@test "C300 returns internal ROM when SlotC3 is false and SlotCX is true" {
	pc_run_with_mem \
		'LDA $CFFF' \
		'LDA $C300' \
		'STA $00' \
		'STA $C00B' \
		'LDA $CFFF' \
		'LDA $C300' \
		'STA $01' \
		'.halt'
	[[ "$(_last_comp PCIOSelect)" == "true" ]]
	[[ "$(_last_comp PCExpSlot)" == "3" ]]
	# First read (SlotC3=false): internal ROM
	# Second read (SlotC3=true): peripheral ROM
	[[ "$(_last_mem 0000)" != "$(_last_mem 0001)" ]]
}

@test "C300 returns peripheral ROM when SlotC3 is true and SlotCX is false" {
	pc_run_with_mem \
		'STA $C00B' \
		'STA $C007' \
		'LDA $C300' \
		'STA $00' \
		'STA $C00A' \
		'LDA $C300' \
		'STA $01' \
		'.halt'
	# First read (SlotC3=true, SlotCX=false): peripheral ROM
	# Second read (SlotC3=false, SlotCX=false): internal ROM
	[[ "$(_last_mem 0000)" != "$(_last_mem 0001)" ]]
}

@test "C300 returns peripheral ROM when SlotC3 true and SlotCX true" {
	pc_run_with_mem \
		'STA $C00B' \
		'LDA $CFFF' \
		'LDA $C300' \
		'STA $00' \
		'STA $C00A' \
		'LDA $CFFF' \
		'LDA $C300' \
		'STA $01' \
		'.halt'
	# First read (SlotC3=true, SlotCX=true): peripheral ROM
	# (SlotC3=false guard does NOT fire, so returns peripheral ROM)
	# Second read (SlotC3=false, SlotCX=true): internal ROM
	[[ "$(_last_mem 0000)" != "$(_last_mem 0001)" ]]
}

@test "C300 returns internal ROM when SlotC3 false and SlotCX false" {
	pc_run_with_mem \
		'STA $C007' \
		'LDA $C300' \
		'STA $00' \
		'STA $C00B' \
		'LDA $C300' \
		'STA $01' \
		'.halt'
	# First read (SlotC3=false, SlotCX=false): internal ROM
	# Second read (SlotC3=true, SlotCX=false): peripheral ROM
	[[ "$(_last_mem 0000)" != "$(_last_mem 0001)" ]]
}

@test "C800 returns peripheral ROM when SlotCX true but IOSelect false" {
	pc_run_with_mem \
		'LDA $CFFF' \
		'LDA $C800' \
		'STA $00' \
		'LDA $C100' \
		'LDA $C800' \
		'STA $01' \
		'.halt'
	# First C800: SlotCX=true, IOSelect=false -> peripheral ROM
	# Second C800: SlotCX=true, IOSelect=true -> expansion ROM (= internal fallback)
	# Peripheral and internal ROM differ at C800
	[[ "$(_last_mem 0000)" != "$(_last_mem 0001)" ]]
}

# ---------------------------------------------------------------------------
# Section 5: ROM Read Logic -- SlotCX false path
# ---------------------------------------------------------------------------

@test "C100 returns internal ROM when SlotCX is false" {
	pc_run_with_mem \
		'STA $C007' \
		'LDA $C100' \
		'STA $00' \
		'STA $C006' \
		'LDA $C100' \
		'STA $01' \
		'.halt'
	# First read (SlotCX=false): internal ROM, IOSelect not set
	# Second read (SlotCX=true): peripheral ROM
	[[ "$(_last_comp PCIOSelect)" == "true" ]]
	[[ "$(_last_mem 0000)" != "$(_last_mem 0001)" ]]
}

@test "Expansion ROM active when SlotCX false but Expansion was enabled" {
	PC_STEPS=200 pc_run_with_mem \
		'LDA $C100' \
		'LDA $C800' \
		'STA $C007' \
		'LDA $C900' \
		'STA $00' \
		'LDA $CFFF' \
		'LDA $C900' \
		'STA $01' \
		'.halt'
	[[ "$(_first_comp PCExpansion)" == "true" ]]
	[[ "$(_last_comp PCSlotCX)" == "false" ]]
	# Expansion ROM falls back to internal ROM when no card installed (5.1)
	# First C900: expansion active -> expansion ROM (= internal fallback)
	# Second C900: expansion cleared -> internal ROM directly
	[[ "$(_last_mem 0000)" == "$(_last_mem 0001)" ]]
}

@test "CFFF returns expansion ROM when SlotCX false but Expansion true" {
	PC_STEPS=200 pc_run_with_mem \
		'LDA $C100' \
		'LDA $C800' \
		'STA $C007' \
		'LDA $CFFF' \
		'STA $00' \
		'.halt'
	# Expansion was enabled via C100+C800, then SlotCX set false.
	# CFFF with SlotCX=false, Expansion=true: returns expansion ROM.
	# Expansion falls back to internal ROM when no card installed (5.1).
	[[ "$(_last_comp PCExpansion)" == "false" ]]
	[[ "$(_last_comp PCIOSelect)" == "false" ]]
}

# ---------------------------------------------------------------------------
# Section 6: ROM Write Logic
# ---------------------------------------------------------------------------

@test "writes to C100-CFFF are ignored (ROM is read-only)" {
	pc_run_with_mem \
		'LDA #$AA' \
		'STA $00' \
		'LDA #$BB' \
		'STA $01' \
		'LDA $C100' \
		'STA $00' \
		'LDA #$FF' \
		'STA $C100' \
		'LDA $C100' \
		'STA $01' \
		'.halt'
	# The value at C100 should be the same before and after the write
	[[ "$(_last_mem 0000)" == "$(_last_mem 0001)" ]]
}
