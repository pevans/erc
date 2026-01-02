# Working with erc, an emulator of retro computers

This repository contains code for a program named "erc", which is an acronym
for "emulator of retro computers". It currently contains code to emulate An
Apple II computer, specifically the Apple //e "enhanced" model.

Erc is a CLI program that is written in Go. For its graphics, it makes use of
the 2D graphics library, ebiten.

Erc's goals are to:

- Reasonably and performantly emulate old computers, but not precisely
  recreate everything a user might have experienced with the physical machine.
  For example, precise cycle-time emulation is not a goal. Being able to run
  the software is a goal.
- Be a readable and searchable code base that others interested in emulation
  may read to learn about emulation generally or older computers specifically.
  Give users the ability to look under the hood and learn more about what's
  happening and, where possible, why it's happening.

## Commands

- `just test`: run every test in every package
- `just test ./some/package`: run every test in one package
- `just lint`: run a linter on every source code file
- `just debug path/to/image/file`: launch erc with the --debug-image flag to
  debug an image file's execution
- `just build`: produce an executable of the software in the root directory,
  which is named `erc`.

## Structure

- `a2`: provides the means to emulate an Apple II computer
- `asm`: code to record and print lines of assembly code -- useful for
  debugging
- `clock`: allows us to emulate the clockspeed of an arbitrary machine
- `cmd`: contains several subcommands to use the various functionality of erc
- `debug`: a debugger that someone can interactively use to investigate the
  state and behavior of the computer
- `emu`: provides an abstract definition of an emulated computer
- `gfx`: the basics necessary to produce graphics on screen, including status
  overlays that fade in/out
- `input`: provides the means to store abstract input events, like keypresses
- `memory`: code to store various kinds of memory in segments and to map
  behavior to certain addresses in segments via "soft switch" maps
- `mos`: provides the means to emulate an MOS 65C02 CPU chip
- `obj`: code that provides object storage for built-in ROM and embedded PNG
  graphics for status overlays
- `render`: the basics necessary to run erc with ebiten
- `shortcut`: code to capture and interpret keyboard shortcuts using a CTRL-A
  prefix (e.g., CTRL-A Q to quit, CTRL-A ESC to pause)
- `work`: a directory that is used for temporary storage of test files, disk
  images, and so forth; these files are never committed to the repository

## Conventions

- Comments should only be added if you believe the intent of the code isn't
  obvious and you must explain why the code is there. Comments that merely
  explain what the code does is not very valuable.
- Test logic, not data. If there is a function with four conditional paths
  that could be taken, you should write 4 test cases for that function. You
  don't need to test every possible data input that may be provided. Focus on
  boundaries and edge cases. A good and fast overview of testing is Kent
  Beck's article, [Test
  Desiderata](https://medium.com/@kentbeck_7670/test-desiderata-94150638a4b3).
- You should use table-driven tests if you have 3 or more behaviors you want
  to test; otherwise, you can use sub-tests without needing a table structure
  to encode all of the inputs for the test.
- There is no need to write a unit test on private, or unexported, code. These
  functions will be tested by the public, exported, code that uses them.
- Think of code organization as if it were provision. Packages provide an API
  that may be used elsewhere. Files provide a collection of related parts of
  the package API.
- erc uses cobra to provide subcommands. The `run` command is used to run an
  image file.

## Debug files

- When using the `just debug` command, you can produce a large number of
  files useful for debugging issues or learning more about what the image
  being emulated is doing. These files conventionally are named for the disk
  image, but have an additional extension added at the end of the filename.
  Here's an example: `image.dsk.ext`.
- An instruction log is a file that looks like `image.dsk.asm`. It contains a
  listing of instructions that were executed, sorted by the address in memory
  where they were executed from. Each line has the instruction's memory
  address, opcode byte, optionally any operand bytes, the instruction name,
  the operand formatted with respect to the instruction's address mode, and
  optionally comments about the code being run. Note that in MOS 6502
  assembly, a comment is anything that follows the operand; some assembly may
  use a comment character like `;`, but not all assemblers required it.
- A disk log is a file that looks like `image.dsk.disklog`. It contains an
  accounting of what bytes were read, and at what positions, on the disk. It
  also shows the instruction that was responsible for reading the byte, along
  with the address in memory. It is important to note that this will show
  bytes that are _encoded_ for the Apple II, and there will not be a 1:1
  correlation to bytes in the disk image that was loaded. For that, you would
  need to look at `image.dsk.physical`, which is a record of the encoded disk
  image.
- A timeset is a file that looks like `image.dsk.time`. It shows how much
  time is spent in total for each instruction as it has been executed, and how
  many times it has been executed. You can use this to build a model of where
  the "hot paths" are in a program, and potentially where there may be a bug.
  It is important to note that some software intentionally has written some
  loops that effectively waste time, because the software may be trying to
  acquire a certain timing based on disk spin. Because erc doesn't emulate
  disk spin, we emulate instructions at full speed when the disk is spinning
  so that those loops go by quickly.
- A metrics file looks like `image.dsk.metrics`. It shows an accounting of the
  total number certain events have been observed while the disk image was
  running; for example, the number of bytes read.
