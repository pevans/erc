---
Specification: 26
Category: Storage
Drafted At: 2026-04-04
Authors:
  - Peter Evans
---

# 1. Overview

A disk set is an ordered collection of disk images that the user can navigate
at runtime. Many classic software packages shipped on multiple disks, so the
disk set lets the user load all of them at once and swap between them with
keyboard shortcuts.

# 2. Construction

## 2.1. Appending Images

Each filename passed on the command line is appended to the disk set in order.
The file must exist at append time; a missing file produces an error before
the emulator starts.

## 2.2. Initial Load

After all images have been appended, the emulator loads the first disk (index
0) into drive 1. This is the disk that is active when the emulator starts.

# 3. Navigation

The keyboard shortcuts for disk navigation are defined in spec 24.

## 3.1. Next

Advancing to the next disk increments the current index by one. If the index
was already at the last disk, it wraps to index 0.

## 3.2. Previous

Retreating to the previous disk decrements the current index by one. If the
index was already at 0, it wraps to the last disk in the set.

## 3.3. Index Tracking

The current disk index is observable as the `DiskIndex` state key. It starts
at 0 after the initial load and updates immediately when a next or previous
operation completes.

# 4. Drive Integration

When a disk swap occurs, the emulator:

1. Saves any pending writes from the current drive image.
2. Loads the new disk image into drive 1.

The drive's mechanical state (motor, track position, etc.) is not reset by a
disk swap -- only the image data changes.

# 5. Single-Disk Behavior

When the disk set contains only one image, next and previous are no-ops. No
load occurs and the index remains 0.

# 6. Test Structure

Tests use the assembler to produce multiple minimal disk images, then invoke
headless mode with `--keys` to trigger disk navigation shortcuts and
`--watch-comp DiskIndex` to observe index changes. This avoids dependence on
external disk files that may not be present.
