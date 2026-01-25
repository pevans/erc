# erc

The _erc_ (**e**mulator of **r**etro **c**omputers) system is designed to
emulate an Apple IIe (enhanced) computer and run software for it. It is
written to be flexible enough to handle multiple machine architectures, and
may be extended to do just that in the future.

Erc's goals are to:
1. Emulate old software in a way that feels as native and natural as can be
   achieved on modern hardware. It is not a goal to precisely emulate every
   facet of the older device.
2. Be written in a clear and straightforward manner, so that others may study
   and learn from how emulation works generally or how some specific older
   machine works in particular.

It's a rewrite of an older project (https://github.com/pevans/erc-c), by
the same name, that was written in C code. This project is instead
written in Go, which has much better tooling for code formatting and
testing--something that I felt was lacking in the previous iteration. I
also just really like Go! I regret nothing.

The soul of the work here, which is emulation of the Apple II, has been
a hobby of mine going back more than a decade.

## Installation

To install Erc, you must have at least Go 1.25 installed on your system.

If you have `just` installed, then you can run `just build`. If not, you can
run `go build -o erc .`. Either way, you'll be left with an executable named
`erc`, which you can install in a path you can execute from (like
`$HOME/bin`).

## Features

Erc can run most Apple II software successfully. It supports:

- Most visual modes, including:
  - Text mode (40 characters)
  - Low resolution graphics
  - High resolution graphics
  - Double high resolution graphics
  - Monochrome color graphics in (green and amber)
- Graphical shaders to simulate the output of a CRT monitor (a soft CRT shader
  is used by default)
- DOS 3.3 (.DSK, .DO) and Nibble (.NIB) disk images
- Basic speaker support
- Save states: load and save the state of your emulation at any time (up to 10
  state slots available)
- Accurate clock cycle emulation: run software at the normal speed of the
  computer, or run it faster (up to 5x the original clock speed).

Plus, lots of features for folks who want to peek under the hood at what the
emulation is doing:

- An inline debugger to peek at memory, change state, etc.
- A large variety of debug files to look at what's being written to disk, the
  assembly that's been executed, how visuals and audio are produced, and more.
  Available when running with the `--debug-image` flag.
- Encode logical disk images to physical images so you can see how those are
  structured with the `erc encode` command.
- Look at disk image metadata with the `erc info` command.

## Running

To run Erc, you need to issue a command like `erc run <somefile>`, where
`<somefile>` is a valid disk image. Erc is validated to run with DOS3.3 disk
images file, which often have a file extension ending in `.dsk`. There are
other commands you can explore with `erc help`.

Many Apple II softwares have more than one disk. If so, you need to tell `erc`
about them when you first run the emulator. You can do this like so:

```
erc run disk1.dsk disk2.dsk disk3.dsk
```

This will load with `disk1.dsk` first. You can use a shortcut to swap disk1
for disk2; disk2 for disk3; disk3 for disk1; and so forth.

Erc is designed to run with the set of disk images you tell it when first
executed. There's no load modal to let you load some other disk image later;
you can only swap disks from within that set.

(If you're truly adventurous, you can still load a disk from outside the
initial disk set if you invoke the inline debugger and use its `load`
command.)

## Floppy disks and disk images

Erc works with _disk images,_ which are files that are a byte-for-byte copy of
what was stored on a floppy disk. There are two kinds of disk images:

- Logical images. These files contain the bytes that represent software data,
  software code, etc. that would have been stored on a disk. Logical disk
  images often have extensions like `.dsk` and `.do`.
- Physical, or _nibblized,_ images. These contain the literal bytes _as they
  would have been written_ on the floppy disks. Apple IIs used an encoding
  scheme so that they could detect errors on disk data, and this scheme made
  the size of the data look different (and larger!) than if they had loaded
  the raw software code and data onto the disk. Physical disk images often
  have an extension like `.nib`.

When loading a logical disk image (which is the most common format), Erc must
encode that data so that it looks like it would have been written on a floppy
disk. After that, the Apple II's disk code would expect to _decode_ that data
into the correct software code to run!

## Save states

Erc can save the state of whatever software you're emulating. Think of this as
a quick save of your progress that you can instantly revert to whenever you
want. Loading up some slow software and you want to skip back to where you
left off? Save states are great for that.

Erc can keep up to 10 save states. It does so with _slots_, numbered 0-9. The
default state slot is 0. Any save-state you create will save to the current
state slot, and text will flash by on the screen to remind you what the slot
would be.

Any state you want to save or load can be done so with keyboard shortcuts.
Read below to learn what those shortcuts are.

## Sound

Erc can emulate sound. On an Apple II, the builtin speaker was very basic, and
the sounds that software might make (like a game) will sound tinny compared to
anything modern.

