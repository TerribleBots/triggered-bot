package bot

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"math"
	"math/rand"
	"os"
	. "triggered-bot/log"
)

type Sampler struct {
	SourceFile, IncludeFile string
	SampleRatio             float64
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

	return removeDuplicates(lines), scanner.Err()
}

func removeDuplicates(lines []string) []string {
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
