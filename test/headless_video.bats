setup_file() { load test_helper; setup_file; }
setup()      { load test_helper; setup; }
teardown()   { load test_helper; teardown; }

@test "capture-video produces video.frame" {
	run_headless --steps 50000 --capture-video 49999 "$DISK"
	[[ $status -eq 0 ]]
	[[ -f "$OUT/video.frame" ]]
}

@test "video.frame has step header and colors legend" {
	run_headless --steps 50000 --capture-video 49999 "$DISK"
	[[ $status -eq 0 ]]
	grep -q '^step 49999: video screen ' "$OUT/video.frame"
	grep -q '^colors: ' "$OUT/video.frame"
}

@test "multiple capture steps appear in video.frame" {
	run_headless --steps 50000 --capture-video 10000,49999 "$DISK"
	[[ $status -eq 0 ]]
	grep -q '^step 10000: video screen ' "$OUT/video.frame"
	grep -q '^step 49999: video screen ' "$OUT/video.frame"
}

@test "no capture-video produces no video.frame" {
	run_headless --steps 1000 "$DISK"
	[[ $status -eq 0 ]]
	[[ ! -f "$OUT/video.frame" ]]
}
