package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joseflores1/rss/internal/database"
)

func handlerRegister(s *state, cmd command) error {
	// Check for valid command input
	if len(cmd.Arguments) != 1 {

		return fmt.Errorf("usage is register <username>")
	}

	// Define variables
	dbQueries := s.db
	userName := cmd.Arguments[0]

	// Check if the user is registered
	_, errGetUser := dbQueries.GetUser(context.Background(), userName)
	if errGetUser == nil {
		log.Fatalf("%s username is already registered\n", userName)
	}

	// Create user into database
	dbUser := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: userName}
	createdUser, errCreatedUser := dbQueries.CreateUser(context.Background(), dbUser)
	if errCreatedUser != nil {
		log.Fatalf("couldn't create user with %s username\n", userName)
	}

	// Set username in .JSON config file
	errSetUser := s.config.SetUser(userName)
	if errSetUser != nil {
		log.Fatalf("couldn't set %s username to .JSON config file\n", userName)
	}
	fmt.Printf("%v was created successfully into the database!\n", createdUser)

	return nil

}
