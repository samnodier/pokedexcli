package main

import (
	"fmt"
	"os"
)

func commandExit(c *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
