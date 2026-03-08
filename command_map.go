package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
