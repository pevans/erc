---
Specification: 2
Category: Tests
Drafted At: 2026-03-05
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes a system for recording and testing audio output during
emulator execution. Where spec 1 (state testing) captures discrete state
transitions -- a register changing value, a memory byte being written -- audio
is a continuous signal that emerges from the interaction of many state changes
over time. A single speaker toggle is not meaningful on its own; what matters
is the waveform that results from a sequence of toggles.

This system records the rendered audio output as a time-series and compares it
against expected properties using fingerprint-based assertions rather than
exact sample matching. This approach works for a simple 1-bit speaker and
scales to multi-channel synthesis hardware.

# 2. Concepts

## 2.1. Audio Recorder

An audio recorder captures the PCM audio stream produced by the emulator's
audio pipeline. It sits at the output of the synthesis stage -- downstream of
all toggle-to-sample conversion, envelope generation, channel mixing, and any
other signal processing. The recorder captures what the user would hear.

The audio recorder operates alongside the state recorder from spec 1. Both
use the same step counter as their shared timeline.

## 2.2. Audio Segment

An audio segment is a contiguous range of audio defined by a start step and
an end step. A segment contains the PCM samples that were generated during
that step range. Segments are the unit of comparison in audio tests -- rather
than asserting properties of individual samples, tests assert properties of
segments.

## 2.3. Audio Fingerprint

An audio fingerprint is a set of derived properties that summarize the
content of an audio segment. Fingerprints abstract away the raw sample data
and describe the audio in terms that are stable across minor timing
variations. A fingerprint includes:

- **Frequency estimate** -- the dominant frequency in the segment, derived
  from zero-crossing analysis
- **Amplitude envelope** -- the minimum, maximum, and mean amplitude across
  the segment
- **Toggle count** -- the number of signal transitions in the segment (most
  relevant for square wave sources)
- **Duty cycle** -- the ratio of high-to-low time within each wave period
- **Silence** -- whether the segment contains no meaningful signal

## 2.4. Audio Assertion

An audio assertion declares the expected properties of a segment. Rather than
requiring exact values, assertions use tolerances. The assertion syntax is:

```
step 100-500: audio freq ~1000Hz +/- 50, amplitude > 0.3
step 501-800: audio silent
step 801-1200: audio freq ~440Hz, duty ~50% +/- 5
```

To assert on a specific channel, include the channel label after `audio`:

```
step 100-500: audio speaker freq ~1000Hz +/- 50
step 100-500: audio mockingboard-A-1 freq ~440Hz +/- 10
```

A test passes if every assertion's properties fall within the specified
tolerances. A test fails if any property is outside tolerance, or if audio is
present where silence was expected (or vice versa).

# 3. Recording

## 3.1. Tap Point

The audio recorder taps the audio stream at the point where PCM samples are
produced for playback. This is after all synthesis and mixing -- the recorder
captures the final output, not intermediate state. This means the recorder
does not need to know how the audio was generated, only what came out.

For systems with multiple audio channels (e.g. a sound card with independent
tone generators), each channel can optionally be recorded separately by
tapping the channel's pre-mix output. The post-mix mixed output is always
recorded.

## 3.2. Step-to-Sample Mapping

The recorder must map step counts to sample positions in the captured audio.
Given the CPU clock rate and the audio sample rate, the mapping is:

```
sample_index = (step_cycle - start_cycle) * sample_rate / clock_rate
```

Where `step_cycle` is the CPU cycle at which a given step occurred,
`start_cycle` is the CPU cycle at the start of recording, `sample_rate` is
the audio output rate (e.g. 44100 Hz), and `clock_rate` is the emulated CPU
clock speed (e.g. 1023000 Hz for the Apple II).

This mapping allows the recorder to extract the sample range corresponding to
any step range, which is then used to compute the segment's fingerprint.

## 3.3. Channel Labels

Each recorded stream is identified by a channel label. For a simple 1-bit
speaker, there is one channel (e.g. `speaker`). For multi-channel hardware,
each channel has its own label (e.g. `mockingboard-A-1`, `mockingboard-A-2`).
The mixed output uses the label `mixed`.

Assertions specify which channel they apply to. If no channel is specified,
the assertion applies to the `mixed` output.

# 4. Fingerprinting

## 4.1. Frequency Estimation

The dominant frequency of a segment is estimated using zero-crossing
analysis. A zero crossing occurs when the signal changes sign (positive to
negative or vice versa). The number of zero crossings in a time interval
gives an estimate of frequency:

