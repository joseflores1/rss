package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joseflores1/rss/internal/config"
	"github.com/joseflores1/rss/internal/database"

	_ "github.com/lib/pq"
)

// Save state in the form of queries, database URL and current user's name
type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {

	// Read config file
	configStruct, errRead := config.Read()
	if errRead != nil {
		log.Fatal(errRead)
	}

	// Open the database connection
	db, errOpen := sql.Open("postgres", configStruct.DBURL)
	if errOpen != nil {
		log.Fatalf("couldn't open a connection to the database: %v", errOpen)
	}

	// Define variables
	commandMap := make(map[string]func(*state, command) error)
	commandsStruct := commands{
		commandList: commandMap,
	}
	dbQueries := database.New(db)
	stateStruct := &state{
		db: dbQueries, 
		config: &configStruct,
	}

	// Register commands
	commandsStruct.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commandsStruct.register("agg", handlerAgg)
	commandsStruct.register("browse", middlewareLoggedIn(handlerBrowse))
	commandsStruct.register("feeds", handlerFeeds)
	commandsStruct.register("follow", middlewareLoggedIn(handlerFollow))
	commandsStruct.register("following", middlewareLoggedIn(handlerFollowing))
	commandsStruct.register("login", handlerLogin)
	commandsStruct.register("register", handlerRegister)
	commandsStruct.register("reset", handlerReset)
	commandsStruct.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commandsStruct.register("users", handlerUsers)

	// Get CLI args
	var commandName string
	var cliArgs []string
	if len(os.Args) < 2 {
		log.Fatal("Usage: rss <command> [args...]\n")
	} else {
		commandName = os.Args[1]
		cliArgs = os.Args[2:]
	}

	// Run command
	commandStruct := command{
		Name: commandName, 
		Arguments: cliArgs,
	}
	errRun := commandsStruct.run(stateStruct, commandStruct)
	if errRun != nil {
		log.Fatalf("couldn't run %s command with %+v arguments: %s\n", commandName, cliArgs, errRun.Error())
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {

		// Get current user
		user, errGetUser := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if errGetUser != nil {
			return fmt.Errorf("couldn't get user: %w", errGetUser)
		}

		// Return normally
		return handler(s, cmd, user)
	}
}