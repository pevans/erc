This file is a list of ideas for erc's development. It is intended to capture only the present and near future scope of work that can be done. This file is used to produce a roadmap of work.

## Vision

Erc is an acronym which stands for "emulator of retro computers", and its long term aim is to do that. It is designed straightforwardly so that people can read the source code and learn how old computers work.

The near-term vision is to reasonably emulate an Apple //e, which is an Apple IIe "enhanced" computer.

## Current state

Erc is able to emulate Apple II DOS-order software, although not well. It can:

- read and execute software
- render high resolution graphics and 40-column text
- reasonably emulate cycles (although precise cycle emulation is not a goal)
- at any point drop into a debugger and record a log of disk access
- produce reasonably complete debugging output to examine what software does

## Problems

- Some software takes a surprisingly long time to load (e.g. Ultima III)
- Some software does not load correctly
- A lot of software has not been tested at all (WOZ, NIB, and PO files particularly)
- Low-resolution graphics have not been well tested

## Opportunities

- 80-column text
- Basic sound support
- Mockingboard sound support
- State file support
- Mouse support
