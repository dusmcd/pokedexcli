package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

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
			break
		}
		if userInput == "exit" {
			break
		}
		command := getCommand(userInput)
		command.callback(&config, cache)
	}

}
