package commands

import (
	"fmt"
)

type CLICommands struct {
	name        string
	description string
	callback    func()
}

func getCommandTypes() map[string]CLICommands {
	return map[string]CLICommands{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    helpMenu,
		},
	}
}

func helpMenu() {
	fmt.Println("Here is a help message")
}

func errorMessage() {
	fmt.Println("Invalid command")
}

func Command(command string) CLICommands {
	// if the string is not in the map, then callback will return an error

	commands := getCommandTypes()
	clicommand, found := commands[command]
	if !found {
		return CLICommands{
			name:        "error",
			description: "invalid command",
			callback:    errorMessage,
		}
	} else {
		return clicommand
	}
}
