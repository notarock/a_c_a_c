package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/notarock/a_c_a_r/pkg/chain"
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

	chain, err := chain.NewChain(chain.ChainConfig{
		Saving:                true,
		SavedMessagesFilepath: SAVED_MESSAGES_FILE,
		SentMessagesFilepath:  MESSAGE_FILE,
	})

	if err != nil {
		log.Fatal(err)
	}

	if ENV != "production" {
		fmt.Println("Environment: ", ENV)
		fmt.Println("Channel: ", CHANNEL)
		fmt.Println("Saving chat messages to: ", MESSAGE_FILE)
		fmt.Println("Saving sent messages to: ", SAVED_MESSAGES_FILE)
	}

	// RunOnTimer(chain, 180)
	RunOnMessageCount(chain, 60)
}

func RunOnMessageCount(chain *chain.Chain, interval int) {
	fmt.Printf("Sending messages every %d chat messages!\n", interval)
	client := twitch.NewClient(TWITCH_USER, TWITCH_OAUTH_STRING)
	// client := twitch.NewAnonymousClient() // for an anonymous user (no write capabilities)
	fmt.Println("Connecting to twitch...")
	client.Join(CHANNEL)
	fmt.Println("Joined " + CHANNEL)

	n_till_next := interval

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		n_till_next = n_till_next - 1
		fmt.Print(".")

		if IgnoreBotMessages(message.User.Name) {
			return
		}

		chain.AddMessage(message.Message)
		err := chain.SaveChatMessage(message.Message)
		if err != nil {
			fmt.Println("Error writing message to file:", err)
			return
		}

		// Don't send a message if we're not ready
		if n_till_next > 0 {
			return
		}

		// Reset the counter
		n_till_next = interval
		fmt.Print("\n")

		response := chain.FilteredMessage()

		// Sleep a bit before sending the message to make it look a bit more natural
		time.Sleep(1 * time.Second)

		if ENV == "production" {
			// Send the message
			fmt.Printf("Sent to %s: %s\n", CHANNEL, response)
			client.Say(CHANNEL, response)
		} else {
			// Present to send the message
			fmt.Printf("I would have sent to %s: %s\n", CHANNEL, response)
		}

		chain.SaveSentMessage(response)

		fmt.Print("Listening")
	})

	fmt.Print("Listening")
	client.Connect()
}

func IgnoreBotMessages(user string) bool {
	return slices.Contains([]string{"oathybot", "funtoon", "cynanbot", TWITCH_USER}, user)
}

// func RunOnTimer(chain *gomarkov.Chain, interval time.Duration) {
// 	fmt.Println("Running on timer...")
// 	client := twitch.NewClient(TWITCH_USER, TWITCH_OAUTH_STRING)
// 	fmt.Println("Connecting to twitch...")

// 	client.Join(CHANNEL) // oats please pepeW
// 	fmt.Println("Joined " + CHANNEL)

// 	go client.Connect()

// 	for {
// 		message := Generate(chain)
// 		fmt.Println("Generated message:", message)

// 		if strings.Contains(message, "@") || strings.Contains(message, "https://") {
// 			fmt.Println("Message contains @, skipping...")
// 			continue
// 		}

// 		// reader := bufio.NewReader(os.Stdin)
// 		// fmt.Print("Yes / No? ")
// 		// input, _ := reader.ReadString('\n')
// 		// input = strings.TrimSpace(strings.ToLower(input))

// 		// if input == "yes" {
// 		// 	    fmt.Println("Sending message")
// 		client.Say(CHANNEL, message)
// 		// }

// 		time.Sleep(interval * time.Second)

// 	}
// }
