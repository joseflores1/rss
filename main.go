package main

import (
	"log"
	"os"

	"github.com/joseflores1/rss/internal/config"
)

type state struct {
	config *config.Config
}

func main() {

	// Read config file
	configStruct, errRead := config.Read()
	if errRead != nil {
		log.Fatal(errRead)
	}

	// Initialize necessary structs
	stateStruct := &state{config: &configStruct}
	commandMap := make(map[string]func(*state, command) error)
	commandsStruct := commands{commandList: commandMap}
	commandsStruct.register("login", handlerLogin)

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
