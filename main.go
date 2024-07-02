package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dusmcd/pokedexcli/cache"
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
		log.Fatal(scanner.Err())
		return "", scanner.Err()
	}

	return input, nil
}

func session() {
	config := config{
		next:     "https://pokeapi.co/api/v2/location/",
		previous: "",
		page:     0,
	}
	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)
	cache := cache.NewCache(10)
	for {
		userInput, err := getUserInput(scanner)
		if err != nil {
			continue
		}
		if userInput == "exit" {
			break
		}

		arguments := strings.Split(userInput, " ")
		if len(arguments) > 1 {
			config.argument = arguments[1]
			userInput = arguments[0]
		}

		command := getCommand(userInput)
		err = command.callback(&config, cache)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

}
