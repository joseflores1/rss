package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {

	if len(cmd.arguments) != 1 {

		return fmt.Errorf("error: gator login <username> expects only an <username> argument")
	}

	userName := cmd.arguments[0]
	errSetUser := s.config.SetUser(userName)
	if errSetUser != nil {
		return fmt.Errorf("error when trying to set username in login handler: %w", errSetUser)
	}

	fmt.Printf("The %s username has been set by the login handler!\n", userName)

	return nil
}
