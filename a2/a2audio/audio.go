package a2audio

import (
	"encoding/binary"
	"sync"
)

const (
	// SampleRate is the audio sample rate in Hz
	SampleRate = 44100

	// BufferSamples is the audio buffer size in samples. Smaller = lower
	// latency but higher risk of glitches. At 44100 Hz: 1024 samples = ~23ms,
	// 2048 = ~46ms, 4096 = ~93ms
	BufferSamples = 1024
)

// ToggleEvent represents a speaker toggle at a specific CPU cycle.
type ToggleEvent struct {
	Cycle uint64
	State bool // true = high, false = low
}

// EventSource provides toggle events for audio generation.
type EventSource interface {
	Pop() *ToggleEvent
	Peek() *ToggleEvent
}

// ClockSource provides the current CPU clock rate and fullspeed status.
type ClockSource interface {
	CPUClockRate() int64
	IsFullSpeed() bool
}

// Stream generates audio samples from speaker toggle events using the
// averaging approach from LinApple/AppleWin. For each audio sample (~23 CPU
// cycles), we calculate the weighted average of high/low speaker time, which
// acts as a natural low-pass filter for rapid speaker toggles.
type Stream struct {
	source EventSource
	clock  ClockSource

	mu sync.Mutex

	speakerHigh bool
	volume      float32

	// Cycle tracking for averaging approach
	currentCycle   uint64 // Current position in CPU cycles
	lastEventCycle uint64 // Cycle of last processed event
	primed         bool
	lastClockRate  int64
}

// NewStream creates a new audio stream from a toggle event source and clock source.
func NewStream(source EventSource, clock ClockSource) *Stream {
	return &Stream{
		source: source,
		clock:  clock,
		volume: 0.5,
	}
}

// SetVolume sets the audio volume (0.0 to 1.0).
func (s *Stream) SetVolume(v float32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.volume = v
}

// Read generates audio samples to fill the buffer. Each sample represents ~23
// CPU cycles (at 1.023 MHz / 44.1 kHz). We average the speaker state across
// those cycles to produce smooth audio output.
func (s *Stream) Read(buf []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.clock.IsFullSpeed() {
		// Discard all events and output silence during fullspeed
		for s.source.Pop() != nil {
		}

		s.primed = false

		for i := range buf {
			buf[i] = 0
		}

		return len(buf), nil
	}

	const bytesPerSample = 4

	numSamples := len(buf) / bytesPerSample
	clockRate := s.clock.CPUClockRate()
	cyclesPerSample := float64(clockRate) / float64(SampleRate)

	// Reset if clock rate changed
	if clockRate != s.lastClockRate && s.lastClockRate != 0 {
		s.primed = false
	}

	s.lastClockRate = clockRate

	amplitude := float64(int16(float32(16384) * s.volume))

	// Maximum gap in cycles before we skip ahead (~50ms)
	maxGapCycles := uint64(clockRate / 20)

	for i := range numSamples {
		// Prime with first event if needed
		if !s.primed {
			ev := s.source.Pop()
			if ev != nil {
				s.currentCycle = ev.Cycle
				s.lastEventCycle = ev.Cycle
				s.speakerHigh = ev.State
				s.primed = true
			}
		}

		// Check if next event is way ahead - if so, skip forward
		// This handles discontinuities like state loads or long pauses
		if s.primed {
			if ev := s.source.Peek(); ev != nil && ev.Cycle > s.currentCycle+maxGapCycles {
				s.currentCycle = ev.Cycle
			}
		}

		// Calculate cycle range for this sample
		sampleStartCycle := s.currentCycle
		sampleEndCycle := sampleStartCycle + uint64(cyclesPerSample)

		// Track time spent high vs low within this sample
		var highCycles, lowCycles float64

		if !s.primed {
			// Not primed - output silence
			lowCycles = cyclesPerSample
		} else {
			// Check if there are any events available
			nextEv := s.source.Peek()
			if nextEv == nil {
				// No events available - output current speaker state but
				// don't advance This prevents racing ahead of the emulator
				// while still reflecting the actual speaker position
				if s.speakerHigh {
					highCycles = cyclesPerSample
				} else {
					lowCycles = cyclesPerSample
				}
				// Don't advance sampleEndCycle - stay at current position
				sampleEndCycle = sampleStartCycle
			} else {
				// Process all events within this sample's cycle range
				currentPos := sampleStartCycle
				currentState := s.speakerHigh

				for {
					ev := s.source.Peek()

					// Determine the end of the current state period
					var periodEnd uint64
					if ev == nil || ev.Cycle >= sampleEndCycle {
						// No more events in this sample, or next event is beyond sample
						periodEnd = sampleEndCycle
					} else if ev.Cycle <= currentPos {
						// Event is at or before current position - consume it
						// immediately This handles gaps where events jumped
						// ahead
						s.source.Pop()
						s.lastEventCycle = ev.Cycle
						currentState = ev.State
						s.speakerHigh = ev.State
						continue
					} else {
						// Event is within this sample
						periodEnd = ev.Cycle
					}

					// Accumulate cycles for current state
					cycles := float64(periodEnd - currentPos)
					if currentState {
						highCycles += cycles
					} else {
						lowCycles += cycles
					}

					currentPos = periodEnd

					// If we've reached the sample end, break
					if currentPos >= sampleEndCycle {
						break
					}

					// Consume the event and switch state
					if ev != nil && ev.Cycle < sampleEndCycle {
						s.source.Pop()
						s.lastEventCycle = ev.Cycle
						currentState = ev.State
						s.speakerHigh = ev.State
					}
				}
			}
		}

		// Calculate weighted average sample value
		totalCycles := highCycles + lowCycles
		var sample int16
		if totalCycles > 0 {
			// Weighted average: high contributes +amplitude, low contributes -amplitude
			avgValue := (highCycles*amplitude - lowCycles*amplitude) / totalCycles
			sample = int16(avgValue)
		}

		// Write stereo sample
		offset := i * bytesPerSample
		binary.LittleEndian.PutUint16(buf[offset:], uint16(sample))
		binary.LittleEndian.PutUint16(buf[offset+2:], uint16(sample))

		// Advance to next sample
		s.currentCycle = sampleEndCycle
	}

	return numSamples * bytesPerSample, nil
}
