package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {

	// Check valid command input
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("%s doesn't expect any arguments", cmd.Name)
	}

	// Define variables
	dbQueries := s.db

	// Delete database users
	errTruncTable := dbQueries.DeleteUsers(context.Background())
	if errTruncTable != nil {
		return fmt.Errorf("couldn't truncate users table: %w", errTruncTable)
	}

	// Return normally
	return nil
}
