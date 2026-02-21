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

// CharacterMap represents a mapping from one character to another
type CharacterMap map[rune]rune

// processCharacter processes a complete character's lines applying hardblank replacement and character mappings
func processCharacter(charLines []string, hardblankChar string, charMap CharacterMap) []string {
	processedLines := make([]string, len(charLines))
	for i, line := range charLines {
		// Process the line by replacing hardblank with space and applying character mappings
		processedLine := strings.ReplaceAll(line, hardblankChar, " ")

		// Apply character mappings
		runes := []rune(processedLine)
		for j, r := range runes {
			if replacement, exists := charMap[r]; exists {
				runes[j] = replacement
			}
		}
		processedLines[i] = string(runes)
	}
	return processedLines
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: flf2bit [options] <input.flf> <output.bit>")
		fmt.Println("Options:")
		fmt.Println("  --name <name>     Set the font name (default: extracted from FLF)")
		fmt.Println("  --author <author> Set the author (default: extracted from FLF)")
		fmt.Println("  --license <license> Set the license (default: 'Converted font, check original license')")
		fmt.Println("  --map-chars <chars> Map first character to second character (can be used multiple times)")
		fmt.Println("  --debug [chars]   Enable debug output for all characters or specific characters")
		fmt.Println("Example: flf2bit --name \"Custom Font\" --author \"John Doe\" --map-chars \"#█\" example.flf example.bit")
		fmt.Println("Example: flf2bit --debug A B C example.flf example.bit")
		os.Exit(1)
	}

	// Parse command line options
	name := ""
	author := ""
	license := ""
	debugEnabled := false
	debugChars := make(map[rune]bool)
	charMaps := []string{} // Store character mapping pairs
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
		} else if arg == "--map-chars" {
			if i+1 < len(args) {
				charMaps = append(charMaps, args[i+1])
				i++ // Skip next argument since it's the value
			} else {
				fmt.Println("Error: --map-chars requires a value")
				os.Exit(1)
			}
		} else if arg == "--debug" {
			debugEnabled = true
			// Check if there are more arguments after --debug
			i++
			for i < len(args) && !strings.HasPrefix(args[i], "--") {
				// Add each character to the debug set
				if len(args[i]) == 1 {
					debugChars[rune(args[i][0])] = true
				} else {
					// If it's a multi-character string, add each character
					for _, char := range args[i] {
						debugChars[char] = true
					}
				}
				i++
			}
			i-- // Adjust index since the loop will increment it
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

	// Process character maps
	charMap := make(CharacterMap)
	for _, mapStr := range charMaps {
		if len(mapStr) >= 2 {
			// Convert string to runes to properly handle multi-byte UTF-8 characters
			runes := []rune(mapStr)
			if len(runes) >= 2 {
				fromChar := runes[0] // First rune (character)
				toChar := runes[1]   // Second rune (character)
				charMap[fromChar] = toChar
			}
		}
	}

	// Read header to get character height
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		fmt.Printf("Error reading header: %v\n", err)
		os.Exit(1)
	}
	header := scanner.Text()
	parts := strings.Split(header, " ")
	charHeight := 0
	if len(parts) >= 2 {
		fmt.Sscanf(parts[1], "%d", &charHeight)
	}
	file.Close()

	font, err := convertFLFToBit(inputFile, name, author, license, charMap, debugEnabled, debugChars, charHeight)
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

