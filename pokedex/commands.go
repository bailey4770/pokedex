package main

import (
	"fmt"
	"os"
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
		locationData, err := cfg.pokeapiClient.GetLocationData(*url)
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
		locationData, err := cfg.pokeapiClient.GetLocationData(*url)
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
