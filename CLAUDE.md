# Working with erc, an emulator of retro computers

This repository contains code for a program named "erc", which is an acronym
for "emulator of retro computers". It currently contains code to emulate An
Apple II computer, specifically the Apple //e "enhanced" model.

Erc is a CLI program that is written in Go. For its graphics, it makes use of
the 2D graphics library, ebiten.

## Commands

- `just lint`: run a linter on every source code file

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
- erc uses cobra to provide subcommands. The `run` command is used to run an
  image file.

## Debug files

- When using the `just debug` command, you can produce a large number of
  files useful for debugging issues or learning more about what the image
  being emulated is doing. These files conventionally are named for the disk
  image, but have an additional extension added at the end of the filename.
  Here's an example: `image.dsk.ext`.
