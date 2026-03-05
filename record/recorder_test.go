package record_test

import (
	"testing"

	"github.com/pevans/erc/memory"
	"github.com/pevans/erc/record"
	"github.com/stretchr/testify/assert"
)

func TestEntryString_byte(t *testing.T) {
	e := record.Entry{Step: 514, Tag: record.TagMem, Name: "$013F", Old: uint8(0x3A), New: uint8(0x3C)}
	assert.Equal(t, "step 514: mem $013F: $3A -> $3C", e.String())
}

func TestEntryString_register(t *testing.T) {
	e := record.Entry{Step: 515, Tag: record.TagReg, Name: "A", Old: uint8(0x00), New: uint8(0x3C)}
	assert.Equal(t, "step 515: reg A: $00 -> $3C", e.String())
}

func TestEntryString_pc(t *testing.T) {
	e := record.Entry{Step: 1, Tag: record.TagReg, Name: "PC", Old: uint16(0x0300), New: uint16(0x0302)}
	assert.Equal(t, "step 1: reg PC: $0300 -> $0302", e.String())
}

func TestEntryString_bool(t *testing.T) {
	e := record.Entry{Step: 518, Tag: record.TagComp, Name: "bank-df-ram", Old: true, New: false}
	assert.Equal(t, "step 518: comp bank-df-ram: true -> false", e.String())
}

func TestRecorder_stepCountStartsAtOne(t *testing.T) {
	var r record.Recorder
	val := uint8(0)
	obs := record.NewObserver("reg", "A", func() any { return val })
	r.Add(obs)

	r.Step(func() { val = 1 })

	entries := r.Entries()
	assert.Len(t, entries, 1)
	assert.Equal(t, 1, entries[0].Step)
}

func TestRecorder_noEntryWhenUnchanged(t *testing.T) {
	var r record.Recorder
	val := uint8(42)
	r.Add(record.NewObserver(record.TagReg, "A", func() any { return val }))

	r.Step(func() {})

	assert.Empty(t, r.Entries())
}

func TestRecorder_multipleSteps(t *testing.T) {
	var r record.Recorder
	val := uint8(0)
	r.Add(record.NewObserver(record.TagReg, "X", func() any { return val }))

	r.Step(func() { val = 1 })
	r.Step(func() { val = 2 })

	entries := r.Entries()
	assert.Len(t, entries, 2)
	assert.Equal(t, 1, entries[0].Step)
	assert.Equal(t, 2, entries[1].Step)
	assert.Equal(t, uint8(0), entries[0].Old)
	assert.Equal(t, uint8(1), entries[0].New)
	assert.Equal(t, uint8(1), entries[1].Old)
	assert.Equal(t, uint8(2), entries[1].New)
}

func TestRecorder_multipleObservers(t *testing.T) {
	var r record.Recorder
	a := uint8(0)
	x := uint8(0)
	r.Add(
		record.NewObserver(record.TagReg, "A", func() any { return a }),
		record.NewObserver(record.TagReg, "X", func() any { return x }),
	)

	r.Step(func() { a = 1; x = 2 })

	entries := r.Entries()
	assert.Len(t, entries, 2)
	assert.Equal(t, "A", entries[0].Name)
	assert.Equal(t, "X", entries[1].Name)
}

func TestRecorder_memObserver(t *testing.T) {
	seg := memory.NewSegment(0x10000)
	seg.DirectSet(0x013F, 0x3A)

	var r record.Recorder
	r.Add(record.MemObserver(seg, 0x013F))

	r.Step(func() {
		seg.DirectSet(0x013F, 0x3C)
	})

	entries := r.Entries()
	assert.Len(t, entries, 1)
	assert.Equal(t, "step 1: mem $013F: $3A -> $3C", entries[0].String())
}

func TestRecorder_string(t *testing.T) {
	var r record.Recorder
	val := uint8(0)
	r.Add(record.NewObserver(record.TagReg, "A", func() any { return val }))

	r.Step(func() { val = 1 })
	r.Step(func() { val = 2 })

	expected := "step 1: reg A: $00 -> $01\nstep 2: reg A: $01 -> $02"
	assert.Equal(t, expected, r.String())
}
