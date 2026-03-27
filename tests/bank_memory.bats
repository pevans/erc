setup_file() { load bank_memory_helper; setup_file; }
setup()      { load bank_memory_helper; setup; }
teardown()   { load bank_memory_helper; teardown; }

# ---------------------------------------------------------------------------
# Section 3.1: Default State
# ---------------------------------------------------------------------------

@test "default state: ReadRAM is false (C012 returns 00)" {
	# C012 returns $80 if ReadRAM is true, $00 if false.
	# Store to $00 which we pre-fill with $FF to detect a $00 write.
	bank_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C012' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "default state: WriteRAM is true" {
	# Verify by writing to bank RAM and reading back.
	# Default: WriteRAM=true, DFBlockBank2=true. Write $AB to $D000,
	# then enable ReadRAM and read it back.
	bank_run_with_mem \
		'LDA #$AB' \
		'STA $D000' \
		'LDA $C080' \
		'LDA $D000' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$AB' ]]
}

@test "default state: DFBlockBank2 is true (C011 returns 80)" {
	# C011 returns $80 if DFBlockBank2 is true.
	bank_run_with_mem \
		'LDA $C011' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$80' ]]
}

# ---------------------------------------------------------------------------
# Section 5.1: Switch Table -- non-write-enable switches
# ---------------------------------------------------------------------------

@test "C080 sets ReadRAM=true WriteRAM=false" {
	bank_run 'LDA $C080' '.halt'
	[[ "$(_last_comp BankReadRAM)" == "true" ]]
	[[ "$(_last_comp BankWriteRAM)" == "false" ]]
}

@test "C080 selects bank 2" {
	# Switch to bank 1 first, then C080 should switch back to bank 2
	bank_run 'LDA $C088' 'LDA $C080' '.halt'
	[[ "$(_last_comp BankDFBlockBank2)" == "true" ]]
}

@test "C082 sets ReadRAM=false WriteRAM=false" {
	bank_run 'LDA $C082' '.halt'
	[[ "$(_last_comp BankWriteRAM)" == "false" ]]
}

@test "C082 selects bank 2" {
	bank_run 'LDA $C088' 'LDA $C082' '.halt'
	[[ "$(_last_comp BankDFBlockBank2)" == "true" ]]
}

@test "C088 sets ReadRAM=true WriteRAM=false DFBlockBank2=false" {
	bank_run 'LDA $C088' '.halt'
	[[ "$(_last_comp BankReadRAM)" == "true" ]]
	[[ "$(_last_comp BankWriteRAM)" == "false" ]]
	[[ "$(_last_comp BankDFBlockBank2)" == "false" ]]
}

@test "C08A sets ReadRAM=false WriteRAM=false DFBlockBank2=false" {
	bank_run 'LDA $C08A' '.halt'
	[[ "$(_last_comp BankWriteRAM)" == "false" ]]
	[[ "$(_last_comp BankDFBlockBank2)" == "false" ]]
}

# ---------------------------------------------------------------------------
# Section 5.1: Switch Table -- duplicates
# ---------------------------------------------------------------------------

@test "C084 duplicates C080" {
	bank_run 'LDA $C088' 'LDA $C084' '.halt'
	[[ "$(_last_comp BankReadRAM)" == "true" ]]
	[[ "$(_last_comp BankWriteRAM)" == "false" ]]
	[[ "$(_last_comp BankDFBlockBank2)" == "true" ]]
}

@test "C086 duplicates C082" {
	bank_run 'LDA $C088' 'LDA $C086' '.halt'
	[[ "$(_last_comp BankWriteRAM)" == "false" ]]
	[[ "$(_last_comp BankDFBlockBank2)" == "true" ]]
}

@test "C08C duplicates C088" {
	bank_run 'LDA $C08C' '.halt'
	[[ "$(_last_comp BankReadRAM)" == "true" ]]
	[[ "$(_last_comp BankWriteRAM)" == "false" ]]
	[[ "$(_last_comp BankDFBlockBank2)" == "false" ]]
}

@test "C08E duplicates C08A" {
	bank_run 'LDA $C08E' '.halt'
	[[ "$(_last_comp BankWriteRAM)" == "false" ]]
	[[ "$(_last_comp BankDFBlockBank2)" == "false" ]]
}

@test "C085 duplicates C081" {
	# Switch to bank 1 first so DFBlockBank2 change is logged
	bank_run \
		'LDA $C08A' \
		'LDA $C085' \
		'LDA $C085' \
		'.halt'
	[[ "$(_last_comp BankReadRAM)" == "false" ]]
	[[ "$(_last_comp BankWriteRAM)" == "true" ]]
	[[ "$(_last_comp BankDFBlockBank2)" == "true" ]]
}

