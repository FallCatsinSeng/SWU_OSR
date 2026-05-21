// Package sanitize provides input sanitization utilities for user-generated content.
// It strips HTML tags and dangerous content before storage to prevent stored XSS
// attacks if content is ever rendered outside React's auto-escaping (emails, mobile
// apps, RSS feeds, markdown renderers, etc.).
package sanitize

import (
	"regexp"
	"strings"
)

// htmlTagRegex matches HTML/XML tags including self-closing and attributes.
var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)

// StripHTML removes all HTML/XML tags from the input string.
// It also decodes common HTML entities that could be used to bypass filters,
// and trims leading/trailing whitespace from the result.
func StripHTML(input string) string {
	// Remove HTML tags
	clean := htmlTagRegex.ReplaceAllString(input, "")

	// Decode common HTML entities that could be used for obfuscation
	clean = strings.ReplaceAll(clean, "&lt;", "<")
	clean = strings.ReplaceAll(clean, "&gt;", ">")
	clean = strings.ReplaceAll(clean, "&amp;", "&")
	clean = strings.ReplaceAll(clean, "&quot;", `"`)
	clean = strings.ReplaceAll(clean, "&#39;", "'")

	// Re-strip in case decoded entities formed new tags
	clean = htmlTagRegex.ReplaceAllString(clean, "")

	return strings.TrimSpace(clean)
}

// StripHTMLPreserveWhitespace removes HTML tags but preserves internal whitespace
// (newlines, indentation). Useful for multi-line content like forum posts and bios.
func StripHTMLPreserveWhitespace(input string) string {
	// Remove HTML tags
	clean := htmlTagRegex.ReplaceAllString(input, "")

	// Decode common HTML entities
	clean = strings.ReplaceAll(clean, "&lt;", "<")
	clean = strings.ReplaceAll(clean, "&gt;", ">")
	clean = strings.ReplaceAll(clean, "&amp;", "&")
	clean = strings.ReplaceAll(clean, "&quot;", `"`)
	clean = strings.ReplaceAll(clean, "&#39;", "'")

	// Re-strip in case decoded entities formed new tags
	clean = htmlTagRegex.ReplaceAllString(clean, "")

	return clean
}

// TruncateUTF8 truncates a string to maxLen runes (not bytes) to safely
// limit user input without breaking multi-byte characters.
func TruncateUTF8(input string, maxLen int) string {
	runes := []rune(input)
	if len(runes) <= maxLen {
		return input
	}
	return string(runes[:maxLen])
}
