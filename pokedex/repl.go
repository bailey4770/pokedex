package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

func startRepl(cfg *config) {
	commands := getCommands()

	completer := readline.NewPrefixCompleter()
	for cmd := range commands {
		completer.Children = append(completer.Children,
			readline.PcItem(cmd),
		)
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "Pokedex > ",
		HistoryLimit:    1000,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		AutoComplete:    completer,
	})
	if err != nil {
		log.Fatalf("readline returned err: %v", err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err != nil {
			log.Fatalf("error reading input: %v", err)
			return
		}

		parts := cleanInput(line)
		if len(parts) == 0 {
			continue
		}

		first := parts[0]
		args := parts[1:]

		if command, ok := commands[first]; ok {
			if err := command.callback(cfg, args); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
