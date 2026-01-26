package a2audio

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockEventSource provides controllable events for testing.
type mockEventSource struct {
	events []ToggleEvent
	index  int
}

func (m *mockEventSource) Push(ev ToggleEvent) {
	m.events = append(m.events, ev)
}

func (m *mockEventSource) Peek() *ToggleEvent {
	if m.index >= len(m.events) {
		return nil
	}
	return &m.events[m.index]
}

func (m *mockEventSource) Pop() *ToggleEvent {
	if m.index >= len(m.events) {
		return nil
	}
	ev := m.events[m.index]
	m.index++
	return &ev
}

func (m *mockEventSource) Len() int {
	return len(m.events) - m.index
}

// mockClockSource provides controllable clock for testing.
type mockClockSource struct {
	clockRate int64
	fullSpeed bool
}

func (m *mockClockSource) CPUClockRate() int64 {
	return m.clockRate
}

func (m *mockClockSource) IsFullSpeed() bool {
	return m.fullSpeed
}

// sampleValue extracts the float32 sample value from a buffer position.
// Buffer is stereo 32-bit float, so each sample is 8 bytes.
func sampleValue(buf []byte, sampleIndex int) float32 {
	offset := sampleIndex * 8
	bits := binary.LittleEndian.Uint32(buf[offset:])
	return math.Float32frombits(bits)
}

// isHigh returns true if the sample represents speaker-high state.
func isHigh(sample float32) bool {
	return sample > 0
}

func TestNoEvents_ConstantOutput(t *testing.T) {
	source := &mockEventSource{}
	clock := &mockClockSource{clockRate: 1_000_000} // 1 MHz

	stream := NewStream(source, clock)
	buf := make([]byte, 800) // 100 stereo samples

	n, err := stream.Read(buf)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if n != 800 {
		t.Fatalf("expected 800 bytes, got %d", n)
	}

	// All samples should be silence (0) when there are no events
	for i := range 100 {
		sample := sampleValue(buf, i)
		if sample != 0 {
			t.Errorf("sample %d should be 0 (silence), got %f", i, sample)
		}
	}
}

func TestSingleToggle_WaveformChanges(t *testing.T) {
	// Generate a square wave and verify transitions appear in the output

	source := &mockEventSource{}
	clock := &mockClockSource{clockRate: 1_000_000}

	// Generate several toggle events to create a sustained waveform Toggle
	// every 500 cycles for a 1000 Hz square wave
	for cycle := uint64(0); cycle < 5000; cycle += 500 {
		source.Push(ToggleEvent{Cycle: cycle, State: (cycle/500)%2 == 0})
	}

	stream := NewStream(source, clock)

	buf := make([]byte, 800) // 100 samples
	_, err := stream.Read(buf)
	assert.NoError(t, err)

	// Count transitions - there should be several
	transitions := 0
	for i := 1; i < 100; i++ {
		prev := isHigh(sampleValue(buf, i-1))
		curr := isHigh(sampleValue(buf, i))
		if prev != curr {
			transitions++
		}
	}

	// With 5000 cycles of events at 22.68 cycles/sample, we cover ~220
	// samples But we only read 100 samples, covering ~2268 cycles At 500
	// cycles per toggle, that's ~4-5 toggles, so ~4-5 transitions
	if transitions < 3 {
		t.Errorf("expected at least 3 transitions, got %d", transitions)
	}
}

func TestSquareWave_CorrectFrequency(t *testing.T) {
	// Generate a 1000 Hz square wave At 1 MHz, 1000 Hz means toggle every 500
	// cycles (1000 cycles per period)

	source := &mockEventSource{}
	clock := &mockClockSource{clockRate: 1_000_000}

	// Generate toggles for several periods
	state := true
	for cycle := uint64(0); cycle < 50000; cycle += 500 {
		source.Push(ToggleEvent{Cycle: cycle, State: state})
		state = !state
	}

	stream := NewStream(source, clock)

	buf := make([]byte, 8000) // 1000 samples
	_, err := stream.Read(buf)
	assert.NoError(t, err)

	// Count zero crossings (transitions)
	transitions := 0
	for i := 1; i < 1000; i++ {
		prev := isHigh(sampleValue(buf, i-1))
		curr := isHigh(sampleValue(buf, i))
		if prev != curr {
			transitions++
		}
	}

	// At 44100 Hz sample rate, 1000 samples = ~22.7ms At 1000 Hz square wave,
	// that's ~22.7 periods = ~45 transitions Allow some tolerance
	expectedTransitions := 45
	tolerance := 5

	if transitions < expectedTransitions-tolerance || transitions > expectedTransitions+tolerance {
		t.Errorf("got %d transitions, expected around %d", transitions, expectedTransitions)
	}
}

