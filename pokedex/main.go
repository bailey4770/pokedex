package main

import (
	"fmt"
	"time"

	"github.com/bailey4770/pokedex/internal/pokeapi"
	"github.com/bailey4770/pokedex/internal/pokecache"
)

type pokedexType map[string]pokeapi.PokemonStatsResponse

type config struct {
	version       string
	pokeapiClient pokeapi.Client
	baseURL       string
	nextURL       string
	prevURL       string
	cache         *pokecache.Cache
	args          []string
	pokedex       pokedexType
}

func main() {
	// initial urls for lcoation data
	baseURL := "https://pokeapi.co/api/v2/"
	nextURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	prevURL := ""

	// create new HTTP client using custom wrapper. Timeout fails any requests that take longer than 5 seconds
	// create new cache using custom wrapper. Cache clears every 5 seconds to save space
	cfg := &config{
		version:       "1.0.0",
		pokeapiClient: pokeapi.NewClient(5 * time.Second),
		baseURL:       baseURL,
		nextURL:       nextURL,
		prevURL:       prevURL,
		cache:         pokecache.NewCache(5 * time.Second),
		pokedex:       make(pokedexType),
	}

	err := loadPokedex(cfg)
	if err != nil {
		fmt.Println(err)
	}

	// initialise repl loop
	startRepl(cfg)
}
