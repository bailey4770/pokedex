package main

import (
	"time"

	"github.com/bailey4770/pokedex/internal/pokeapi"
	"github.com/bailey4770/pokedex/internal/pokecache"
)

type config struct {
	pokeapiClient pokeapi.Client
	baseURL       *string
	currentURL    *string
	nextURL       *string
	prevURL       *string
	cache         *pokecache.Cache
	args          []string
}

func main() {
	// initial urls for lcoation data
	baseURL := "https://pokeapi.co/api/v2/location-area/"
	currentURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	nextURL := "https://pokeapi.co/api/v2/location-area?offset=20&limit=20"

	// create new HTTP client using custom wrapper. Timeout fails any requests that take longer than 5 seconds
	cfg := &config{
		pokeapiClient: pokeapi.NewClient(5 * time.Second),
		baseURL:       &baseURL,
		currentURL:    &currentURL,
		nextURL:       &nextURL,
		// also create cache to attach to config
		cache: pokecache.NewCache(5 * time.Second),
	}

	// initialise repl loop
	replLoop(cfg)
}
