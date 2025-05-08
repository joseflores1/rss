package main

import (
	"fmt"
	"log"

	"github.com/joseflores1/rss/internal/config"
)

const (
	CURRENT_USER = "josei"
)

func main() {
	// Read file and get Config struct
	configStruct, errRead := config.Read()
	if errRead != nil {
		log.Fatal(errRead)
	}

	fmt.Printf("The contents of the initial config struct are: %+v\n", configStruct)

	// Set username of Config struct and write it to disk
	errSetUser := configStruct.SetUser(CURRENT_USER)
	if errSetUser != nil {
		log.Fatal(errSetUser)
	}

	// Read updated file and get new Config struct
	newConfigStruct, errNewRead := config.Read()
	if errNewRead != nil {
		log.Fatal(errNewRead)
	}

	// Print config struct's contents
	fmt.Printf("The contents of the config struct are: %+v\n", newConfigStruct)
}
