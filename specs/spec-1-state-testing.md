---
Specification: 1
Drafted At: 2026-03-04
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes a system for recording state transitions during emulator
execution. A "state transition" is any change to a singular piece of
observable state -- a byte in memory, a CPU register, a flag, or a piece of
emulator state (e.g. whether bank memory is active for a given block). These
transitions are recorded at each execution step and can be compared against
expected results to enable black box testing of emulated programs.

# 2. Concepts

## 2.1. State Entry

A state entry is a record of a single value changing from one state to
another. It contains:

- **Step** -- the instruction execution count at which the change occurred
  (e.g. the 514th instruction since boot)
- **Tag** -- a short label identifying the category of state (e.g. `mem`,
  `reg`, `comp`)
- **Name** -- an identifier for the specific item that changed (e.g. `$013F`
  for a memory address, `A` for the accumulator, `bank-df-ram` for a computer
  state flag)
- **Old value** -- the value before the change
- **New value** -- the value after the change

## 2.2. Recorder

A recorder is the component responsible for collecting state entries. It is
attached to the emulator before execution begins and accumulates entries as
instructions execute. The recorder itself holds no knowledge of what is being
recorded -- it simply accepts state entries from any source.

## 2.3. Observer

An observer is a function or component that watches a specific piece of state
and reports changes to the recorder. Observers are registered with the
recorder before execution begins. Each observer is responsible for exactly one
piece of state -- one memory address, one register, one flag, etc.

# 3. State Entry Format

When rendered as text, a state entry uses the following format:

```
step <step>: <tag> <name>: <old> -> <new>
```

Examples:

```
step 514: mem $013F: $3A -> $3C
step 515: reg A: $00 -> $3C
step 515: reg P: $30 -> $B0
step 518: comp bank-df-ram: true -> false
```

Values are formatted according to their type:

- Byte values are printed as two-digit hex with a `$` prefix (e.g. `$3A`)
- 16-bit values are printed as four-digit hex with a `$` prefix (e.g. `$013F`)
- Boolean values are printed as `true` or `false`

# 4. Observer Types

## 4.1. Memory Observer

Watches a single memory address. Before each step, the observer reads the
current value at the address. After the step executes, it reads the value
again. If the value differs, it emits a state entry with tag `mem` and the
address as the name.

Multiple memory observers can be registered to watch a range of addresses.
Each observer is independent and watches exactly one address.

## 4.2. Register Observer

Watches a single CPU register (A, X, Y, S, P, or PC). Before each step, the
observer captures the register value. After the step, it compares. If the
value differs, it emits a state entry with tag `reg` and the register name.

## 4.3. Computer State Observer

Watches a single entry in the emulator's state map. Before each step, the
observer reads the value. After the step, it reads the value again. If the
value differs, it emits a state entry with tag `comp` and a descriptive name
for the state key.

# 5. Recording Lifecycle

## 5.1. Setup

1. Create a recorder.
2. Register observers for every piece of state that should be tracked. This
   can be done by the test harness -- the emulator itself does not need to
   know which observers exist.
3. Boot the emulator normally.

## 5.2. Per-Step Recording

For each execution step:

1. All observers capture their "before" snapshot.
2. The emulator executes one instruction.
3. The step counter increments.
4. All observers capture their "after" snapshot and compare against the
   "before" value. Any differences are emitted as state entries to the
   recorder.

## 5.3. Retrieval

After execution completes (or at any point during), the recorder provides
access to the full ordered list of state entries. This list can be:

- Compared against an expected list for test assertions
- Serialized to text for human inspection
- Filtered by tag, name, step range, or other criteria

# 6. Testing Model

## 6.1. Black Box Test Structure

A black box test consists of:

1. A program or disk image to load into the emulator
2. A set of observers to register (what to watch)
3. A number of steps to execute
4. An expected list of state entries

The test loads the program, registers the observers, runs the emulator for the
specified number of steps, and then compares the recorded state entries
against the expected list.

## 6.2. Comparison

Two state entry lists are compared entry-by-entry. A test passes if and only
if every entry matches in step, tag, name, old value, and new value. Missing
or extra entries cause a failure. The comparison reports the first mismatch to
aid debugging.

## 6.3. Test-Only Concern

The recorder and observers are only instantiated during testing. Production
execution paths do not create a recorder and thus incur no overhead. The
observer mechanism should not require changes to the core emulation loop --
observers wrap the existing execution call rather than being embedded within
it.

# 7. Design Considerations

## 7.1. Generality

The recorder accepts any state entry regardless of source. New observer types
can be added without modifying the recorder. This makes the system open to
extension -- if future emulator features introduce new kinds of state, new
observers can be written to track them.

## 7.2. Granularity

Each observer tracks exactly one item. To watch a range of memory (e.g.
`$0400`--`$07FF`), the test harness registers one observer per address. This
keeps each observer trivial and avoids complex multi-address logic.

## 7.3. Performance

Since this system is only active during testing, performance is secondary to
correctness. However, the before/after snapshot approach avoids instrumenting
the memory or register write paths, which means the core emulation code
remains unchanged and unaware of observation.

## 7.4. Step Counting

The step counter is maintained by the recorder, not the CPU's cycle counter.
Steps count instructions executed, not clock cycles. This provides a stable,
deterministic reference point regardless of instruction timing.
