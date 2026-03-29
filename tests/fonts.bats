setup_file() { load fonts_helper; setup_file; }
setup()      { load fonts_helper; setup; }
teardown()   { load fonts_helper; teardown; }

# ---------------------------------------------------------------------------
# Glyph bitmaps -- uppercase set (section 3.1)
# ---------------------------------------------------------------------------

@test "uppercase '@' glyph matches font data" {
	fill_screen_40 C0
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/uppercase.txt" "@")
	[[ "$actual" == "$expected" ]]
}

@test "uppercase 'A' glyph matches font data" {
	fill_screen_40 C1
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/uppercase.txt" "A")
	[[ "$actual" == "$expected" ]]
}

@test "uppercase '_' glyph matches font data" {
	fill_screen_40 DF
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/uppercase.txt" "_")
	[[ "$actual" == "$expected" ]]
}

# ---------------------------------------------------------------------------
# Glyph bitmaps -- special set (section 3.2)
# ---------------------------------------------------------------------------

@test "space glyph is entirely blank" {
	fill_screen_40 A0
	[[ $status -eq 0 ]]
	local actual
	actual=$(glyph_pattern_40)
	local expected
	expected=$(expected_glyph "$FONTDIR/special.txt" "SP")
	[[ "$actual" == "$expected" ]]
}

@test "special '!' glyph matches font data" {
	fill_screen_40 A1
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/special.txt" "!")
	[[ "$actual" == "$expected" ]]
}

# ---------------------------------------------------------------------------
# Glyph bitmaps -- lowercase set (section 3.3)
# ---------------------------------------------------------------------------

@test "lowercase 'a' glyph matches font data" {
	fill_screen_40 E1
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/lowercase.txt" "a")
	[[ "$actual" == "$expected" ]]
}

@test "lowercase 'g' glyph matches font data" {
	fill_screen_40 E7
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/lowercase.txt" "g")
	[[ "$actual" == "$expected" ]]
}

@test "DEL glyph is entirely blank" {
	fill_screen_40 FF
	[[ $status -eq 0 ]]
	local actual
	actual=$(glyph_pattern_40)
	local expected
	expected=$(expected_glyph "$FONTDIR/lowercase.txt" "DEL")
	[[ "$actual" == "$expected" ]]
}

# ---------------------------------------------------------------------------
# Border convention (section 2.4)
# ---------------------------------------------------------------------------

@test "'{' has column 0 on in row 3" {
	fill_screen_40 FB
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/lowercase.txt" "{")
	[[ "$actual" == "$expected" ]]
	local row3
	row3=$(echo "$actual" | sed -n '4p')
	[[ "${row3:0:1}" == "#" ]]
}

@test "'}' has column 6 on in row 3" {
	fill_screen_40 FD
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/lowercase.txt" "}")
	[[ "$actual" == "$expected" ]]
	local row3
	row3=$(echo "$actual" | sed -n '4p')
	[[ "${row3:6:1}" == "#" ]]
}

# ---------------------------------------------------------------------------
# 40-column scaling (section 5.1)
# ---------------------------------------------------------------------------

@test "40-col: each dot row is emitted twice (vertical doubling)" {
	fill_screen_40 C1
	[[ $status -eq 0 ]]
	local r
	for ((r = 0; r < 16; r += 2)); do
		local row_a row_b
		row_a=$(pixel_row "$r")
		row_b=$(pixel_row "$((r + 1))")
		[[ "$row_a" == "$row_b" ]]
	done
}

@test "40-col: each dot column is emitted twice (horizontal doubling)" {
	fill_screen_40 C1
	[[ $status -eq 0 ]]
	local row
	row=$(pixel_row 0)
	local c
	for ((c = 0; c < 14; c += 2)); do
		[[ "${row:$c:1}" == "${row:$((c + 1)):1}" ]]
	done
}

# ---------------------------------------------------------------------------
# Display modes (section 4.2)
# ---------------------------------------------------------------------------

@test "inverse \$21 has white background" {
	fill_screen_40 21
	[[ $status -eq 0 ]]
	[[ "$(bg_color)" == "FFFFFF" ]]
}

@test "normal \$C1 has black background" {
	fill_screen_40 C1
	[[ $status -eq 0 ]]
	[[ "$(bg_color)" == "000000" ]]
}

@test "inverse \$01 has white background" {
	fill_screen_40 01
	[[ $status -eq 0 ]]
	[[ "$(bg_color)" == "FFFFFF" ]]
}

@test "inverse \$01 glyph shape matches normal \$C1" {
	fill_screen_40 01
	[[ $status -eq 0 ]]
	local pattern_01
	pattern_01=$(glyph_pattern_40)

	reset_out
	fill_screen_40 C1
	[[ $status -eq 0 ]]
	local pattern_c1
	pattern_c1=$(glyph_pattern_40)

	[[ "$pattern_01" == "$pattern_c1" ]]
}

