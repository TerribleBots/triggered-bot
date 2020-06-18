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

func RemoveDuplicates(lines []string) []string {
	var out []string
	seen := make(map[string]interface{})
	var v struct{}
	for _, s := range lines {
		if _, ok := seen[s]; !ok {
			out = append(out, s)
			seen[s] = v
		}
	}

	return out
}
