package elog

import (
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	// AudioSampleRate is the audio sample rate in Hz
	AudioSampleRate = 44100

	// AudioFrameDuration is how long each frame captures (in seconds)
	AudioFrameDuration = 1.0

	// SamplesPerFrame is the number of samples in each frame
	SamplesPerFrame = int(AudioSampleRate * AudioFrameDuration)
)

// AudioFrame represents one second of captured audio data.
type AudioFrame struct {
	// Timestamp is the time since boot when this frame was captured
	Timestamp float64

	// Samples contains the audio sample values (16-bit signed integers)
	// This is mono data - one sample per time point
	Samples []int16
}

// AudioLog is a collection of AudioFrames, representing all audio captured
// during emulation.
type AudioLog struct {
	// frames are all the AudioFrames that we're holding
	frames []AudioFrame

	// currentFrame is the frame being built
	currentFrame AudioFrame

	// sampleCount tracks samples added to current frame
	sampleCount int
}

// NewAudioLog returns a newly allocated AudioLog ready to record audio.
func NewAudioLog() *AudioLog {
	return &AudioLog{
		frames: make([]AudioFrame, 0),
		currentFrame: AudioFrame{
			Samples: make([]int16, 0, SamplesPerFrame),
		},
	}
}

// AddSamples adds audio samples to the log. When a full second of samples has
// been collected, it creates a new frame.
func (a *AudioLog) AddSamples(samples []int16, timestamp float64) {
	for _, sample := range samples {
		a.currentFrame.Samples = append(a.currentFrame.Samples, sample)
		a.sampleCount++

		if a.sampleCount >= SamplesPerFrame {
			// Complete the current frame
			a.currentFrame.Timestamp = timestamp
			a.frames = append(a.frames, a.currentFrame)

			// Start a new frame
			a.currentFrame = AudioFrame{
				Samples: make([]int16, 0, SamplesPerFrame),
			}
			a.sampleCount = 0
		}
	}
}

// renderWaveform converts audio samples to a text waveform visualization.
// Each line represents a portion of time, and uses characters to show
// amplitude.
func renderWaveform(samples []int16, width int) []string {
	if len(samples) == 0 {
		return []string{}
	}

	// How many samples per column
	samplesPerCol := max(len(samples)/width, 1)

	// Height of the waveform (in lines)
	const height = 20

	lines := make([]string, height)
	for i := range lines {
		lines[i] = strings.Repeat(" ", width)
	}

	// For each column, calculate the RMS (root mean square) amplitude
	// and plot it
	for col := 0; col < width && col*samplesPerCol < len(samples); col++ {
		startIdx := col * samplesPerCol
		endIdx := min(startIdx+samplesPerCol, len(samples))

		// Calculate RMS amplitude for this column
		var sumSquares float64
		for i := startIdx; i < endIdx; i++ {
			val := float64(samples[i])
			sumSquares += val * val
		}
		rms := math.Sqrt(sumSquares / float64(endIdx-startIdx))

		// Normalize to 0-1 range (assuming max amplitude is 16384)
		normalized := rms / 16384.0
		if normalized > 1.0 {
			normalized = 1.0
		}

		// Map to height (0 at bottom, height-1 at top)
		// We'll use the middle as zero, and show amplitude symmetrically
		midHeight := height / 2
		amplitudeHeight := int(normalized * float64(midHeight))

		// Plot symmetrically above and below middle
		for h := midHeight - amplitudeHeight; h <= midHeight+amplitudeHeight; h++ {
			if h >= 0 && h < height {
				// Use different characters for different amplitudes
				char := '·'
				if normalized > 0.8 {
					char = '█'
				} else if normalized > 0.6 {
					char = '▓'
				} else if normalized > 0.4 {
					char = '▒'
				} else if normalized > 0.2 {
					char = '░'
				}

				// Convert line to rune slice, modify, convert back
				runes := []rune(lines[h])
				runes[col] = char
				lines[h] = string(runes)
			}
		}
	}

	// Reverse lines so highest amplitude is at top
	for i := 0; i < len(lines)/2; i++ {
		j := len(lines) - 1 - i
		lines[i], lines[j] = lines[j], lines[i]
	}

	return lines
}

// FrameAnalysis contains dropout and pop detection metrics for an audio frame.
type FrameAnalysis struct {
	ZeroCrossings    int     // Number of times the waveform crosses zero
	MaxRunLength     int     // Longest sequence of identical samples
	ActivityRate     float64 // Percentage of time with varying samples (0-100)
	ActivityTimeline string  // Visual representation of activity
}

