package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joseflores1/rss/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {	

	// Check for right number of arguments
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
	}

	// Initialize appropiate variables
	currentUser := s.config.CurrentUserName
	dbQueries := s.db
	feedName := cmd.Arguments[0]
	feedURL := cmd.Arguments[1]

	// Get current user
	user, errGetUser := dbQueries.GetUser(context.Background(), currentUser)
	if errGetUser != nil {
		return fmt.Errorf("couldn't get user: %w", errGetUser)
	}

	// Create feed into database
	dbFeed := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: feedName,
		Url: feedURL,
		UserID: user.ID,
	}

	createdFeed, errCreatedFeed := dbQueries.CreateFeed(context.Background(), dbFeed)
	if errCreatedFeed != nil {
		return fmt.Errorf("couldn't create '%s' feed from %s: %w", createdFeed.Name, createdFeed.Url, errCreatedFeed)
	}

	// Print feed to stdout
	fmt.Println("Feed created successfully!")
	printFeed(createdFeed)

	return nil
}

func printFeed(feed database.Feed) {
	// Print user's ID and Name
	fmt.Printf(" * ID:      %v\n", feed.ID)
	fmt.Printf(" * Name:    %s\n", feed.Name)
	fmt.Printf(" * URL:     %s\n", feed.Url)


}