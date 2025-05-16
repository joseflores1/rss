package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joseflores1/rss/internal/config"
	"github.com/joseflores1/rss/internal/database"
	_ "github.com/lib/pq"
)

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
		log.Fatalf("error when trying to open a connection to the database: %v", errOpen)
	}
	dbQueries := database.New(db)

	// Initialize necessary structs
	stateStruct := &state{db: dbQueries, config: &configStruct}
	commandMap := make(map[string]func(*state, command) error)
	commandsStruct := commands{commandList: commandMap}

	// Register commands
	commandsStruct.register("register", handlerRegister)
	commandsStruct.register("login", handlerLogin)
	commandsStruct.register("reset", handlerReset)
	commandsStruct.register("users", handlerUsers)
	commandsStruct.register("agg", handlerAgg)
	commandsStruct.register("addfeed", handlerAddFeed)
	commandsStruct.register("feeds", handlerFeeds)

	// Get CLI args
	var commandName string
	var cliArgs []string
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]\n")
	} else {
		commandName = os.Args[1]
		cliArgs = os.Args[2:]
	}

	// Run command
	commandStruct := command{Name: commandName, Arguments: cliArgs}
	errRun := commandsStruct.run(stateStruct, commandStruct)
	if errRun != nil {
		log.Fatalf("error when trying to run %s command with %+v arguments: %s\n", commandName, cliArgs, errRun.Error())
	}

}
