package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

func TextValidateProcess(text string) (string, error) {
	re := regexp.MustCompile(`[^\p{L}\d\s]`)
	matched := re.MatchString(text)
	if matched {
		return "", errors.New("text is not valid")
	}
	s := strings.TrimSpace(text)

	// Split string into words
	words := strings.Fields(s)

	// Loop over each word and convert the first letter to uppercase, others to lowercase
	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			for j := 1; j < len(runes); j++ {
				runes[j] = unicode.ToLower(runes[j])
			}
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " "), nil

}
