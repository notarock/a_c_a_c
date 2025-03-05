package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/mb-14/gomarkov"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

var messages []string

const MESSAGE_FILE = "messages.txt"

func main() {
	// Listen()
	// fmt.Println("Listening for messages...")
	// time.Sleep(10 * time.Second)
	// fmt.Println("Generating Markov chain...")
	// Markov()

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

	RunOnTimer(chain, 5) // says things like "oats you are not so washed!"
}

func RunOnTimer(chain *gomarkov.Chain, interval time.Duration) {
	for {
		time.Sleep(interval * time.Second)

		fmt.Println(Generate(chain))

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
	fmt.Println("File contents as array of strings:")
	return lines, err
}

func LoadModel() (*gomarkov.Chain, error) {
	chain := gomarkov.NewChain(1)
	fmt.Println("Adding saved messages to model...")

	filename := "messages.txt"
	lines, err := ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, line := range lines {
		chain.Add(strings.Split(line, " "))
	}
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
	client := twitch.NewClient("notarock95", "oauth:123123123")
	// client := twitch.NewAnonymousClient() // for an anonymous user (no write capabilities)

	filename := "messages.txt"

	// Open the file in append mode, create if not exists
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println("Adding message:", message.User.Name+": "+message.Message)
		messages = append(messages, message.Message)
		// Write the new line
		_, err = file.WriteString(message.Message + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	})

	client.Join("oatsngoats") // oats please pepeW

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}

func Markov() {
	//Create a chain of order 2
	chain := gomarkov.NewChain(10)
	var err error

	//Feed in training data
	for _, message := range messages {
		chain.Add(strings.Split(message, " "))
	}

	// From teh readme
	// //Get transition probability of a sequence
	// prob, _ := chain.TransitionProbability("a", []string{"I"})
	// fmt.Println(prob)
	// //Output: 0.6666666666666666

	rand.Seed(time.Now().UnixNano())
	randInt := rand.Intn(len(messages))
	randomFirstWord := strings.Split(messages[randInt], " ")[0]

	randLen := rand.Intn(12)
	generated := []string{randomFirstWord}
	for i := 0; i < randLen; i++ {
		next, err := chain.Generate(generated)
		if err != nil {
			fmt.Println(err)
		}
		generated = append(generated, next)
	}

	fmt.Println(generated)

	//The chain is JSON serializable
	jsonObj, _ := json.Marshal(chain)
	err = ioutil.WriteFile("model.json", jsonObj, 0644)
	if err != nil {
		fmt.Println(err)
	}

}
