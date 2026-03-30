# text80_helper.bash -- shared helpers for 80-column text mode tests (spec 19).
# Inherits setup, teardown, asm_video, asm_video_mono, asm_video_long, and
# all pixel-inspection helpers from text40_helper.bash.
source "${BASH_SOURCE[0]%/*}/text40_helper.bash"

# any_col_range_has_mixed_colors START_ROW END_ROW COL_START COL_LEN --
# succeed if any pixel row in [START_ROW, END_ROW] has mixed colors within
# columns [COL_START, COL_START+COL_LEN-1].
any_col_range_has_mixed_colors() {
	local r
	for ((r = $1; r <= $2; r++)); do
		if col_range_has_mixed_colors "$r" "$3" "$4"; then
			return 0
		fi
	done
	return 1
}
