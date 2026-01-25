package a2speaker

import (
	"sync"
	"time"

	"github.com/pevans/erc/a2/a2audio"
)

// Speaker is an interface for accessing speaker-related methods.
type Speaker interface {
	Push(cycle uint64, state bool)
	Pop() *a2audio.ToggleEvent
	Peek() *a2audio.ToggleEvent
	Len() int
	IsActive() bool
}

const (
	// speakerActivityTimeout is how long after the last speaker toggle we
	// consider the speaker "active". During this time, fullspeed mode should
	// be inhibited to prevent audio gaps. This needs to be longer than the
	// longest delay a sound routine might have between toggles (e.g., the
	// Apple II WAIT routine with A=$C0 is ~180ms).
	speakerActivityTimeout = 300 * time.Millisecond
)

// SpeakerBuffer holds recent speaker toggle events for audio generation. It's
// designed as a ring buffer.
type SpeakerBuffer struct {
	mu     sync.Mutex
	events []a2audio.ToggleEvent
	head   int
	tail   int
	size   int

	// Activity tracking for fullspeed inhibition
	lastActivity time.Time

	// Debug counters
	Pushed  uint64
	Popped  uint64
	Dropped uint64
}

// NewSpeakerBuffer creates a new speaker buffer for toggle events.
func NewSpeakerBuffer(size int) *SpeakerBuffer {
	return &SpeakerBuffer{
		events: make([]a2audio.ToggleEvent, size),
		size:   size,
	}
}

// Push adds a toggle event to the buffer.
func (sb *SpeakerBuffer) Push(cycle uint64, state bool) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	sb.events[sb.head] = a2audio.ToggleEvent{
		Cycle: cycle,
		State: state,
	}
	sb.head = (sb.head + 1) % sb.size
	sb.Pushed++
	sb.lastActivity = time.Now()

	if sb.head == sb.tail {
		// Buffer full, advance tail (drop oldest)
		sb.tail = (sb.tail + 1) % sb.size
		sb.Dropped++
	}
}

// IsActive returns true if the speaker has been toggled recently (within
// speakerActivityTimeout) OR if there are pending events in the buffer
// waiting to be processed. This is used to inhibit fullspeed mode during
// sound playback.
func (sb *SpeakerBuffer) IsActive() bool {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	// If there are events waiting to be processed, speaker is active
	if sb.head != sb.tail {
		return true
	}

	// Also check recent activity (handles the case where buffer was just
	// drained but sound is still playing)
	if sb.lastActivity.IsZero() {
		return false
	}

	return time.Since(sb.lastActivity) < speakerActivityTimeout
}

// Pop removes and returns the oldest toggle event, or nil if empty.
func (sb *SpeakerBuffer) Pop() *a2audio.ToggleEvent {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if sb.head == sb.tail {
		return nil
	}
	ev := sb.events[sb.tail]
	sb.tail = (sb.tail + 1) % sb.size
	sb.Popped++
	return &ev
}

// Peek returns the oldest toggle event without removing it, or nil if empty.
func (sb *SpeakerBuffer) Peek() *a2audio.ToggleEvent {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if sb.head == sb.tail {
		return nil
	}
	ev := sb.events[sb.tail]
	return &ev
}

// Len returns the number of events in the buffer.
func (sb *SpeakerBuffer) Len() int {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if sb.head >= sb.tail {
		return sb.head - sb.tail
	}
	return sb.size - sb.tail + sb.head
}

// Stats returns the debug counters (pushed, popped, dropped).
func (sb *SpeakerBuffer) Stats() (pushed, popped, dropped uint64) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.Pushed, sb.Popped, sb.Dropped
}
