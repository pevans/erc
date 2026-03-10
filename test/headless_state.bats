setup_file() { load test_helper; setup_file; }
setup()      { load test_helper; setup; }
teardown()   { load test_helper; teardown; }

@test "watch-mem produces state.log with mem entries" {
	run_headless --steps 50000 --watch-mem 0400-0403 "$DISK"
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q '^step [0-9]*: mem ' "$OUT/state.log"
}

@test "watch-reg produces state.log with reg entries" {
	run_headless --steps 50000 --watch-reg PC,A "$DISK"
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q '^step [0-9]*: reg A ' "$OUT/state.log"
	grep -q '^step [0-9]*: reg PC ' "$OUT/state.log"
}

@test "watch-comp produces state.log with comp entries" {
	run_headless --steps 50000 --watch-comp SpeakerState "$DISK"
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q '^step [0-9]*: comp SpeakerState ' "$OUT/state.log"
}

@test "combined watch flags produce all tag types" {
	run_headless --steps 50000 \
		--watch-mem 0400 \
		--watch-reg A \
		--watch-comp SpeakerState \
		"$DISK"
	[[ $status -eq 0 ]]
	[[ -f "$OUT/state.log" ]]
	grep -q '^step [0-9]*: mem ' "$OUT/state.log"
	grep -q '^step [0-9]*: reg ' "$OUT/state.log"
	grep -q '^step [0-9]*: comp ' "$OUT/state.log"
}

@test "no observers produces no state.log" {
	run_headless --steps 1000 "$DISK"
	[[ $status -eq 0 ]]
	[[ ! -f "$OUT/state.log" ]]
}

@test "state entries use old -> new format" {
	run_headless --steps 50000 --watch-reg A "$DISK"
	[[ $status -eq 0 ]]
	grep -q ' -> ' "$OUT/state.log"
}
