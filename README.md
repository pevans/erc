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
- Code for the MOS 65c02 processor is in `mos65c02`.
- Code for the Apple II architecture is in `a2`.
- There's random other code for things generally within other subdirectories.

## What can I do here?

You can:

- Build the executable `erc` by running `go build`.
- Run various tests, including:
  - `task lint` for linting (requires [Task](https://taskfile.dev/))
  - `go test ./...` for unit tests
