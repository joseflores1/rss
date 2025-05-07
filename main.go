package main

import (
	"fmt"
	"github.com/joseflores1/rss/internal/config"
)

const (
	CURRENT_USER = "josei"
)
func main() {
	// Read file and get Config struct
	configStruct, errRead := config.Read()
	if errRead != nil {
		fmt.Println(errRead)
	}
	// Set username of Config struct and write it to disk	
	errSetUser := configStruct.SetUser(CURRENT_USER)
	if errSetUser != nil {
		fmt.Println(errSetUser)
	}
	// Read updated file and get new Config struct
	newConfigStruct, errNewRead := config.Read()
	if errNewRead != nil {
		fmt.Println(errNewRead)
	}
	// Print config struct's contents
	fmt.Printf("The contents of the config struct are: %+v\n", newConfigStruct)
}
