package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dusmcd/pokedexcli/cache"
	"github.com/dusmcd/pokedexcli/pokeapi"
)

type config struct {
	next     string
	previous string
	page     int
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
	callback    func(config *config, cache *cache.Cache)
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

func showPreviousLocations(config *config, cacheStruct *cache.Cache) {
	if config.previous == "" {
		fmt.Println("Previous page does not exist")
		return
	}
	config.setPage("previous")
	location := pokeapi.Location{}
	var err error
	var rawData []byte

	// checking cache
	ch := make(chan cache.CacheData)
	go cacheStruct.GetEntry(fmt.Sprintf("Page %d", config.page), ch)
	cacheData := <-ch
	if cacheData.Found {
		err = json.Unmarshal(cacheData.Val, &location)
	} else {
		location, rawData, err = pokeapi.GetLocationData(config.previous)
		go cacheStruct.AddEntry(fmt.Sprintf("Page %d", config.page), rawData)
	}

	if err != nil {
		log.Fatal(err)
		return
	}
	if location.Previous != nil {
		config.previous = location.Previous.(string)
	} else {
		config.previous = ""
	}
	config.next = location.Next
	fmt.Printf("Page %d\n", config.page)
	for _, result := range location.Results {
		fmt.Println(result.Name)
	}
	fmt.Print("\n")
}

func showNextLocations(config *config, cacheStruct *cache.Cache) {
	var err error
	location := pokeapi.Location{}
	var rawData []byte
	config.setPage("next")

	// checking cache
	ch := make(chan cache.CacheData)
	go cacheStruct.GetEntry(fmt.Sprintf("Page %d", config.page), ch)
	cacheData := <-ch
	if cacheData.Found {
		err = json.Unmarshal(cacheData.Val, &location)
	} else {
		location, rawData, err = pokeapi.GetLocationData(config.next)
		go cacheStruct.AddEntry(fmt.Sprintf("Page %d", config.page), rawData)
	}

	if err != nil {
		log.Fatal(err)
		return
	}
	if location.Previous == nil {
		config.previous = ""
	} else {
		config.previous = location.Previous.(string)
	}
	config.next = location.Next
	fmt.Printf("Page %d\n", config.page)
	for _, result := range location.Results {
		fmt.Println(result.Name)
	}
	fmt.Print("\n")
}

func helpMenu(config *config, cache *cache.Cache) {
	fmt.Print("Usage:\n\n")
	commands := getCommandTypes()
	for command := range commands {
		fmt.Printf("%s: %s\n", commands[command].name, commands[command].description)
	}
	fmt.Println("exit: Exits the pokedex")
	fmt.Print("\n")
}

func errorMessage(config *config, cache *cache.Cache) {
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
