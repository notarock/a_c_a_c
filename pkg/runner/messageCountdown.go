package runner

import (
	"fmt"

	gotwitch "github.com/gempir/go-twitch-irc/v4"

	"github.com/notarock/a_c_a_r/pkg/chain"
	"github.com/notarock/a_c_a_r/pkg/twitch"
)

type MessageCountdownRunner struct {
	client    *twitch.TwitchClient
	chain     *chain.Chain
	interval  int
	countdown int
}

type MessageCountdownConfig struct {
	Client   *twitch.TwitchClient
	Chain    *chain.Chain
	Interval int
}

func NewMessageCountdownRunner(config MessageCountdownConfig) *MessageCountdownRunner {
	runner := MessageCountdownRunner{
		client:    config.Client,
		chain:     config.Chain,
		interval:  config.Interval,
		countdown: config.Interval,
	}
	fmt.Println("Adding message hook...")

	runner.client.AddMessageHook(func(message gotwitch.PrivateMessage) {

		runner.countdown = runner.countdown - 1 // Decrement countdown
		if runner.client.IsUserIgnored(message.User.Name) {
			return
		}

		runner.chain.AddMessage(message.Message)             // Learn
		err := runner.chain.SaveChatMessage(message.Message) // Save to file
		if err != nil {
			fmt.Println("Error writing message to file:", err)
			return
		}

		fmt.Println(runner.client.Channel, ":", message.Message)

		// Don't send a message if we have not reached the countdown yet
		if runner.countdown > 0 {
			return
		}

		runner.countdown = runner.interval // Reset countdown

		response := runner.chain.FilteredMessage() // Generate a response

		runner.client.SendMessage(response)    // Send the message
		runner.chain.SaveSentMessage(response) // Save the sent message
	})

	return &runner
}

func (m *MessageCountdownRunner) Run() error {
	return m.client.Connect()
}
