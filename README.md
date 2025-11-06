# flf2bit

A command-line tool for converting [FIGlet font files](https://github.com/cmatsuoka/figlet-fonts) (.flf) to a JSON-based ANSI fonts (.bit) used by the [`bit`]((https://github.com/superstarryeyes/bit)) project.

## Description

flf2bit converts FIGlet font files (FLF) to a custom JSON-based format that can be used with the `bit` terminal font renderer. This allows you to use some FIGlet fonts that .

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
- `--author <author>`: Set the author (default: extracted from FLF comments or "Converted from FLF")
- `--license <license>`: Set the license (default: "Converted font, check original license")

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

After creating the new .bit font you need to copy if to the `bit` fonts directory and rebuild bit.

```bash
git clone https://github.com/superstarryeyes/bit
cd bit
go mod tidy
# copy the generated font file
cp <path_to>/example.bit fonts/
# rebuild bit
go build -o bit ./cmd/bit
```

## License

`flf2bit` is licensed under the MIT License.

## AI

This project was initially created with assistance from QWEN code.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.
