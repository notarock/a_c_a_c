package chain

import (
	"fmt"
	"strings"

	"github.com/mb-14/gomarkov"
)

/**
 * Generate a message that was filtered for prohibited content
 * */
func (c *Chain) FilteredMessage() string {
	response := c.generateMessage()

	for !c.validMessage(response) {
		fmt.Printf("Message '%s' prohibited content, skipping.../n", response)
		response = c.generateMessage()
	}

	c.lastMessage = response

	return response
}

/**
 * Generate a message from the chain
 * */
func (c *Chain) generateMessage() string {
	tokens := []string{gomarkov.StartToken}

	for tokens[len(tokens)-1] != gomarkov.EndToken {
		next, _ := c.chain.Generate(tokens[(len(tokens) - 1):])
		tokens = append(tokens, next)
	}

	return strings.Join(tokens[1:len(tokens)-1], " ")
}

/**
 * Validate a message against prohibited patterns
 * */
func (c *Chain) validMessage(message string) bool {
	for _, prohibitedMessage := range c.ProhibitedMessages {
		if message == prohibitedMessage {
			return false
		}
	}

	for _, pattern := range c.ProhibitedStrings {
		if strings.Contains(message, pattern) {
			return false
		}
	}

	return true
}

/**
 * Check if a message is similar to what was just sent
 * */
func (c *Chain) IsParrot(message string) bool {
	return c.lastMessage == message
}
