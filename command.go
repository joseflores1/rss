package main

import (
	"fmt"

	"github.com/joseflores1/rss/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commandList map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {

	commandName := cmd.name

	handler, ok := c.commandList[commandName]
	if ok {
		errCommand := handler(s, cmd)
		if errCommand != nil {
			return fmt.Errorf("error when trying to run %s command: %w", commandName, errCommand)
		}
	} else {
		return fmt.Errorf("error: %s does not exist", commandName)
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {

	c.commandList[name] = f

}
