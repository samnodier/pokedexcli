package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *Config) error
}

type Config struct {
	Next     string
	Previous string
}

type PokeAPIResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	commands := getCommands()
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMapNext(c *Config) error {
	if c.Next == "" {
		fmt.Println("You are on the last page")
		return nil
	}
	res, err := http.Get(c.Next)
	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	locationData := PokeAPIResponse{}

	if err := decoder.Decode(&locationData); err != nil {
		return fmt.Errorf("Error decoding response: %w", err)
	}
	for _, area := range locationData.Results {
		fmt.Println(area.Name)
	}
	c.Next = locationData.Next
	c.Previous = locationData.Previous
	return nil
}

func commandMapPrevious(c *Config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	res, err := http.Get(c.Previous)
	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	locationData := PokeAPIResponse{}

	if err := decoder.Decode(&locationData); err != nil {
		return fmt.Errorf("Error decoding response: %w", err)
	}
	for _, area := range locationData.Results {
		fmt.Println(area.Name)
	}
	c.Next = locationData.Next
	c.Previous = locationData.Previous
	return nil
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of next 20 location areas in the Pokemon world",
			callback:    commandMapNext,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of previous 20 location areas in the Pokemon world",
			callback:    commandMapPrevious,
		},
	}
	return commands
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)
	return words
}
