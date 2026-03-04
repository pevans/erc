# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.1.1]: https://github.com/pevans/erc/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/pevans/erc/releases/tag/v0.1.0
