package bot

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"math"
	"math/rand"
	"os"
	. "triggered-bot/log"
	"triggered-bot/text"
)

type Sampler struct {
	SourceFile, IncludeFile string
	SampleRatio             float64
}

func NewSampler(SourceFile, IncludeFile string, SampleRatio float64) Sampler {
	return Sampler{SourceFile, IncludeFile, SampleRatio}
}

func (s *Sampler) SampleWords() []string {
	words, err := readLines(s.SourceFile)
	if err != nil {
		Log.Errorf("unable to sample words: %s", err)
	}

	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })
	n := int(math.Round(s.SampleRatio * float64(len(words))))
	words = words[:n]
	Log.Info(zap.Strings("words", words))
	return append(words, include(s.IncludeFile)...)
}

func readLines(file string) ([]string, error) {
	var lines []string
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %s", file, err)
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return text.RemoveDuplicates(lines), scanner.Err()
}

func include(includeFile string) []string {
	if includeFile != "" {
		included, err := readLines(includeFile)
		if err != nil {
			Log.Fatal(err)
		}

		return included
	}

	return []string{}
}
