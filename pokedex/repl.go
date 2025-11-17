package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/bailey4770/pokedex/internal/pokeapi"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

// to add a new command, add an entry to the map return by getCommands() then add the commandFunc to the end of the file
type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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

	return commands
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			return
		}

		input := scanner.Text()
		parts := cleanInput(input)
		if len(parts) == 0 {
			continue
		}
		first := parts[0]
		cfg.args = parts[1:]

		if command, ok := commands[first]; ok {
			if err := command.callback(cfg); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

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

func commandMapb(cfg *config) error {
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

func commandExplore(cfg *config) error {
	args := cfg.args
	if len(args) > 1 {
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

func commandCatch(cfg *config) error {
	args := cfg.args
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

func commandInspect(cfg *config) error {
	args := cfg.args
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

func commandPokedex(cfg *config) error {
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

func commandVersion(cfg *config) error {
	fmt.Printf("Version %s\n", cfg.version)
	return nil
}
