package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {

	// Check ausence of args
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("%s doesn't expect any arguments", cmd.Name)
	}

	// Delete database users
	dbQueries := s.db
	errTruncTable := dbQueries.DeleteUsers(context.Background())
	if errTruncTable != nil {
		return fmt.Errorf("couldn't truncate users table: %w", errTruncTable)
	}
	return nil
}