@test "C087 duplicates C083" {
	# Switch to bank 1 first so DFBlockBank2 change is logged
	bank_run \
		'LDA $C08A' \
		'LDA $C087' \
		'LDA $C087' \
		'.halt'
	[[ "$(_last_comp BankReadRAM)" == "true" ]]
	[[ "$(_last_comp BankWriteRAM)" == "true" ]]
	[[ "$(_last_comp BankDFBlockBank2)" == "true" ]]
}

@test "C08D duplicates C089" {
	bank_run \
		'LDA $C08A' \
		'LDA $C08D' \
		'LDA $C08D' \
		'.halt'
	[[ "$(_last_comp BankReadRAM)" == "false" ]]
	[[ "$(_last_comp BankWriteRAM)" == "true" ]]
	[[ "$(_last_comp BankDFBlockBank2)" == "false" ]]
}

@test "C08F duplicates C08B" {
	bank_run \
		'LDA $C08A' \
		'LDA $C08F' \
		'LDA $C08F' \
		'.halt'
	[[ "$(_last_comp BankReadRAM)" == "true" ]]
	[[ "$(_last_comp BankWriteRAM)" == "true" ]]
	[[ "$(_last_comp BankDFBlockBank2)" == "false" ]]
}

# ---------------------------------------------------------------------------
# Section 5.3: Double-Read Write Protection
# ---------------------------------------------------------------------------

@test "single read of C083 does not enable WriteRAM" {
	bank_run \
		'LDA $C082' \
		'LDA $C083' \
		'.halt'
	[[ "$(_last_comp BankWriteRAM)" == "false" ]]
}

@test "double read of C083 enables WriteRAM" {
	bank_run \
		'LDA $C082' \
		'LDA $C083' \
		'LDA $C083' \
		'.halt'
	[[ "$(_last_comp BankWriteRAM)" == "true" ]]
}

@test "single read of C08B does not enable WriteRAM" {
	bank_run \
		'LDA $C08A' \
		'LDA $C08B' \
		'.halt'
	[[ "$(_last_comp BankWriteRAM)" == "false" ]]
}

@test "double read of C08B enables WriteRAM" {
	bank_run \
		'LDA $C08A' \
		'LDA $C08B' \
		'LDA $C08B' \
		'.halt'
	[[ "$(_last_comp BankWriteRAM)" == "true" ]]
}

@test "two different write-enable switches satisfy double-read" {
	bank_run \
		'LDA $C08A' \
		'LDA $C081' \
		'LDA $C083' \
		'.halt'
	[[ "$(_last_comp BankWriteRAM)" == "true" ]]
}

@test "non-write-enable switch between two write-enable reads resets counter" {
	bank_run \
		'LDA $C082' \
		'LDA $C083' \
		'LDA $C080' \
		'LDA $C083' \
		'.halt'
	[[ "$(_last_comp BankWriteRAM)" == "false" ]]
}

# ---------------------------------------------------------------------------
# Section 5.1: Write-enable switches with double read (C081, C089)
# ---------------------------------------------------------------------------

@test "double read of C081 enables WriteRAM with ReadRAM=false" {
	bank_run \
		'LDA $C082' \
		'LDA $C081' \
		'LDA $C081' \
		'.halt'
	[[ "$(_last_comp BankReadRAM)" == "false" ]]
	[[ "$(_last_comp BankWriteRAM)" == "true" ]]
}

@test "double read of C089 sets ReadRAM=false WriteRAM=true bank2=false" {
	bank_run \
		'LDA $C08A' \
		'LDA $C089' \
		'LDA $C089' \
		'.halt'
	[[ "$(_last_comp BankReadRAM)" == "false" ]]
	[[ "$(_last_comp BankWriteRAM)" == "true" ]]
	[[ "$(_last_comp BankDFBlockBank2)" == "false" ]]
}

# ---------------------------------------------------------------------------
# Section 5: Soft switch reads return 00
# ---------------------------------------------------------------------------

@test "reading a bank soft switch returns 00" {
	bank_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C080' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "writing to a bank soft switch has no effect" {
	# Start with default state (ReadRAM=false). Writing to C080 should
	# not trigger the switch; ReadRAM should remain false.
	bank_run \
		'STA $C080' \
		'.halt'
	[[ "$(_last_comp BankReadRAM)" != "true" ]]
}

# ---------------------------------------------------------------------------
# Section 6: Status Switches
# ---------------------------------------------------------------------------

@test "C011 returns 80 when DFBlockBank2 is true" {
	bank_run_with_mem \
		'LDA $C011' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$80' ]]
}

@test "C011 returns 00 when DFBlockBank2 is false" {
	bank_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C088' \
		'LDA $C011' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "C012 returns 00 when ReadRAM is false" {
	bank_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C012' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "C012 returns 80 when ReadRAM is true" {
	bank_run_with_mem \
		'LDA $C080' \
		'LDA $C012' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$80' ]]
}

