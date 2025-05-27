package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joseflores1/rss/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {

	// Check for valid command input
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	// Define variables
	dbQueries := s.db
	feedURL := cmd.Arguments[0]

	// Get feed by URL
	feed, errGetFeedURL := dbQueries.GetFeedByURL(context.Background(), feedURL)
	if errGetFeedURL != nil {
		return fmt.Errorf("couldn't get feed by URL: %w", errGetFeedURL)
	}

	// Get feed follow by user and feed IDs to check for duplicates
	_, errGetFeedFollow := dbQueries.GetFeedFollowByIDS(context.Background(), database.GetFeedFollowByIDSParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if errGetFeedFollow == sql.ErrNoRows {
		fmt.Println("Feed follow not found, registering!")
	} else {
		return fmt.Errorf("feed follow is already registered")
	}

	// Create feed follow record
	dbFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, errCreateFeedFollow := dbQueries.CreateFeedFollow(context.Background(), dbFeedFollow)
	if errCreateFeedFollow != nil {
		return fmt.Errorf("couldn't create feed follow: %w", errCreateFeedFollow)
	}

	// Print and return normally
	fmt.Println("Feed follow created successfully!")
	printFeedFollow(feedFollow)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	// Check for valid command input
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("%s doesn't accept any arguments", cmd.Name)
	}

	// Define variables
	currentUser := s.config.CurrentUserName
	dbQueries := s.db

	// Get all feed follow from an user
	feedFollowsSlice, errGetFeedFollowSlice := dbQueries.GetFeedFollowsForUser(context.Background(), user.ID)
	if errGetFeedFollowSlice != nil {
		return fmt.Errorf("couldn't get feed follows for %s user: %w", currentUser, errGetFeedFollowSlice)
	}

	// Print feed follow slice
	if len(feedFollowsSlice) == 0 {
		fmt.Println("There are no registered feed follows for this user!")
		return nil
	}
	printFeedFollowsByUser(feedFollowsSlice, currentUser)

	// Return normally
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {

	// Check for valid command input
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage %s <feed_url>", cmd.Name)
	}

	// Define variables
	dbQueries := s.db
	feedURL := cmd.Arguments[0]

	// Get feed by URL
	feed, errGetFeed := dbQueries.GetFeedByURL(context.Background(), feedURL)
	if errGetFeed != nil {
		return fmt.Errorf("couldn't get feed by url: %w", errGetFeed)
	}

	// Unfollow feed
	errUnfollow := dbQueries.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if errUnfollow != nil {
		return fmt.Errorf("couldn't unfollow feed: %w", errUnfollow)
	}

	// Return normally
	return nil
}

func printFeedFollow(feed database.CreateFeedFollowRow) {

	// Print a single feed follow's info
	fmt.Printf("*****************************\n")
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* User ID:       %s\n", feed.UserID)
	fmt.Printf("* Feed ID:       %s\n", feed.FeedID)
	fmt.Printf("* Feed:          %s\n", feed.FeedName)
	fmt.Printf("* User:          %s\n", feed.UserName)
	fmt.Printf("*****************************\n")
}

func printFeedFollowsByUser(feedsSlice []database.GetFeedFollowsForUserRow, currentUsername string) {

	// Print all of the feed follow of an user
	for i, feedFollow := range feedsSlice {
		fmt.Printf("Feed %d:\n\n", i + 1)
		fmt.Printf("* ID:            %s\n", feedFollow.ID)
		fmt.Printf("* Created:       %v\n", feedFollow.CreatedAt)
		fmt.Printf("* Updated:       %v\n", feedFollow.UpdatedAt)
		fmt.Printf("* User ID:       %s\n", feedFollow.UserID)
		fmt.Printf("* Feed ID:       %s\n", feedFollow.FeedID)
		fmt.Printf("* Feed:          %s\n", feedFollow.FeedName)
		fmt.Printf("* User:          %s\n", currentUsername)
		fmt.Printf("*****************************************************************\n")
	}
}