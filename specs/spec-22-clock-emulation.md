---
Specification: 22
Drafted At: 2026-03-31
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes how the emulator throttles CPU execution to match the
clock speed of the original hardware. Without clock emulation, the host
machine would execute instructions as fast as the host CPU allows, causing the
emulated system to run far faster than real hardware.

The clock emulator measures elapsed wall-clock time, calculates how many
cycles should have completed in that time, and runs instructions until the
cycle count catches up. When the host is faster than the target clock (the
normal case), the emulator naturally idles between bursts. When the host falls
behind, it processes extra cycles to catch up.

# 2. Configuration

## 2.1. Hertz

The clock emulator is initialized with a target clock rate in hertz (cycles
per second). For example, the Apple //e runs at 1,023,000 Hz.

From the hertz value, the emulator precomputes the wall-clock duration of a
single cycle:

    time_per_cycle = 1 second / hertz

This value is used throughout the process loop to convert elapsed time into a
target cycle count.

## 2.2. Changing the Clock Rate

The target clock rate can be changed at runtime. When the rate changes:

1. The hertz and time-per-cycle values are updated.
2. The timing state is reset (section 3.3) so the emulator does not think it
   is behind or ahead on cycles.

# 3. Process Loop

The process loop is the main execution loop. It runs indefinitely on its own
thread and does not return unless the emulated computer's process function
reports an error.

## 3.1. Outer Loop

Each iteration of the outer loop proceeds as follows:

1. Call the breakpoint check function (section 5.1).
2. If the debugger state is active, call the debugger entry function (section
   5.2), reset timing, and restart the loop.
3. If the paused state is active, sleep for 100 ms and restart the loop.
   When the emulator transitions from paused to unpaused, reset timing before
   continuing.
4. Read the current timing state: elapsed time since the last resume,
   time-per-cycle, total cycles executed, full-speed flag, and timing epoch.
5. Compute the target cycle count:

       wanted_cycles = elapsed / time_per_cycle

6. Enter the inner loop (section 3.2).

The debugger and paused states are read from the computer's state map, not
stored by the clock emulator itself. The clock emulator treats them as
external inputs.

## 3.2. Inner Loop

The inner loop calls the computer's process function repeatedly. Each call
returns the number of cycles consumed by the instruction that was executed.
The emulator accumulates these into a running total.

The inner loop continues while either:

- `total_cycles < wanted_cycles` (the emulator is behind and needs to catch
  up), or
- Full-speed mode is active (section 4).

After each instruction, the inner loop acquires the timing mutex to update
`total_cycles` and re-read the full-speed flag and timing epoch. It does not
hold the mutex while executing the instruction itself.

The inner loop breaks early in two cases:

- **Timing epoch change**: if the timing epoch (section 3.3) has changed since
  the inner loop began, another thread has reset the timing state. The inner
  loop breaks so the outer loop can recompute `wanted_cycles` with the new
  values.
- **Full-speed exit**: if the emulator was in full-speed mode at the start of
  the inner loop iteration but is no longer, timing is reset and the inner
  loop breaks. This prevents a burst of catch-up cycles after a period of
  unthrottled execution.

## 3.3. Timing Resets

Several events disrupt the relationship between wall-clock time and cycle
count. After any of these, the emulator must reset its timing state to avoid
thinking it is behind and running a burst of catch-up cycles:

- Exiting the debugger
- Transitioning from paused to unpaused
- Changing the clock rate
- Transitioning from full-speed to normal speed

A timing reset:

1. Sets the resume time to now.
2. Zeros the total cycle count.
3. Increments the timing epoch (a monotonic counter that the inner loop uses
   to detect concurrent resets).

# 4. Full-Speed Mode

Full-speed mode disables clock throttling. The emulator runs instructions as
fast as the host allows, ignoring the target cycle count. This is useful when
the emulated machine is performing a long operation where timing fidelity does
not matter, such as loading data from a disk.

Full-speed mode is set by external code. It is the responsibility of the
emulated computer's instructions to clear full-speed mode when the operation
that triggered it is complete.

When full-speed mode is cleared, the emulator resets timing (section 3.3)
before resuming normal throttled execution.

# 5. Debugger and Breakpoint Hooks

The clock emulator accepts two callback functions that integrate the debugger
with the process loop.

## 5.1. Breakpoint Check

A breakpoint check function is called at the start of every outer loop
iteration. Its job is to inspect the current machine state and, if a
breakpoint condition is met, set the debugger state so that the next iteration
enters the debugger.

## 5.2. Debugger Entry

A debugger entry function is called when the debugger state is active. This
function should block until the user exits the debugger. After it returns, the
emulator resets timing and resumes the outer loop.

# 6. Queries

The clock emulator exposes two read-only queries:

- **Full-speed status**: returns whether the emulator is currently in
  full-speed mode.
- **Time per cycle**: returns the precomputed duration of a single cycle
  (section 2.1).

# 7. Concurrency

The process loop runs on its own thread. External code may concurrently
change the clock rate, set full-speed mode, or set state flags (paused,
debugger). All timing-related fields -- hertz, time-per-cycle, resume time,
total cycles, full-speed flag, and timing epoch -- are protected by a mutex.

The inner loop acquires the mutex briefly to read and update timing fields,
then releases it. It does not hold the mutex while executing instructions.
The timing epoch allows the inner loop to detect concurrent resets without
holding the lock for the entire inner loop body.
