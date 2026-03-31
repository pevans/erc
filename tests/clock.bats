setup_file() { load test_helper; setup_file; }
setup()      { load test_helper; setup; }
teardown()   { load test_helper; teardown; }

@test "speed increase changes the Speed state" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:+" \
		--watch-comp Speed \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp Speed .* -> 2' "$OUT/state.log"
}

@test "speed decrease changes the Speed state" {
	run_headless --steps 1000 \
		--keys "100:ctrl-a,101:+,200:ctrl-a,201:-" \
		--watch-comp Speed \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp Speed .* -> 2' "$OUT/state.log"
	grep -q 'comp Speed .* -> 1' "$OUT/state.log"
}

@test "disk motor on triggers full-speed" {
	run_headless --steps 100000 \
		--watch-comp FullSpeed \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp FullSpeed .* -> true' "$OUT/state.log"
}
