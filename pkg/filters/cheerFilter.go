package filters

import (
	"fmt"
	"regexp"
	"strings"
)

type CheerFilter struct {
	cheermoteRegexp regexp.Regexp
}

func NewCheerFilter(cheermotes []string) *CheerFilter {
	return &CheerFilter{
		cheermoteRegexp: *buildCheermoteRegex(cheermotes),
	}
}

func buildCheermoteRegex(cheermotes []string) *regexp.Regexp {
	pattern := `(?i)\b(?:` + strings.Join(cheermotes, "|") + `)(\d+)\b`
	return regexp.MustCompile(pattern)
}

func (cf *CheerFilter) Filter(message string) bool {
	if cf.cheermoteRegexp.MatchString(message) {
		fmt.Println("CheerFilter matched:", message)
		return true
	}
	return false
}
