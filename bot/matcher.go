package bot

import "strings"

type Matcher interface {
	Match(content string) string
}

type SimpleMatcher struct {
	words map[string]interface{}
}

func NewSimpleMatcher(source []string) *SimpleMatcher {
	words := make(map[string]interface{})

	var v struct{}
	for _, s := range source {
		words[s] = v
	}

	return &SimpleMatcher{words}
}

func (m *SimpleMatcher) Match(content string) string {
	for _, word := range strings.Fields(content) {
		if _, ok := m.words[word]; ok {
			return word
		}
	}
	return ""
}