// analyzeFrame computes dropout and pop detection metrics.
func analyzeFrame(samples []int16) FrameAnalysis {
	if len(samples) == 0 {
		return FrameAnalysis{}
	}

	var zeroCrossings int
	var maxRunLength int
	var currentRunLength int
	var lastSample int16

	// Count active windows (100ms windows with >10 sample changes)
	const windowSize = AudioSampleRate / 10 // 100ms windows
	activeWindows := 0
	totalWindows := 0

	timeline := strings.Builder{}
	const timelineWidth = 80

	for i, sample := range samples {
		// Zero crossing detection
		if i > 0 && ((samples[i-1] < 0 && sample >= 0) || (samples[i-1] >= 0 && sample < 0)) {
			zeroCrossings++
		}

		// Run length detection (consecutive identical samples)
		if i > 0 && sample == lastSample {
			currentRunLength++
			if currentRunLength > maxRunLength {
				maxRunLength = currentRunLength
			}
		} else {
			currentRunLength = 1
		}
		lastSample = sample

		// Activity timeline (per window)
		if i > 0 && i%windowSize == 0 {
			// Check how active this window was
			windowStart := i - windowSize
			uniqueValues := make(map[int16]bool)
			for j := windowStart; j < i && j < len(samples); j++ {
				uniqueValues[samples[j]] = true
			}

			totalWindows++
			if len(uniqueValues) > 10 {
				activeWindows++
				timeline.WriteRune('█')
			} else if len(uniqueValues) > 5 {
				timeline.WriteRune('▒')
			} else if len(uniqueValues) > 2 {
				timeline.WriteRune('░')
			} else {
				timeline.WriteRune('·')
			}
		}
	}

	// Pad timeline to fixed width
	timelineStr := timeline.String()
	if len(timelineStr) < timelineWidth {
		timelineStr += strings.Repeat(" ", timelineWidth-len(timelineStr))
	}
	if len(timelineStr) > timelineWidth {
		timelineStr = timelineStr[:timelineWidth]
	}

	activityRate := 0.0
	if totalWindows > 0 {
		activityRate = float64(activeWindows) / float64(totalWindows) * 100.0
	}

	return FrameAnalysis{
		ZeroCrossings:    zeroCrossings,
		MaxRunLength:     maxRunLength,
		ActivityRate:     activityRate,
		ActivityTimeline: timelineStr,
	}
}

// WriteToFile writes the audio log to a file in text format.
func (a *AudioLog) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:errcheck

	_, err = fmt.Fprintf(file, "Audio Log - Sample Rate: %d Hz\n", AudioSampleRate)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "Each frame represents %.1f second of audio\n", AudioFrameDuration)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "\nActivity Timeline Legend:\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "  █ = Active audio (>10 unique values per 100ms)\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "  ▒ = Moderate activity (6-10 unique values)\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "  ░ = Low activity (3-5 unique values)\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "  · = Idle/dropout (≤2 unique values)\n\n")
	if err != nil {
		return err
	}

	const waveformWidth = 80

	for _, frame := range a.frames {
		_, err := fmt.Fprintf(file, "FRAME %.6f\n", frame.Timestamp)
		if err != nil {
			return err
		}

		// Calculate statistics for this frame
		var minSample, maxSample int16
		var sumAbs int64
		if len(frame.Samples) > 0 {
			minSample = frame.Samples[0]
			maxSample = frame.Samples[0]
			for _, sample := range frame.Samples {
				if sample < minSample {
					minSample = sample
				}
				if sample > maxSample {
					maxSample = sample
				}
				if sample < 0 {
					sumAbs += int64(-sample)
				} else {
					sumAbs += int64(sample)
				}
			}
		}
		avgAbs := float64(sumAbs) / float64(len(frame.Samples))

		_, err = fmt.Fprintf(file, "  Samples: %d, Min: %d, Max: %d, Avg Amplitude: %.1f\n",
			len(frame.Samples), minSample, maxSample, avgAbs)
		if err != nil {
			return err
		}

		// Analyze frame for dropouts and pops
		analysis := analyzeFrame(frame.Samples)
		_, err = fmt.Fprintf(file, "  Zero Crossings: %d, Max Run: %d samples, Activity: %.1f%%\n",
			analysis.ZeroCrossings, analysis.MaxRunLength, analysis.ActivityRate)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(file, "  Timeline: %s\n", analysis.ActivityTimeline)
		if err != nil {
			return err
		}

		// Render waveform
		waveform := renderWaveform(frame.Samples, waveformWidth)
		for _, line := range waveform {
			_, err := fmt.Fprintln(file, "  "+line)
			if err != nil {
				return err
			}
		}

		_, err = fmt.Fprintln(file)
		if err != nil {
			return err
		}
	}

	return nil
}
