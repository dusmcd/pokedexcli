package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Print("Welcome to the Pokedex! Type help for valid commands.\n\n")
	session()
	fmt.Println("Thank you for using Pokedex CLI.")
}

func getUserInput() (string, error) {
	fmt.Print("pokedex > ")
	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	input := scanner.Text()

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
		return "", scanner.Err()
	}

	return input, nil
}

func session() {
	config := config{
		next:     "https://pokeapi.co/api/v2/location/",
		previous: "",
	}
	for {
		userInput, err := getUserInput()
		if err != nil {
			break
		}
		if userInput == "exit" {
			break
		}
		command := getCommand(userInput)
		command.callback(&config)
	}

}
