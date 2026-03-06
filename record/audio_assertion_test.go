package record_test

import (
	"bytes"
	"testing"

	"github.com/pevans/erc/record"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAudioAssertion_freq(t *testing.T) {
	a, err := record.ParseAudioAssertion("step 100-500: audio freq ~1000Hz +/- 50")
	require.NoError(t, err)

	assert.Equal(t, 100, a.StartStep)
	assert.Equal(t, 500, a.EndStep)
	assert.Equal(t, "", a.Channel)
	require.Len(t, a.Checks, 1)
	assert.Equal(t, "freq", a.Checks[0].Property)
	assert.Equal(t, "~", a.Checks[0].Op)
	assert.InDelta(t, 1000, a.Checks[0].Value, 0.01)
	assert.InDelta(t, 50, a.Checks[0].Tolerance, 0.01)
}

func TestParseAudioAssertion_silent(t *testing.T) {
	a, err := record.ParseAudioAssertion("step 501-800: audio silent")
	require.NoError(t, err)

	assert.Equal(t, 501, a.StartStep)
	assert.Equal(t, 800, a.EndStep)
	require.Len(t, a.Checks, 1)
	assert.Equal(t, "silent", a.Checks[0].Property)
}

func TestParseAudioAssertion_withChannel(t *testing.T) {
	a, err := record.ParseAudioAssertion("step 100-500: audio speaker freq ~1000Hz +/- 50")
	require.NoError(t, err)

	assert.Equal(t, "speaker", a.Channel)
	require.Len(t, a.Checks, 1)
	assert.Equal(t, "freq", a.Checks[0].Property)
}

func TestParseAudioAssertion_multipleChecks(t *testing.T) {
	a, err := record.ParseAudioAssertion("step 100-500: audio freq ~1000Hz +/- 50, amplitude > 0.3")
	require.NoError(t, err)

	require.Len(t, a.Checks, 2)
	assert.Equal(t, "freq", a.Checks[0].Property)
	assert.Equal(t, "amplitude", a.Checks[1].Property)
	assert.Equal(t, ">", a.Checks[1].Op)
	assert.InDelta(t, 0.3, a.Checks[1].Value, 0.001)
}

func TestParseAudioAssertion_dutyWithTolerance(t *testing.T) {
	a, err := record.ParseAudioAssertion("step 801-1200: audio freq ~440Hz, duty ~50% +/- 5")
	require.NoError(t, err)

	require.Len(t, a.Checks, 2)
	assert.Equal(t, "duty", a.Checks[1].Property)
	assert.InDelta(t, 50, a.Checks[1].Value, 0.01)
	assert.InDelta(t, 5, a.Checks[1].Tolerance, 0.01)
}

func TestParseAudioAssertion_errors(t *testing.T) {
	cases := []string{
		"",
		"not a step",
		"step 100: audio freq ~1000Hz",    // missing end step
		"step 100-500 audio freq ~1000Hz", // missing colon
	}

	for _, c := range cases {
		_, err := record.ParseAudioAssertion(c)
		assert.Error(t, err, "expected error for: %q", c)
	}
}

func TestEvaluateAudioAssertions_pass(t *testing.T) {
	var cycle uint64

	mono := squareWave(1000, 44100, 44100, 0.5)
	reader := bytes.NewReader(stereoFloat32Bytes(mono))
	ar := record.NewAudioRecorder("speaker", reader, func() uint64 { return cycle }, 1023000, 44100)

	// totalCycles * 44100 / 1023000 = totalSamples We want 44100 samples, so
	// totalCycles = 1023000, at 23 cycles/step ≈ 44478 steps
	cyclesPerStep := uint64(23)
	stepsNeeded := 1023000/int(cyclesPerStep) + 1

	for i := range stepsNeeded {
		ar.Before()
		cycle += cyclesPerStep
		ar.Observe(i + 1)
	}

	assertions := []record.AudioAssertion{
		{
			StartStep: 1,
			EndStep:   stepsNeeded,
			Checks: []record.AudioCheck{
				{Property: "freq", Op: "~", Value: 1000, Tolerance: 50},
				{Property: "amplitude", Op: ">", Value: 0.1},
			},
		},
	}

	recorders := map[string]*record.AudioRecorder{"speaker": ar}
	results := record.EvaluateAudioAssertions(assertions, recorders)

	require.Len(t, results, 1)
	assert.True(t, results[0].Passed, "expected pass, got failures: %+v", results[0].Failures)
}

func TestEvaluateAudioAssertions_silentPass(t *testing.T) {
	var cycle uint64

	mono := make([]float32, 1000)
	reader := bytes.NewReader(stereoFloat32Bytes(mono))
	ar := record.NewAudioRecorder("speaker", reader, func() uint64 { return cycle }, 1000, 100)

	for i := range 1000 {
		ar.Before()
		cycle += 10
		ar.Observe(i + 1)
	}

	assertions := []record.AudioAssertion{
		{
			StartStep: 1,
			EndStep:   1000,
			Checks:    []record.AudioCheck{{Property: "silent"}},
		},
	}

	recorders := map[string]*record.AudioRecorder{"speaker": ar}
	results := record.EvaluateAudioAssertions(assertions, recorders)

	require.Len(t, results, 1)
	assert.True(t, results[0].Passed)
}

func TestEvaluateAudioAssertions_fail(t *testing.T) {
	var cycle uint64

	// 440Hz wave, but assert 1000Hz -- should fail
	mono := squareWave(440, 44100, 44100, 0.5)
	reader := bytes.NewReader(stereoFloat32Bytes(mono))
	ar := record.NewAudioRecorder("speaker", reader, func() uint64 { return cycle }, 1023000, 44100)

	cyclesPerStep := uint64(23)
	stepsNeeded := 1023000/int(cyclesPerStep) + 1

	for i := range stepsNeeded {
		ar.Before()
		cycle += cyclesPerStep
		ar.Observe(i + 1)
	}

	assertions := []record.AudioAssertion{
		{
			StartStep: 1,
			EndStep:   stepsNeeded,
			Checks: []record.AudioCheck{
				{Property: "freq", Op: "~", Value: 1000, Tolerance: 50},
			},
		},
	}

	recorders := map[string]*record.AudioRecorder{"speaker": ar}
	results := record.EvaluateAudioAssertions(assertions, recorders)

	require.Len(t, results, 1)
	assert.False(t, results[0].Passed)
	assert.Len(t, results[0].Failures, 1)
	assert.Equal(t, "freq", results[0].Failures[0].Property)
}

func TestEvaluateAudioAssertions_channelLookup(t *testing.T) {
	assertions := []record.AudioAssertion{
		{
			StartStep: 1,
			EndStep:   100,
			Channel:   "nonexistent",
			Checks:    []record.AudioCheck{{Property: "silent"}},
		},
	}

	recorders := map[string]*record.AudioRecorder{}
	results := record.EvaluateAudioAssertions(assertions, recorders)

	require.Len(t, results, 1)
	assert.False(t, results[0].Passed)
	assert.Equal(t, "channel", results[0].Failures[0].Property)
}
