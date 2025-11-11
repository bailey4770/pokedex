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

		if command, ok := commands[first]; ok {
			if err := command.callback(cfg); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
