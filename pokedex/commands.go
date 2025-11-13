package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bailey4770/pokedex/internal/pokeapi"
)

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands map[string]cliCommand) func(cfg *config) error {
	return func(cfg *config) error {
		fmt.Println("Welcome to the Pokedex!\nUsage:")
		fmt.Println()

		for _, cmd := range commands {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		return nil
	}
}

func commandMap(cfg *config) error {
	if cfg.nextURL != "" {
		locationData, err := pokeapi.GetData[pokeapi.LocationResponse](&cfg.pokeapiClient, cfg.nextURL, cfg.cache)
		if err != nil {
			return err
		}

		for _, item := range locationData.Results {
			fmt.Println(item.Name)
		}

		if locationData.Previous == nil {
			cfg.prevURL = ""
		} else {
			cfg.prevURL = *locationData.Previous
		}
		if locationData.Next == nil {
			cfg.nextURL = ""
		} else {
			cfg.nextURL = *locationData.Next
		}

		log.Printf("Next URL: %s, prevURL: %s", cfg.nextURL, cfg.prevURL)

		return nil

	} else {
		return fmt.Errorf("you're on the last page")
	}
}

func commandMapb(cfg *config) error {
	if cfg.prevURL != "" {
		locationData, err := pokeapi.GetData[pokeapi.LocationResponse](&cfg.pokeapiClient, cfg.prevURL, cfg.cache)
		if err != nil {
			return err
		}

		for _, item := range locationData.Results {
			fmt.Println(item.Name)
		}

		if locationData.Next == nil {
			cfg.nextURL = ""
		} else {
			cfg.nextURL = *locationData.Next
		}
		if locationData.Previous == nil {
			cfg.prevURL = ""
		} else {
			cfg.prevURL = *locationData.Next
		}

		return nil

	} else {
		return fmt.Errorf("you're on the first page")
	}
}

func commandExplore(cfg *config) error {
	args := cfg.args
	if len(args) > 1 {
		return fmt.Errorf("explore command takes one location area argument")
	}

	url := cfg.baseURL + args[0]
	pokemonData, err := pokeapi.GetData[pokeapi.PokemonListResponse](&cfg.pokeapiClient, url, cfg.cache)
	if err != nil {
		return err
	}

	for _, pokemon := range pokemonData.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}
