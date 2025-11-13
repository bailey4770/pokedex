package main

import (
	"fmt"
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
	if url := cfg.currentURL; url != nil {
		locationData, err := pokeapi.GetData[pokeapi.LocationResponse](&cfg.pokeapiClient, *url, cfg.cache)
		if err != nil {
			return err
		}

		for _, item := range locationData.Results {
			fmt.Println(item.Name)
		}

		cfg.prevURL = locationData.Previous
		cfg.currentURL = cfg.nextURL
		cfg.nextURL = locationData.Next

		return nil

	} else {
		return fmt.Errorf("you're on the last page")
	}
}

func commandMapb(cfg *config) error {
	if url := cfg.prevURL; url != nil {
		locationData, err := pokeapi.GetData[pokeapi.LocationResponse](&cfg.pokeapiClient, *url, cfg.cache)
		if err != nil {
			return err
		}

		for _, item := range locationData.Results {
			fmt.Println(item.Name)
		}

		cfg.nextURL = cfg.currentURL
		cfg.currentURL = cfg.prevURL
		cfg.prevURL = locationData.Previous

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

	url := *cfg.baseURL + args[0]
	pokemonData, err := pokeapi.GetData[pokeapi.PokemonListResponse](&cfg.pokeapiClient, url, cfg.cache)
	if err != nil {
		return err
	}

	for _, pokemon := range pokemonData.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}