If you would prefer not to hear sound, you can toggle it off using the
**CTRL-A V** shortcut (see more on keyboard shortcuts below). You can also
adjust the volume up or down using other shortcuts.

## Monochrome

You can emulate software in a monochrome color by passing the CLI flag,
`--monochrome=x` where `x` is either `green` or `amber`. This feature is
designed to provide the same feeling of running under such a monitor.

## Shaders

Erc uses shaders to try provide different visual effects to the software being
emulated. Those shaders are:

- **softcrt.** The default shader. Adds a very light scanline effect to the
  graphics in an effort to mimic the look that graphics had on monitors of the
  time.
- **curvedcrt.** Like softcrt, this shader adds a light scanline effect, but
  also adds a kind of curvature to the edges of the screen. Use this shader if
  you want to pretend that you were plugged into a cheap television.
- **hardcrt.** This adds a harder scanline effect to the graphics. Oddly,
  "better" CRT screens had more pronounced scanlines than did cheaper
  displays.

You can use the CLI flag, `--shader=x`, where `x` is one of the above shaders,
to configure Erc to use that shader. See `erc run help` for more information
on what flags are available when running a disk image.

## Keyboard shortcuts

All keyboard shortcuts are a combination of keys. You must hit Control-A, then
some other key, which is the shortcut. You don't need to hit these keys all at
the same time.

When you hit Control-A, the screen will highlight the edge of the screen to
indicate that it's waiting for the next key. Once a shortcut is completed,
you'll see a graphic in the top left that will fade out to indicate what's
happening.

Shortcuts available:

- **CTRL-A Escape: Pause emulation until you hit Escape again.** Any other
  keypress will show the Pause icon to indicate it's still paused. You don't
  need to hit CTRL-A Escape to resume; just Escape.
- **CTRL-A 0-9: Set the current state slot to N, where N is some number key on
  your keyboard.** This will tell Erc what slot to use when loading and saving
  state. See more information in the Save State section of this file.
- **CTRL-A +: Increase the speed so emulation runs faster.** Up to a maximum
  of 5x the normal speed of emulation.
- **CTRL-A -: Decrease the speed of emulation so things run slower.** The
  emulator will not go any slower than 1x the normal speed.
- **CTRL-A ]: Increase the sound volume by 10%.** Up to 100%.
- **CTRL-A [: Decrease the sound volume by 10%.** Down to 0% (muted).
- **CTRL-A B: Start the debugger in the console where you ran Erc from.** From
  the debugger, type `help` to see a list of commands available there, or type
  `resume` to resume emulation and leave the debugger.
- **CTRL-A L: Load a saved state from the current slot into the emulator.**
  See more information in the Save State section of this file.
- **CTRL-A N: Swap the disk currently in the drive with the _next_ disk
  available.** If you're on disk 1, you'll swap to disk 2; if you're on
  disk 2, you'll swap to disk 3. And so forth until you're back to disk 1.
  Any changes made to the disk being swapped out will be written back to the
  image file.
- **CTRL-A P: Swap the disk currently in the drive with the _previous_ disk
  available.** If you're on disk 3, you'll swap to disk 2, and so forth. Any
  changes made to the disk being swapped out will be written back to the image
  file.
- **CTRL-A Q: Quit the emulator and save all changes to the disk image.** (If
  you're on a Mac, CMD-Q should also work.)
- **CTRL-A S: Save the state the emulator in the current state slot.** See
  more information in the Save State section of this file.
- **CTRL-A V: Toggle volume on or off.**
- **CTRL-A W: Toggle write-protect for the disk.** Write-protection will prevent
  data on the disk from being overwritten. Some software cannot run with
  write-protect enabled, and some software cannot run _without_ write-protect
  enabled -- exactly which is the case depends on the software. By default,
  write-protect is always off.

_Note: if you really need to send a literal control-A to the software
being emulated, you can do so by typing CTRL-A twice._

## What's in here?

- The main command (erc) is located in this directory (`erc.go`).
- Code for the MOS 65c02 processor is in `mos`.
- Code for the Apple II architecture is in `a2`.

Most other packages are there to support the work to run/render/etc. the
supporting elements that make emulation possible.

## Credits

- I am enormously grateful for a bunch of books:
  - Rodney Zaks for his book _Programming the 6502_;
  - Jim Sather for his work _Understanding the Apple II_;
  - Don Worth and Pieter Lechner's book, _Beneath Apple DOS_;
  - And Apple's fantastic own technical reference manuals, which were a huge
    help.
- To Linapple/AppleWin and everyone who contributed to those -- their work was
  a great resource for me while trying to figure out why things weren't
  working.
- To Steve Wozniak, whose Apple II is the device that introduced me to
  computers and to my love of computing generally. Apple could never be what
  it is today without Steve's dedication and craftsmenship.

