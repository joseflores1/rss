package main

import (
	"fmt"
)

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	commandList map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {

	commandName := cmd.Name

	handler, ok := c.commandList[commandName]
	if ok {
		errCommand := handler(s, cmd)
		if errCommand != nil {
			return errCommand
		}
	} else {
		return fmt.Errorf("%s does not exist", commandName)
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandList[name] = f
}
