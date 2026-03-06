package record

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// AudioAssertion represents an assertion about audio over a step range.
type AudioAssertion struct {
	StartStep int
	EndStep   int
	Channel   string
	Checks    []AudioCheck
}

// AudioCheck is a single property check within an assertion.
type AudioCheck struct {
	Property  string // "freq", "amplitude", "duty", "silent"
	Op        string // "~" (approximate) or ">" (greater than)
	Value     float64
	Tolerance float64
}

// AudioCheckFailure describes why a single check failed.
type AudioCheckFailure struct {
	Property string
	Expected string
	Actual   float64
}

// AudioAssertionResult holds the outcome of evaluating one AudioAssertion.
type AudioAssertionResult struct {
	Assertion AudioAssertion
	Passed    bool
	Failures  []AudioCheckFailure
}

// ParseAudioAssertion parses lines like:
// - step 100-500: audio freq ~1000Hz +/- 50, amplitude > 0.3
// - step 501-800: audio silent
// - step 100-500: audio speaker freq ~1000Hz +/- 50
func ParseAudioAssertion(line string) (AudioAssertion, error) {
	var a AudioAssertion

	line = strings.TrimSpace(line)

	// Parse "step N-M:"
	if !strings.HasPrefix(line, "step ") {
		return a, fmt.Errorf("expected 'step' prefix")
	}

	line = line[5:]

	colonIdx := strings.Index(line, ":")
	if colonIdx < 0 {
		return a, fmt.Errorf("expected ':' after step range")
	}

	rangePart := strings.TrimSpace(line[:colonIdx])
	rest := strings.TrimSpace(line[colonIdx+1:])

	// Parse range
	parts := strings.SplitN(rangePart, "-", 2)
	if len(parts) != 2 {
		return a, fmt.Errorf("expected step range N-M, got %q", rangePart)
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return a, fmt.Errorf("invalid start step: %w", err)
	}

	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return a, fmt.Errorf("invalid end step: %w", err)
	}

	a.StartStep = start
	a.EndStep = end

	// Expect "audio" keyword
	if !strings.HasPrefix(rest, "audio") {
		return a, fmt.Errorf("expected 'audio' keyword")
	}

	rest = strings.TrimSpace(rest[5:])

	// Check if next token is a channel label (not a property keyword)
	if rest != "" && !isPropertyKeyword(firstToken(rest)) {
		token := firstToken(rest)
		a.Channel = token
		rest = strings.TrimSpace(rest[len(token):])
	}

	// Parse checks
	if rest == "" {
		return a, nil
	}

	checks, err := parseChecks(rest)
	if err != nil {
		return a, err
	}

	a.Checks = checks
	return a, nil
}

var propertyKeywords = map[string]bool{
	"freq":      true,
	"amplitude": true,
	"duty":      true,
	"silent":    true,
}

func isPropertyKeyword(token string) bool {
	return propertyKeywords[token]
}

func firstToken(s string) string {
	s = strings.TrimSpace(s)
	for i, ch := range s {
		if ch == ' ' || ch == ',' {
			return s[:i]
		}
	}
	return s
}

func parseChecks(s string) ([]AudioCheck, error) {
	var checks []AudioCheck

	for part := range strings.SplitSeq(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		check, err := parseSingleCheck(part)
		if err != nil {
			return nil, err
		}

		checks = append(checks, check)
	}

	return checks, nil
}

