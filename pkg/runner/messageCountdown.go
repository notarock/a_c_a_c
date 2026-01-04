package runner

import (
	"fmt"
	"time"

	gotwitch "github.com/gempir/go-twitch-irc/v4"

	"github.com/notarock/a_c_a_c/pkg/chain"
	"github.com/notarock/a_c_a_c/pkg/filters"
	"github.com/notarock/a_c_a_c/pkg/twitch"
)

const RED = "\033[31m"
const BLUE = "\033[34m"
const RESET = "\033[0m"
const PARROT = "🦜"
const TALKING_HEAD = "🗣️"

type MessageCountdownRunner struct {
	client    *twitch.TwitchClient
	chain     *chain.Chain
	filters   []filters.Filter
	interval  int
	countdown int
}

type MessageCountdownConfig struct {
	Client   *twitch.TwitchClient
	Chain    *chain.Chain
	Filters  []filters.Filter
	Interval int
}

func NewMessageCountdownRunner(config MessageCountdownConfig) *MessageCountdownRunner {
	runner := MessageCountdownRunner{
		client:    config.Client,
		chain:     config.Chain,
		interval:  config.Interval,
		countdown: config.Interval,
		filters:   config.Filters,
	}
	fmt.Println("Adding message hook for", runner.client.Channel)

	runner.client.AddMessageHook(func(message gotwitch.PrivateMessage) {

		// Don't learn messages from ignored users (bots)
		if runner.client.IsUserIgnored(message.User.Name) {
			return
		}

		if runner.client.IsUserModerator(message.User.Name) || message.User.Name == runner.client.Channel {
			if message.Message == "!acac" {
				fmt.Println("Moderator", message.User.Name, "made me speak!", runner.client.Channel)
				response := runner.chain.GenerateValidMessage(runner.filters) // Generate a valid message

				fmt.Println(TALKING_HEAD, BLUE, runner.client.Channel, ":", response, RESET)
				runner.client.SendMessage(response)    // Send the message
				runner.chain.SaveSentMessage(response) // Save the sent message
				return
			}
		}

		// Don't learn parroted messages (if enabled)
		if runner.chain.IsParrot(message.Message) && runner.chain.IgnoreParrots {
			fmt.Println(PARROT, RED, runner.client.Channel, ":", message.Message, RESET)
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
		fmt.Println(runner.client.Channel, "(", runner.countdown, ")", ":", message.Message)

		// Don't send a message if we have not reached the countdown yet
		if runner.countdown > 0 {
			return
		}

		// Countdown reached, send a message and reset countdown

		runner.countdown = runner.interval // Reset countdown

		response := runner.chain.GenerateValidMessage(runner.filters) // Generate a valid message

		fmt.Println(TALKING_HEAD, BLUE, runner.client.Channel, ":", response, RESET)
		runner.delayAndSend(response)          // Send the message with delay
		runner.chain.SaveSentMessage(response) // Save the sent message
	})

	return &runner
}

func (m *MessageCountdownRunner) delayAndSend(message string) {
	delay := m.client.GetResponseDelay()
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Second)
	}
	m.client.SendMessage(message)
}

func (m *MessageCountdownRunner) Run() error {
	return m.client.Connect()
}
