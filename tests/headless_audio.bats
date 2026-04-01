setup_file() { load test_helper; setup_file; }
setup()      { load test_helper; setup; }
teardown()   { load test_helper; teardown; }

@test "record-audio produces audio.pcm" {
	run_headless --steps 100000 --record-audio "$DISK"
	[[ $status -eq 0 ]]
	[[ -f "$OUT/audio.pcm" ]]
	[[ -s "$OUT/audio.pcm" ]]
}

@test "audio.pcm size is a multiple of 4 bytes" {
	run_headless --steps 100000 --record-audio "$DISK"
	[[ $status -eq 0 ]]
	size=$(wc -c < "$OUT/audio.pcm")
	[[ $((size % 4)) -eq 0 ]]
}

@test "no record-audio produces no audio.pcm" {
	run_headless --steps 1000 "$DISK"
	[[ $status -eq 0 ]]
	[[ ! -f "$OUT/audio.pcm" ]]
}

@test "volume level clamps at maximum of 100" {
	# Default is 50; six increases of 10 would reach 110 without clamping
	run_headless --steps 2000 \
		--keys "100:ctrl-a,101:],200:ctrl-a,201:],300:ctrl-a,301:],400:ctrl-a,401:],500:ctrl-a,501:],600:ctrl-a,601:]" \
		--watch-comp VolumeLevel \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp VolumeLevel .* -> 100' "$OUT/state.log"
	! grep -q 'comp VolumeLevel .* -> 110' "$OUT/state.log"
}

@test "volume level clamps at minimum of 0" {
	# Default is 50; six decreases of 10 would reach -10 without clamping.
	# VolumeDown to 0 sets VolumeMuted=true while preserving VolumeLevel
	# at the last non-zero value, so we check VolumeMuted instead.
	run_headless --steps 2000 \
		--keys "100:ctrl-a,101:[,200:ctrl-a,201:[,300:ctrl-a,301:[,400:ctrl-a,401:[,500:ctrl-a,501:[,600:ctrl-a,601:[" \
		--watch-comp VolumeLevel,VolumeMuted \
		"$DISK"
	[[ $status -eq 0 ]]
	grep -q 'comp VolumeMuted .* -> true' "$OUT/state.log"
}
