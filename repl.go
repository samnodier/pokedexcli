package main

import (
	"encoding/json"
	"fmt"
	"github.com/samnodier/pokedexcli/internal/pokecache"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *Config, args ...string) error
}

type Config struct {
	Next     string
	Previous string
	cache    *pokecache.Cache
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

type PokeAPIAreaResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func commandExit(c *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	commands := getCommands()
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMapNext(c *Config, args ...string) error {
	locationData := PokeAPIResponse{}
	if c.Next == "" {
		fmt.Println("You are on the last page")
		return nil
	}
	if val, ok := c.cache.Get(c.Next); ok {
		err := json.Unmarshal(val, &locationData)
		if err != nil {
			return fmt.Errorf("Error reading the cache: %w\n", err)
		}
	} else {
		res, err := http.Get(c.Next)
		if err != nil {
			return fmt.Errorf("Error creating request: %w\n", err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Status error: %d\n", res.StatusCode)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading the body: %w\n", err)
		}
		c.cache.Add(c.Next, body)
		if err := json.Unmarshal(body, &locationData); err != nil {
			return fmt.Errorf("Error decoding response: %w\n", err)
		}
	}
	for _, area := range locationData.Results {
		fmt.Println(area.Name)
	}
	c.Next = locationData.Next
	c.Previous = locationData.Previous
	return nil
}

func commandMapPrevious(c *Config, args ...string) error {
	locationData := PokeAPIResponse{}
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	if val, ok := c.cache.Get(c.Previous); ok {
		err := json.Unmarshal(val, &locationData)
		if err != nil {
			return fmt.Errorf("Error reading the cache: %w\n", err)
		}
	} else {
		res, err := http.Get(c.Previous)
		if err != nil {
			return fmt.Errorf("Error creating request: %w\n", err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Status error: %d\n", res.StatusCode)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading the body: %w\n", err)
		}
		c.cache.Add(c.Previous, body)
		if err := json.Unmarshal(body, &locationData); err != nil {
			return fmt.Errorf("Error decoding response: %w\n", err)
		}
	}
	for _, area := range locationData.Results {
		fmt.Println(area.Name)
	}
	c.Next = locationData.Next
	c.Previous = locationData.Previous
	return nil
}

func commandExplore(c *Config, args ...string) error {
	locationData := PokeAPIAreaResponse{}
	if len(args) == 0 {
		fmt.Println("you must provide a location name")
		return nil
	}
	locationName := args[0]
	locationURL := "https://pokeapi.co/api/v2/location-area/" + locationName + "/"
	if val, ok := c.cache.Get(locationURL); ok {
		err := json.Unmarshal(val, &locationData)
		if err != nil {
			return fmt.Errorf("Error reading the cache: %w\n", err)
		}
	} else {
		res, err := http.Get(locationURL)
		if err != nil {
			return fmt.Errorf("Error creating request: %w\n", err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Status error: %d\n", res.StatusCode)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading the body: %w\n", err)
		}
		c.cache.Add(locationURL, body)
		if err := json.Unmarshal(body, &locationData); err != nil {
			return fmt.Errorf("Error decoding response: %w\n", err)
		}
	}
	fmt.Printf("Exploring %s...\n", locationName)
	if len(locationData.PokemonEncounters) == 0 {
		fmt.Println("Found no Pokemon")
		return nil
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range locationData.PokemonEncounters {
		fmt.Println(" - " + encounter.Pokemon.Name)
	}
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
		"explore": {
			name:        "explore",
			description: "Displays a list of all Pokémon locate in an area",
			callback:    commandExplore,
		},
	}
	return commands
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)
	return words
}
