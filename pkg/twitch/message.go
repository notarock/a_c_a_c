package twitch

import (
	"slices"
	"strings"
)

func (t *TwitchClient) IsUserIgnored(user string) bool {
	return slices.Contains(t.ignoredUsers, strings.ToLower(user))
}

func (t *TwitchClient) IsUserModerator(user string) bool {
	return slices.Contains(t.moderatorUsers, strings.ToLower(user))
}
