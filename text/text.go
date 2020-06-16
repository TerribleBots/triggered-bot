package text

import (
	"strings"
	"unicode"
)

const OverrideText = "trigger warning"

func Normalize(candidate string) string {
	candidate = strings.ToLower(candidate)
	candidate = strings.TrimSpace(candidate)
	candidate = strings.TrimFunc(candidate, unicode.IsPunct)
	return candidate
}

func Overriden(msg string) bool {
	return strings.HasPrefix(msg, OverrideText)
}
