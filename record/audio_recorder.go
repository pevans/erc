package record

import (
	"encoding/binary"
	"io"
	"math"
)

// AudioRecorder implements Observer to capture PCM audio samples in lockstep
// with the Recorder's step lifecycle. It reads stereo float32 PCM from a
// reader (such as a2audio.Stream) and extracts mono samples.
type AudioRecorder struct {
	label           string
	reader          io.Reader
	getCycle        func() uint64
	samplesPerCycle float64
	sampleRate      int

	// Cycle tracking per step
	startCycle uint64
	started    bool
	totalRead  int // total mono samples read so far

	// Per-step cycle log for mapping steps to sample indices
	stepCycles []uint64 // stepCycles[i] = CPU cycle at end of step i+1
	samples    []float32
	readBuf    []byte
}

// NewAudioRecorder returns an AudioRecorder that reads PCM from reader. The
// getCycle function returns the current CPU cycle count, clockRate is the CPU
// clock frequency, and sampleRate is the audio sample rate.
func NewAudioRecorder(
	label string,
	reader io.Reader,
	getCycle func() uint64,
	clockRate int64,
	sampleRate int,
) *AudioRecorder {
	return &AudioRecorder{
		label:           label,
		reader:          reader,
		getCycle:        getCycle,
		samplesPerCycle: float64(sampleRate) / float64(clockRate),
		sampleRate:      sampleRate,
	}
}

func (a *AudioRecorder) Label() string {
	return a.label
}

// Before records the starting cycle on the first call. Subsequent calls are
// no-ops because sample counts are derived from the cumulative cycle distance
// from startCycle, not from per-step deltas.
func (a *AudioRecorder) Before() {
	if !a.started {
		a.startCycle = a.getCycle()
		a.started = true
	}
}

// Observe is called after a step. It computes how many total samples should
// have been produced since startCycle, reads the delta from the reader, and
// stores them. Returns nil (audio changes are not tracked as state entries).
func (a *AudioRecorder) Observe(step int) []Entry {
	cycle := a.getCycle()
	a.stepCycles = append(a.stepCycles, cycle)

	totalSamples := a.cyclesToSamples(cycle - a.startCycle)

	delta := totalSamples - a.totalRead
	if delta <= 0 {
		return nil
	}

	// Read stereo float32 PCM: 8 bytes per sample (4 bytes * 2 channels)
	needed := delta * 8
	if cap(a.readBuf) < needed {
		a.readBuf = make([]byte, needed)
	} else {
		a.readBuf = a.readBuf[:needed]
	}

	n, _ := io.ReadFull(a.reader, a.readBuf)

	// Extract mono (left channel) from stereo pairs
	completeSamples := n / 8
	for i := range completeSamples {
		offset := i * 8
		bits := binary.LittleEndian.Uint32(a.readBuf[offset : offset+4])
		a.samples = append(a.samples, math.Float32frombits(bits))
	}

	a.totalRead += completeSamples

	return nil
}

// Segment returns the AudioSegment for the given step range. Step 0 is not a
// valid step because it implies no execution has occurred.
func (a *AudioRecorder) Segment(startStep, endStep int) AudioSegment {
	startIdx := a.sampleIndexForStep(startStep)
	endIdx := a.sampleIndexForStep(endStep)

	endIdx = min(endIdx, len(a.samples))
	if startIdx >= endIdx {
		return AudioSegment{
			SampleRate: a.sampleRate,
			StartStep:  startStep,
			EndStep:    endStep,
		}
	}

	return AudioSegment{
		Samples:    a.samples[startIdx:endIdx],
		SampleRate: a.sampleRate,
		StartStep:  startStep,
		EndStep:    endStep,
	}
}

func (a *AudioRecorder) Samples() []float32 {
	return a.samples
}

func (a *AudioRecorder) cyclesToSamples(cycles uint64) int {
	return int(float64(cycles) * a.samplesPerCycle)
}

func (a *AudioRecorder) sampleIndexForStep(step int) int {
	if step <= 0 || len(a.stepCycles) == 0 {
		return 0
	}

	idx := min(step-1, len(a.stepCycles)-1)

	return a.cyclesToSamples(a.stepCycles[idx] - a.startCycle)
}
