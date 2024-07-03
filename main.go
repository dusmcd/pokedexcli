package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dusmcd/pokedexcli/cache"
	"github.com/dusmcd/pokedexcli/pokeapi"
)

func main() {
	fmt.Print("Welcome to the Pokedex! Type help for valid commands.\n\n")
	session()
	fmt.Println("Thank you for using Pokedex CLI.")
}

func getUserInput(scanner *bufio.Scanner) (string, error) {
	fmt.Print("pokedex > ")
	scan := scanner.Scan()
	input := scanner.Text()

	err := scanner.Err()
	if !scan && err != nil {
		return "", err
	}

	return input, nil
}

func session() {
	config := config{
		next:     "https://pokeapi.co/api/v2/location/",
		previous: "",
		page:     0,
		pokedex: pokedex{
			data: make(map[string]pokeapi.PokemonStats),
		},
	}
	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)
	cache := cache.NewCache(120)
	for {
		userInput, err := getUserInput(scanner)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if userInput == "exit" {
			break
		}

		inputs := strings.Split(userInput, " ")
		if len(inputs) > 1 {
			config.argument = inputs[1]
			userInput = inputs[0]
		}

		command := getCommand(userInput)
		err = command.callback(&config, cache)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

}
