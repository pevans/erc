package clog

import (
	"fmt"
	"io"
)

// A Channel is able to receive messages in a channel and write those
// messages to some io.Writer.
type Channel struct {
	mesgs  chan string
	Writer io.Writer
}

// NewChannel returns a newly allocated channel, ready to receive
// messages.
func NewChannel(w io.Writer) *Channel {
	c := new(Channel)
	c.mesgs = make(chan string)
	c.Writer = w
	return c
}

// Printf will send a formatted message to the channel following the
// conventions of the fmt package.
func (c *Channel) Printf(format string, vals ...interface{}) {
	c.mesgs <- fmt.Sprintf(format, vals...)
}

// WriteLoop will iterate roughly forever, writing messages to its
// writer. You can exit the write loop by sending any bool (value
// doesn't matter) to its shutdown channel.
func (c *Channel) WriteLoop(shutdown chan bool) {
	for {
		select {
		case mesg := <-c.mesgs:
			c.Writer.Write([]byte(mesg + "\n"))
		case _ = <-shutdown:
			return
		}
	}
}
