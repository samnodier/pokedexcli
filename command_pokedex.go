package main

import (
	"fmt"
)

func commandPokedex(c *Config, args ...string) error {
	if len(c.Pokedex) == 0 {
		fmt.Println("you have not caught any pokemon")
		return nil
	}
	for _, pokemon := range c.Pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}