func TestMultipleReads_Continuity(t *testing.T) {
	// Verify that multiple Read calls produce continuous output

	source := &mockEventSource{}
	clock := &mockClockSource{clockRate: 1_000_000}

	// Prime at cycle 0, toggle at cycle 1000
	source.Push(ToggleEvent{Cycle: 0, State: false})
	source.Push(ToggleEvent{Cycle: 1000, State: true})

	stream := NewStream(source, clock)

	// First read: 20 samples (covers ~453 cycles at 22.68 cycles/sample)
	buf1 := make([]byte, 160)
	_, err := stream.Read(buf1)
	assert.NoError(t, err)

	// All should be low (toggle is at cycle 1000, we only covered ~453)
	for i := range 20 {
		if isHigh(sampleValue(buf1, i)) {
			t.Errorf("sample %d in first read should be low", i)
		}
	}

	// Second read: 40 samples (covers cycles ~453 to ~1360)
	buf2 := make([]byte, 320)
	_, err = stream.Read(buf2)
	assert.NoError(t, err)

	// Should see transition somewhere in this buffer
	foundTransition := false
	for i := 1; i < 40; i++ {
		if !isHigh(sampleValue(buf2, i-1)) && isHigh(sampleValue(buf2, i)) {
			foundTransition = true
			break
		}
	}

	if !foundTransition {
		t.Error("expected transition in second read buffer")
	}
}

func TestFullSpeed_OutputsSilenceAndDiscardsEvents(t *testing.T) {
	source := &mockEventSource{}
	source.Push(ToggleEvent{Cycle: 100, State: true})

	clock := &mockClockSource{clockRate: 1_000_000, fullSpeed: true}

	stream := NewStream(source, clock)
	buf := make([]byte, 800)
	_, err := stream.Read(buf)
	assert.NoError(t, err)

	// All bytes should be zero (silence)
	for i, b := range buf {
		if b != 0 {
			t.Errorf("byte %d should be 0, got %d", i, b)
		}
	}

	// Events should be discarded during fullspeed
	if source.Peek() != nil {
		t.Error("events should be discarded during fullspeed")
	}
}

func TestSpeedChange_StillProducesOutput(t *testing.T) {
	// Simulate: generate events at 1 MHz, then change to 2 MHz mid-stream
	source := &mockEventSource{}
	clock := &mockClockSource{clockRate: 1_000_000}

	// Events at 1 MHz timing - toggle every 500 cycles for 1000 Hz
	for cycle := uint64(0); cycle < 10000; cycle += 500 {
		source.Push(ToggleEvent{Cycle: cycle, State: (cycle/500)%2 == 0})
	}

	stream := NewStream(source, clock)

	// First read at 1 MHz
	buf1 := make([]byte, 800) // 100 samples
	_, err := stream.Read(buf1)
	assert.NoError(t, err)

	transitions1 := countTransitions(buf1, 100)

	// Change to 2 MHz - this will reset primed state
	clock.clockRate = 2_000_000

	// Add more events at 2 MHz timing
	for cycle := uint64(10000); cycle < 50000; cycle += 1000 {
		source.Push(ToggleEvent{Cycle: cycle, State: (cycle/1000)%2 == 0})
	}

	// Read after speed change
	buf2 := make([]byte, 800) // 100 samples
	_, err = stream.Read(buf2)
	assert.NoError(t, err)

	transitions2 := countTransitions(buf2, 100)

	t.Logf("transitions at 1MHz: %d, at 2MHz: %d", transitions1, transitions2)

	// We mainly want to verify no crash and some reasonable output
	if transitions1 == 0 {
		t.Error("no transitions at 1 MHz")
	}
	if transitions2 == 0 {
		t.Error("no transitions at 2 MHz")
	}
}

