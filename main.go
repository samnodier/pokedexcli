package main

import (
	"bufio"
	"fmt"
	"github.com/samnodier/pokedexcli/internal/pokecache"
	"os"
	"time"
)

func main() {
	cache := pokecache.NewCache(5 * time.Minute)
	c := &Config{
		Next:  "https://pokeapi.co/api/v2/location-area",
		cache: cache,
	}
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
		cmdName := cleanText[0]
		args := cleanText[1:]
		if cmd, ok := commands[cmdName]; !ok {
			fmt.Println("Unknown cmdName")
			continue
		} else {
			if err := cmd.callback(c, args...); err != nil {
				fmt.Printf("Error encountered: %v", err)
			}
		}
	}
}
