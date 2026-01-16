package input

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPushEvent(t *testing.T) {
	el := len(eventChannel)

	go func() {
		PushEvent(Event{Key: 'a', Modifier: ModOption})
	}()

	e := <-eventChannel
	assert.Equal(t, len(eventChannel), el)
	assert.Equal(t, 'a', e.Key)
	assert.Equal(t, ModOption, e.Modifier)
}

func TestListen(t *testing.T) {
	key := 'f'

	// This will run indefinitely on a goroutine waiting for events, executing
	// the listener as necessary
	go Listen(func(e Event) {
		key = e.Key
	})

	PushEvent(Event{Key: 'g', Modifier: ModShift})

	// We're not really sure when Listen will get to the event, so use
	// Eventually to set up a loop that tests for what we want
	assert.Eventually(t, func() bool {
		return key == 'g'
	}, 1*time.Second, 10*time.Millisecond)

	shutdown <- true
}

func TestShutdown(t *testing.T) {
	sl := len(shutdown)

	go func() {
		Shutdown()
	}()

	<-shutdown
	assert.Equal(t, len(shutdown), sl)
}
