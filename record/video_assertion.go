package record

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// VideoAssertionKind identifies the type of video assertion.
type VideoAssertionKind string

const (
	VideoScreen VideoAssertionKind = "screen"
	VideoRegion VideoAssertionKind = "region"
	VideoRow    VideoAssertionKind = "row"
)

// ColorLegend maps single-character codes to RGB colors.
type ColorLegend struct {
	Forward map[byte]color.RGBA
	Reverse map[color.RGBA]byte
}

// VideoAssertion represents an assertion about video content at a specific
// execution step.
type VideoAssertion struct {
	Step     int
	Kind     VideoAssertionKind
	GridW    int
	GridH    int
	RegionX  int
	RegionY  int
	RegionW  int
	RegionH  int
	RowIndex int
	Legend   ColorLegend
	Expected []string
}

// VideoAssertionResult holds the outcome of evaluating one VideoAssertion.
type VideoAssertionResult struct {
	Assertion VideoAssertion
	Passed    bool
	Failure   *VideoMismatch
}

// VideoMismatch describes the first mismatched cell in a video assertion.
type VideoMismatch struct {
	X, Y          int
	Expected      byte
	Actual        byte
	ActualRGB     color.RGBA
	ExpectedRow   string
	ActualRow     string
	MissingFrame  bool
	UnmappedColor bool
}

// ParseVideoAssertion parses a complete video assertion block from a slice of
// lines. Line 0 is the header, line 1 is the color legend, and the remaining
// lines are the expected grid rows.
func ParseVideoAssertion(lines []string) (VideoAssertion, error) {
	if len(lines) < 3 {
		return VideoAssertion{}, fmt.Errorf("video assertion requires at least 3 lines, got %d", len(lines))
	}

	a, err := parseVideoHeader(lines[0])
	if err != nil {
		return VideoAssertion{}, fmt.Errorf("header: %w", err)
	}

	legend, err := parseColorLegend(lines[1])
	if err != nil {
		return VideoAssertion{}, fmt.Errorf("colors: %w", err)
	}

	a.Legend = legend
	a.Expected = lines[2:]

	// Validate grid row count
	expectedRows := a.gridHeight()
	if len(a.Expected) != expectedRows {
		return VideoAssertion{}, fmt.Errorf(
			"expected %d grid rows, got %d", expectedRows, len(a.Expected))
	}

	// Validate grid row widths
	expectedWidth := a.gridWidth()
	for i, row := range a.Expected {
		if len(row) != expectedWidth {
			return VideoAssertion{}, fmt.Errorf(
				"row %d: expected width %d, got %d", i, expectedWidth, len(row))
		}
	}

	return a, nil
}

func (a *VideoAssertion) gridHeight() int {
	switch a.Kind {
	case VideoRegion:
		return a.RegionH
	case VideoRow:
		return 1
	default:
		return a.GridH
	}
}

func (a *VideoAssertion) gridWidth() int {
	switch a.Kind {
	case VideoRegion:
		return a.RegionW
	default:
		return a.GridW
	}
}

func parseVideoHeader(line string) (VideoAssertion, error) {
	var a VideoAssertion

	line = strings.TrimSpace(line)

	if !strings.HasPrefix(line, "step ") {
		return a, fmt.Errorf("expected 'step' prefix")
	}

	line = line[5:]

	colonIdx := strings.Index(line, ":")
	if colonIdx < 0 {
		return a, fmt.Errorf("expected ':' after step number")
	}

	step, err := strconv.Atoi(strings.TrimSpace(line[:colonIdx]))
	if err != nil {
		return a, fmt.Errorf("invalid step number: %w", err)
	}

	a.Step = step

	rest := strings.TrimSpace(line[colonIdx+1:])

	if !strings.HasPrefix(rest, "video ") {
		return a, fmt.Errorf("expected 'video' keyword")
	}

	rest = strings.TrimSpace(rest[6:])
	tokens := strings.Fields(rest)

	if len(tokens) == 0 {
		return a, fmt.Errorf("expected assertion kind (screen, region, row)")
	}

	switch tokens[0] {
	case "screen":
		return parseScreenHeader(a, tokens[1:])
	case "region":
		return parseRegionHeader(a, tokens[1:])
	case "row":
		return parseRowHeader(a, tokens[1:])
	default:
		return a, fmt.Errorf("unknown video assertion kind: %q", tokens[0])
	}
}

