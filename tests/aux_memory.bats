setup_file() { load aux_memory_helper; setup_file; }
setup()      { load aux_memory_helper; setup; }
teardown()   { load aux_memory_helper; teardown; }

# ---------------------------------------------------------------------------
# Section 3.1: Default State
# ---------------------------------------------------------------------------

@test "default state: ReadAux is false (C013 returns 00)" {
	aux_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C013' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "default state: WriteAux is false (C014 returns 00)" {
	aux_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C014' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "default state: 80STORE is off (C018 returns 00)" {
	aux_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C018' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "default state: PAGE2 is off (C01C returns 00)" {
	aux_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C01C' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

@test "default state: HIRES is off (C01D returns 00)" {
	aux_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C01D' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
}

# ---------------------------------------------------------------------------
# Section 4: Soft Switches -- RAMRD / RAMWRT
# ---------------------------------------------------------------------------

@test "C003 write sets ReadAux true" {
	aux_run \
		'STA $C003' \
		'.halt'
	[[ "$(_first_comp MemReadAux)" == "true" ]]
}

@test "C005 write sets WriteAux true" {
	aux_run \
		'STA $C005' \
		'.halt'
	[[ "$(_last_comp MemWriteAux)" == "true" ]]
}

@test "C002 write sets ReadAux false" {
	aux_run_80store \
		'STA $C003' \
		'STA $C002' \
		'.halt'
	[[ "$(_last_comp MemReadAux)" == "false" ]]
}

@test "C004 write sets WriteAux false" {
	aux_run \
		'STA $C005' \
		'STA $C004' \
		'.halt'
	[[ "$(_last_comp MemWriteAux)" == "false" ]]
}

@test "reading RAMRD/RAMWRT on switches returns 00 and has no effect" {
	aux_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'STA $01' \
		'LDA $C003' \
		'STA $00' \
		'LDA $C005' \
		'STA $01' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
	[[ "$(_last_mem 0001)" == '$00' ]]
	[[ "$(_last_comp MemReadAux)" != "true" ]]
	[[ "$(_last_comp MemWriteAux)" != "true" ]]
}

@test "reading RAMRD/RAMWRT off switches returns 00 and has no effect" {
	aux_run_80store_with_mem \
		'STA $C003' \
		'STA $C005' \
		'LDA #$FF' \
		'STA $00' \
		'STA $01' \
		'LDA $C002' \
		'STA $00' \
		'LDA $C004' \
		'STA $01' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
	[[ "$(_last_mem 0001)" == '$00' ]]
	[[ "$(_last_comp MemReadAux)" == "true" ]]
	[[ "$(_last_comp MemWriteAux)" == "true" ]]
}

# ---------------------------------------------------------------------------
# Section 4.1: Independence
# ---------------------------------------------------------------------------

@test "ReadAux and WriteAux are independent" {
	aux_run \
		'STA $C005' \
		'.halt'
	[[ "$(_last_comp MemWriteAux)" == "true" ]]
	[[ "$(_last_comp MemReadAux)" != "true" ]]
}

# ---------------------------------------------------------------------------
# Section 5: Status Switches
# ---------------------------------------------------------------------------

@test "C013 returns 80 when ReadAux is true" {
	aux_run_80store_with_mem \
		'STA $C003' \
		'LDA $C013' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$80' ]]
}

@test "C014 returns 80 when WriteAux is true" {
	aux_run_with_mem \
		'STA $C005' \
		'LDA $C014' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$80' ]]
}

@test "C018 returns 80 when 80STORE is on" {
	aux_run_with_mem \
		'STA $C001' \
		'LDA $C018' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$80' ]]
}

@test "C01C returns 80 when PAGE2 is on" {
	aux_run_with_mem \
		'STA $C055' \
		'LDA $C01C' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$80' ]]
}

@test "C01D returns 80 when HIRES is on" {
	aux_run_with_mem \
		'STA $C057' \
		'LDA $C01D' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$80' ]]
}

@test "writing to status addresses has no effect" {
	aux_run \
		'STA $C013' \
		'STA $C014' \
		'STA $C018' \
		'STA $C01C' \
		'STA $C01D' \
		'.halt'
	[[ "$(_last_comp MemReadAux)" != "true" ]]
	[[ "$(_last_comp MemWriteAux)" != "true" ]]
	[[ "$(_last_comp DisplayStore80)" != "true" ]]
	[[ "$(_last_comp DisplayPage2)" != "true" ]]
	[[ "$(_last_comp DisplayHires)" != "true" ]]
}

# ---------------------------------------------------------------------------
# Section 6.1: 80STORE Soft Switch
# ---------------------------------------------------------------------------

@test "C001 write enables 80STORE" {
	aux_run \
		'STA $C001' \
		'.halt'
	[[ "$(_last_comp DisplayStore80)" == "true" ]]
}

@test "C000 write disables 80STORE" {
	aux_run \
		'STA $C001' \
		'STA $C000' \
		'.halt'
	[[ "$(_last_comp DisplayStore80)" == "false" ]]
}

@test "reading C000 does not trigger 80STORE" {
	aux_run \
		'LDA $C000' \
		'.halt'
	[[ "$(_last_comp DisplayStore80)" != "true" ]]
}

@test "reading C001 does not trigger 80STORE" {
	aux_run \
		'LDA $C001' \
		'.halt'
	[[ "$(_last_comp DisplayStore80)" != "true" ]]
}

# ---------------------------------------------------------------------------
# Section 6.2: PAGE2 Soft Switch
# ---------------------------------------------------------------------------

@test "C055 write enables PAGE2" {
	aux_run \
		'STA $C055' \
		'.halt'
	[[ "$(_last_comp DisplayPage2)" == "true" ]]
}

@test "C054 write disables PAGE2" {
	aux_run \
		'STA $C055' \
		'STA $C054' \
		'.halt'
	[[ "$(_last_comp DisplayPage2)" == "false" ]]
}

@test "reading C055 triggers PAGE2 and returns 00" {
	aux_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C055' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
	[[ "$(_last_comp DisplayPage2)" == "true" ]]
}

@test "reading C054 triggers PAGE2 off and returns 00" {
	aux_run_with_mem \
		'STA $C055' \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C054' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
	[[ "$(_last_comp DisplayPage2)" == "false" ]]
}

# ---------------------------------------------------------------------------
# Section 6.3: HIRES Soft Switch
# ---------------------------------------------------------------------------

@test "C057 write enables HIRES" {
	aux_run \
		'STA $C057' \
		'.halt'
	[[ "$(_last_comp DisplayHires)" == "true" ]]
}

@test "C056 write disables HIRES" {
	aux_run \
		'STA $C057' \
		'STA $C056' \
		'.halt'
	[[ "$(_last_comp DisplayHires)" == "false" ]]
}

@test "reading C057 triggers HIRES and returns 00" {
	aux_run_with_mem \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C057' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
	[[ "$(_last_comp DisplayHires)" == "true" ]]
}

@test "reading C056 triggers HIRES off and returns 00" {
	aux_run_with_mem \
		'STA $C057' \
		'LDA #$FF' \
		'STA $00' \
		'LDA $C056' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$00' ]]
	[[ "$(_last_comp DisplayHires)" == "false" ]]
}

# ---------------------------------------------------------------------------
# Section 6.3: 80STORE Override Rules -- Text Page
# ---------------------------------------------------------------------------

@test "80STORE+PAGE2 routes text page to aux" {
	# Write $CC to aux text page via 80STORE+PAGE2, then $DD to main.
	# Reading with PAGE2 on should return aux value.
	aux_run_with_mem \
		'STA $C001' \
		'STA $C055' \
		'LDA #$CC' \
		'STA $0400' \
		'STA $C054' \
		'LDA #$DD' \
		'STA $0400' \
		'STA $C055' \
		'LDA $0400' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$CC' ]]
}

@test "80STORE on PAGE2 off routes text page to main" {
	# Write $DD to main, then $CC to aux. Read with PAGE2 off -> main.
	aux_run_with_mem \
		'STA $C001' \
		'LDA #$DD' \
		'STA $0400' \
		'STA $C055' \
		'LDA #$CC' \
		'STA $0400' \
		'STA $C054' \
		'LDA $0400' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$DD' ]]
}

# ---------------------------------------------------------------------------
# Section 6.3: 80STORE Override Rules -- Hires Page
# ---------------------------------------------------------------------------

@test "80STORE+HIRES+PAGE2 routes hires page to aux" {
	aux_run_with_mem \
		'STA $C001' \
		'STA $C057' \
		'STA $C055' \
		'LDA #$EE' \
		'STA $2000' \
		'STA $C054' \
		'LDA #$DD' \
		'STA $2000' \
		'STA $C055' \
		'LDA $2000' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$EE' ]]
}

@test "80STORE+HIRES on PAGE2 off routes hires page to main" {
	aux_run_with_mem \
		'STA $C001' \
		'STA $C057' \
		'LDA #$DD' \
		'STA $2000' \
		'STA $C055' \
		'LDA #$EE' \
		'STA $2000' \
		'STA $C054' \
		'LDA $2000' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$DD' ]]
}

@test "80STORE without HIRES does not override hires page" {
	# Write $AA to aux hires page (80STORE+HIRES+PAGE2), then turn off
	# HIRES and write $BB to $2000 (should follow RAMWRT -> main).
	# Re-enable HIRES and read from aux: should still be $AA, not $BB.
	aux_run_with_mem \
		'STA $C001' \
		'STA $C057' \
		'STA $C055' \
		'LDA #$AA' \
		'STA $2000' \
		'STA $C056' \
		'LDA #$BB' \
		'STA $2000' \
		'STA $C057' \
		'LDA $2000' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$AA' ]]
}

# ---------------------------------------------------------------------------
# Section 8/9: Read and Write Logic
# ---------------------------------------------------------------------------

@test "WriteAux routes writes to aux segment" {
	# Write $AA to main $0300, then $BB to aux $0300 via WriteAux.
	# Read back from main (ReadAux off) -> should be $AA.
	aux_run_with_mem \
		'LDA #$AA' \
		'STA $0300' \
		'STA $C005' \
		'LDA #$BB' \
		'STA $0300' \
		'STA $C004' \
		'LDA $0300' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$AA' ]]
}

@test "ReadAux routes reads to aux segment" {
	# Write $AA to main $0300, then $BB to aux $0300 via WriteAux.
	# Read back from aux (ReadAux on) -> should be $BB.
	aux_run_80store_with_mem \
		'LDA #$AA' \
		'STA $0300' \
		'STA $C005' \
		'LDA #$BB' \
		'STA $0300' \
		'STA $C004' \
		'STA $C003' \
		'LDA $0300' \
		'STA $C002' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$BB' ]]
}

@test "without 80STORE text page follows RAMWRT" {
	# With 80STORE off, text page writes follow RAMWRT.
	# Write $AA to main $0400, then $BB to aux via WriteAux.
	# Read back from main -> $AA.
	aux_run_with_mem \
		'LDA #$AA' \
		'STA $0400' \
		'STA $C005' \
		'LDA #$BB' \
		'STA $0400' \
		'STA $C004' \
		'LDA $0400' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$AA' ]]
}

@test "without 80STORE hires page follows RAMWRT" {
	aux_run_with_mem \
		'LDA #$AA' \
		'STA $2000' \
		'STA $C005' \
		'LDA #$BB' \
		'STA $2000' \
		'STA $C004' \
		'LDA $2000' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$AA' ]]
}

# ---------------------------------------------------------------------------
# 80STORE override ignores RAMRD/RAMWRT
# ---------------------------------------------------------------------------

@test "80STORE text override ignores RAMRD" {
	# Write $CC to aux $07F0 via WriteAux (80STORE off), then $DD to main.
	# Enable 80STORE. With RAMRD=aux and PAGE2 off, read $07F0 -> main ($DD).
	# Code runs from $0400 with 80STORE on; PAGE2 stays off throughout.
	aux_run_80store_with_mem \
		'STA $C000' \
		'STA $C005' \
		'LDA #$CC' \
		'STA $07F0' \
		'STA $C004' \
		'LDA #$DD' \
		'STA $07F0' \
		'STA $C001' \
		'STA $C003' \
		'LDA $07F0' \
		'STA $C002' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$DD' ]]
}

@test "80STORE hires override ignores RAMRD" {
	# Write $CC to aux $2000 via WriteAux (80STORE off), then $DD to main.
	# Enable 80STORE+HIRES. With RAMRD=aux and PAGE2 off, read $2000 -> main.
	# Code runs from $0400 with 80STORE on; PAGE2 stays off throughout.
	aux_run_80store_with_mem \
		'STA $C000' \
		'STA $C005' \
		'LDA #$CC' \
		'STA $2000' \
		'STA $C004' \
		'LDA #$DD' \
		'STA $2000' \
		'STA $C001' \
		'STA $C057' \
		'STA $C003' \
		'LDA $2000' \
		'STA $C002' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$DD' ]]
}

# ---------------------------------------------------------------------------
# Page 2 ranges unaffected by 80STORE
# ---------------------------------------------------------------------------

@test "80STORE does not affect text page 2 range" {
	# $0800-$0BFF follows RAMWRT even with 80STORE+PAGE2 on.
	aux_run_with_mem \
		'STA $C001' \
		'STA $C055' \
		'LDA #$AA' \
		'STA $0800' \
		'STA $C005' \
		'LDA #$BB' \
		'STA $0800' \
		'STA $C004' \
		'LDA $0800' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$AA' ]]
}

@test "80STORE does not affect hires page 2 range" {
	# $4000-$5FFF follows RAMWRT even with 80STORE+HIRES+PAGE2 on.
	aux_run_with_mem \
		'STA $C001' \
		'STA $C057' \
		'STA $C055' \
		'LDA #$AA' \
		'STA $4000' \
		'STA $C005' \
		'LDA #$BB' \
		'STA $4000' \
		'STA $C004' \
		'LDA $4000' \
		'STA $00' \
		'.halt'
	[[ "$(_last_mem 0000)" == '$AA' ]]
}
