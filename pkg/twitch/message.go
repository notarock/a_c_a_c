package twitch

import "slices"

func (t *TwitchClient) IsUserIgnored(user string) bool {
	return slices.Contains(t.ignoredUsers, user)
}

func (t *TwitchClient) IsUserModerator(user string) bool {
	return slices.Contains(t.moderatorUsers, user)
}
