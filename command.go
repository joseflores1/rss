package main

import (
	"fmt"
)

// Struct to save command's name and its arguments
type command struct {
	Name      string
	Arguments []string
}

// Struct to save command's list
type commands struct {
	commandList map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {

	// Check if the command exists already
	c.commandList[name] = f
}

func (c *commands) run(s *state, cmd command) error {

	// Get command's name
	commandName := cmd.Name

	// Run command's handler if the command exists
	handler, ok := c.commandList[commandName]
	if ok {
		errCommand := handler(s, cmd)
		if errCommand != nil {
			return errCommand
		}
	} else {
		return fmt.Errorf("%s does not exist", commandName)
	}

	// Return normally
	return nil
}