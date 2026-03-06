package record

import "math"

// AudioSegment holds a slice of PCM samples captured over a range of
// execution steps.
type AudioSegment struct {
	Samples    []float32
	SampleRate int
	StartStep  int
	EndStep    int
}

// AudioFingerprint summarizes the acoustic properties of a segment.
type AudioFingerprint struct {
	Frequency     float64
	AmplitudeMin  float64
	AmplitudeMax  float64
	AmplitudeMean float64
	ToggleCount   int
	DutyCycle     float64
	Silent        bool
}

// Fingerprint computes an AudioFingerprint from the segment's PCM data.
func (s AudioSegment) Fingerprint() AudioFingerprint {
	if len(s.Samples) == 0 {
		return AudioFingerprint{Silent: true}
	}

	duration := float64(len(s.Samples)) / float64(s.SampleRate)

	// Count zero crossings and compute duty cycle
	var toggleCount int
	var positiveSamples int

	for i := 1; i < len(s.Samples); i++ {
		prev := s.Samples[i-1]
		curr := s.Samples[i]

		if (prev < 0 && curr >= 0) || (prev >= 0 && curr < 0) {
			toggleCount++
		}

		if curr > 0 {
			positiveSamples++
		}
	}

	// Include first sample in positive count
	if s.Samples[0] > 0 {
		positiveSamples++
	}

	frequency := float64(toggleCount) / (2.0 * duration)
	dutyCycle := float64(positiveSamples) / float64(len(s.Samples)) * 100.0

	// Amplitude envelope: divide into sub-windows and compute peak per window
	ampMin, ampMax, ampMean := amplitudeEnvelope(s.Samples)

	silent := ampMean < 0.01

	return AudioFingerprint{
		Frequency:     frequency,
		AmplitudeMin:  ampMin,
		AmplitudeMax:  ampMax,
		AmplitudeMean: ampMean,
		ToggleCount:   toggleCount,
		DutyCycle:     dutyCycle,
		Silent:        silent,
	}
}

// amplitudeEnvelope divides samples into sub-windows and computes the peak
// absolute amplitude of each window, returning the min, max, and mean of
// those peaks.
func amplitudeEnvelope(samples []float32) (ampMin, ampMax, mean float64) {
	if len(samples) == 0 {
		return 0, 0, 0
	}

	windowCount := min(100, len(samples))
	windowSize := max(len(samples)/windowCount, 1)

	var sum float64
	var count int

	for start := 0; start < len(samples); start += windowSize {
		end := min(start+windowSize, len(samples))

		var peak float64
		for i := start; i < end; i++ {
			peak = max(peak, math.Abs(float64(samples[i])))
		}

		if count == 0 {
			ampMin = peak
			ampMax = peak
		} else {
			ampMin = min(ampMin, peak)
			ampMax = max(ampMax, peak)
		}

		sum += peak
		count++
	}

	mean = sum / float64(count)
	return ampMin, ampMax, mean
}
