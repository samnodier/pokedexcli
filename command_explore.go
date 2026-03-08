package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
