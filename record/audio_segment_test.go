package record_test

import (
	"math"
	"testing"

	"github.com/pevans/erc/record"
	"github.com/stretchr/testify/assert"
)

// squareWave generates a square wave at the given frequency and sample rate.
func squareWave(freq float64, sampleRate int, numSamples int, amplitude float32) []float32 {
	samples := make([]float32, numSamples)
	samplesPerCycle := float64(sampleRate) / freq

	for i := range samples {
		pos := math.Mod(float64(i), samplesPerCycle)
		if pos < samplesPerCycle/2 {
			samples[i] = amplitude
		} else {
			samples[i] = -amplitude
		}
	}

	return samples
}

func TestFingerprint_squareWave1000Hz(t *testing.T) {
	seg := record.AudioSegment{
		Samples:    squareWave(1000, 44100, 44100, 0.5),
		SampleRate: 44100,
		StartStep:  1,
		EndStep:    100,
	}

	fp := seg.Fingerprint()

	assert.InDelta(t, 1000, fp.Frequency, 25)
	assert.InDelta(t, 0.5, fp.AmplitudeMax, 0.01)
	assert.InDelta(t, 0.5, fp.AmplitudeMean, 0.05)
	assert.InDelta(t, 50, fp.DutyCycle, 2)
	assert.False(t, fp.Silent)
	assert.Greater(t, fp.ToggleCount, 1900)
}

func TestFingerprint_squareWave440Hz(t *testing.T) {
	seg := record.AudioSegment{
		Samples:    squareWave(440, 44100, 44100, 0.8),
		SampleRate: 44100,
		StartStep:  1,
		EndStep:    100,
	}

	fp := seg.Fingerprint()

	assert.InDelta(t, 440, fp.Frequency, 10)
	assert.InDelta(t, 0.8, fp.AmplitudeMax, 0.01)
	assert.False(t, fp.Silent)
}

func TestFingerprint_silence(t *testing.T) {
	samples := make([]float32, 44100)
	seg := record.AudioSegment{
		Samples:    samples,
		SampleRate: 44100,
		StartStep:  1,
		EndStep:    100,
	}

	fp := seg.Fingerprint()

	assert.True(t, fp.Silent)
	assert.Equal(t, 0, fp.ToggleCount)
	assert.InDelta(t, 0, fp.Frequency, 0.01)
	assert.InDelta(t, 0, fp.AmplitudeMean, 0.01)
}

func TestFingerprint_emptySamples(t *testing.T) {
	seg := record.AudioSegment{
		Samples:    nil,
		SampleRate: 44100,
		StartStep:  1,
		EndStep:    100,
	}

	fp := seg.Fingerprint()

	assert.True(t, fp.Silent)
}

func TestFingerprint_asymmetricWave(t *testing.T) {
	// 75% duty cycle: 3/4 positive, 1/4 negative per cycle
	sampleRate := 44100
	freq := 500.0
	samplesPerCycle := float64(sampleRate) / freq
	numSamples := 44100
	samples := make([]float32, numSamples)

	for i := range samples {
		pos := math.Mod(float64(i), samplesPerCycle)
		if pos < samplesPerCycle*0.75 {
			samples[i] = 0.5
		} else {
			samples[i] = -0.5
		}
	}

	seg := record.AudioSegment{
		Samples:    samples,
		SampleRate: sampleRate,
		StartStep:  1,
		EndStep:    100,
	}

	fp := seg.Fingerprint()

	assert.InDelta(t, 500, fp.Frequency, 15)
	assert.InDelta(t, 75, fp.DutyCycle, 3)
	assert.False(t, fp.Silent)
}
