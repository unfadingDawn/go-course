package usecase

import (
	"strings"
	"unicode"
)

func tokenize(text string) []string {
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	result := make([]string, 0, len(words))
	seen := make(map[string]struct{}, len(words))
	for _, word := range words {
		if _, exists := seen[word]; exists {
			continue
		}
		seen[word] = struct{}{}
		result = append(result, word)
	}
	return result
}