func convertFLFToBit(inputFile string, name string, author string, license string, charMap CharacterMap, debugEnabled bool, debugChars map[rune]bool, charHeight int) (*FontData, error) {
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

	// Read the first line after comments
	var firstDataLine string = ""

	// Read the first line after comments
	if scanner.Scan() {
		firstDataLine = scanner.Text()

		// If it's an empty line, continue reading until we find a non-empty line
		for firstDataLine == "" && scanner.Scan() {
			firstDataLine = scanner.Text()
		}

	}

	// Process characters
	characters := make(map[string][]string)

	// Process the first character line
	// currentLine is firstDataLine
	currentLine := firstDataLine

	// Process characters until EOF
	// Count all characters in file order, then map to appropriate ASCII codes
	charIndex := 0 // Start at 0, will adjust to ASCII later

	// Process characters one at a time
	for {
		// Skip lines that consist only of a single repeated character (like "@@" or "##")
		// These are separators between characters, not actual character data.
		// IMPORTANT: We check BEFORE trimming whitespace, because a line like "  ##"
		// should NOT be skipped - it's the end marker for the previous character.
		// Only lines like "@@" or "##" (with no leading whitespace) should be skipped.
		isOnlyDelimiter := false
		if len(currentLine) >= 2 {
			lastChar := currentLine[len(currentLine)-1:]
			if strings.Trim(currentLine, lastChar) == "" {
				isOnlyDelimiter = true
			}
		}
		if isOnlyDelimiter {
			// Skip this line and read the next one
			if !scanner.Scan() {
				break
			}
			currentLine = scanner.Text()
			continue
		}

		// Determine the line end character from the last character of the first line
		// This needs to be done for EACH character as TLF files can have different terminators
		trimmedFirstCharLine := strings.TrimRight(currentLine, " \t\r\n")
		lineEndChar := "@"
		if len(trimmedFirstCharLine) > 0 {
			lineEndChar = string(trimmedFirstCharLine[len(trimmedFirstCharLine)-1])
		}

		// 1. Determine marker character for current character block from first line of character
		if currentLine == "" {
			// Skip empty lines
			if !scanner.Scan() {
				break // End of file
			}
			currentLine = scanner.Text()
			continue
		}

		// 2. Read ahead to the end of character marker and collect all lines
		var charLines []string

		// For fonts with height=1, each line is a complete character
		// So we just process the current line and don't look for double terminator
		if charHeight == 1 {
			trimmedLine := strings.TrimRight(currentLine, " \t\r\n")
			if len(trimmedLine) > 0 {
				lineEndChar := string(trimmedLine[len(trimmedLine)-1])
				charLine := strings.TrimSuffix(trimmedLine, lineEndChar)
				charLines = append(charLines, charLine)
			}
		} else {
			// For fonts with height>1, look for double terminator to mark end of character
			for {
				trimmedLine := strings.TrimRight(currentLine, " \t\r\n")

				// Check if this is the end of the character (two consecutive end-of-line markers)
				endOfCharDelimiter := lineEndChar + lineEndChar
				if len(trimmedLine) >= 2 && strings.HasSuffix(trimmedLine, endOfCharDelimiter) {
					// Remove the end-of-character delimiter from the end
					charData := strings.TrimSuffix(trimmedLine, endOfCharDelimiter)

					// Split the data and collect each part
					parts := strings.Split(charData, lineEndChar)
					for _, part := range parts {
						charLines = append(charLines, part)
					}
					break
				} else if strings.HasSuffix(trimmedLine, lineEndChar) {
					// This is a line of the current character
					// Remove the line end character and add the line to current character
					charLine := strings.TrimSuffix(trimmedLine, lineEndChar)
					charLines = append(charLines, charLine)
				} else if currentLine == "" {
					// Empty line after character data may indicate end of character
					break
				}

				// Read next line
				if !scanner.Scan() {
					// End of file
					break
				}

				currentLine = scanner.Text()
			}
		}

		// 3. Process the block of character data to convert
		if len(charLines) > 0 {
			// Process the character lines with hardblank replacement and character mappings
			processedLines := processCharacter(charLines, hardblankChar, charMap)

			// Calculate the ASCII value for this character position
			asciiValue := charIndex + 32

			// 4. Update the `characters` with the current character as each character is processed
			// Only add characters in the standard ASCII range (space to tilde, 32-126)
			if asciiValue >= 32 && asciiValue <= 126 {
				char := string(rune(asciiValue))

				// Debug output - only show if debug is enabled
				if debugEnabled {
					// Check if we should show debug for this specific character
					showDebug := len(debugChars) == 0 // If no specific chars specified, show all
					if len(debugChars) > 0 {
						showDebug = debugChars[rune(asciiValue)] // Only show if this char is in the list
					}

					if showDebug {
						fmt.Printf("Processing character '%c' (ASCII %d)\n", asciiValue, asciiValue)
						fmt.Println("FLF input:")
						for _, line := range charLines {
							fmt.Printf("  %q\n", line)
						}
						fmt.Println("BIT output:")
						for _, line := range processedLines {
							fmt.Printf("  %q\n", line)
						}
						fmt.Println()
					}
				}

				characters[char] = processedLines
			}
			charIndex++
		}

		// 5. Repeat for next character until end of file
		if !scanner.Scan() {
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
