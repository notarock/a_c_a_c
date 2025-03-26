package twitch

import "slices"

func (t *TwitchClient) IsUserIgnored(user string) bool {
	return slices.Contains(t.ignoredUsers, user)
}
