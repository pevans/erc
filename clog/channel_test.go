package clog

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChannel(t *testing.T) {
	assert.NotNil(t, NewChannel(os.Stdout))
}

func TestChannelPrintf(t *testing.T) {
	c := NewChannel(os.Stdout)

	go func() {
		c.Printf("hello %v", "world")
	}()

	s := <-c.mesgs
	assert.Equal(t, "hello world", s)
}

func TestChannelWriteLoop(t *testing.T) {
	b := strings.Builder{}
	c := NewChannel(&b)
	shutdown := make(chan bool)

	go func() {
		c.Printf("hello %v", "world")
		shutdown <- true
	}()

	c.WriteLoop(shutdown)

	assert.Zero(t, len(c.mesgs))
	assert.Equal(t, "hello world\n", b.String())
}
