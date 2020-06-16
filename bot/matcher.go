package bot

import (
	"github.com/sajari/fuzzy"
	"strings"
	"triggered-bot/text"
)

type Matcher interface {
	Match(content string) string
	SetWords(source []string)
}

type SimpleMatcher struct {
	words   map[string]interface{}
	model   *fuzzy.Model
	sampler Sampler
}

func NewSimpleMatcher(sampler Sampler) *SimpleMatcher {
	matcher := &SimpleMatcher{sampler: sampler}
	source := sampler.SampleWords()
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
		n := text.Normalize(word)
		if m.isMatch(n) {
			return n
		}
	}
	return ""
}

func (m *SimpleMatcher) isMatch(s string) bool {
	return m.inWords(s) ||
		(len(s) > 5 && m.inWords(m.model.SpellCheck(s)))
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
