package bot

import (
	"github.com/sajari/fuzzy"
	"strings"
)

type Matcher interface {
	Match(content string) string
	SetWords(source []string)
}

type SimpleMatcher struct {
	words map[string]interface{}
	model *fuzzy.Model
}

func NewSimpleMatcher(source []string) *SimpleMatcher {
	matcher := &SimpleMatcher{}
	matcher.SetWords(source)
	return matcher
}

func (m *SimpleMatcher) SetWords(source []string) {
	words := make(map[string]interface{})

	var v struct{}
	for _, s := range source {
		words[s] = v
	}

	m.words = words
	m.model = newModel(source)
}

func (m *SimpleMatcher) Match(content string) string {
	for _, word := range strings.Fields(content) {
		if m.isMatch(word) {
			return word
		}
	}
	return ""
}

func (m *SimpleMatcher) isMatch(s string) bool {
	return m.inWords(s) || m.inWords(m.model.SpellCheck(s))
}

func (m *SimpleMatcher) inWords(s string) bool {
	_, ok := m.words[s]
	return ok
}

func newModel(words []string) *fuzzy.Model {
	model := fuzzy.NewModel()
	model.SetThreshold(1)
	model.SetDepth(1)
	model.Train(words)
	return model
}
