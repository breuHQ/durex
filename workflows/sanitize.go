package workflows

import (
	"regexp"
	"strings"
)

// Sanitize converts a string to snake_case, preserving hyphens and handling various edge cases. It is intended for use
// in generating unique identifiers for elements, modifiers, and blocks.
func Sanitize(input string) string {
	// Normalize whitespace and special characters (excluding hyphens) to underscores.
	re := regexp.MustCompile(`[^a-zA-Z0-9-]+`)
	output := re.ReplaceAllString(input, "_")

	// Convert to lower case.
	output = strings.ToLower(output)

	// Normalize multiple underscores to a single underscore.
	output = regexp.MustCompile(`_+`).ReplaceAllString(output, "_")

	// Normalize multiple dashes to a single dash.  (Crucial order)
	output = regexp.MustCompile(`-+`).ReplaceAllString(output, "-")

	// Handle empty input (crucially, trim whitespace first).
	output = strings.TrimSpace(output)
	if output == "" {
		return ""
	}

	if len(output) > 1 {
		// Remove leading/trailing underscores.
		output = strings.Trim(output, "_")

		// Remove leading/trailing dashes.
		output = strings.Trim(output, "-")
	}

	return output
}
