package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// FontData represents the structure of the .bit font file
type FontData struct {
	Name       string              `json:"name"`
	Author     string              `json:"author"`
	License    string              `json:"license"`
	Characters map[string][]string `json:"characters"`
}

// processCharacterLine processes a character line by replacing the hardblank character with space and # with █
func processCharacterLine(line string, hardblankChar string) string {
	processedLine := strings.ReplaceAll(line, hardblankChar, " ")  // Replace hardblank with space
	processedLine = strings.ReplaceAll(processedLine, "#", "█")    // Replace # with █
	return processedLine
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: flf2bit [options] <input.flf> <output.bit>")
		fmt.Println("Options:")
		fmt.Println("  --name <name>     Set the font name (default: extracted from FLF)")
		fmt.Println("  --author <author> Set the author (default: extracted from FLF)")
		fmt.Println("  --license <license> Set the license (default: 'Converted font, check original license')")
		fmt.Println("Example: flf2bit -name \"Custom Font\" -author \"John Doe\" example.flf example.bit")
		os.Exit(1)
	}

	// Parse command line options
	name := ""
	author := ""
	license := ""
	args := os.Args[1:]
	var inputFile, outputFile string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--name" {
			if i+1 < len(args) {
				name = args[i+1]
				i++ // Skip next argument since it's the value
			} else {
				fmt.Println("Error: --name requires a value")
				os.Exit(1)
			}
		} else if arg == "--author" {
			if i+1 < len(args) {
				author = args[i+1]
				i++ // Skip next argument since it's the value
			} else {
				fmt.Println("Error: --author requires a value")
				os.Exit(1)
			}
		} else if arg == "--license" {
			if i+1 < len(args) {
				license = args[i+1]
				i++ // Skip next argument since it's the value
			} else {
				fmt.Println("Error: --license requires a value")
				os.Exit(1)
			}
		} else if inputFile == "" {
			inputFile = arg
		} else if outputFile == "" {
			outputFile = arg
		}
	}

	if inputFile == "" || outputFile == "" {
		fmt.Println("Usage: flf2bit [options] <input.flf> <output.bit>")
		fmt.Println("Example: flf2bit example.flf example.bit")
		os.Exit(1)
	}

	font, err := convertFLFToBit(inputFile, name, author, license)
	if err != nil {
		fmt.Printf("Error converting font: %v\n", err)
		os.Exit(1)
	}

	err = saveFontData(font, outputFile)
	if err != nil {
		fmt.Printf("Error saving font: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}

func convertFLFToBit(inputFile string, name string, author string, license string) (*FontData, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read the header line
	if !scanner.Scan() {
		return nil, fmt.Errorf("could not read header line")
	}

	header := scanner.Text()
	parts := strings.Split(header, " ")
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid FLF header format")
	}

	// Extract the hardblank character from the 6th character of the header
	var hardblankChar string = "$"
	if len(header) >= 6 {
		hardblankChar = string(header[5]) // 6th character is the hardblank
	}

	// Extract comment line count from the header (6th field in the space-separated parts)
	commentLineCount := 0
	if len(parts) >= 6 {
		// The 6th field contains the comment line count
		fmt.Sscanf(parts[5], "%d", &commentLineCount)
	}

	// Skip the exact number of comment lines specified in the header
	commentLines := []string{}
	for i := 0; i < commentLineCount; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected end of file while reading comment lines")
		}
		line := scanner.Text()
		commentLines = append(commentLines, line)
	}

	// Use provided values or extract from FLF file if not provided
	fontName := name
	fontAuthor := author
	fontLicense := license

	// Only extract from file if not provided via command line
	if fontName == "" {
		fontName = "Converted Font"
		if strings.Contains(header, "\"") {
			quoteSplit := strings.Split(header, "\"")
			if len(quoteSplit) >= 2 {
				fontName = quoteSplit[1]
			}
		}
	}

	if fontAuthor == "" {
		fontAuthor = "Converted from FLF"
		// Look for author information in comments
		for _, line := range commentLines {
			if strings.Contains(line, "by") {
				fontAuthor = line
				break
			}
		}
	}

	if fontLicense == "" {
		fontLicense = "Converted font, check original license"
	}

	// First, find the line end character by looking at the first non-comment line
	var lineEndChar string = "@"
	var firstDataLine string = ""
	
	// Read the first line after comments to determine the line end character
	if scanner.Scan() {
		firstDataLine = scanner.Text()
		
		// If it's an empty line, continue reading until we find a non-empty line
		for firstDataLine == "" && scanner.Scan() {
			firstDataLine = scanner.Text()
		}
		
		// Determine the line end character from the last character of the first non-comment line
		if len(firstDataLine) > 0 {
			// Trim right whitespace to get the actual last character
			trimmedLine := strings.TrimRight(firstDataLine, " \t\r\n")
			if len(trimmedLine) > 0 {
				lineEndChar = string(trimmedLine[len(trimmedLine)-1])
			}
		}
	}
	
	// Create the end-of-character delimiter (two of the line end characters)
	endOfCharDelimiter := lineEndChar + lineEndChar

	// Process characters
	characters := make(map[string][]string)

	// Process the first character line
	// currentLine is firstDataLine
	currentLine := firstDataLine

	// Process characters until EOF
	// Count all characters in file order, then map to appropriate ASCII codes
	charIndex := 0 // Start at 0, will adjust to ASCII later
	inCharacter := false
	var currentCharLines []string

	// Process the first character line
	for {
		line := currentLine

		// If line contains the end-of-character delimiter, it's the end of a character
		if strings.Contains(line, endOfCharDelimiter) {
			parts := strings.Split(line, endOfCharDelimiter)

			// Process the first part (end of current character)
			if parts[0] != "" {
				charParts := strings.Split(parts[0], lineEndChar)
				for _, part := range charParts {
					processedPart := processCharacterLine(part, hardblankChar)
					currentCharLines = append(currentCharLines, processedPart)
				}
			}

			// Add completed character
			if len(currentCharLines) > 0 {
				// Process each line to replace hardblank with space and # with █
				processedLines := make([]string, len(currentCharLines))
				for i, line := range currentCharLines {
					processedLines[i] = processCharacterLine(line, hardblankChar)
				}

				// Calculate the ASCII value for this character position
				asciiValue := charIndex + 32

				// Only add characters in the standard ASCII range (space to tilde, 32-126)
				if asciiValue >= 32 && asciiValue <= 126 {
					char := string(rune(asciiValue))
					characters[char] = processedLines
				}
				charIndex++
				currentCharLines = []string{}
				inCharacter = false
			}

			// Process the second part if it exists (start of next character)
			if len(parts) > 1 && parts[1] != "" {
				nextParts := strings.Split(parts[1], lineEndChar)
				for _, part := range nextParts {
					processedPart := processCharacterLine(part, hardblankChar)
					currentCharLines = append(currentCharLines, processedPart)
					inCharacter = true
				}
			}
		} else if strings.HasSuffix(line, lineEndChar) {
			// This is a line of the current character
			// Split by the line end character to get the actual character data
			parts := strings.Split(line, lineEndChar)
			// The last element after splitting is usually empty, so we process all but the last
			for i := 0; i < len(parts)-1; i++ {
				processedPart := processCharacterLine(parts[i], hardblankChar)
				currentCharLines = append(currentCharLines, processedPart)
			}
			inCharacter = true
		} else if line == "" {
			// Empty line after character data may indicate end of character
			if inCharacter && len(currentCharLines) > 0 {
				processedLines := make([]string, len(currentCharLines))
				for i, line := range currentCharLines {
					processedLines[i] = processCharacterLine(line, hardblankChar)
				}

				// Calculate the ASCII value for this character position
				asciiValue := charIndex + 32

				// Only add characters in the standard ASCII range (space to tilde, 32-126)
				if asciiValue >= 32 && asciiValue <= 126 {
					char := string(rune(asciiValue))
					characters[char] = processedLines
				}
				charIndex++
				currentCharLines = []string{}
				inCharacter = false
			}
		}

		// Read next line
		if !scanner.Scan() {
			// End of file
			if len(currentCharLines) > 0 {
				// Process each line to replace hardblank with space and # with █
				processedLines := make([]string, len(currentCharLines))
				for i, line := range currentCharLines {
					processedLines[i] = processCharacterLine(line, hardblankChar)
				}

				// Calculate the ASCII value for this character position
				asciiValue := charIndex + 32

				// Only add the last character if it's in the standard ASCII range (space to tilde, 32-126)
				if asciiValue >= 32 && asciiValue <= 126 {
					// Add the last character if any data remains
					char := string(rune(asciiValue))
					characters[char] = processedLines
				}
			}
			break
		}

		currentLine = scanner.Text()
	}

	fontData := &FontData{
		Name:       fontName,
		Author:     fontAuthor,
		License:    fontLicense,
		Characters: characters,
	}

	return fontData, nil
}

func saveFontData(font *FontData, filename string) error {
	data, err := json.MarshalIndent(font, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
