package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/notarock/a_c_a_r/pkg/chain"
	"github.com/notarock/a_c_a_r/pkg/runner"
	"github.com/notarock/a_c_a_r/pkg/twitch"
)

var BASE_PATH = os.Getenv("BASE_PATH")
var CHANNEL = os.Getenv("TWITCH_CHANNEL")
var TWITCH_USER = os.Getenv("TWITCH_USER")
var TWITCH_OAUTH_STRING = os.Getenv("TWITCH_OAUTH_STRING")
var ENV = os.Getenv("ENV")

var MESSAGE_FILE_PATTERN = "%s/%s.txt"             // BASEPATH-CHANNEL.txt
var SAVED_MESSAGES_FILE_PATTERN = "%s/%s-sent.txt" // BASEPATH-CHANNEL-sent.txt

const GREEN = "\033[32m"
const RED = "\033[31m"
const RESET = "\033[0m"

func main() {
	if BASE_PATH == "" || CHANNEL == "" || TWITCH_USER == "" || TWITCH_OAUTH_STRING == "" {
		fmt.Println("Missing environment variables")
		return
	}

	channels := strings.Split(CHANNEL, ",")
	var runners []*runner.MessageCountdownRunner

	for _, channel := range channels {

		savedMessagesFilepath := fmt.Sprintf(MESSAGE_FILE_PATTERN, BASE_PATH, channel)
		sentMessagesFilepath := fmt.Sprintf(SAVED_MESSAGES_FILE_PATTERN, BASE_PATH, channel)

		if ENV != "production" {
			fmt.Println("Environment: ", ENV)
			fmt.Println("Channel: ", CHANNEL)
			fmt.Println("Base path: ", BASE_PATH)
			fmt.Println("Saving chat messages to: ", savedMessagesFilepath)
			fmt.Println("Saving sent messages to: ", sentMessagesFilepath)
		}

		chain, err := chain.NewChain(chain.ChainConfig{
			Saving:                true,
			SavedMessagesFilepath: savedMessagesFilepath,
			SentMessagesFilepath:  sentMessagesFilepath,
		})

		if err != nil {
			log.Fatal(err)
		}

		client := twitch.NewClient(twitch.ClientConfig{
			Username: TWITCH_USER,
			OAuth:    TWITCH_OAUTH_STRING,
			Channel:  channel,
			Sending:  ENV == "production",
		})

		r := runner.NewMessageCountdownRunner(runner.MessageCountdownConfig{
			Client:   client,
			Chain:    chain,
			Interval: 5,
		})

		runners = append(runners, r)
	}

	for _, runner := range runners {
		go runner.Run()
	}

	select {}
}
