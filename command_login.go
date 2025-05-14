package main

import (
	"context"
	"fmt"
	"log"
)

func handlerLogin(s *state, cmd command) error {

	// Check valid command input
	if len(cmd.Arguments) != 1 {

		return fmt.Errorf("usage is login <username>")
	}

	// Define variables
	userName := cmd.Arguments[0]
	dbQueries := s.db

	// Check if the user is registered
	_, errGetUser := dbQueries.GetUser(context.Background(), userName)
	if errGetUser != nil {
		log.Fatalf("can't not login with %s username because it is not registered!\n", userName)
	}

	return nil
}
