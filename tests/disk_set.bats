setup_file() { load disk_set_helper; setup_file; }
setup()      { load disk_set_helper; setup; }
teardown()   { load disk_set_helper; teardown; }

# ---------------------------------------------------------------------------
# Appending Images
# ---------------------------------------------------------------------------

@test "missing file produces an error before the emulator starts" {
	make_disk a
	run "$ERC_BIN" headless --steps 100 "$TMP/a.dsk" "$TMP/nonexistent.dsk"
	[[ $status -ne 0 ]]
	[[ $output == *"nonexistent.dsk"* ]]
}

# ---------------------------------------------------------------------------
# Next
# ---------------------------------------------------------------------------

@test "next advances from disk 0 to disk 1" {
	make_disk a && make_disk b
	diskset_run 500 "100:ctrl-a,101:n" "$TMP/a.dsk" "$TMP/b.dsk"
	[[ $status -eq 0 ]]
	grep -q 'comp DiskIndex .* -> 1' "$OUT/state.log"
}

@test "next wraps from last disk back to disk 0" {
	make_disk a && make_disk b
	diskset_run 500 "100:ctrl-a,101:n,200:ctrl-a,201:n" \
		"$TMP/a.dsk" "$TMP/b.dsk"
	[[ $status -eq 0 ]]
	grep -q 'comp DiskIndex .* -> 1' "$OUT/state.log"
	grep -q 'comp DiskIndex .* -> 0' "$OUT/state.log"
}

# ---------------------------------------------------------------------------
# Previous
# ---------------------------------------------------------------------------

@test "previous wraps from disk 0 to the last disk" {
	make_disk a && make_disk b
	diskset_run 500 "100:ctrl-a,101:p" "$TMP/a.dsk" "$TMP/b.dsk"
	[[ $status -eq 0 ]]
	grep -q 'comp DiskIndex .* -> 1' "$OUT/state.log"
}

@test "previous after next returns to disk 0" {
	make_disk a && make_disk b
	diskset_run 500 "100:ctrl-a,101:n,200:ctrl-a,201:p" \
		"$TMP/a.dsk" "$TMP/b.dsk"
	[[ $status -eq 0 ]]
	grep -q 'comp DiskIndex .* -> 1' "$OUT/state.log"
	grep -q 'comp DiskIndex .* -> 0' "$OUT/state.log"
}

# ---------------------------------------------------------------------------
# Index Tracking
# ---------------------------------------------------------------------------

@test "initial disk index is 0" {
	make_disk a && make_disk b
	# No key injection -- just verify no DiskIndex change from 0
	diskset_run 200 "" "$TMP/a.dsk" "$TMP/b.dsk"
	[[ $status -eq 0 ]]
	! grep -q 'comp DiskIndex' "$OUT/state.log"
}

@test "sequential next visits each disk in a three-disk set" {
	make_disk a && make_disk b && make_disk c
	diskset_run 800 \
		"100:ctrl-a,101:n,300:ctrl-a,301:n,500:ctrl-a,501:n" \
		"$TMP/a.dsk" "$TMP/b.dsk" "$TMP/c.dsk"
	[[ $status -eq 0 ]]
	grep -q 'comp DiskIndex .* -> 1' "$OUT/state.log"
	grep -q 'comp DiskIndex .* -> 2' "$OUT/state.log"
	grep -q 'comp DiskIndex .* -> 0' "$OUT/state.log"
}

# ---------------------------------------------------------------------------
# Single-Disk Behavior
# ---------------------------------------------------------------------------

@test "next with a single disk keeps index at 0" {
	make_disk a
	diskset_run 500 "100:ctrl-a,101:n" "$TMP/a.dsk"
	[[ $status -eq 0 ]]
	# Index goes from 0 -> 0 (wraps), so no state change is logged
	! grep -q 'comp DiskIndex' "$OUT/state.log"
}

@test "previous with a single disk keeps index at 0" {
	make_disk a
	diskset_run 500 "100:ctrl-a,101:p" "$TMP/a.dsk"
	[[ $status -eq 0 ]]
	! grep -q 'comp DiskIndex' "$OUT/state.log"
}
