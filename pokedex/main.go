package main

import (
	"time"

	"github.com/bailey4770/pokedex/internal/pokeapi"
	"github.com/bailey4770/pokedex/internal/pokecache"
)

type Pokemon struct {
	name    string
	baseExp int
}

type config struct {
	pokeapiClient pokeapi.Client
	baseURL       string
	nextURL       string
	prevURL       string
	cache         *pokecache.Cache
	args          []string
	pokedex       map[string]Pokemon
}

func main() {
	// initial urls for lcoation data
	baseURL := "https://pokeapi.co/api/v2/"
	nextURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	prevURL := ""

	pokedex := make(map[string]Pokemon)

	// create new HTTP client using custom wrapper. Timeout fails any requests that take longer than 5 seconds
	// create new cache using custom wrapper. Cache clears every 5 seconds to save space
	cfg := &config{
		pokeapiClient: pokeapi.NewClient(5 * time.Second),
		baseURL:       baseURL,
		nextURL:       nextURL,
		prevURL:       prevURL,
		cache:         pokecache.NewCache(5 * time.Second),
		pokedex:       pokedex,
	}

	// initialise repl loop
	replLoop(cfg)
}
