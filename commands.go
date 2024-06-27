package main

import (
	"fmt"
)

type cliCommand struct {
	name        string
	description string
	Callback    func()
}

func getCommandTypes() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    helpMenu,
		},
	}
}

func helpMenu() {
	fmt.Print("Usage:\n\n")
	commands := getCommandTypes()
	for command := range commands {
		fmt.Printf("%s: %s\n", commands[command].name, commands[command].description)
	}
	fmt.Println("exit: Exits the pokedex")
}

func errorMessage() {
	fmt.Println("Invalid command")
}

func getCommand(command string) cliCommand {

	commands := getCommandTypes()
	clicommand, found := commands[command]
	if !found {
		return cliCommand{
			name:        "error",
			description: "invalid command",
			Callback:    errorMessage,
		}
	} else {
		return clicommand
	}
}
