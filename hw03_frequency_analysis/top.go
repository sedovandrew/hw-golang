package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var wordRegexp = regexp.MustCompile(`([-а-я]{2,}|[а-я])`)

func Top10(text string) []string {
	words := make([]string, 0)
	wordsMap := make(map[string]int)

	// Find words
	lowerText := strings.ToLower(text)
	byteWords := wordRegexp.FindAll([]byte(lowerText), -1)
	for _, byteWord := range byteWords {
		word := string(byteWord)
		if _, ok := wordsMap[word]; !ok {
			words = append(words, word)
		}
		wordsMap[word]++
	}

	// Sort words
	sort.Slice(words, func(i, j int) bool {
		if wordsMap[words[i]] == wordsMap[words[j]] {
			return words[i] < words[j]
		}
		return wordsMap[words[i]] > wordsMap[words[j]]
	})

	// Truncate words
	if len(words) > 10 {
		return words[:10]
	}
	return words
}
