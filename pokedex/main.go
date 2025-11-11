package main

import (
	"time"

	"github.com/bailey4770/pokedex/internal/pokeapi"
)

type config struct {
	pokeapiClient pokeapi.Client
	currentURL    *string
	nextURL       *string
	prevURL       *string
}

func main() {
	// create new HTTP client using custom wrapper. Timeout fails any requests that take longer than 5 seconds
	pokeClient := pokeapi.NewClient(5 * time.Second)
	currentURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	nextURL := "https://pokeapi.co/api/v2/location-area?offset=20&limit=20"

	cfg := &config{
		pokeapiClient: pokeClient,
		currentURL:    &currentURL,
		nextURL:       &nextURL,
	}

	replLoop(cfg)
}
