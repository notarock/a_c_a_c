package filters

import (
	"strings"
)

// MessageFilter filters messages that exactly match any of the specified messages.
type MessageFilter struct {
	Messages []string
}

// Filter returns true if the message exactly matches any of the specified messages.
func (f *MessageFilter) Filter(message string) bool {
	for _, m := range f.Messages {
		if strings.EqualFold(message, m) {
			return true
		}
	}
	return false
}
