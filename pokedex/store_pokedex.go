package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

func savePokedex(cfg *config) error {
	file, err := os.Create("pokedex.txt")
	if err != nil {
		return fmt.Errorf("error creating pokedex file: %v", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(cfg.pokedex)
}

func loadPokedex(cfg *config) error {
	file, err := os.Open("pokedex.txt")
	if err != nil {
		return nil
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err = decoder.Decode(&cfg.pokedex); err != nil {
		return err
	}

	return nil
}

func resetPokedex(cfg *config) error {
	cfg.pokedex = nil
	fmt.Println("Pokedex successfully cleared.")
	return nil
}
