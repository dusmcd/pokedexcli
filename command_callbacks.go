package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dusmcd/pokedexcli/cache"
	"github.com/dusmcd/pokedexcli/pokeapi"
)

/*
callback function for mapb command
*/
func showPreviousLocations(config *config, cacheStruct *cache.Cache) error {
	if config.previous == "" {
		return errors.New("previous page does not exist")
	}
	config.setPage("previous")
	ch := make(chan cache.CacheData)
	location, err := getPreviousLocations(ch, config.page, config.previous, cacheStruct)
	if err != nil {
		return err
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
	return nil
}

func getPreviousLocations(ch chan cache.CacheData, page int, prevUrl string, cacheStruct *cache.Cache) (pokeapi.Location, error) {
	location := pokeapi.Location{}
	var err error
	var rawData []byte

	// checking cache
	go cacheStruct.GetEntry(fmt.Sprintf("Page %d", page), ch)
	cacheData := <-ch
	if cacheData.Found {
		err = json.Unmarshal(cacheData.Val, &location)
	} else {
		location, rawData, err = pokeapi.GetLocationData(prevUrl)
		go cacheStruct.AddEntry(fmt.Sprintf("Page %d", page), rawData)
	}

	if err != nil {
		return location, err
	}
	return location, nil

}

/*
callback function for map command
*/
func showNextLocations(config *config, cacheStruct *cache.Cache) error {
	config.setPage("next")

	ch := make(chan cache.CacheData)
	location, err := getNextLocations(ch, config.page, config.next, cacheStruct)
	if err != nil {
		return err
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
	return nil
}

func getNextLocations(ch chan cache.CacheData, page int, nextUrl string, cacheStruct *cache.Cache) (pokeapi.Location, error) {
	var err error
	location := pokeapi.Location{}
	var rawData []byte

	// check cache
	go cacheStruct.GetEntry(fmt.Sprintf("Page %d", page), ch)
	cacheData := <-ch
	if cacheData.Found {
		err = json.Unmarshal(cacheData.Val, &location)
	} else {
		location, rawData, err = pokeapi.GetLocationData(nextUrl)
		go cacheStruct.AddEntry(fmt.Sprintf("Page %d", page), rawData)
	}

	if err != nil {
		return location, err
	}
	return location, nil

}

/*
callback function for help command
*/
func helpMenu(config *config, cache *cache.Cache) error {
	fmt.Print("Usage:\n\n")
	commands := getCommandTypes()
	for command := range commands {
		fmt.Printf("%s: %s\n", commands[command].name, commands[command].description)
	}
	fmt.Println("exit: Exits the pokedex")
	fmt.Print("\n")
	return nil
}

/*
callback function for invalid command
*/
func errorMessage(config *config, cache *cache.Cache) error {
	fmt.Println("Invalid command")
	return nil
}

func getPokemonInLocation(config *config, cacheStruct *cache.Cache) (pokeapi.Pokemon, error) {
	url := "https://pokeapi.co/api/v2/location/" + config.argument
	var err error
	//var rawData []byte
	ch := make(chan cache.CacheData)
	go cacheStruct.GetEntry(config.argument, ch)
	cacheData := <-ch
	pokemon := pokeapi.Pokemon{}
	if cacheData.Found {
		err = json.Unmarshal(cacheData.Val, &pokemon)
	} else {
		pokemon, _, err = pokeapi.GetPokemonInLocation(url)
		// go cacheStruct.AddEntry(config.argument, rawData)
	}
	if err != nil {
		return pokeapi.Pokemon{}, err
	}
	return pokemon, nil
}

/*
callback function for explore <location> command
*/
func showPokemonInLocation(config *config, cacheStruct *cache.Cache) error {
	pokemon, err := getPokemonInLocation(config, cacheStruct)
	if err != nil {
		return err
	}
	fmt.Println("Exploring " + pokemon.Location.Name + "...")
	for _, pokemonEncounter := range pokemon.PokemonEncounters {
		fmt.Println("-" + pokemonEncounter.Pokemon.Name)
	}
	fmt.Print("\n")
	return nil
}
