# Pokedex CLI

A command-line Pokedex REPL built in Go, powered by the [PokéAPI](https://pokeapi.co/).

## About

This project is a REPL (Read-Eval-Print Loop) that lets you look up Pokémon info — names, types, stats, and more — directly from your terminal.

## Learning Goals

- Parse JSON in Go
- Make HTTP GET requests
- Build a CLI tool that interacts with a backend API
- Practice local Go development and tooling
- Implement caching to improve performance

## Usage

```bash
go run main.go
```

Once the REPL starts, type a command and hit enter. Type `help` to see available commands.

## Data Source

All Pokémon data is fetched from the [PokéAPI](https://pokeapi.co/) — a free, open REST API for Pokémon data.
