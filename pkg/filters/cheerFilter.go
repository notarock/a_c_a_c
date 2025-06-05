package filters

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nicklaw5/helix"
)

type CheerFilter struct {
	cheermoteRegexp *regexp.Regexp
}

func NewCheerFilter(twitchApiClient *helix.Client, channelName string) (cf *CheerFilter, err error) {
	// Fetch channel ID using twitchApiClient and channel name
	usersResp, err := twitchApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{channelName},
	})
	if err != nil || usersResp.StatusCode != 200 || len(usersResp.Data.Users) == 0 {
		return cf, fmt.Errorf("failed to fetch user ID for channel %s: %v", channelName, err)
	}
	channelID := usersResp.Data.Users[0].ID

	resp, err := twitchApiClient.GetCheermotes(&helix.CheermotesParams{
		BroadcasterID: channelID,
	})

	if err != nil {
		return cf, fmt.Errorf("failed to fetch cheermotes for channel %s: %v", channelName, err)
	} else if resp.StatusCode != 200 {
		return cf, fmt.Errorf("failed to fetch cheermotes for channel %s: http response code was %d", channelName, resp.StatusCode)
	}
	var cheerKeywords []string
	for _, cheermote := range resp.Data.Cheermotes {
		cheerKeywords = append(cheerKeywords, cheermote.Prefix)
	}
	return &CheerFilter{cheermoteRegexp: buildCheermoteRegex(cheerKeywords)}, nil
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
