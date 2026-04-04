# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Added

- There's a help modal that shows all available shortcut commands (hit CTRL-A
  ? or CTRL-A H to see it).

### Fixed

- Added the missing shortcut to select save state slot 0 (CTRL-A 0), which is
  the default state slot. (So if you went away from slot 0, you wouldn't have
  been able to go back.)

## [0.2.0] - 2026-04-01

### Added

- We now support a caps lock state. Press CTRL-A c to toggle caps lock. When
  toggled on, all letters will be uppercased with key presses, regardless of
  whether shift was held down. You may enable caps lock at boot time with a
  command line flag, `--caps-lock`, for the `run` subcommand.
- Support for 80-column text mode
- Support for [MouseText](https://en.wikipedia.org/wiki/MouseText) characters.
  Now you can see the running man (if you wish).
- Support for flashing characters, which are text that alternate between
  normal video and inverse video.

### Fixed

- Lowercase letters are now used for non-shifted key presses. (Early versions
  of the Apple II only supported uppercase letters, but the Apple //e both
  supported lowercase and defaulted to it.)
- Keyboard latch is now returned for all soft switches in the $C000-$C00F
  range.
- Alternate character set now shows lowercase letters in the $60-$7F character
  code range (changed from the special character set).

### Removed

- The `info` subcommand has been removed. In my experience, it rarely worked
  as software either didn't support the standardized VTOC or would write
  changes that would invalidate the VTOC.

## [0.1.3] - 2026-03-28

### Added

- A new subcommand, `decode`, is available. You can run this command to decode
  a physically-formatted (nibble) disk image to a DOS 3.3 or ProDOS
  logically-formatted disk image.

### Fixed

- Keyboard data is returned on reads to soft switches in the 0x range of the
  $C0 page.
- The correct ROM's byte is returned (expansion or PC) when reading $CFFF.
  Previously, we were always returning the PC ROM byte.
- Expansion slot is always set to zero on machine reset
- Expansion slot is set on reads to $C100-$CFFF when SlotCX is true
- The speaker toggle is now be triggered on both reads and writes to its soft
  switch
- Bank 2 RAM is now enabled by default when the Apple II machine resets.
- The internal state variable that tracks consecutive read attempts to enable
  writing to bank 1 and bank 2 RAM is now reset to zero when the Apple II
  machine resets.
- Address modes no longer inadvertantly trigger soft switch reads
  when resolving for non-read instructions. (E.g., STA $C088 should not count
  as both a read and a write on $C088; it should only count as a write.) This
  was happening because address modes call the `Get()` method of a segment to
  fetch an effective value of an address regardless of what the instruction
  intends to do.

## [0.1.2] - 2026-03-24

### Added

- A headless mode has been added. This allows Erc to run without graphics or
  sound. In headless mode, Erc is able to watch and logs various computer
  states, memory addresses, registers and so forth. Headless mode is intended
  to support black box testing.
- A miniature assembler was built as a standalone command (erc-assembler)
  which is used to build one-off disk images for black box testing. The
  assembler compiles an input program and writes the program data into the
  first track of the disk. Anything more than that won't work, as the Apple's
  first stage boot loader will only read the first track into memory.
- Added a large cohort of black box tests to support the new headless mode;
  the MOS 65c02 CPU; keyboard shortcuts; the inline interactive debugger; and
  40-column TEXT mode. In all we have over 280 black box tests now.
- Support ZPI (zero page indirect) address mode for the MOS 65c02 CPU.

### Fixed

- Carry flag is now set if the decimal accumulator is not negative after a
  subtraction. The previous behavior (setting carry if the result was >= 0)
  wasn't correct in that context.
- Text page 2 ($0800-$0BFF) and display page 2 ($4000-$5FFF) are properly
  mapped for graphic display. Previously they had not been, which meant that
  updates to those areas might not cause a screen rerender. One way I've seen
  this problem is a static text screen where some characters are missing.

## [0.1.1] - 2026-03-04

### Added

- `--start-in-debugger` flag that tells erc to boot into the debugger
  immediately after starting.
- `runfor` debugger command to tell erc to execute for a given number of
  seconds before reentering the debugger prompt.

### Changed

- [Ebitengine](https://ebitengine.org/) was upgraded to 2.8
- Adopt new text-render API for system messages rendered on screen. (E.g. when
  the volume is changed, erc will print the new volume setting.) Replaces a
  deprecated text API from an earlier version of ebitengine.

### Removed

- The MCP server experiment ended and it was removed.

## [0.1.0] - 2026-01-25

Initial versioned release of Erc.

[0.2.0]: https://github.com/pevans/erc/compare/v0.1.3...v0.2.0
[0.1.3]: https://github.com/pevans/erc/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/pevans/erc/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/pevans/erc/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/pevans/erc/releases/tag/v0.1.0
