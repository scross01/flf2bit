# AGENTS.MD

## Project Overview

### flf2bit Project
The `flf2bit` project is a Go-based utility for converting [FIGlet font files](https://github.com/cmatsuoka/figlet-fonts) (.flf) to a custom JSON-based format (.bit) used by the [Bit](https://github.com/superstarryeyes/bit) terminal logo designer. The project serves as a converter tool that transforms FIGlet font files into a structured JSON format that can be used by the Bit application.

**Key Files:**
- `main.go`: Main application file containing the conversion logic
- `go.mod`: Go module file defining the project dependencies
- `Makefile`: Contains build and conversion commands

**Project Purpose:**
The flf2bit tool converts FIGlet font files (.flf) to .bit files (JSON format) with the following structure:
- Each character is represented as an array of strings
- Each font contains metadata (name, author, license) and character data
- Characters are stored as ASCII art using "█" for solid blocks
- Supports command-line options for setting font name, author, and license

**Usage:**
```
flf2bit [options] <input.flf> <output.bit>
```
- `--name`: Set the font name
- `--author`: Set the author
- `--license`: Set the license

### Bit Terminal ANSI Logo Designer
The `bit` project https://github.com/superstarryeyes/bit is a comprehensive terminal ANSI logo designer and font library built in Go. It serves as the main application that consumes the .bit font files created by the flf2bit tool.

**Key Components:**
- Interactive TUI (Tea) application for designing ANSI logos
- Standalone Go library (`ansifonts`) for ANSI text rendering
- 100+ built-in bitmap fonts in .bit format
- Support for gradients, shadows, scaling, and alignment
- Multi-format export (TXT, Go, JavaScript, Python, Rust, Bash)

**Key Directories:**
- `ansifonts/`: Core library for font loading and rendering
- `cmd/bit/`: Interactive CLI application
- `cmd/ansifonts/`: Command-line tool for text rendering
- `images/`: Project assets including icon and screenshots

**Building and Running:**
```
# Build the interactive UI
go build -o bit ./cmd/bit

# Build the command line tool
go build -o ansifonts-cli ./cmd/ansifonts

# Run interactive UI
./bit

# Use command line tool
./ansifonts-cli -font dogica -color 32 "Hello"
```

## Relationship Between Projects
The `flf2bit` project serves as a converter tool that enables the `bit` project to use FIGlet fonts by converting them to the JSON-based .bit format that `bit` natively supports. This allows the Bit application to leverage the large collection of existing FIGlet fonts.

## Development Conventions
Both projects follow Go conventions:
- Standard project structure with go.mod
- CLI applications with proper error handling
- Structured logging and user-friendly error messages
- Modular code organization

## Building and Running flf2bit
```
# Build
go build .

# Convert FIGlet font to .bit format
go run main.go input.flf output.bit

# Using Makefile to work with C64 fonts
make flf2bit
```

## Key Technologies
- **Go 1.25+**: Both projects are written in Go
- **TUI Framework**: Bit uses Bubble Tea framework for interactive UI
- **JSON**: Font data is stored in JSON format
- **ANSI Terminal**: Both projects work with ANSI terminal output
