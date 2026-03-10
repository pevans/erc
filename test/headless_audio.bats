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
