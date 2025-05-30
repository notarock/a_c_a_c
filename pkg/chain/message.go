package chain

import (
	"fmt"
	"strings"

	"github.com/mb-14/gomarkov"
	"github.com/notarock/a_c_a_c/pkg/filters"
)

const RED = "\033[31m"
const BLUE = "\033[34m"
const RESET = "\033[0m"

/**
 * Generate a message that was filtered for prohibited content
 * */
func (c *Chain) GenerateValidMessage(filters []filters.Filter) string {
	response := c.generateMessage()

	for runFilters(response, filters) {
		fmt.Println(RED, "Message \"", response, "\" filtered, generating new response...", RESET)
		response = c.generateMessage()
	}

	c.lastMessage = response

	return response
}

// runFilters checks if any of the filters return true for the given message.
func runFilters(message string, filters []filters.Filter) bool {
	for _, f := range filters {
		if f.Filter(message) {
			return true
		}
	}
	return false
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
 * Check if a message is similar to what was just sent
 * */
func (c *Chain) IsParrot(message string) bool {
	return c.lastMessage == message
}
