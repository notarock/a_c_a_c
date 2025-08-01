package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nicklaw5/helix"
	"github.com/notarock/a_c_a_c/pkg/chain"
	"github.com/notarock/a_c_a_c/pkg/config"
	"github.com/notarock/a_c_a_c/pkg/filters"
	"github.com/notarock/a_c_a_c/pkg/runner"
	"github.com/notarock/a_c_a_c/pkg/twitch"
)

var IGNORE_PARROTS = os.Getenv("IGNORE_PARROTS") == "true"
var BASE_PATH = os.Getenv("BASE_PATH")
var TWITCH_USER = os.Getenv("TWITCH_USER")
var TWITCH_OAUTH_STRING = os.Getenv("TWITCH_OAUTH_STRING")
var ENV = os.Getenv("ENV")

var TWITCH_API_CLIENTID = os.Getenv("TWITCH_API_CLIENTID")
var TWITCH_API_TOKEN = os.Getenv("TWITCH_API_TOKEN")

var PROHIBITED_STRINGS = strings.Split(os.Getenv("PROHIBITED_STRINGS"), ",")
var PROHIBITED_MESSAGES = strings.Split(os.Getenv("PROHIBITED_MESSAGES"), ",")

var GLOBAL_MODERATORS = strings.Split(os.Getenv("GLOBAL_MODERATORS"), ",")

var MESSAGE_FILE_PATTERN = "%s/%s.txt"             // BASEPATH-CHANNEL.txt
var SAVED_MESSAGES_FILE_PATTERN = "%s/%s-sent.txt" // BASEPATH-CHANNEL-sent.txt

const GREEN = "\033[32m"
const RED = "\033[31m"
const RESET = "\033[0m"

func main() {
	messagesFile := flag.String("from-file", "", "A file to read messages from and spit out one generated message")
	flag.Parse()

	if *messagesFile != "" {
		message := loadAndGenerate(*messagesFile)
		fmt.Println(message)
		os.Exit(0)
	}

	if BASE_PATH == "" || TWITCH_USER == "" || TWITCH_OAUTH_STRING == "" {
		log.Panic("Missing environment variables")
	}

	twitchApiClient, err := helix.NewClient(&helix.Options{
		ClientID:        TWITCH_API_CLIENTID,
		UserAccessToken: TWITCH_API_TOKEN,
	})

	if err != nil {
		log.Fatal(err)
	}

	channelConfig, err := config.LoadChannelConfig(os.Getenv("CHANNEL_CONFIG"))

	if err != nil {
		log.Panic("Error loading channel config:", err)
	}

	baseFilters := []filters.Filter{
		&filters.NaughtyWordsFilter{
			ProhibitedWords: PROHIBITED_STRINGS,
		},
		&filters.MessageFilter{
			Messages: PROHIBITED_MESSAGES,
		},
	}

	var runners []*runner.MessageCountdownRunner

	for _, channel := range channelConfig.Channels {

		channelName := strings.ToLower(channel.Name)

		savedMessagesFilepath := fmt.Sprintf(MESSAGE_FILE_PATTERN, BASE_PATH, channelName)
		sentMessagesFilepath := fmt.Sprintf(SAVED_MESSAGES_FILE_PATTERN, BASE_PATH, channelName)

		if ENV != "production" {
			fmt.Println("Environment: ", ENV)
			fmt.Println("Channel: ", GREEN, fmt.Sprintf("%+v", channel), RESET)
			fmt.Println("Base path: ", BASE_PATH)
			fmt.Println("Saving chat messages to: ", savedMessagesFilepath)
			fmt.Println("Saving sent messages to: ", sentMessagesFilepath)
		}

		chain, err := chain.NewChain(chain.ChainConfig{
			Saving:                true,
			IgnoreParrots:         IGNORE_PARROTS,
			SavedMessagesFilepath: savedMessagesFilepath,
			SentMessagesFilepath:  sentMessagesFilepath,
		})

		if err != nil {
			log.Fatal(err)
		}

		client := twitch.NewClient(twitch.ClientConfig{
			Username:      TWITCH_USER,
			OAuth:         TWITCH_OAUTH_STRING,
			Channel:       channelName,
			Bots:          append(channelConfig.Bots, channel.ExtraBots...),
			Sending:       ENV == "production",
			BotModerators: GLOBAL_MODERATORS,
		})

		channelFilters := baseFilters
		if !channel.AllowBits {
			cheerFilter, err := filters.NewCheerFilter(twitchApiClient, channelName)
			if err != nil {
				log.Fatal(err)
			}
			channelFilters = append(channelFilters, cheerFilter)
		}

		r := runner.NewMessageCountdownRunner(runner.MessageCountdownConfig{
			Client:   client,
			Chain:    chain,
			Interval: channel.Frequency,
			Filters:  channelFilters,
		})

		runners = append(runners, r)
	}

	fmt.Println("Starting runners... Currently", len(runners), "channels configured.")

	for _, runner := range runners {
		go runner.Run()
		time.Sleep(1 * time.Second) // Rate limit the startup (max 20 per 10 seconds)
	}

	select {}
}

func loadAndGenerate(messagesFile string) string {
	chain, err := chain.NewChain(chain.ChainConfig{
		Saving:                false,
		IgnoreParrots:         false,
		SavedMessagesFilepath: messagesFile,
	})
	if err != nil {
		log.Panic(err)
	}

	if ENV != "production" {
		fmt.Println("Prohibited strings:", PROHIBITED_STRINGS)
		fmt.Println("Prohibited messages:", PROHIBITED_MESSAGES)
	}

	return chain.GenerateValidMessage([]filters.Filter{
		&filters.NaughtyWordsFilter{
			ProhibitedWords: PROHIBITED_STRINGS,
		},
		&filters.MessageFilter{
			Messages: PROHIBITED_MESSAGES,
		},
	}) // Generate a valid message from the loaded messages
}
