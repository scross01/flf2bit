# flf2bit

A command-line tool for converting
[FIGlet font files](https://github.com/cmatsuoka/figlet-fonts) (.flf) and [TOIlet font files]() (.tlf) to a
JSON-based (.bit) fonts used by [`bit`](<(https://github.com/superstarryeyes/bit)>) project.

## Description

`flf2bit` converts FIGlet and TOIlet font file to the JSON-based .bit format that
can be used with the `bit` terminal font renderer. This allows you to convert
FIGlet/TOIlet fonts for use the with bit.

Most FIGlet fonts should now work with flf2bit. The included Makefile shows how
to convert many the fonts from the figlet-fonts repository to .bit format.

## Installation

### Install from source

```bash
git clone https://github.com/scross01/flf2bit
cd flf2bit
go install .
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
- `--map-chars <chars>`: Map the first character to the second character during
  font conversion (can be used multiple times)
- `--debug [chars]`:  Enable debug output for all characters or specific characters

## Examples

Convert a FIGlet font file to the .bit format:

```bash
flf2bit example.flf example.bit
```

Convert a TOIlet font file to the .bit format:

```bash
flf2bit example.tlf example.bit
```


Convert with custom metadata:

```bash
flf2bit --name "Custom Font" --author "John Doe" --license "MIT" example.flf example.bit
```

Convert with character mapping (replaces # with █):

```bash
flf2bit --map-chars "#█" example.flf example.bit
```

Convert with multiple character mappings (replaces # with █ and . with space):

```bash
flf2bit --map-chars "#█" --map-chars ". " example.flf example.bit
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

## Converting figlet-fonts to .bit format

The included `Makefile` demonstrates how to convert many of the fonts from the
[figlet-fonts](https://github.com/cmatsuoka/figlet-fonts) repository to .bit
format using `flf2bit`.

This will download the figlet-fonts repository and the bit repository, convert
the C64-Fonts and bdffonts, and place them in the `bit/ansifont/fonts`
directory, and rebuilt bit with the new fonts.

```bash
make
make all-fonts
make install
```

Not all fonts include both upper and lower case letters, and some fonts
transpose the the case.

## License

`flf2bit` is licensed under the MIT License.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on
GitHub.

## AI

This project has been coded with assistance from QWEN code.
