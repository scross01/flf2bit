#!/usr/bin/env python3
"""
Test script to verify that each character in .bit font files has 
consistent width across all rows.
"""

import json
import re
import sys
from pathlib import Path


def check_font_consistency(font_file):
    """Check if all characters in a font file have consistent row widths."""
    
    # Pattern to match ANSI escape sequences
    ansi_escape = re.compile(r'\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])')
    
    print(f"\nChecking {font_file}...")
    
    with open(font_file, 'r', encoding='utf-8') as f:
        font_data = json.load(f)
    
    characters = font_data.get('characters', {})
    errors = []
    
    for char_name, char_lines in characters.items():
        if not char_lines:
            continue
            
        # Get the width of each row (in characters, not bytes)
        # Remove ANSI escape sequences before counting
        cleaned_lines = [ansi_escape.sub('', line) for line in char_lines]
        row_widths = [len(line) for line in cleaned_lines]
        
        # Check if all rows have the same width
        if len(set(row_widths)) > 1:
            errors.append({
                'character': char_name,
                'char_code': ord(char_name) if len(char_name) == 1 else None,
                'row_widths': row_widths,
                'lines': char_lines
            })
    
    if errors:
        print(f"  ❌ Found {len(errors)} characters with inconsistent widths:")
        for error in errors[:5]:  # Show first 5 errors
            char = error['character']
            if error['char_code']:
                char_display = f"'{char}' (ASCII {error['char_code']})"
            else:
                char_display = f"'{char}'"
            print(f"    - {char_display}: widths {error['row_widths']}")
            # Show the actual lines
            for i, line in enumerate(error['lines']):
                cleaned_line = ansi_escape.sub('', line)
                print(f"      Row {i}: '{line}' (len={len(line)}, clean_len={len(cleaned_line)})")
        
        if len(errors) > 5:
            print(f"    ... and {len(errors) - 5} more errors")
        return False
    else:
        print(f"  ✅ All {len(characters)} characters have consistent widths")
        return True


def main():
    fonts_dir = Path("fonts")
    
    if not fonts_dir.exists():
        print(f"Error: fonts directory not found at {fonts_dir}")
        sys.exit(1)
    
    # Find all .bit files
    bit_files = list(fonts_dir.glob("*.bit"))
    
    if not bit_files:
        print("No .bit files found in fonts/ directory")
        sys.exit(1)
    
    print(f"Found {len(bit_files)} .bit font files")
    
    all_consistent = True
    inconsistent_fonts = []
    
    for font_file in sorted(bit_files):
        is_consistent = check_font_consistency(font_file)
        if not is_consistent:
            all_consistent = False
            inconsistent_fonts.append(font_file.name)
    
    print("\n" + "="*60)
    if all_consistent:
        print("✅ SUCCESS: All fonts have consistent character widths")
    else:
        print(f"❌ FAILURE: {len(inconsistent_fonts)} fonts have "
              f"inconsistent widths:")
        for font in inconsistent_fonts:
            print(f"  - {font}")
        sys.exit(1)


if __name__ == "__main__":
    main()