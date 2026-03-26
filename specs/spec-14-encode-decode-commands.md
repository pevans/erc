---
Specification: 14
Drafted At: 2026-03-26
Authors:
  - Peter Evans
---

# 1. Overview

This spec describes the `erc encode` and `erc decode` CLI subcommands. These
commands convert disk images between logical format (DOS 3.3 or ProDOS) and
physical format (6-and-2 nibblized). They are standalone utilities that
operate on files without booting the emulator. Though they are used in black
box testing for spec 13 (disk images), they are designed to work on their own
as utilities to work with logical and physical disk images.

# 2. Encode Command

## 2.1. Usage

```
erc encode [image] -o [output]
```

The encode command reads a logical disk image and writes a physically encoded
(6-and-2 nibblized) file.

## 2.2. Input

The input file must be a logical disk image in one of the following formats,
determined by file extension:

- `.dsk` -- DOS 3.3
- `.do` -- DOS 3.3
- `.po` -- ProDOS

The input file must be exactly 143,360 bytes. If the file is a different size,
the command fails with an error indicating the expected vs. actual size.

## 2.3. Rejected Input

- `.nib` files are rejected because they are already in physical (nibble)
  format. The error message indicates the file is already in nibble format.
- Files with unrecognized extensions are rejected.

## 2.4. Output

The `-o` (or `--output`) flag is required and specifies the output file path.
The output is a physically encoded file of exactly 223,440 bytes (35 tracks x
6,384 bytes per track).

On success, the command prints a message naming both the input and output paths
and exits with status 0.

# 3. Decode Command

## 3.1. Usage

```
erc decode [encoded-file] -o [output]
```

The decode command reads a physically encoded file and writes a logical disk
image.

## 3.2. Input

The input file must be a physically encoded disk image. Physical disk images
do not need to be exactly sized; nibble disk images found in the wild may have
variable sizes depending on the tool that produced them. The decode command
accepts any input that the underlying decoder can process.

The image type (DOS 3.3 or ProDOS) is determined by the **output** file's
extension, not the input file's extension.

## 3.3. Output

The `-o` (or `--output`) flag is required. The output file extension determines
the interleave used for decoding:

- `.dsk` -- DOS 3.3
- `.do` -- DOS 3.3
- `.po` -- ProDOS

The output file is exactly 143,360 bytes. `.nib` is rejected as an output
extension because decoding to nibble format is not supported.

On success, the command prints a message naming both the input and output paths
and exits with status 0.

## 3.4. Rejected Output Formats

- `.nib` output is rejected with an error.
- Unrecognized output extensions are rejected.

# 4. Error Handling

Both commands fail immediately and exit with a non-zero status on any error.
Error conditions include:

- Missing `-o` flag
- Unrecognized or unsupported file extension
- Input file does not exist or cannot be read
- Input file has the wrong size
- Encoding or decoding fails internally

# 5. Round-Trip Invariant

Encoding a logical image and then decoding the result with the same image type
must produce a file identical to the original input. That is:

```
erc encode input.dsk -o encoded
erc decode encoded -o output.dsk
diff input.dsk output.dsk  # no difference
```

The same holds for `.po` files. This property is the primary correctness
guarantee for both commands.

# 6. Relationship to Spec 13

The encode and decode commands expose the `a2enc.Encode` and `a2enc.Decode`
functions (specified in spec 13, section 9) as CLI tools. The encoding format,
sector interleaving, and physical layout are all defined in spec 13. This spec
covers only the CLI interface and its observable behavior.
