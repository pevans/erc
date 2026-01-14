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
	Len() int
}

// ClockSource provides the current CPU clock rate and fullspeed status.
type ClockSource interface {
	CPUClockRate() int64
	IsFullSpeed() bool
}

// AudioLogger is an interface for logging audio samples.
type AudioLogger interface {
	AddSamples(samples []int16, timestamp float64)
}

// StreamStats contains diagnostic counters for audio stream health
// monitoring.
type StreamStats struct {
	SamplesGenerated uint64 // Total samples produced
	EventsProcessed  uint64 // Toggle events consumed
	GapsDetected     uint64 // Timeline resyncs to event stream
	FullSpeedSamples uint64 // Samples output as silence during fullspeed
	LastSample       int16  // Most recent sample value output
	CurrentBufferLen int    // Current event buffer length
}

// Stream generates audio samples by syncing to the speaker toggle event stream.
// Rather than maintaining an independent timeline, we follow the events and only
// advance when we have data to process. This avoids gaps when the event stream
// pauses temporarily.
type Stream struct {
	source EventSource
	clock  ClockSource

	mu sync.Mutex

	speakerHigh bool
	volume      float32
	lastSample  int16 // Most recent sample value, for diagnostics

	// Cycle tracking (synced to event stream)
	currentCycle uint64

	// Diagnostic counters
	samplesGenerated uint64
	eventsProcessed  uint64
	gapsDetected     uint64
	fullSpeedSamples uint64

	// Optional audio logger for debugging
	audioLogger AudioLogger
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

// SetAudioLogger sets an optional audio logger for debugging.
func (s *Stream) SetAudioLogger(logger AudioLogger) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.audioLogger = logger
}

// Stats returns diagnostic counters for monitoring stream health.
func (s *Stream) Stats() StreamStats {
	s.mu.Lock()
	defer s.mu.Unlock()
	return StreamStats{
		SamplesGenerated: s.samplesGenerated,
		EventsProcessed:  s.eventsProcessed,
		GapsDetected:     s.gapsDetected,
		FullSpeedSamples: s.fullSpeedSamples,
		LastSample:       s.lastSample,
		CurrentBufferLen: s.source.Len(),
	}
}

// Read generates audio samples based on speaker toggle events. We sync our
// timeline to the event stream and only advance when we have events to
// process.
func (s *Stream) Read(buf []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	const bytesPerSample = 4
	numSamples := len(buf) / bytesPerSample
	clockRate := s.clock.CPUClockRate()
	cyclesPerSample := float64(clockRate) / float64(SampleRate)

	// Track samples for logging (mono only)
	var logSamples []int16
	if s.audioLogger != nil {
		logSamples = make([]int16, 0, numSamples)
	}

	amplitude := float64(int16(float32(16384) * s.volume))

	// If in fullspeed mode, consume events and output silence
	if s.clock.IsFullSpeed() {
		for s.source.Pop() != nil {
			s.eventsProcessed++
		}

		s.currentCycle = 0 // Reset timeline and resync when fullspeed ends

		// Output silence
		for i := range buf {
			buf[i] = 0
		}

		s.fullSpeedSamples += uint64(numSamples)
		s.samplesGenerated += uint64(numSamples)

		return numSamples * bytesPerSample, nil
	}

	// Peek at next event to sync our timeline
	if ev := s.source.Peek(); ev != nil {
		// If we haven't started yet, or there's a large gap, sync to the
		// event stream
		if s.currentCycle == 0 || ev.Cycle > s.currentCycle+uint64(clockRate/10) {
			// Start our timeline at the first event's cycle
			s.currentCycle = ev.Cycle
			s.gapsDetected++
		}
	}

	for i := range numSamples {
		sampleStartCycle := s.currentCycle
		sampleEndCycle := sampleStartCycle + uint64(cyclesPerSample)

		// Check if we have any events to process for this sample period
		ev := s.source.Peek()
		hasEventsInRange := ev != nil && ev.Cycle < sampleEndCycle

		if !hasEventsInRange && s.source.Len() == 0 {
			// Just output based on current speaker state
			var sample int16

			if s.speakerHigh {
				sample = int16(amplitude)
			} else {
				sample = int16(-amplitude)
			}

			if s.audioLogger != nil {
				logSamples = append(logSamples, sample)
			}

			offset := i * bytesPerSample
			binary.LittleEndian.PutUint16(buf[offset:], uint16(sample))
			binary.LittleEndian.PutUint16(buf[offset+2:], uint16(sample))
			s.lastSample = sample
			s.samplesGenerated++

			continue
		}

		// Process events and generate sample
		var highCycles, lowCycles float64

		currentPos := sampleStartCycle
		currentState := s.speakerHigh

		for {
			ev := s.source.Peek()

			var periodEnd uint64
			if ev == nil || ev.Cycle >= sampleEndCycle {
				periodEnd = sampleEndCycle
			} else if ev.Cycle <= currentPos {
				// If there's an event at or before current position, we
				// consume and update state
				s.source.Pop()
				s.eventsProcessed++
				currentState = ev.State
				s.speakerHigh = ev.State
				continue
			} else {
				periodEnd = ev.Cycle
			}

			cycles := float64(periodEnd - currentPos)
			if currentState {
				highCycles += cycles
			} else {
				lowCycles += cycles
			}

			currentPos = periodEnd

			if currentPos >= sampleEndCycle {
				break
			}

			if ev != nil && ev.Cycle < sampleEndCycle {
				s.source.Pop()
				s.eventsProcessed++
				currentState = ev.State
				s.speakerHigh = ev.State
			}
		}

		// Calculate sample value
		totalCycles := highCycles + lowCycles
		var sample int16
		if totalCycles > 0 {
			avgValue := (highCycles*amplitude - lowCycles*amplitude) / totalCycles
			sample = int16(avgValue)
		} else {
			// Just use the current speaker state
			if s.speakerHigh {
				sample = int16(amplitude)
			} else {
				sample = int16(-amplitude)
			}
		}

		if s.audioLogger != nil {
			logSamples = append(logSamples, sample)
		}

		offset := i * bytesPerSample
		binary.LittleEndian.PutUint16(buf[offset:], uint16(sample))
		binary.LittleEndian.PutUint16(buf[offset+2:], uint16(sample))

		s.lastSample = sample
		s.samplesGenerated++
		s.currentCycle = sampleEndCycle
	}

	if s.audioLogger != nil && len(logSamples) > 0 {
		timestamp := float64(s.currentCycle) / float64(clockRate)
		s.audioLogger.AddSamples(logSamples, timestamp)
	}

	return numSamples * bytesPerSample, nil
}
