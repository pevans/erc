# erc

The _erc_ (**e**mulator of **r**etro **c**omputers) system is designed to
emulate an Apple IIe (enhanced) computer and run software for it. It is
written to be flexible enough to handle multiple machine architectures,
and may be extended to do just that in the future.

It's a rewrite of an older project (https://github.com/pevans/erc-c), by
the same name, that was written in C code. This project is instead
written in Go, which has much better tooling for code formatting and
testing--something that I felt was lacking in the previous iteration. I
also just really like Go! I regret nothing.

The soul of the work here, which is emulation of the Apple II, has been
a hobby of mine going back more than a decade. I am enormously grateful
to Rodney Zaks for his book _Programming the 6502_, Jim Sather for his
work _Understanding the Apple II_, Apple--beyond building the machine,
they published a wealth of technical reference material for it--and, of
course, Steve Wozniak, without whom the Apple II and Apple-as-we-know-it
would not exist.

## What's in here?

- The main command (erc) is located in this directory (`erc.go`).
- Code for the MOS 65c02 processor is in `mos`.
- Code for the Apple II architecture is in `a2`.

## What can I do here?

Erc uses `just` as a command runner. With it, you can:

- Build the executable with `just build`
- Test with images by running `just run` or `just debug`. Useful if you're poking around in the source code and want to see what changes.