func parseSingleCheck(s string) (AudioCheck, error) {
	s = strings.TrimSpace(s)

	// Handle "silent" as a bare keyword
	if s == "silent" {
		return AudioCheck{Property: "silent"}, nil
	}

	tokens := strings.Fields(s)
	if len(tokens) < 2 {
		return AudioCheck{}, fmt.Errorf("incomplete check: %q", s)
	}

	prop := tokens[0]
	if !isPropertyKeyword(prop) {
		return AudioCheck{}, fmt.Errorf("unknown property: %q", prop)
	}

	var check AudioCheck
	check.Property = prop

	valueStr := tokens[1]

	// Determine operator from value prefix
	tolStart := 2
	if strings.HasPrefix(valueStr, "~") {
		check.Op = "~"
		valueStr = valueStr[1:]
	} else if valueStr == ">" {
		check.Op = ">"
		if len(tokens) < 3 {
			return AudioCheck{}, fmt.Errorf("missing value after '>'")
		}
		valueStr = tokens[2]
		tolStart = 3
	} else {
		return AudioCheck{}, fmt.Errorf("expected operator '~' or '>' in %q", s)
	}

	// Strip unit suffixes
	valueStr = strings.TrimRight(valueStr, "Hz%")

	val, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return AudioCheck{}, fmt.Errorf("invalid value %q: %w", valueStr, err)
	}

	check.Value = val

	// Look for "+/- N" tolerance
	for i := tolStart; i < len(tokens); i++ {
		if tokens[i] == "+/-" && i+1 < len(tokens) {
			tolStr := strings.TrimRight(tokens[i+1], "Hz%")
			tol, err := strconv.ParseFloat(tolStr, 64)
			if err != nil {
				return AudioCheck{}, fmt.Errorf("invalid tolerance %q: %w", tokens[i+1], err)
			}
			check.Tolerance = tol
			break
		}
	}

	return check, nil
}

// EvaluateAudioAssertions evaluates a set of assertions against the given
// recorders. If an assertion specifies a channel, the recorder is looked up
// by that label; otherwise the first recorder is used.
func EvaluateAudioAssertions(
	assertions []AudioAssertion,
	recorders map[string]*AudioRecorder,
) []AudioAssertionResult {
	results := make([]AudioAssertionResult, len(assertions))

	for i, a := range assertions {
		result := AudioAssertionResult{Assertion: a, Passed: true}

		rec := findRecorder(a.Channel, recorders)
		if rec == nil {
			result.Passed = false
			result.Failures = append(result.Failures, AudioCheckFailure{
				Property: "channel",
				Expected: a.Channel,
			})
			results[i] = result
			continue
		}

		seg := rec.Segment(a.StartStep, a.EndStep)
		fp := seg.Fingerprint()

		for _, check := range a.Checks {
			if failure, ok := evaluateCheck(check, fp); !ok {
				result.Passed = false
				result.Failures = append(result.Failures, failure)
			}
		}

		results[i] = result
	}

	return results
}

func findRecorder(channel string, recorders map[string]*AudioRecorder) *AudioRecorder {
	if channel != "" {
		return recorders[channel]
	}

	// Return the first (or only) recorder
	for _, r := range recorders {
		return r
	}

	return nil
}

func evaluateCheck(check AudioCheck, fp AudioFingerprint) (AudioCheckFailure, bool) {
	switch check.Property {
	case "silent":
		if !fp.Silent {
			return AudioCheckFailure{
				Property: "silent",
				Expected: "true",
				Actual:   fp.AmplitudeMean,
			}, false
		}
		return AudioCheckFailure{}, true

	case "freq":
		return evaluateOp(check, fp.Frequency)

	case "amplitude":
		return evaluateOp(check, fp.AmplitudeMean)

	case "duty":
		return evaluateOp(check, fp.DutyCycle)

	default:
		return AudioCheckFailure{
			Property: check.Property,
			Expected: "known property",
		}, false
	}
}

func evaluateOp(check AudioCheck, actual float64) (AudioCheckFailure, bool) {
	switch check.Op {
	case "~":
		tol := check.Tolerance
		if tol == 0 {
			tol = check.Value * 0.1 // default 10% tolerance
		}

		if math.Abs(actual-check.Value) > tol {
			return AudioCheckFailure{
				Property: check.Property,
				Expected: fmt.Sprintf("~%.2f +/- %.2f", check.Value, tol),
				Actual:   actual,
			}, false
		}

	case ">":
		if actual <= check.Value {
			return AudioCheckFailure{
				Property: check.Property,
				Expected: fmt.Sprintf("> %.2f", check.Value),
				Actual:   actual,
			}, false
		}

	default:
		return AudioCheckFailure{
			Property: check.Property,
			Expected: "known operator",
		}, false
	}

	return AudioCheckFailure{}, true
}
