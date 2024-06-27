package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/dusmcd/pokedexcli/commands"
)

func main() {
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
	for {
		userInput, err := getUserInput()
		if err != nil {
			break
		}
		if userInput == "exit" {
			break
		}
		command := commands.Command(userInput)
		action := command.callback
		action()
	}

}
