package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

func commandCatch(c *Config, args ...string) error {
	pokemonData := Pokemon{}
	if len(args) == 0 {
		fmt.Println("you must provide a pokemon name")
		return nil
	}
	pokemonName := args[0]
	pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)
	if val, ok := c.cache.Get(pokemonURL); ok {
		err := json.Unmarshal(val, &pokemonData)
		if err != nil {
			return fmt.Errorf("Error reading the pokemon cache: %w\n", err)
		}
	} else {
		res, err := http.Get(pokemonURL)
		if err != nil {
			return fmt.Errorf("Error creating request: %w\n", err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Pokemon doesn't exist: %d\n", res.StatusCode)
		}
		if err != nil {
			return fmt.Errorf("Error reading the body: %w\n", err)
		}
		c.cache.Add(pokemonURL, body)
		if err := json.Unmarshal(body, &pokemonData); err != nil {
			return fmt.Errorf("Error decoding response: %w\n", err)
		}
	}
	if _, ok := c.Pokedex[pokemonName]; ok {
		fmt.Printf("You already caught %s!\n", pokemonName)
		return nil
	}
	// Random number between 0 and the pokemon base experience
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	maxRoll := pokemonData.BaseExperience + (pokemonData.BaseExperience / 2)
	randNum := rand.Intn(maxRoll)
	if randNum > pokemonData.BaseExperience {
		c.Pokedex[pokemonName] = pokemonData
		fmt.Printf("%s was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	return nil
}
