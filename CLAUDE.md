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
- `gfx`: the basics necessary to produce graphics on screen
- `input`: provides the means to store abstract input events, like keypresses
- `memory`: code to store various kinds of memory in segments and to map
  behavior to certain addresses in segments via "soft switch" maps
- `mos`: provides the means to emulate an MOS 65C02 CPU chip
- `obj`: code that provides object storage for built-in ROM
- `render`: the basics necessary to run erc with ebiten
- `shortcut`: code to capture and interpret keyboard shortcuts (like quitting
  the program)
- `work`: a directory that is used for temporary storage of test files, disk
  images, and so forth; these files are never committed to the repository

## Conventions

- Comments should only be added if you believe the intent of the code isn't
  obvious and you must explain why the code is there. Comments that merely
  explain what the code does is not very valuable.
- Tests should effectively encode the behavior of the functions being tested.
  For example, if there was a factorial function, you don't need a test for
  every integer you could possibly pass in; you only need a test for every
  conditional behavior the factorial function may exhibit. A good and fast
  overview of testing is Kent Beck's article, [Test
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
