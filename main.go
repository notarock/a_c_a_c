package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mb-14/gomarkov"
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

	listen := flag.Bool("listen", false, "Listen for chat messages and saves them to "+MESSAGE_FILE)

	flag.Parse()
	if *listen {
		fmt.Println("Listening for chat messages...")
		Listen()
		os.Exit(0)
	}

	chain, err := LoadModel()
	if err != nil {
		fmt.Println(err)
		return
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

func RunOnMessageCount(chain *gomarkov.Chain, interval int) {
	fmt.Printf("Sending messages every %d chat messages!\n", interval)
	client := twitch.NewClient(TWITCH_USER, TWITCH_OAUTH_STRING)
	// client := twitch.NewAnonymousClient() // for an anonymous user (no write capabilities)
	fmt.Println("Connecting to twitch...")
	client.Join(CHANNEL)
	fmt.Println("Joined " + CHANNEL)

	n_till_next := interval

	sendFile, err := os.OpenFile(SAVED_MESSAGES_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening "+SAVED_MESSAGES_FILE+":", err)
		return
	}
	defer sendFile.Close()

	messagesFile, err := os.OpenFile(MESSAGE_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening "+MESSAGE_FILE+":", err)
		return
	}
	defer messagesFile.Close()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		n_till_next = n_till_next - 1
		fmt.Print(".")

		if IgnoreBotMessages(message.User.Name) {
			return
		}

		// Learn from the message
		chain.Add(strings.Split(message.Message, " "))

		// Write the message to the file for future training
		_, err = messagesFile.WriteString(message.Message + "\n")
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

		// Generate a message compliant with basic filtering
		response := FilteredMessage(chain)

		// Sleep a bit before sending the message to make it look a bit more natural
		time.Sleep(2 * time.Second)

		// Send the message
		fmt.Printf("Sent to %s: %s\n", CHANNEL, response)

		if ENV == "production" {
			client.Say(CHANNEL, response)
		}

		_, err = sendFile.WriteString(response + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		fmt.Print("Listening")
	})

	fmt.Print("Listening")
	client.Connect()
}

func IgnoreBotMessages(user string) bool {
	return slices.Contains([]string{"oathybot", "funtoon", "cynanbot", TWITCH_USER}, user)
}

func FilteredMessage(chain *gomarkov.Chain) string {
	response := Generate(chain)

	for strings.Contains(response, "@") || strings.Contains(response, "https://") {
		fmt.Printf("Message %s contains @, skipping.../n", response)
		response = Generate(chain) // try again
	}

	return response
}

func RunOnTimer(chain *gomarkov.Chain, interval time.Duration) {
	fmt.Println("Running on timer...")
	client := twitch.NewClient(TWITCH_USER, TWITCH_OAUTH_STRING)
	fmt.Println("Connecting to twitch...")

	client.Join(CHANNEL) // oats please pepeW
	fmt.Println("Joined " + CHANNEL)

	go client.Connect()

	for {
		message := Generate(chain)
		fmt.Println("Generated message:", message)

		if strings.Contains(message, "@") || strings.Contains(message, "https://") {
			fmt.Println("Message contains @, skipping...")
			continue
		}

		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Yes / No? ")
		// input, _ := reader.ReadString('\n')
		// input = strings.TrimSpace(strings.ToLower(input))

		// if input == "yes" {
		// 	    fmt.Println("Sending message")
		client.Say(CHANNEL, message)
		// }

		time.Sleep(interval * time.Second)

	}
}

func ReadFile(filename string) ([]string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Read file line by line
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	// Print the lines (for demonstration)
	return lines, err
}

func LoadModel() (*gomarkov.Chain, error) {
	chain := gomarkov.NewChain(1)
	fmt.Println("Adding saved messages to model...")

	lines, err := ReadFile(MESSAGE_FILE)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Loading", len(lines), "messages...")

	for _, line := range lines {
		chain.Add(strings.Split(line, " "))
	}

	fmt.Println("Done!")
	return chain, nil
}

func Generate(chain *gomarkov.Chain) string {
	tokens := []string{gomarkov.StartToken}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := chain.Generate(tokens[(len(tokens) - 1):])
		tokens = append(tokens, next)
	}
	return strings.Join(tokens[1:len(tokens)-1], " ")
}

func Listen() {
	// client := twitch.NewClient(TWITCH_USER, TWITCH_OAUTH_STRING)
	client := twitch.NewAnonymousClient() // for an anonymous user (no write capabilities)

	// Open the file in append mode, create if not exists
	file, err := os.OpenFile(MESSAGE_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		var prefix string
		if message.User.Name == TWITCH_USER {
			prefix = RED + message.User.Name + ": "
		} else if strings.Contains(message.Message, "@"+TWITCH_USER) {
			prefix = GREEN + message.User.Name + ": "
		} else {
			prefix = message.User.Name + ": "
		}

		if message.User.Name == TWITCH_USER ||
			message.User.Name == "oathybot" ||
			message.User.Name == "funtoon" ||
			message.User.Name == "cynanbot" {
			return
		}

		fmt.Println("Adding message:", prefix+message.Message+RESET)

		// messages = append(messages, message.Message)
		//
		// Write the new line
		_, err = file.WriteString(message.Message + "\n")

		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	})

	client.Join(CHANNEL) // oats please pepeW

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}

// func Markov() {
// 	//Create a chain of order 2
// 	chain := gomarkov.NewChain(10)
// 	var err error

// 	//Feed in training data
// 	for _, message := range messages {
// 		chain.Add(strings.Split(message, " "))
// 	}

// 	// From teh readme
// 	// //Get transition probability of a sequence
// 	// prob, _ := chain.TransitionProbability("a", []string{"I"})
// 	// fmt.Println(prob)
// 	// //Output: 0.6666666666666666

// 	rand.Seed(time.Now().UnixNano())
// 	randInt := rand.Intn(len(messages))
// 	randomFirstWord := strings.Split(messages[randInt], " ")[0]

// 	randLen := rand.Intn(12)
// 	generated := []string{randomFirstWord}
// 	for i := 0; i < randLen; i++ {
// 		next, err := chain.Generate(generated)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		generated = append(generated, next)
// 	}

// 	fmt.Println(generated)

// 	//The chain is JSON serializable
// 	jsonObj, _ := json.Marshal(chain)
// 	err = ioutil.WriteFile("model.json", jsonObj, 0644)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// }
