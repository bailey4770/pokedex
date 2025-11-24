package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/bailey4770/pokedex/internal/pokeapi"
)

// to add a new command, add an entry to the map return by getCommands() then add the commandFunc to the end of the file
type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func getCommands() map[string]cliCommand {
	commands := make(map[string]cliCommand)

	// register exit command
	commands["exit"] = cliCommand{
		"exit",
		"Exit the pokedex",
		commandExit,
	}
	commands["help"] = cliCommand{
		"help",
		"Displays a help message",
		commandHelp(commands),
	}
	commands["map"] = cliCommand{
		"map",
		"Displays the next location areas",
		commandMap,
	}
	commands["mapb"] = cliCommand{
		"mapb",
		"Displays the previous location areas",
		commandMapb,
	}
	commands["explore"] = cliCommand{
		"explore",
		"Lists all pokemon found in this location area",
		commandExplore,
	}
	commands["catch"] = cliCommand{
		"catch",
		"Attempts to catch a named pokemon to add to pokedex",
		commandCatch,
	}
	commands["inspect"] = cliCommand{
		"inspect",
		"Takes the name of a pokemon and prints some statistics",
		commandInspect,
	}
	commands["pokedex"] = cliCommand{
		"pokedex",
		"Lists caught pokemon in the pokedex",
		commandPokedex,
	}
	commands["version"] = cliCommand{
		"version",
		"Displays current pokedex version",
		commandVersion,
	}
	commands["reset"] = cliCommand{
		"reset",
		"Resets pokedex to empty",
		commandReset,
	}

	return commands
}

func commandExit(cfg *config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("exit command was not expecting any arguments")
	}

	if err := savePokedex(cfg); err != nil {
		// print rather than return since we don't want to be stuck unable to exit if there is a saving problem
		fmt.Printf("Error saving the pokedex: %v/n", err)
	} else {
		fmt.Println("Successfully saved pokedex.")
	}
	fmt.Println("Closing the Pokedex... Goodbye!")

	os.Exit(0)
	return nil
}

func commandHelp(commands map[string]cliCommand) func(cfg *config, args []string) error {
	return func(cfg *config, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("help command was not expecting any arguments")
		}

		fmt.Println("Welcome to the Pokedex!\nUsage:")
		fmt.Println()

		for _, cmd := range commands {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		return nil
	}
}

func commandMap(cfg *config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("map command was not expecting any arguments")
	}

	if cfg.nextURL != "" {
		locationData, err := pokeapi.GetData[pokeapi.LocationResponse](&cfg.pokeapiClient, cfg.nextURL, cfg.cache)
		if err != nil {
			return err
		}

		// fmt.Printf("showing results for %s\n", cfg.nextURL)
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

		return nil

	} else {
		return fmt.Errorf("you're on the last page")
	}
}

func commandMapb(cfg *config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("mapb command was not expecting any arguments")
	}

	if cfg.prevURL != "" {
		locationData, err := pokeapi.GetData[pokeapi.LocationResponse](&cfg.pokeapiClient, cfg.prevURL, cfg.cache)
		if err != nil {
			return err
		}

		// fmt.Printf("showing results for %s\n", cfg.prevURL)
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
			cfg.prevURL = *locationData.Previous
		}

		return nil

	} else {
		return fmt.Errorf("you're on the first page")
	}
}

func commandExplore(cfg *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("explore command takes one location area argument")
	}

	url := cfg.baseURL + "location-area/" + args[0]
	pokemonData, err := pokeapi.GetData[pokeapi.PokemonListResponse](&cfg.pokeapiClient, url, cfg.cache)
	if err != nil {
		return err
	}

	for _, pokemon := range pokemonData.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("catch command takes one pokemon name as an argument")
	}
	pokemon := args[0]

	if _, exists := cfg.pokedex[pokemon]; exists {
		return fmt.Errorf("%s is already in your pokedex", pokemon)
	}

	url := cfg.baseURL + "pokemon/" + pokemon
	pokemonStats, err := pokeapi.GetData[pokeapi.PokemonStatsResponse](&cfg.pokeapiClient, url, cfg.cache)
	if err != nil {
		return err
	}

	const k = 0.008
	chance := math.Exp(-k * float64(pokemonStats.BaseExperience))
	fmt.Printf("%s has base experience: %d. There is a %.1f%% chance of success\n", pokemon, pokemonStats.BaseExperience, chance*100)
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	diceRoll := rand.Float64()
	if diceRoll < chance {
		fmt.Printf("You successfully caught %s!\n", pokemon)
		cfg.pokedex[pokemon] = pokemonStats
	} else {
		fmt.Printf("You were unsuccessful in catching %s\n", pokemon)
	}

	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("inspect command takes one pokemon name as an argument")
	}
	pokemon := args[0]

	pokemonStats, exists := cfg.pokedex[pokemon]
	if !exists {
		return fmt.Errorf("%s is not in your pokedex", pokemon)
	}

	fmt.Printf("Name: %s\n", pokemon)
	fmt.Printf("Height: %d\n", pokemonStats.Height)
	fmt.Printf("Weight: %d\n", pokemonStats.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemonStats.Stats {
		fmt.Printf("- %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemonStats.Types {
		fmt.Printf("- %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("pokedex command was not expecting any arguments")
	}
	pokedex := cfg.pokedex

	if len(pokedex) == 0 {
		return fmt.Errorf("pokedex is empty")
	}

	fmt.Println("Your pokedex:")
	for name := range pokedex {
		fmt.Printf("- %s\n", name)
	}

	return nil
}

func commandVersion(_ *config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("version command was not expecting any arguments")
	}

	fmt.Printf("Version %s\n", VERSION)
	return nil
}

func commandReset(cfg *config, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("reset command was not expecting any arguments")
	}

	fmt.Println("Clearing pokedex...")
	_ = resetPokedex(cfg)

	return nil
}
