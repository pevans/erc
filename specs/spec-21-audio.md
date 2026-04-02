---
Specification: 21
Category: Sound
Drafted At: 2026-03-30
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes how the Apple II produces audio. The Apple II has no
dedicated sound chip. Instead, it has a single-bit speaker that software
controls by toggling a soft switch. Each access to the switch physically
deflects the speaker cone, and a rapid series of toggles produces an audible
waveform. The frequency and shape of the sound are entirely determined by how
often -- and with what timing -- software accesses the toggle address.

The emulator must convert these toggle events into a PCM audio stream suitable
for playback on the host system.

# 2. Hardware Model

## 2.1. The Speaker

The Apple II speaker is a 1-bit output device. It has two physical positions:
high and low. Accessing the speaker toggle switch flips the speaker from its
current position to the opposite one. There is no way to set the speaker to a
specific position -- each access simply inverts the current state.

The emulated speaker maintains two pieces of state:

- **SpeakerState**: a boolean indicating the current position of the speaker
  (high or low). Default: low (false).
- **ToggleCycle**: the CPU cycle count at which the most recent toggle
  occurred. Used by the audio stream to determine timing.

SpeakerState is initialized to false (low) at power-on.

## 2.2. No Waveform Generator

Unlike machines with dedicated sound hardware (SID, AY-3-8910, etc.), the
Apple II has no oscillator, envelope generator, or volume register. The
speaker is either fully deflected or fully at rest. Any waveform -- square,
triangle-approximation, or noise -- must be constructed by software through
carefully timed toggle sequences.

# 3. Soft Switch

## 3.1. $C030 -- Speaker Toggle (Read and Write)

Any access to address $C030 -- read or write, with any value -- toggles the
speaker state. The data bus value is irrelevant; only the address decode
matters.

    Address   Trigger       Side Effect
    -------   -------       -----------
    $C030     Read          Toggle speaker state
    $C030     Write         Toggle speaker state

On each toggle, the emulator:

1. Inverts SpeakerState.
2. Records the current CPU cycle count and the new state as a toggle event.
3. Pushes the toggle event into the speaker buffer (section 4).

The read returns 0 (floating bus behavior is not modeled for this address).

# 4. Speaker Buffer

Toggle events must cross from the CPU thread to the audio thread. A ring
buffer serves as the intermediary.

## 4.1. Toggle Event

Each event records:

- **Cycle**: the CPU cycle at which the toggle occurred (uint64).
- **State**: the new speaker position after the toggle (true = high, false =
  low).

## 4.2. Ring Buffer

The buffer is a fixed-size circular queue of 1,024 entries. If the buffer
fills (the audio thread is not consuming events fast enough), the oldest event
is dropped to make room for the new one. This prevents unbounded memory growth
at the cost of occasional audio glitches under heavy load.

## 4.3. Thread Safety

All buffer operations (push, pop, peek, length query) must be protected by a
mutex. The CPU thread pushes events and the audio thread pops them
concurrently.

## 4.4. Activity Tracking

The buffer tracks the wall-clock time of the most recent push. The speaker is
considered "active" when either:

- The buffer is non-empty (events are pending), or
- The last push occurred within an activity timeout (300 ms).

This activity signal is used to inhibit full-speed mode (section 7) so that
audio playback is not disrupted by the emulator running faster than real time.

The timeout value (300 ms) must be longer than the longest delay a typical
Apple II sound routine inserts between toggles. The firmware WAIT routine with
A=$C0 produces a delay of roughly 180 ms, so 300 ms provides adequate margin.

# 5. Audio Stream

The audio stream converts toggle events from the speaker buffer into PCM
samples for the host audio system.

## 5.1. Output Format

The stream produces stereo interleaved 32-bit float PCM at 44,100 Hz. Each
sample frame consists of two identical float32 values (left and right channels
are the same, since the Apple II speaker is mono).

    Sample rate:       44,100 Hz
    Channels:          2 (stereo, mono-duplicated)
    Sample format:     float32, little-endian
    Bytes per frame:   8 (4 bytes x 2 channels)
    Buffer size:       1,024 sample frames (~23 ms)

The buffer size determines how many sample frames the host audio system
requests per callback. Smaller buffers reduce latency but increase the risk of
underruns (audible glitches). 1,024 frames at 44,100 Hz provides a good
balance.

## 5.2. Stream Interface

The audio stream implements `io.Reader`. The host audio system calls `Read(buf
[]byte)` to pull PCM data. Each call fills `buf` with as many complete sample
frames as fit (`len(buf) / 8`), and returns the number of bytes written. The
stream never returns an error.

## 5.3. Cycle-to-Sample Mapping

The Apple II CPU runs at 1,023,000 Hz. Each audio sample spans a fraction of
the CPU timeline:

    cycles_per_sample = cpu_clock_rate / sample_rate

