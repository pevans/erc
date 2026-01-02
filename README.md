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

## Running

To run Erc, you need to issue a command like `erc run <somefile>`, where
`<somefile>` is a valid disk image. Erc is validated to run with DOS3.3 disk
images file, which often have a file extension ending in `.dsk`. There are
other commands you can explore with `erc help`.

Many Apple II softwares have more than one disk. If so, you need to tell `erc`
about them when you first run the emulator. You can do this like so:

```
erc run disk1.dsk,disk2.dsk,disk3.dsk
```

This will load with `disk1.dsk` first. You can use a shortcut to swap disk1
for disk2; disk2 for disk3; disk3 for disk1; and so forth.

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
- **CTRL-A B: Start the debugger in the console where you ran Erc from.** From
  the debugger, type `help` to see a list of commands available there, or type
  `resume` to resume emulation and leave the debugger.
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
- **CTRL-A W: Toggle write-protect for the disk.** Write-protection will prevent
  data on the disk from being overwritten. Some software cannot run with
  write-protect enabled, and some software cannot run _without_ write-protect
  enabled -- exactly which is the case depends on the software. By default,
  write-protect is always off.

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