# ---------------------------------------------------------------------------
# Character space mapping (section 4.1)
# ---------------------------------------------------------------------------

@test "\$81 and \$C1 map to the same glyph (uppercase A)" {
	fill_screen_40 81
	[[ $status -eq 0 ]]
	local pattern_81
	pattern_81=$(glyph_pattern_40)

	reset_out
	fill_screen_40 C1
	[[ $status -eq 0 ]]
	local pattern_c1
	pattern_c1=$(glyph_pattern_40)

	[[ "$pattern_81" == "$pattern_c1" ]]
}

@test "\$21 and \$A1 map to the same glyph (special '!')" {
	fill_screen_40 21
	[[ $status -eq 0 ]]
	local pattern_21
	pattern_21=$(glyph_pattern_40)

	reset_out
	fill_screen_40 A1
	[[ $status -eq 0 ]]
	local pattern_a1
	pattern_a1=$(glyph_pattern_40)

	[[ "$pattern_21" == "$pattern_a1" ]]
}

# ---------------------------------------------------------------------------
# Glyph bitmaps -- MouseText set (section 3.4)
# ---------------------------------------------------------------------------

@test "mousetext \$40 glyph matches font data (alt charset)" {
	fill_screen_40_alt 40
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/mousetext.txt" "0x40")
	[[ "$actual" == "$expected" ]]
}

@test "mousetext \$48 glyph matches font data (alt charset)" {
	fill_screen_40_alt 48
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/mousetext.txt" "0x48")
	[[ "$actual" == "$expected" ]]
}

@test "mousetext \$53 glyph matches font data (alt charset)" {
	fill_screen_40_alt 53
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/mousetext.txt" "0x53")
	[[ "$actual" == "$expected" ]]
}

@test "mousetext \$5A glyph matches font data (alt charset)" {
	fill_screen_40_alt 5A
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/mousetext.txt" "0x5A")
	[[ "$actual" == "$expected" ]]
}

@test "mousetext \$5F glyph matches font data (alt charset)" {
	fill_screen_40_alt 5F
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/mousetext.txt" "0x5F")
	[[ "$actual" == "$expected" ]]
}

@test "mousetext \$53 has black background in alt charset" {
	fill_screen_40_alt 53
	[[ $status -eq 0 ]]
	[[ "$(bg_color)" == "000000" ]]
}

@test "\$48 in primary charset is flash uppercase, not mousetext" {
	fill_screen_40 48
	[[ $status -eq 0 ]]
	local primary
	primary=$(glyph_pattern_40)

	reset_out
	fill_screen_40_alt 48
	[[ $status -eq 0 ]]
	local alt
	alt=$(glyph_pattern_40)

	[[ "$primary" != "$alt" ]]
}

@test "\$61 matches special '!' glyph from font data" {
	fill_screen_40 61
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/special.txt" "!")
	[[ "$actual" == "$expected" ]]
}

@test "flash \$61 renders as inverse (same as normal \$A1)" {
	fill_screen_40 61
	[[ $status -eq 0 ]]
	local pattern_61
	pattern_61=$(glyph_pattern_40)

	reset_out
	fill_screen_40 A1
	[[ $status -eq 0 ]]
	local pattern_a1
	pattern_a1=$(glyph_pattern_40)

	[[ "$pattern_61" == "$pattern_a1" ]]
}

# ---------------------------------------------------------------------------
# Character space mapping -- additional range coverage (section 4.1)
# ---------------------------------------------------------------------------

@test "\$00 maps to inverse uppercase '@'" {
	fill_screen_40 00
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/uppercase.txt" "@")
	[[ "$actual" == "$expected" ]]
}

@test "\$00 has white background (inverse)" {
	fill_screen_40 00
	[[ $status -eq 0 ]]
	[[ "$(bg_color)" == "FFFFFF" ]]
}

@test "\$80 and \$C0 both map to normal uppercase '@'" {
	fill_screen_40 80
	[[ $status -eq 0 ]]
	local pattern_80
	pattern_80=$(glyph_pattern_40)

	reset_out
	fill_screen_40 C0
	[[ $status -eq 0 ]]
	local pattern_c0
	pattern_c0=$(glyph_pattern_40)

	[[ "$pattern_80" == "$pattern_c0" ]]
}

@test "\$48 in primary charset renders as inverse uppercase 'H'" {
	fill_screen_40 48
	[[ $status -eq 0 ]]
	local actual expected
	actual=$(glyph_pattern_40)
	expected=$(expected_glyph "$FONTDIR/uppercase.txt" "H")
	[[ "$actual" == "$expected" ]]
	[[ "$(bg_color)" == "FFFFFF" ]]
}
