package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joseflores1/rss/internal/database"
)

func handlerFollow(s *state, cmd command) error {

	// Check for right number of arguments
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	// Initialize appropiate variables
	currentUser := s.config.CurrentUserName
	dbQueries := s.db
	feedURL := cmd.Arguments[0]

	// Get current user
	user, errGetUser := dbQueries.GetUser(context.Background(), currentUser)
	if errGetUser != nil {
		return fmt.Errorf("couldn't get user: %w", errGetUser)
	}

	// Get feed by URL
	feed, errGetFeedURL := dbQueries.GetFeedByURL(context.Background(), feedURL)
	if errGetFeedURL != nil {
		return fmt.Errorf("couldn't get feed by uRL: %w", errGetFeedURL)
	}

	_, errGetFeedFollow := dbQueries.GetFeedFollowByIDS(context.Background(), database.GetFeedFollowByIDSParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if errGetFeedFollow == sql.ErrNoRows {
		fmt.Println("feed follow not found, recording!")
	} else {
		return fmt.Errorf("feed follow is already registered")
	}

	// Initialize feed follow struct for further creation
	dbFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	// Create feed follow record
	feedFollow, errCreateFeedFollow := dbQueries.CreateFeedFollow(context.Background(), dbFeedFollow)
	if errCreateFeedFollow != nil {
		return fmt.Errorf("couldn't create feed follow: %w", errCreateFeedFollow)
	}

	fmt.Println("Feed follow created successfully!")
	printFeedFollow(feedFollow)

	return nil
}

func printFeedFollow(feed database.CreateFeedFollowRow) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* User ID:       %s\n", feed.UserID)
	fmt.Printf("* Feed ID:       %s\n", feed.FeedID)
	fmt.Printf("* Feed:          %s\n", feed.FeedName)
	fmt.Printf("* User:          %s\n", feed.UserName)

}