At the standard clock rate, this is approximately 23.2 cycles per sample.

The stream maintains a `currentCycle` counter that advances by
`cycles_per_sample` for each output sample. This counter is synced to the
toggle event stream rather than to wall-clock time, which prevents drift
between the audio and emulation timelines.

## 5.4. Sample Generation

The stream maintains a `lastState` value representing the speaker position
carried forward from the most recently processed toggle event. This is the
state assumed to be in effect at the start of each sample interval until a new
event says otherwise.

For each output sample, the stream examines the CPU cycle interval
`[currentCycle, currentCycle + cycles_per_sample)` and computes how many
cycles within that interval the speaker spent in the high state vs. the low
state.

The sample value is a weighted average:

    sample = ((high_cycles - low_cycles) / total_cycles) * amplitude

This produces values in the range `[-amplitude, +amplitude]`. When the speaker
is fully high for the entire sample, the output is `+amplitude`. When fully
low, `-amplitude`. When the speaker toggles mid-sample, the value falls
somewhere in between, producing implicit anti-aliasing at toggle boundaries.

## 5.5. Silence

When the speaker buffer is empty and no events fall within the current sample
interval, the stream outputs silence (0.0) rather than holding the last
speaker state as a DC offset. Sound is produced by transitions, not by holding
a static position. Outputting a constant non-zero value when the speaker is
idle would produce an audible pop when sound begins or ends.

This means a single isolated toggle (with no subsequent toggle to return the
speaker to its original position) will not produce a sustained DC level in the
output. On real hardware, such a toggle would physically deflect the speaker
cone and hold it there, but the resulting constant offset is inaudible and
suppressing it avoids the pop artifacts described above.

## 5.6. Timeline Synchronization

The stream syncs its `currentCycle` to the event stream in two cases:

1. **Initial sync**: when `currentCycle` is 0 (stream has just started), it is
   set to the cycle of the first available event.
2. **Gap recovery**: when the next event's cycle is more than 1/10th of a
   second ahead of `currentCycle`, the stream jumps forward to the event's
   cycle. This handles cases where the emulator was paused, loading a disk, or
   otherwise not producing toggle events for an extended period.

In both cases, the stream resyncs to the event stream rather than generating
a burst of silence to "catch up", which would introduce audible latency.

# 6. Volume Control

The emulator provides volume control that scales the output amplitude.

## 6.1. Volume Level

Volume is stored as an integer percentage (0-100). The default is 50%.

When applying volume to audio samples, the percentage is converted to a float
(0.0-1.0) and further scaled by 0.5 to prevent the output from being
excessively loud:

    amplitude = (volume_level / 100.0) * 0.5

At the default volume of 50%, the effective amplitude is 0.25.

## 6.2. Mute

The mute state is independent of the volume level. When muted:

- The audio stream receives a volume of 0.0.
- The stored volume level is preserved so it can be restored on unmute.

This matches the behavior users expect from operating system volume controls:
muting and unmuting returns to the previous volume rather than resetting it.

## 6.3. Volume Adjustment

Volume is adjusted in steps of 10 percentage points. Adjustments are clamped
to the 0-100 range. When volume is decreased to 0%, the mute flag is
automatically set. When volume is increased from a muted state, the mute flag
is cleared and the new level takes effect.

# 7. Full-Speed Mode

The emulator may run faster than real time (e.g., during disk loading). In
full-speed mode, the CPU clock runs as fast as possible, and the relationship
between CPU cycles and wall-clock time breaks down.

When the audio stream detects full-speed mode:

1. All pending toggle events are consumed and discarded.
2. The stream outputs silence (0.0) for all samples.
3. The `currentCycle` counter is reset to 0, so that when full-speed mode
   ends, the stream will resync to the event stream (section 5.6).

This prevents the audio stream from trying to process a burst of thousands of
toggle events that accumulated during full-speed execution, which would
produce garbled noise.

The speaker's activity tracking (section 4.4) inhibits full-speed mode when
the speaker has been recently toggled, so that audio playback in normal
operation is not disrupted by a full-speed transition.

# 8. Interaction with Other Subsystems

## 8.1. CPU

The speaker soft switch is accessed during normal instruction execution. Each
access to $C030 is a standard memory-mapped I/O operation with no special
cycle cost beyond the instruction that performs it. The CPU provides the cycle
counter that timestamps each toggle event.

## 8.2. Graphics

Audio and graphics are independent subsystems. The audio stream runs on its
own thread and reads from the speaker buffer asynchronously. There is no
synchronization between audio and video frames.

## 8.3. Shortcuts

Volume up, volume down, and mute toggle are controlled through the emulator's
shortcut system (spec 7). These shortcuts adjust the volume state described in
section 6 and do not interact with the speaker toggle mechanism.
