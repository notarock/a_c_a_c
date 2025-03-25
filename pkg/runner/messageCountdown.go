package runner

import (
	"fmt"

	gotwitch "github.com/gempir/go-twitch-irc/v4"

	"github.com/notarock/a_c_a_r/pkg/chain"
	"github.com/notarock/a_c_a_r/pkg/twitch"
)

const RED = "\033[31m"
const RESET = "\033[0m"

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

		// Don't learn messages from ignored users (bots)
		if runner.client.IsUserIgnored(message.User.Name) {
			return
		}

		// Don't learn parroted messages (if enabled)
		if runner.chain.IsParrot(message.Message) && runner.chain.IgnoreParrots {
			fmt.Println(RED, runner.client.Channel, ":", message.Message, RESET)
			return
		}

		runner.countdown = runner.countdown - 1              // Decrement countdown
		runner.chain.AddMessage(message.Message)             // Learn
		err := runner.chain.SaveChatMessage(message.Message) // Save to file
		if err != nil {
			fmt.Println("Error writing message to file:", err)
			return
		}

		// Log message
		fmt.Println(runner.client.Channel, ":", message.Message)

		// Don't send a message if we have not reached the countdown yet
		if runner.countdown > 0 {
			return
		}

		// Countdown reached, send a message and reset countdown

		runner.countdown = runner.interval         // Reset countdown
		response := runner.chain.FilteredMessage() // Generate a response

		runner.client.SendMessage(response)    // Send the message
		runner.chain.SaveSentMessage(response) // Save the sent message
	})

	return &runner
}

func (m *MessageCountdownRunner) Run() error {
	return m.client.Connect()
}
