package chain

import (
	"fmt"
	"os"
)

/**
 * Save a chat message to the chat message file
 * */
func (c *Chain) SaveChatMessage(message string) error {
	err := c.Save(c.savedMessagesFilepath, message)
	if err != nil {
		return fmt.Errorf("error while saving chat message: %v", err)
	}

	return nil
}

/**
 * Save a sent message to the sent message file
 * */
func (c *Chain) SaveSentMessage(message string) error {
	err := c.Save(c.sentMessagesFilepath, message)
	if err != nil {
		return fmt.Errorf("error while saving sent message: %v", err)
	}

	return nil
}

/**
 * Handle saving a message to a file
 * */
func (c *Chain) Save(path, message string) error {
	if c.Saving {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file %s while saving message: %v", path, err)
		}
		defer file.Close()

		_, err = file.WriteString(message + "\n")
		if err != nil {
			return fmt.Errorf("failed to save message %s to %s: %v", message, path, err)
		}
	}

	return nil
}
