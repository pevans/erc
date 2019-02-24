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

* The main command (erc) is located in the `cmd/erc` subdirectory.
* Machine-related code (for running an architecture) is located in
	`pkg/mach`.
* Processor-related code (for emulating processor chips) is located in
	`pkg/proc`.
* Code that reads static object data (fonts, system roms, etc.) is in
	`pkg/obj`.

## How can I run it?

At the moment, all you can or would need to do is run `go build
./cmd/erc`. This will build the _erc_ executable and you can run the
software from there.
