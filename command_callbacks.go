package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

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
	var rawData []byte
	ch := make(chan cache.CacheData)
	go cacheStruct.GetEntry(config.argument, ch)
	cacheData := <-ch
	pokemon := pokeapi.Pokemon{}
	if cacheData.Found {
		err = json.Unmarshal(cacheData.Val, &pokemon)
	} else {
		pokemon, rawData, err = pokeapi.GetPokemonInLocation(url)
		go cacheStruct.AddEntry(config.argument, rawData)
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
	fmt.Println("Exploring " + config.argument + "...")
	pokemon, err := getPokemonInLocation(config, cacheStruct)
	if err != nil {
		return err
	}
	for _, pokemonEncounter := range pokemon.PokemonEncounters {
		fmt.Println("-" + pokemonEncounter.Pokemon.Name)
	}
	fmt.Print("\n")
	return nil
}

/*
callback function for catch <pokemon> command
*/
func catchPokemon(config *config, cacheStruct *cache.Cache) error {
	_, found := config.pokedex.data[config.argument]
	if found {
		return errors.New("pokemon already caught")
	}
	fmt.Println("Throwing a Pokeball at " + config.argument + "...")
	randomNumber := rand.Float64()

	ch := make(chan cache.CacheData)
	pokemonStats := pokeapi.PokemonStats{}
	var err error
	var rawData []byte

	go cacheStruct.GetEntry(config.argument, ch)
	cacheData := <-ch

	if cacheData.Found {
		err = json.Unmarshal(cacheData.Val, &pokemonStats)
	} else {
		url := "https://pokeapi.co/api/v2/pokemon/" + config.argument
		pokemonStats, rawData, err = pokeapi.GetPokemonStats(url)
		go cacheStruct.AddEntry(config.argument, rawData)
	}

	if err != nil {
		return err
	}
	threshold := calculateThreshold(pokemonStats.BaseExperience)
	if randomNumber < threshold {
		fmt.Println(config.argument + " caught and added to pokedex!")
		config.pokedex.add(config.argument, pokemonStats)
	} else {
		fmt.Println(config.argument + " escaped!")
	}

	return nil
}

func calculateThreshold(baseExperience int) float64 {
	if baseExperience < 100 {
		return 0.50
	} else if baseExperience >= 100 && baseExperience < 150 {
		return 0.40
	}
	return 0.25
}

/*
callback function for the inspect <pokemon> command
*/
func showPokemonInfo(config *config, cacheStruct *cache.Cache) error {
	pokemon, found := config.pokedex.data[config.argument]
	if !found {
		return errors.New("you have not caught that pokemon")
	}
	fmt.Println("Name: " + pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("\t-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pokemonType := range pokemon.Types {
		fmt.Println("\t- " + pokemonType.Type.Name)
	}

	return nil
}

/*
callback function for pokedex command
*/
func showPokedex(config *config, cacheStruct *cache.Cache) error {
	fmt.Println("Your Pokedex:")
	for pokemon := range config.pokedex.data {
		fmt.Println("\t- " + config.pokedex.data[pokemon].Name)
	}
	return nil
}
