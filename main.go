package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	commands := getCommands()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanned := scanner.Scan()
		if !scanned {
			fmt.Println("Failed to scan text")
			break
		}
		text := scanner.Text()
		cleanText := cleanInput(text)
		if len(cleanText) == 0 {
			continue
		}
		firstWord := cleanText[0]
		if command, ok := commands[firstWord]; !ok {
			fmt.Println("Unknown command")
			continue
		} else {
			if err := command.callback(); err != nil {
				fmt.Printf("Error encountered: %v", err)
			}
		}
	}
}
