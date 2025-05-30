package filters

import (
	"strings"
)

// NaughtyWordsFilter filters messages containing prohibited words.
type NaughtyWordsFilter struct {
	ProhibitedWords []string
}

// Filter returns true if the message contains any prohibited word.
func (f *NaughtyWordsFilter) Filter(message string) bool {
	lowerMsg := strings.ToLower(message)
	for _, word := range f.ProhibitedWords {
		if strings.Contains(lowerMsg, strings.ToLower(word)) {
			return true
		}
	}
	return false
}
