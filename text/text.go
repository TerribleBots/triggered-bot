package text

import (
	"strings"
	"unicode"
)

const OverrideText = "trigger warning"

func Normalize(candidate string) string {
	candidate = strings.ToLower(candidate)
	candidate = strings.TrimSpace(candidate)
	candidate = strings.TrimFunc(candidate, isTerminator)
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

var otherTerminators = map[rune]bool{
	'`':  true,
	'1':  true,
	'2':  true,
	'3':  true,
	'4':  true,
	'5':  true,
	'6':  true,
	'7':  true,
	'8':  true,
	'9':  true,
	'0':  true,
	'-':  true,
	'=':  true,
	'~':  true,
	'!':  true,
	'@':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'^':  true,
	'&':  true,
	'*':  true,
	'(':  true,
	')':  true,
	'_':  true,
	'+':  true,
	'\\': true,
	']':  true,
	'[':  true,
	'|':  true,
	'}':  true,
	'{':  true,
	'\'': true,
	';':  true,
	'"':  true,
	':':  true,
	'/':  true,
	'.':  true,
	',':  true,
	'?':  true,
	'>':  true,
	'<':  true,
}

func isTerminator(r rune) bool {
	return unicode.IsPunct(r) || otherTerminators[r]
}
