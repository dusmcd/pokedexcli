package main

import (
	"github.com/dusmcd/pokedexcli/cache"
	"github.com/dusmcd/pokedexcli/pokeapi"
)

type config struct {
	next     string
	previous string
	page     int
	argument string
	pokedex  pokedex
}

func (c *config) setPage(command string) {
	if command == "next" {
		c.page++
	} else if command == "previous" {
		c.page--
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *config, cache *cache.Cache) error
}

type pokedex struct {
	data map[string]pokeapi.PokemonStats
}

func (p *pokedex) add(pokemonName string, stats pokeapi.PokemonStats) {
	p.data[pokemonName] = stats
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
		"explore": {
			name:        "explore <location>",
			description: "Shows the pokemon found in the given <location>",
			callback:    showPokemonInLocation,
		},
		"catch": {
			name:        "catch <pokemon>",
			description: "Catch the specified <pokemon> according to random chance and the pokemon's base experience",
			callback:    catchPokemon,
		},
		"inspect": {
			name:        "inspect <pokemon>",
			description: "View information about a given <pokemon>",
			callback:    showPokemonInfo,
		},
	}
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
