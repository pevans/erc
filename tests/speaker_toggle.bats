setup_file() { load speaker_helper; setup_file; }
setup()      { load speaker_helper; setup; }
teardown()   { load speaker_helper; teardown; }

# --- $C030 Speaker Toggle ---

@test "reading C030 toggles speaker state to true" {
	spk_run \
		'LDA $C030' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'comp SpeakerState .* -> true' "$OUT/state.log"
}

@test "reading C030 twice toggles speaker state back to false" {
	spk_run \
		'LDA $C030' \
		'LDA $C030' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'comp SpeakerState .* -> true' "$OUT/state.log"
	grep -q 'comp SpeakerState .* -> false' "$OUT/state.log"
}

@test "writing C030 toggles speaker state" {
	spk_run \
		'STA $C030' \
		'.halt'
	[[ $status -eq 0 ]]
	grep -q 'comp SpeakerState .* -> true' "$OUT/state.log"
}

@test "multiple toggles alternate speaker state" {
	spk_run \
		'LDA $C030' \
		'LDA $C030' \
		'LDA $C030' \
		'.halt'
	[[ $status -eq 0 ]]
	local count
	count=$(grep -c 'comp SpeakerState' "$OUT/state.log")
	[[ $count -eq 3 ]]
	[[ "$(_last_comp SpeakerState)" == "true" ]]
}

# --- Audio Recording ---

@test "toggling C030 produces non-silent audio" {
	spk_run_audio \
		'loop: LDA $C030' \
		'NOP' \
		'NOP' \
		'NOP' \
		'NOP' \
		'JMP loop'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/audio.pcm" ]]
	[[ -s "$OUT/audio.pcm" ]]
	# PCM should contain non-zero bytes (actual audio, not silence)
	local nonzero
	nonzero=$(LC_ALL=C tr -d '\0' < "$OUT/audio.pcm" | wc -c)
	[[ $nonzero -gt 0 ]]
}

@test "speaker state unchanged when C030 is never accessed" {
	spk_run \
		'NOP' \
		'NOP' \
		'NOP' \
		'.halt'
	[[ $status -eq 0 ]]
	# No SpeakerState transitions should appear (state.log may not exist)
	[[ ! -f "$OUT/state.log" ]] || ! grep -q 'comp SpeakerState' "$OUT/state.log"
}

@test "audio.pcm size is a multiple of 4 bytes" {
	spk_run_audio \
		'loop: LDA $C030' \
		'JMP loop'
	[[ $status -eq 0 ]]
	[[ -f "$OUT/audio.pcm" ]]
	local size
	size=$(wc -c < "$OUT/audio.pcm")
	[[ $((size % 4)) -eq 0 ]]
}
