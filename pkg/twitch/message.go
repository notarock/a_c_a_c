package twitch

import "slices"

func (t *TwitchClient) IsUserIgnored(user string) bool {
	return slices.Contains([]string{"oathybot", "funtoon", "cynanbot", "mandoobot", t.Username}, user)
}