```
frequency = zero_crossings / (2 * duration_in_seconds)
```

This method is well-suited for square waves and simple tonal signals. For
segments containing multiple simultaneous frequencies (e.g. a chord from
multiple channels), frequency estimation applies to each channel
independently rather than to the mixed output.

## 4.2. Amplitude Envelope

The amplitude envelope captures the signal's volume characteristics over the
segment. It is computed by dividing the segment into sub-windows and
measuring the peak absolute amplitude in each window. The number and size of
sub-windows is an implementation detail, but should be consistent within a
test run so that min/max/mean values are comparable. The fingerprint records:

- **Min amplitude** -- the lowest sub-window peak
- **Max amplitude** -- the highest sub-window peak
- **Mean amplitude** -- the average of all sub-window peaks

This captures volume changes over time, such as an envelope decay or a
sudden onset.

## 4.3. Silence Detection

A segment is considered silent if the mean amplitude (on a normalized scale
where samples range from -1.0 to 1.0) is below a threshold (e.g. 0.01).
This is a distinct property from "no toggle events occurred" --
a speaker that is held high produces a DC offset that is not silence in the
toggle sense but produces no audible tone.

## 4.4. Duty Cycle

For square wave signals, the duty cycle is the fraction of each period spent
in the high state. A 50% duty cycle produces a symmetric square wave. The
duty cycle is estimated from the PCM waveform by measuring the ratio of
positive-to-negative sample durations per cycle.

# 5. Testing Model

## 5.1. Audio Test Structure

An audio test consists of:

1. A program or disk image to load into the emulator
2. A set of state observers to register (from spec 1, if needed)
3. An audio recorder attached to the audio stream
4. A number of steps to execute
5. A list of audio assertions, each specifying a step range, an optional
   channel, and expected fingerprint properties with tolerances

The audio recorder begins recording at step 0 and continues for the
duration of the test. There is no explicit start/stop control -- the
recorder captures the entire execution.

## 5.2. Assertion Matching

Each assertion is evaluated independently. For a given assertion:

1. Extract the PCM samples corresponding to the assertion's step range using
   the step-to-sample mapping.
2. Compute the fingerprint of the extracted segment.
3. Compare each asserted property against the fingerprint, respecting
   tolerances.

A test passes if all assertions pass. A failing assertion reports the
expected and actual values for the property that was out of tolerance.

## 5.3. Combining with State Tests

Audio assertions and state entry assertions (from spec 1) can coexist in the
same test. Both systems share the same step counter and the same execution
run. This allows a test to verify, for example, that a program writes a
particular value to a memory address *and* produces a tone of a given
frequency.

# 6. Raw Sample Capture

## 6.1. Debugging Failed Tests

When an audio assertion fails, the raw PCM samples for the failing segment
are available for inspection. The recorder can write these samples to a file
for offline analysis (e.g. importing into an audio editor or plotting the
waveform). This is analogous to printing the recorded state entry list when a
state test fails.

## 6.2. Not Used for Comparison

Raw PCM samples are never used directly for test comparison. Sample-level
comparison is fragile -- minor changes in timing, interpolation, or synthesis
can shift individual samples without changing the audible result.
Fingerprint-based comparison captures the perceptually meaningful properties
while tolerating implementation-level variation.

# 7. Design Considerations

## 7.1. Scalability

The fingerprint approach scales to complex audio sources. Adding a new sound
chip (e.g. a Mockingboard with six tone channels and two noise channels) does
not require rewriting the test framework -- it requires adding new channel
labels and tapping each channel's output. The same assertion language and
fingerprinting logic apply regardless of the source.

## 7.2. Independence from Synthesis Method

Because the recorder taps the rendered PCM output, it is independent of how
audio is synthesized. A 1-bit speaker toggling via the Apple II's `$C030` soft
switch, a PSG chip with frequency registers, and a wavetable synthesizer all
produce PCM samples that can be fingerprinted the same way. Tests written
against the output remain valid even if the internal synthesis implementation
changes.

## 7.3. Tolerance and Precision

Frequency estimation via zero-crossing has limited precision, especially for
short segments or low frequencies. Tests should use tolerances appropriate to
the segment duration and expected frequency. As a guideline, a segment should
contain at least several complete cycles of the expected waveform for a
reliable frequency estimate.

## 7.4. Test-Only Concern

As with the state recording system in spec 1, the audio recorder is only
instantiated during testing. Production execution paths do not create an
audio recorder. The recorder taps the existing audio stream interface and
does not require changes to the synthesis or playback code.