// parseScreenHeader parses: 280x192
func parseScreenHeader(a VideoAssertion, tokens []string) (VideoAssertion, error) {
	if len(tokens) != 1 {
		return a, fmt.Errorf("screen assertion expects 1 argument (grid size), got %d", len(tokens))
	}

	w, h, err := parseDimensions(tokens[0])
	if err != nil {
		return a, fmt.Errorf("grid size: %w", err)
	}

	a.Kind = VideoScreen
	a.GridW = w
	a.GridH = h

	return a, nil
}

// parseRegionHeader parses: 10,5 20x10 280x192
func parseRegionHeader(a VideoAssertion, tokens []string) (VideoAssertion, error) {
	if len(tokens) != 3 {
		return a, fmt.Errorf("region assertion expects 3 arguments, got %d", len(tokens))
	}

	// Origin: x,y
	parts := strings.SplitN(tokens[0], ",", 2)
	if len(parts) != 2 {
		return a, fmt.Errorf("expected region origin x,y, got %q", tokens[0])
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return a, fmt.Errorf("region x: %w", err)
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return a, fmt.Errorf("region y: %w", err)
	}

	// Region size: wxh
	rw, rh, err := parseDimensions(tokens[1])
	if err != nil {
		return a, fmt.Errorf("region size: %w", err)
	}

	// Grid size: wxh
	gw, gh, err := parseDimensions(tokens[2])
	if err != nil {
		return a, fmt.Errorf("grid size: %w", err)
	}

	a.Kind = VideoRegion
	a.RegionX = x
	a.RegionY = y
	a.RegionW = rw
	a.RegionH = rh
	a.GridW = gw
	a.GridH = gh

	return a, nil
}

// parseRowHeader parses: 96 280x192
func parseRowHeader(a VideoAssertion, tokens []string) (VideoAssertion, error) {
	if len(tokens) != 2 {
		return a, fmt.Errorf("row assertion expects 2 arguments, got %d", len(tokens))
	}

	row, err := strconv.Atoi(tokens[0])
	if err != nil {
		return a, fmt.Errorf("row index: %w", err)
	}

	gw, gh, err := parseDimensions(tokens[1])
	if err != nil {
		return a, fmt.Errorf("grid size: %w", err)
	}

	a.Kind = VideoRow
	a.RowIndex = row
	a.GridW = gw
	a.GridH = gh

	return a, nil
}

func parseDimensions(s string) (int, int, error) {
	parts := strings.SplitN(s, "x", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected WxH, got %q", s)
	}

	w, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("width: %w", err)
	}

	h, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("height: %w", err)
	}

	return w, h, nil
}

func parseColorLegend(line string) (ColorLegend, error) {
	line = strings.TrimSpace(line)

	if !strings.HasPrefix(line, "colors:") {
		return ColorLegend{}, fmt.Errorf("expected 'colors:' prefix")
	}

	line = strings.TrimSpace(line[7:])

	legend := ColorLegend{
		Forward: make(map[byte]color.RGBA),
		Reverse: make(map[color.RGBA]byte),
	}

	for entry := range strings.SplitSeq(line, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			return ColorLegend{}, fmt.Errorf("expected 'char = hex', got %q", entry)
		}

		charPart := strings.TrimSpace(parts[0])
		hexPart := strings.TrimSpace(parts[1])

		if len(charPart) != 1 {
			return ColorLegend{}, fmt.Errorf("legend key must be a single character, got %q", charPart)
		}

		clr, err := parseHexColor(hexPart)
		if err != nil {
			return ColorLegend{}, fmt.Errorf("color for %q: %w", charPart, err)
		}

		ch := charPart[0]
		legend.Forward[ch] = clr
		legend.Reverse[clr] = ch
	}

	return legend, nil
}

func parseHexColor(s string) (color.RGBA, error) {
	if len(s) != 6 {
		return color.RGBA{}, fmt.Errorf("expected 6-digit hex color, got %q", s)
	}

	val, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid hex color %q: %w", s, err)
	}

	return color.RGBA{
		R: uint8(val >> 16),
		G: uint8(val >> 8),
		B: uint8(val),
		A: 0xff,
	}, nil
}
