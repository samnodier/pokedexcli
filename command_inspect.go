package main

import (
	"errors"
	"fmt"
)

func commandInspect(c *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}
	pokemonName := args[0]
	if _, ok := c.Pokedex[pokemonName]; !ok {
		fmt.Println("You have not caught that pokemon")
		return nil
	}
	// Print caught Pokemon data
	pokemon := c.Pokedex[pokemonName]
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, stat := range pokemon.Types {
		fmt.Printf("  - %s\n", stat.Type.Name)
	}
	return nil
}
