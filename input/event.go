package input

// An Event is a record of something that happened with user input.
type Event struct {
	Key      rune
	Modifier int
}

// An EventListener is a function which receives an event
type EventListener func(e Event)

// These are key modifiers. Note that they represent a real keyboard
// you'd use today, not a virtual keyboard from any emulated device.
const (
	ModNone    = iota
	ModShift   // shift key
	ModOption  // option key
	ModControl // control key
	ModCommand // command key
)

const (
	KeyNone rune = 0
)

var (
	eventChannel = make(chan Event)
	shutdown     = make(chan bool)
	listening    = false
)

// PushEvent adds a new event to the event channel, which something else
// can listen for and act on.
func PushEvent(e Event) {
	eventChannel <- e
}

// Listen registers an EventListener and starts a loop in a goroutine,
// which listens for new events (via PushEvent) and calls the listen
// function with them. Alternatively, if a shutdown event is received,
// it ends.
func Listen(listen EventListener) {
	listening = true
	for {
		select {
		case e := <-eventChannel:
			listen(e)
		case <-shutdown:
			return
		}
	}
}

// Shutdown sends a message to the shutdown channel, which would end any
// Listen goroutine.
func Shutdown() {
	if !listening {
		return
	}

	shutdown <- true
}
