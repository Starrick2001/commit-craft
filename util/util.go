// Package util provides shared constants and helpers.
package util

import "strings"

const (
	GeneralPrompt = "You are an expert at creating a git commit message for a set of changes. Here is a diff of changes we need a commit message for (response with title and description): "
)

func SanitizeJSONResponse(raw string) string {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)
	return cleaned
}
