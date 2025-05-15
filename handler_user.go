package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joseflores1/rss/internal/config"
	"github.com/joseflores1/rss/internal/database"
)

func handlerRegister(s *state, cmd command) error {
	// Check for valid command input
	if len(cmd.Arguments) != 1 {

		return fmt.Errorf("usage is %s <username>", cmd.Name)
	}

	// Define variables
	dbQueries := s.db
	userName := cmd.Arguments[0]

	// Check if the user is registered
	_, errGetUser := dbQueries.GetUser(context.Background(), userName)
	if errGetUser == nil {
		return fmt.Errorf("'%s' username is already registered", userName)
	}

	// Create user into database
	dbUser := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Name: userName}
	createdUser, errCreatedUser := dbQueries.CreateUser(context.Background(), dbUser)
	if errCreatedUser != nil {
		return fmt.Errorf("couldn't create '%s' user: %w", userName, errCreatedUser)
	}

	// Set username in .JSON config file
	errSetUser := s.config.SetUser(userName)
	if errSetUser != nil {
		return fmt.Errorf("couldn't set '%s' username to .JSON config file: %w", userName, errSetUser)
	}
	fmt.Println("User created successfully:")
	printUser(createdUser)
	return nil

}

func handlerLogin(s *state, cmd command) error {

	// Check valid command input
	if len(cmd.Arguments) != 1 {

		return fmt.Errorf("usage is %s <username>", cmd.Name)
	}

	// Define variables
	userName := cmd.Arguments[0]
	dbQueries := s.db

	// Check if the user is registered
	_, errGetUser := dbQueries.GetUser(context.Background(), userName)
	if errGetUser != nil {
		return fmt.Errorf("can't login with %s username because it is not registered: %w", userName, errGetUser)
	}

	// Set username in .JSON config file
	errSetUser := s.config.SetUser(userName)
	if errSetUser != nil {
		return fmt.Errorf("couldn't set %s username to .JSON config file: %w", userName, errSetUser)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handlerUsers(s *state, cmd command) error {

	// Check for no arguments
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("%s doesn't expect any arguments", cmd.Name)
	}

	// Read config file to get active user's name
	dbQueries := s.db
	configStruct, errRead := config.Read()
	if errRead != nil {
		return fmt.Errorf("couldn't read config file within users handler: %w", errRead)
	}
	currentName := configStruct.CurrentUserName

	// Get Users slice
	usersSlice, errGetUsers := dbQueries.GetUsers(context.Background())
	if errGetUsers != nil {
		return fmt.Errorf("couldn't get users: %s", errGetUsers)
	}

	// Print slice of users
	if len(usersSlice) == 0 {
		fmt.Println("There are no registered users!")
	}
	for _, user := range usersSlice {
		if user.Name == currentName {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}

	return nil
}

func printUser(user database.User) {

	// Print user's ID and Name
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}

