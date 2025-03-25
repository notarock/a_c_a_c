package chain

import (
	"fmt"
	"strings"

	"github.com/mb-14/gomarkov"
)

/**
 * Chain string holds the markov chain and information related to
 * the chain such as where to save incomming chat messages, sent messages and wether or not we are saving them.
 * */
type Chain struct {
	chain                 *gomarkov.Chain
	sentMessagesFilepath  string
	savedMessagesFilepath string
	Saving                bool
	IgnoreParrots         bool

	// To keep track of the last message
	lastMessage string
}

type ChainConfig struct {
	SentMessagesFilepath  string
	SavedMessagesFilepath string
	Saving                bool
	IgnoreParrots         bool
}

func NewChain(config ChainConfig) (*Chain, error) {
	chain := gomarkov.NewChain(1) // Not sure what the magic 1 is for

	c := &Chain{
		chain:                 chain,
		sentMessagesFilepath:  config.SentMessagesFilepath,
		savedMessagesFilepath: config.SavedMessagesFilepath,
		Saving:                config.Saving,
		IgnoreParrots:         config.IgnoreParrots,
	}

	err := c.LoadModel()

	return c, err
}

func (c *Chain) AddMessage(message string) {
	c.chain.Add(strings.Split(message, " "))
}

func (c *Chain) LoadModel() error {
	lines, err := ReadFile(c.savedMessagesFilepath)
	if err != nil {
		return fmt.Errorf("failed to load previous messages from file %s: %v:", c.savedMessagesFilepath, err)
	}

	for _, line := range lines {
		c.chain.Add(strings.Split(line, " "))
	}
	fmt.Println("Loaded", len(lines), "messages...")

	return nil
}
