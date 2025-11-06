# flf2bit

A command-line tool for converting
[FIGlet font files](https://github.com/cmatsuoka/figlet-fonts) (.flf) to a
JSON-based (.bit) used by ansifont and
[`bit`](<(https://github.com/superstarryeyes/bit)>) project.

## Description

`flf2bit` converts FIGlet font files (FLF) to the JSON-based .bit format that
can be used with the `bit` terminal font renderer. This allows you to convert
some FIGlet bitmap style fonts for use the with bit.

The initial version of this tool has been tested with the figlet C64-Fonts and
bdffonts. More fonts may work, but most FIGlet fonts use features not supported
by `bit`, so results may vary. Additional details on creating the C64 fonts and
BDF fonts can be found below.

## Installation

### Install from source

```bash
git clone https://github.com/scross01/flf2bit
cd flf2bit
go install
```

Or install directly using Go:

```bash
go install github.com/scross01/flf2bit@latest
```

## Usage

Basic usage:

```bash
flf2bit [options] <input.flf> <output.bit>
```

### Command-line Options

- `--name <name>`: Set the font name (default: extracted from FLF header)
- `--author <author>`: Set the author (default: extracted from FLF comments or
  "Converted from FLF")
- `--license <license>`: Set the license (default: "Converted font, check
  original license")

## Examples

Convert a FIGlet font file to the .bit format:

```bash
flf2bit example.flf example.bit
```

Convert with custom metadata:

```bash
flf2bit --name "Custom Font" --author "John Doe" --license "MIT" example.flf example.bit
```

## Adding fonts to `bit`

After creating the new .bit font you need to copy if to the `bit` fonts
directory and rebuild bit.

```bash
git clone https://github.com/superstarryeyes/bit
cd bit
go mod tidy
# copy the generated font file
cp <path_to>/example.bit fonts/
# rebuild bit
go build -o bit ./cmd/bit
```

## Creating C64 and BDF Fonts

The included `Makefile` demonstrates how to convert the C64-font and bdffonts
from the figlet-fonts repository to .bit format using `flf2bit`.

This will download the figlet-fonts repository and the bit repository, convert
the C64-Fonts and bdffonts, and place them in the `bit/ansifont/fonts`
directory, and rebuilt bit with the new fonts.

```bash
make c64fonts
make bdffonts
make install
```

Note that the Figlet C64-Fonts provided by Figlet where originally extracted
from Commodore 64 character set file and converted using Commodore2Figlet v1.00
by David Proper. Some characters are different in the original
[PETSCII](https://en.wikipedia.org/wiki/PETSCII) than in ASCII, certain
charactors will be different or missing. Not all fonts include both upper and
lower case letters, and some fonts transpose the the case.

## License

`flf2bit` is licensed under the MIT License.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on
GitHub.

## AI

This project was initially created with assistance from QWEN code.
