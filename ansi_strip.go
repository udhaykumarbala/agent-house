package main

import (
	"regexp"
	"strings"
)

// ANSI escape code patterns
var (
	// Standard ANSI escape sequences: ESC[...m (colors, styles)
	ansiColorRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

	// Cursor movement and control: ESC[...H, ESC[...A, ESC[...K, etc.
	ansiCursorRegex = regexp.MustCompile(`\x1b\[[0-9;]*[A-HJKSTfhlmnsu]`)

	// Private mode sequences: ESC[?...h, ESC[?...l
	ansiPrivateRegex = regexp.MustCompile(`\x1b\[\?[0-9;]*[hl]`)

	// Operating system commands: ESC]...BEL or ESC]...ST
	ansiOSCRegex = regexp.MustCompile(`\x1b\][^\x07\x1b]*(?:\x07|\x1b\\)`)

	// Remaining ESC sequences
	ansiEscRegex = regexp.MustCompile(`\x1b[^\x1b]*`)
)

// StripANSI removes ANSI escape codes from text
func StripANSI(input string) string {
	result := input

	// Remove in order of specificity
	result = ansiOSCRegex.ReplaceAllString(result, "")
	result = ansiPrivateRegex.ReplaceAllString(result, "")
	result = ansiColorRegex.ReplaceAllString(result, "")
	result = ansiCursorRegex.ReplaceAllString(result, "")

	// Clean up any remaining escape sequences
	result = ansiEscRegex.ReplaceAllString(result, "")

	// Clean up carriage returns and excessive newlines
	result = strings.ReplaceAll(result, "\r\n", "\n")
	result = strings.ReplaceAll(result, "\r", "\n")

	// Collapse multiple blank lines
	multiNewline := regexp.MustCompile(`\n{3,}`)
	result = multiNewline.ReplaceAllString(result, "\n\n")

	return strings.TrimSpace(result)
}

// ExtractClaudeResponse tries to extract just Claude's response text
func ExtractClaudeResponse(input string) string {
	// First strip ANSI codes
	clean := StripANSI(input)

	// Common patterns in Claude output to filter
	lines := strings.Split(clean, "\n")
	var output []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip common UI elements
		if line == "" ||
			strings.HasPrefix(line, ">") ||
			strings.HasPrefix(line, "─") ||
			strings.HasPrefix(line, "│") ||
			strings.HasPrefix(line, "╭") ||
			strings.HasPrefix(line, "╰") ||
			strings.Contains(line, "for shortcuts") ||
			strings.Contains(line, "Puzzling") ||
			strings.Contains(line, "Noodling") ||
			strings.Contains(line, "esc to interrupt") ||
			strings.Contains(line, "Welcome back") ||
			strings.Contains(line, "Claude Code v") {
			continue
		}

		output = append(output, line)
	}

	return strings.Join(output, "\n")
}
