package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	gotwitch "github.com/gempir/go-twitch-irc/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/notarock/a_c_a_r/pkg/chain"
	"github.com/notarock/a_c_a_r/pkg/twitch"
)

var BASE_PATH = os.Getenv("BASE_PATH")
var CHANNEL = os.Getenv("TWITCH_CHANNEL")
var TWITCH_USER = os.Getenv("TWITCH_USER")
var TWITCH_OAUTH_STRING = os.Getenv("TWITCH_OAUTH_STRING")
var ENV = os.Getenv("ENV")

var MESSAGE_FILE = BASE_PATH + CHANNEL + ".txt"
var SAVED_MESSAGES_FILE = BASE_PATH + "sent.txt"

const GREEN = "\033[32m"
const RED = "\033[31m"
const RESET = "\033[0m"

func main() {
	if BASE_PATH == "" || CHANNEL == "" || TWITCH_USER == "" || TWITCH_OAUTH_STRING == "" {
		fmt.Println("Missing environment variables")
		return
	}

	if ENV != "production" {
		fmt.Println("Environment: ", ENV)
		fmt.Println("Channel: ", CHANNEL)
		fmt.Println("Saving chat messages to: ", MESSAGE_FILE)
		fmt.Println("Saving sent messages to: ", SAVED_MESSAGES_FILE)
	}

	chain, err := chain.NewChain(chain.ChainConfig{
		Saving:                true,
		SavedMessagesFilepath: MESSAGE_FILE,
		SentMessagesFilepath:  SAVED_MESSAGES_FILE,
	})

	if err != nil {
		log.Fatal(err)
	}

	client := twitch.NewClient(twitch.ClientConfig{
		Username: TWITCH_USER,
		OAuth:    TWITCH_OAUTH_STRING,
		Channel:  CHANNEL,
		Sending:  ENV == "production",
	})

	RunOnMessageCount(chain, client, 60)
}

func RunOnMessageCount(chain *chain.Chain, client *twitch.TwitchClient, interval int) {
	messageCountdown := interval

	fmt.Println("Adding message hook...")
	client.AddMessageHook(func(message gotwitch.PrivateMessage) {

		messageCountdown = messageCountdown - 1 // Decrement countdown
		if IgnoreMessagesFromBots(message.User.Name) {
			return
		}

		chain.AddMessage(message.Message)             // Learn
		err := chain.SaveChatMessage(message.Message) // Save to file
		if err != nil {
			fmt.Println("Error writing message to file:", err)
			return
		}

		fmt.Println("Saved:", message.Message)

		// Don't send a message if we have not reached the countdown yet
		if messageCountdown > 0 {
			return
		}

		messageCountdown = interval // Reset countdown

		response := chain.FilteredMessage() // Generate a response

		client.SendMessage(response)    // Send the message
		chain.SaveSentMessage(response) // Save the sent message
	})

	fmt.Println("Hook added, now listening.")
	err := client.Connect()
	if err != nil {
		log.Fatalln(err)
	}
}

func IgnoreMessagesFromBots(user string) bool {
	return slices.Contains([]string{"oathybot", "funtoon", "cynanbot", "mandoobot", TWITCH_USER}, user)
}
