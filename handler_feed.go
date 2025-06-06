package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joseflores1/rss/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {

	// Check valid command input
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
	}

	// Define variables
	dbQueries := s.db
	feedName := cmd.Arguments[0]
	feedURL := cmd.Arguments[1]

	// Get feed by URL
	_, errGetFeedBYURL := dbQueries.GetFeedByURL(context.Background(), feedURL)
	if errGetFeedBYURL == sql.ErrNoRows {
		fmt.Printf("Feed not found, adding!\n")
	} else {
		return fmt.Errorf("feed is already added")
	}

	// Create feed into database
	dbFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}
	createdFeed, errCreatedFeed := dbQueries.CreateFeed(context.Background(), dbFeed)
	if errCreatedFeed != nil {
		return fmt.Errorf("couldn't create '%s' feed from %s: %w", createdFeed.Name, createdFeed.Url, errCreatedFeed)
	}

	// Create feed follow into database
	dbFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    createdFeed.ID,
	}
	_, errCreateFeedFollow := dbQueries.CreateFeedFollow(context.Background(), dbFeedFollow)
	if errCreateFeedFollow != nil {
		return fmt.Errorf("couldn't create feed follow: %w", errCreateFeedFollow)
	}

	// Print and return normally
	fmt.Println("Feed created successfully!")
	printFeed(createdFeed, user)
	return nil
}

func handlerFeeds(s *state, cmd command) error {

	// Check for valid command input
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("%s doesn't accept any arguments", cmd.Name)
	}

	// Define variables
	dbQueries := s.db

	// Get Feeds slice
	feedsSlice, errGetFeeds := dbQueries.GetFeeds(context.Background())
	if errGetFeeds != nil {
		return fmt.Errorf("couldn't get feeds: %w", errGetFeeds)
	}

	// Print list of feeds with their name, URL and creator's name
	if len(feedsSlice) == 0 {
		fmt.Println("There are no registered feeds!")
		return nil
	}
	for i, feed := range feedsSlice {
		fmt.Printf("Feed %d:\n", i + 1)
		user, errGetUser := dbQueries.GetUserById(context.Background(), feed.UserID)
		if errGetUser != nil {
			return fmt.Errorf("couldn't get feed's user: %w", errGetUser)
		}
		printFeed(feed, user)
		fmt.Println("------------------------------------")
	}

	// Return normally
	return nil
}

func printFeed(feed database.Feed, user database.User) {

	// Print feed's info and related username
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