func TestEventGap_RecoversGracefully(t *testing.T) {
	// Simulate a gap in events (like after a state load)
	source := &mockEventSource{}
	clock := &mockClockSource{clockRate: 1_000_000}

	// First batch of events: cycles 0-5000
	for cycle := uint64(0); cycle < 5000; cycle += 500 {
		source.Push(ToggleEvent{Cycle: cycle, State: (cycle/500)%2 == 0})
	}

	// Gap: no events from 5000-500000 (way more than 50ms threshold)

	// Second batch starts at cycle 500000 (simulating state load)
	for cycle := uint64(500000); cycle < 600000; cycle += 500 {
		source.Push(ToggleEvent{Cycle: cycle, State: (cycle/500)%2 == 0})
	}

	stream := NewStream(source, clock)

	// Read through the first batch
	buf := make([]byte, 8000) // 1000 samples covers ~22680 cycles
	_, err := stream.Read(buf)
	assert.NoError(t, err)

	transitions1 := countTransitions(buf, 1000)

	// Now we're past the first batch, next event is at cycle 500000 This is a
	// huge gap - the stream should cap it and continue

	buf2 := make([]byte, 8000)
	_, err = stream.Read(buf2)
	assert.NoError(t, err)

	transitions2 := countTransitions(buf2, 1000)
	t.Logf("transitions before gap: %d, after gap: %d", transitions1, transitions2)

	// Should have recovered and be producing output
	if transitions2 == 0 {
		t.Error("no transitions after event gap - recovery failed")
	}
}

func countTransitions(buf []byte, numSamples int) int {
	transitions := 0
	for i := 1; i < numSamples; i++ {
		prev := isHigh(sampleValue(buf, i-1))
		curr := isHigh(sampleValue(buf, i))
		if prev != curr {
			transitions++
		}
	}
	return transitions
}

func TestLongRunning_NoDrift(t *testing.T) {
	// Simulate several minutes of audio to check for accumulated drift.
	// Generate a 500 Hz square wave and verify frequency stays consistent.

	source := &mockEventSource{}
	clock := &mockClockSource{clockRate: 1_023_000} // Real Apple II speed

	// 500 Hz at 1.023 MHz = toggle every 1023 cycles Generate 5 minutes worth
	// of events = 300 seconds * 500 Hz * 2 toggles = 300,000 toggles
	toggleInterval := uint64(1023)
	numToggles := 300000
	for i := range numToggles {
		cycle := uint64(i) * toggleInterval
		source.Push(ToggleEvent{Cycle: cycle, State: i%2 == 0})
	}

	stream := NewStream(source, clock)

	// Log how many events remain after setup
	t.Logf("Events queued: %d", len(source.events))

	// At 44100 Hz, 5 minutes = 300 * 44100 = 13,230,000 samples We'll read in
	// chunks and check periodically
	samplesPerRead := 4410 // 0.1 seconds worth
	totalSamples := 0
	totalTransitions := 0

	// Track transitions at different points
	type checkpoint struct {
		samples     int
		transitions int
	}
	var checkpoints []checkpoint

	buf := make([]byte, samplesPerRead*8)

	// Calculate actual expected based on cycles covered
	cyclesPerSample := float64(1_023_000) / float64(44100)

	// Track state across buffer boundaries
	var lastSampleHigh *bool

	// Run for simulated 5 minutes
	for totalSamples < 13_230_000 {
		_, err := stream.Read(buf)
		assert.NoError(t, err)

		// Check for transition at buffer boundary
		if lastSampleHigh != nil {
			firstHigh := isHigh(sampleValue(buf, 0))
			if firstHigh != *lastSampleHigh {
				totalTransitions++
			}
		}

		// Count transitions within buffer
		transitions := countTransitions(buf, samplesPerRead)
		totalTransitions += transitions
		totalSamples += samplesPerRead

		// Track last sample for next iteration
		lastHigh := isHigh(sampleValue(buf, samplesPerRead-1))
		lastSampleHigh = &lastHigh

		// Record checkpoint every simulated minute
		if totalSamples%(44100*60) < samplesPerRead {
			checkpoints = append(checkpoints, checkpoint{
				samples:     totalSamples,
				transitions: totalTransitions,
			})
		}
	}

	// Calculate expected transitions based on cycles actually covered
	totalCycles := float64(totalSamples) * cyclesPerSample
	expectedTransitions := int(totalCycles / float64(toggleInterval))

	t.Logf("Total samples: %d, cycles covered (expected): %.0f", totalSamples, totalCycles)
	t.Logf("Events: consumed=%d, remaining=%d", source.index, len(source.events)-source.index)
	t.Logf("Total transitions: %d, expected: %d", totalTransitions, expectedTransitions)

	tolerance := expectedTransitions / 100 // 1% tolerance
	if totalTransitions < expectedTransitions-tolerance ||
		totalTransitions > expectedTransitions+tolerance {
		t.Errorf("drift detected: got %d transitions, expected %d ± %d",
			totalTransitions, expectedTransitions, tolerance)
	}

	// Check that rate is consistent across checkpoints
	t.Log("Checkpoints (cumulative):")
	for i, cp := range checkpoints {
		cyclesCovered := float64(cp.samples) * cyclesPerSample
		expectedAtPoint := int(cyclesCovered / float64(toggleInterval))
		drift := cp.transitions - expectedAtPoint
		driftPercent := float64(drift) / float64(expectedAtPoint) * 100
		t.Logf("  Minute %d: %d samples, %d transitions (expected %d, drift %.2f%%)",
			i+1, cp.samples, cp.transitions, expectedAtPoint, driftPercent)

		// Alert if drift exceeds 1% at any checkpoint
		if driftPercent > 1 || driftPercent < -1 {
			t.Errorf("excessive drift at minute %d: %.2f%%", i+1, driftPercent)
		}
	}
}

