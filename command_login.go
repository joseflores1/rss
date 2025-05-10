package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {

	if len(cmd.Arguments) != 1 {

		return fmt.Errorf("usage is login <username>")
	}

	userName := cmd.Arguments[0]
	errSetUser := s.config.SetUser(userName)
	if errSetUser != nil {
		return fmt.Errorf("error when trying to set username in login handler: %w", errSetUser)
	}

	fmt.Printf("The %s username has been set by the login handler!\n", userName)

	return nil
}
