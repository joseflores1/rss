package main

import (
	"log"
	"os"

	"github.com/joseflores1/rss/internal/config"
)

const (
	CURRENT_USER = "josei"
)

func main() {
	configStruct, errRead := config.Read()
	if errRead != nil {
		log.Fatal(errRead)
	}

	stateStruct := state{config: &configStruct}

	commandMap := make(map[string]func(*state, command) error)

	commandsStruct := commands{commandList: commandMap}

	commandsStruct.register("login", handlerLogin)

	var commandName string
	var cliArgs []string

	if len(os.Args) < 3 {
		log.Fatal("error: at least 2 arguments must be provided to the CLI\n")
	} else {
		commandName = os.Args[1]
		cliArgs = os.Args[2:]
	}

	errRun := commandsStruct.run(&stateStruct, command{name: commandName, arguments: cliArgs})
	if errRun != nil {
		log.Fatalf("error when trying to run %s command with %+v arguments: %s\n", commandName, cliArgs, errRun.Error())
	}

}
