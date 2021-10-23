package cli

import (
	"fmt"
	"strings"
)

// ToSentence converts a slice of terms into sentence.
func ToSentence(words []string, andOrOr string) string {
	l := len(words)

	if l == 1 {
		return fmt.Sprintf("'%s'", words[0])
	} else if l == 2 {
		return fmt.Sprintf("'%s' or '%s'", words[0], words[1])
	}

	wordsForSentence := []string{}
	for _, w := range words {
		wordsForSentence = append(wordsForSentence, fmt.Sprintf("'%s'", w))
	}

	wordsForSentence[l-1] = andOrOr + " " + wordsForSentence[l-1]
	return strings.Join(wordsForSentence, ", ")
}
