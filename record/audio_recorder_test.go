package record_test

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/pevans/erc/record"
	"github.com/stretchr/testify/assert"
)

// stereoFloat32Bytes encodes mono samples as stereo float32 little-endian
// PCM.
func stereoFloat32Bytes(mono []float32) []byte {
	buf := new(bytes.Buffer)

	for _, s := range mono {
		bits := math.Float32bits(s)
		_ = binary.Write(buf, binary.LittleEndian, bits) // left
		_ = binary.Write(buf, binary.LittleEndian, bits) // right
	}

	return buf.Bytes()
}

func TestAudioRecorder_basic(t *testing.T) {
	// Simulate: clock rate 1000, sample rate 100 -- 10 cycles per sample
	var cycle uint64

	mono := squareWave(10, 100, 100, 0.5)
	reader := bytes.NewReader(stereoFloat32Bytes(mono))

	ar := record.NewAudioRecorder("speaker", reader, func() uint64 { return cycle }, 1000, 100)

	assert.Equal(t, "speaker", ar.Label())

	// Run 100 steps, each advancing 10 cycles (= 1 sample per step)
	for i := range 100 {
		ar.Before()
		cycle += 10
		entries := ar.Observe(i + 1)
		assert.Nil(t, entries)
	}

	samples := ar.Samples()
	assert.Len(t, samples, 100)
}

func TestAudioRecorder_segment(t *testing.T) {
	var cycle uint64

	// 50 samples with known values
	mono := make([]float32, 50)
	for i := range mono {
		mono[i] = float32(i) / 50.0
	}

	reader := bytes.NewReader(stereoFloat32Bytes(mono))
	ar := record.NewAudioRecorder("speaker", reader, func() uint64 { return cycle }, 1000, 100)

	// 50 steps, 10 cycles each = 1 sample per step
	for i := range 50 {
		ar.Before()
		cycle += 10
		ar.Observe(i + 1)
	}

	seg := ar.Segment(11, 20)
	assert.Equal(t, 11, seg.StartStep)
	assert.Equal(t, 20, seg.EndStep)
	assert.Equal(t, 100, seg.SampleRate)
	assert.Len(t, seg.Samples, 9)

	// Samples should correspond to mono[11:20]
	for i, s := range seg.Samples {
		assert.InDelta(t, float32(i+11)/50.0, s, 0.001)
	}
}

func TestAudioRecorder_subSampleSteps(t *testing.T) {
	// Simulate Apple II-like: clockRate=1023000, sampleRate=44100 Most steps
	// consume ~3 cycles, meaning many steps share a sample
	var cycle uint64

	// Provide enough samples for the total cycle count we'll accumulate
	mono := make([]float32, 100)
	for i := range mono {
		mono[i] = 0.25
	}

	reader := bytes.NewReader(stereoFloat32Bytes(mono))
	ar := record.NewAudioRecorder("speaker", reader, func() uint64 { return cycle }, 1023000, 44100)

	// 500 steps of 3 cycles each = 1500 total cycles 1500 * 44100 / 1023000
	// ≈ 64 samples
	for i := range 500 {
		ar.Before()
		cycle += 3
		ar.Observe(i + 1)
	}

	samples := ar.Samples()
	assert.Greater(t, len(samples), 50)
	assert.Less(t, len(samples), 80)
}

func TestAudioRecorder_withRecorder(t *testing.T) {
	var cycle uint64
	var r record.Recorder

	mono := squareWave(10, 100, 20, 0.5)
	reader := bytes.NewReader(stereoFloat32Bytes(mono))

	ar := record.NewAudioRecorder("speaker", reader, func() uint64 { return cycle }, 1000, 100)
	r.Add(ar)

	for range 20 {
		r.Step(func() { cycle += 10 })
	}

	samples := ar.Samples()
	assert.Len(t, samples, 20)
}
