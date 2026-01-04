package twitch

import (
	"fmt"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

type TwitchClient struct {
	client               *twitch.Client
	oauth                string
	Channel              string
	Username             string
	Sending              bool
	ignoredUsers         []string
	moderatorUsers       []string
	ResponseDelaySeconds int
}

type ClientConfig struct {
	Username             string
	OAuth                string
	Channel              string
	Sending              bool
	Bots                 []string
	BotModerators        []string
	ResponseDelaySeconds int
}

func NewClient(config ClientConfig) *TwitchClient {
	var client *twitch.Client

	if config.Username != "" && config.OAuth != "" {
		client = twitch.NewClient(config.Username, config.OAuth)
	} else {
		fmt.Println("No username or OAuth provided, using anonymous client")
		client = twitch.NewAnonymousClient()
	}

	client.OnConnect(func() {
		fmt.Println("Connected to Twitch as", config.Username, "in channel", config.Channel)
	})

	client.Join(config.Channel)

	return &TwitchClient{
		client:               client,
		oauth:                config.OAuth,
		Channel:              config.Channel,
		Username:             config.Username,
		Sending:              config.Sending,
		ignoredUsers:         append(config.Bots, config.Username, config.Channel),
		moderatorUsers:       config.BotModerators,
		ResponseDelaySeconds: config.ResponseDelaySeconds,
	}
}

func (t *TwitchClient) Connect() error {
	return t.client.Connect()
}

func (t *TwitchClient) AddMessageHook(hook func(twitch.PrivateMessage)) {
	t.client.OnPrivateMessage(hook)
}

func (t *TwitchClient) SendMessage(message string) {
	if t.ResponseDelaySeconds > 0 {
		time.Sleep(time.Duration(t.ResponseDelaySeconds) * time.Second)
	}
	if t.Sending {
		t.client.Say(t.Channel, message)
	} else {
		fmt.Println("Not sending message to", t.Channel, ":", message)
	}
}
