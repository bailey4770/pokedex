package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func replLoop(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
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