@test "C016 returns 00 when SETSTDZP is active" {
	bank_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C016' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "C016 returns 80 when SETALTZP is active" {
	bank_run_with_mem \
		'STA $C009' \
		'LDA $C016' \
		'STA $C008' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$80' ]]
}

# ---------------------------------------------------------------------------
# Section 7: Zero Page and Stack Page Switching
# ---------------------------------------------------------------------------

@test "C008 write sets SETSTDZP" {
	bank_run \
		'STA $C009' \
		'STA $C008' \
		'.halt'
	[[ "$(_last_comp BankSysBlockAux)" == "false" ]]
}

@test "C009 write sets SETALTZP" {
	bank_run \
		'STA $C009' \
		'.halt'
	[[ "$(_last_comp BankSysBlockAux)" == "true" ]]
}

@test "reading C008 or C009 has no effect and returns 00" {
	bank_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'STA $01' \
		'LDA $C009' \
		'STA $00' \
		'LDA $C008' \
		'STA $01' \
		'.halt'
	# Both reads should return $00
	[[ "$(_last_mem 0000)" == '$00' ]]
	[[ "$(_last_mem 0001)" == '$00' ]]
	# SETALTZP should not have been activated by the reads
	[[ "$(_last_comp BankSysBlockAux)" != "true" ]]
}

# ---------------------------------------------------------------------------
# Section 8: Read Logic -- ROM vs RAM
# ---------------------------------------------------------------------------

@test "default state reads ROM in D000-FFFF range" {
	# ReadRAM is false by default, so reading D000 should return the ROM
	# byte at offset $1000 (= $D000 - $C000). Read both D000 and E000
	# to verify ROM is served from the expected offset. Write known values
	# to RAM first so we can confirm ROM -- not RAM -- is returned.
	bank_run_with_mem \
		'LDA #$00' \
		'STA $D000' \
		'STA $E000' \
		'LDA $D000' \
		'STA $00' \
		'LDA $E000' \
		'STA $01' \
		'.halt'
	# With WriteRAM=true (default), we wrote $00 to RAM. If we're reading
	# ROM, the values should differ from $00 (ROM is not all zeroes).
	[[ "$(_last_mem 0000)" != '$00' ]]
	[[ "$(_last_mem 0001)" != '$00' ]]
}

@test "C080 enables RAM reads in D000-FFFF range" {
	bank_run_with_mem \
		'LDA #$A5' \
		'STA $D000' \
		'LDA $C080' \
		'LDA $D000' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$A5' ]]
}

# ---------------------------------------------------------------------------
# Section 8/9: Bank 1 vs Bank 2 independence
# ---------------------------------------------------------------------------

@test "bank 1 and bank 2 hold independent values in D000-DFFF" {
	# Default: WriteRAM=true, DFBlockBank2=true (bank 2).
	# Write $AA to bank 2 D000.
	# Switch to bank 1 (C08B x2 = ReadRAM+WriteRAM, bank1).
	# Write $55 to bank 1 D000.
	# Switch to bank 2 (C083 x2 = ReadRAM+WriteRAM, bank2).
	# Read D000 -> should be $AA.
	# Switch to bank 1 (C08B x2).
	# Read D000 -> should be $55.
	bank_run_with_mem \
		'LDA #$AA' \
		'STA $D000' \
		'LDA $C08B' \
		'LDA $C08B' \
		'LDA #$55' \
		'STA $D000' \
		'LDA $C083' \
		'LDA $C083' \
		'LDA $D000' \
		'STA $00' \
		'LDA $C08B' \
		'LDA $C08B' \
		'LDA $D000' \
		'STA $01' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$AA' ]]
	[[ "$(_last_mem 0001)" == '$55' ]]
}

# ---------------------------------------------------------------------------
# Section 8/9: E000-FFFF shared between banks
# ---------------------------------------------------------------------------

@test "E000-FFFF is shared regardless of bank selection" {
	bank_run_with_mem \
		'LDA #$42' \
		'STA $E000' \
		'LDA $C08B' \
		'LDA $C08B' \
		'LDA $E000' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$42' ]]
}

# ---------------------------------------------------------------------------
# Section 9: Write Logic -- WriteRAM=false discards writes
# ---------------------------------------------------------------------------

@test "writes are discarded when WriteRAM is false" {
	# C083 x2 = ReadRAM+WriteRAM, bank2. Write $BE to D000.
	# C080 = ReadRAM, WriteRAM=false, bank2. Try to overwrite D000.
	# Read D000 -> should still be $BE.
	bank_run_with_mem \
		'LDA $C083' \
		'LDA $C083' \
		'LDA #$BE' \
		'STA $D000' \
		'LDA $C080' \
		'LDA #$00' \
		'STA $D000' \
		'LDA $D000' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$BE' ]]
}

@test "write and read back in E000-FFFF shared region" {
	bank_run_with_mem \
		'LDA #$CD' \
		'STA $E000' \
		'LDA $C080' \
		'LDA $E000' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$CD' ]]
}