func TestPauseResume_RecoversGracefully(t *testing.T) {
	// When paused, no new events are generated. When resumed, events continue
	// from a later cycle. This is similar to an event gap.

	source := &mockEventSource{}
	clock := &mockClockSource{clockRate: 1_000_000}

	// Events before pause: cycles 0-10000
	for cycle := uint64(0); cycle < 10000; cycle += 500 {
		source.Push(ToggleEvent{Cycle: cycle, State: (cycle/500)%2 == 0})
	}

	stream := NewStream(source, clock)

	// First read - normal operation
	buf1 := make([]byte, 800)
	_, err := stream.Read(buf1)
	assert.NoError(t, err)
	transitions1 := countTransitions(buf1, 100)

	// Simulate pause: no events for a while, then resume with later cycles
	// This creates an event gap which gets capped
	for cycle := uint64(200000); cycle < 300000; cycle += 500 {
		source.Push(ToggleEvent{Cycle: cycle, State: (cycle/500)%2 == 0})
	}

	// Read after "pause" - should handle the gap gracefully
	buf2 := make([]byte, 800)
	_, err = stream.Read(buf2)
	assert.NoError(t, err)

	transitions2 := countTransitions(buf2, 100)

	t.Logf("transitions before pause: %d, after: %d", transitions1, transitions2)

	if transitions1 == 0 {
		t.Error("no transitions before pause")
	}
	if transitions2 == 0 {
		t.Error("no transitions after pause - recovery failed")
	}
}

func TestVolumeZero_ProducesSilence(t *testing.T) {
	source := &mockEventSource{}
	clock := &mockClockSource{clockRate: 1_000_000}

	// Generate events that would normally produce sound
	for cycle := uint64(0); cycle < 5000; cycle += 500 {
		source.Push(ToggleEvent{Cycle: cycle, State: (cycle/500)%2 == 0})
	}

	stream := NewStream(source, clock)
	stream.SetVolume(0.0) // Mute the audio

	buf := make([]byte, 800)
	_, err := stream.Read(buf)
	assert.NoError(t, err)

	// All samples should be silence (0.0) when volume is 0
	for i := range 100 {
		sample := sampleValue(buf, i)
		if sample != 0.0 {
			t.Errorf("sample %d should be 0.0 (muted), got %f", i, sample)
		}
	}
}
