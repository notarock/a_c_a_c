package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/notarock/a_c_a_c/pkg/chain"
	"github.com/notarock/a_c_a_c/pkg/runner"
	"github.com/notarock/a_c_a_c/pkg/twitch"
)

var COUNTDOWN = os.Getenv("COUNTDOWN")
var IGNORE_PARROTS = os.Getenv("IGNORE_PARROTS") == "true"
var BASE_PATH = os.Getenv("BASE_PATH")
var TWITCH_USER = os.Getenv("TWITCH_USER")
var TWITCH_OAUTH_STRING = os.Getenv("TWITCH_OAUTH_STRING")
var ENV = os.Getenv("ENV")
var CHANNELS = strings.Split(os.Getenv("TWITCH_CHANNELS"), ",")
var BOTS = strings.Split(os.Getenv("TWITCH_BOT_USERNAMES"), ",")

var PROHIBITED_STRINGS = strings.Split(os.Getenv("PROHIBITED_STRINGS"), ",")
var PROHIBITED_MESSAGES = strings.Split(os.Getenv("PROHIBITED_MESSAGES"), ",")

var MESSAGE_FILE_PATTERN = "%s/%s.txt"             // BASEPATH-CHANNEL.txt
var SAVED_MESSAGES_FILE_PATTERN = "%s/%s-sent.txt" // BASEPATH-CHANNEL-sent.txt

const GREEN = "\033[32m"
const RED = "\033[31m"
const RESET = "\033[0m"

func main() {
	if BASE_PATH == "" || TWITCH_USER == "" || TWITCH_OAUTH_STRING == "" || COUNTDOWN == "" {
		fmt.Println("Missing environment variables")
		return
	}

	countdownInterval, err := strconv.Atoi(COUNTDOWN)
	if err != nil {
		log.Fatalf("Failed to read COUNTDOWN as int: %v", err)
	}

	var runners []*runner.MessageCountdownRunner

	for _, channel := range CHANNELS {

		savedMessagesFilepath := fmt.Sprintf(MESSAGE_FILE_PATTERN, BASE_PATH, channel)
		sentMessagesFilepath := fmt.Sprintf(SAVED_MESSAGES_FILE_PATTERN, BASE_PATH, channel)

		if ENV != "production" {
			fmt.Println("Environment: ", ENV)
			fmt.Println("Channel: ", GREEN, channel, RESET)
			fmt.Println("Base path: ", BASE_PATH)
			fmt.Println("Saving chat messages to: ", savedMessagesFilepath)
			fmt.Println("Saving sent messages to: ", sentMessagesFilepath)
		}

		chain, err := chain.NewChain(chain.ChainConfig{
			Saving:                true,
			IgnoreParrots:         IGNORE_PARROTS,
			SavedMessagesFilepath: savedMessagesFilepath,
			SentMessagesFilepath:  sentMessagesFilepath,
			ProhibitedStrings:     PROHIBITED_STRINGS,
			ProhibitedMessages:    PROHIBITED_MESSAGES,
		})

		if err != nil {
			log.Fatal(err)
		}

		client := twitch.NewClient(twitch.ClientConfig{
			Username: TWITCH_USER,
			OAuth:    TWITCH_OAUTH_STRING,
			Channel:  channel,
			Bots:     BOTS,
			Sending:  ENV == "production",
		})

		r := runner.NewMessageCountdownRunner(runner.MessageCountdownConfig{
			Client:   client,
			Chain:    chain,
			Interval: countdownInterval,
		})

		runners = append(runners, r)
	}

	for _, runner := range runners {
		go runner.Run()
	}

	select {}
}
