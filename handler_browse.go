package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/joseflores1/rss/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {

	// Set default limit
	limit := 2

	// Check for valid command input
	if len(cmd.Arguments) > 1 {
		return fmt.Errorf("usage: %s [limit]", cmd.Name)
	} else if len(cmd.Arguments) == 1 {
		parsedInt, errAtoi := strconv.Atoi(cmd.Arguments[0])
		if errAtoi != nil {
			return fmt.Errorf("couldn't parse limit argument: %w", errAtoi)
		}
		limit = parsedInt
	}

	// Define variables
	dbQueries := s.db

	// Get posts for usesr
	posts, errGetPosts := dbQueries.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID, 
		Limit: int32(limit)},
	)
	if errGetPosts != nil {
		return fmt.Errorf("couldn't get posts for user: %w", errGetPosts)
	}

	// Print posts and return normally
	for i, post := range posts {
		fmt.Printf("POST %v\n", i+1)
		printPost(post)
		fmt.Println("--------------------------------")
	}

	return nil
}