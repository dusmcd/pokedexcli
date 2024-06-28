package main

import (
	"fmt"
	"log"

	"github.com/dusmcd/pokedexcli/pokeapi"
)

type config struct {
	next     string
	previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *config)
}

/*
expose the map of the different commands that are available to the user
*/
func getCommandTypes() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    helpMenu,
		},
		"map": {
			name:        "map",
			description: "Displays the name of 20 location areas in the Pokemon world. Each subsequent call will display the next 20 locations, and so on.",
			callback:    showNextLocations,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the name of the previous 20 location areas.",
			callback:    showPreviousLocations,
		},
	}
}

func showPreviousLocations(config *config) {
	if config.previous == "" {
		fmt.Println("Previous page does not exist")
		return
	}
	location, err := pokeapi.GetLocationData(config.previous)
	if err != nil {
		log.Fatal(err)
	}
	if location.Previous != nil {
		config.previous = location.Previous.(string)
	} else {
		config.previous = ""
	}
	config.next = location.Next
	for _, result := range location.Results {
		fmt.Printf("Name: %s, URL: %s\n", result.Name, result.URL)
	}
	fmt.Print("\n")
}

func showNextLocations(config *config) {
	location, err := pokeapi.GetLocationData(config.next)
	if err != nil {
		log.Fatal(err)
	}
	config.previous = config.next
	config.next = location.Next
	for _, result := range location.Results {
		fmt.Printf("Name: %s, URL: %s\n", result.Name, result.URL)
	}
	fmt.Print("\n")
}

func helpMenu(config *config) {
	fmt.Print("Usage:\n\n")
	commands := getCommandTypes()
	for command := range commands {
		fmt.Printf("%s: %s\n", commands[command].name, commands[command].description)
	}
	fmt.Println("exit: Exits the pokedex")
	fmt.Print("\n")
}

func errorMessage(config *config) {
	fmt.Println("Invalid command")
}

func getCommand(command string) cliCommand {

	commands := getCommandTypes()
	clicommand, found := commands[command]
	if !found {
		return cliCommand{
			name:        "error",
			description: "invalid command",
			callback:    errorMessage,
		}
	} else {
		return clicommand
	}
}
