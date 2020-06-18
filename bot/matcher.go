package bot

import (
	"github.com/sajari/fuzzy"
	"strings"
	"triggered-bot/text"
)

type Matcher interface {
	Match(content string) MatchResult
	SetWords(source []string)
}

type MatchResult struct {
	matches, approximates []string
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

func (m *SimpleMatcher) Match(content string) MatchResult {
	var matches, approximate []string
	for _, word := range strings.Fields(content) {
		n := text.Normalize(word)
		if m.inWords(n) {
			matches = append(matches, n)
		} else if len(n) > 5 {
			corrected := m.model.SpellCheck(n)
			if m.inWords(corrected) {
				approximate = append(approximate, corrected)
			}
		}
	}
	return MatchResult{
		text.RemoveDuplicates(matches),
		text.RemoveDuplicates(approximate),
	}
}

func (m *SimpleMatcher) inWords(s string) bool {
	_, ok := m.words[s]
	return ok
}

func (m MatchResult) AnyMatch() bool {
	return len(m.matches) > 0 || len(m.approximates) > 0
}

func newModel(words []string) *fuzzy.Model {
	model := fuzzy.NewModel()
	model.SetThreshold(1)
	model.SetDepth(1)
	model.Train(words)
	return model
}
