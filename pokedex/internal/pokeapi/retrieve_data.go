package pokeapi

import (
	"encoding/json"
	"fmt"

	"github.com/bailey4770/pokedex/internal/pokecache"
)

type LocationResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokemonListResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type PokemonStatsResponse struct {
	BaseExperience int `json:"base_experience"`
}

type Response interface {
	LocationResponse | PokemonListResponse | PokemonStatsResponse
}

func GetData[T Response](c *Client, url string, cache *pokecache.Cache) (T, error) {
	// Typical pattern:
	//     Create a client with timeouts.
	//     Build a request (optionally attach headers, context).
	//     Do the request with client.Do(req).
	//     Always close resp.Body when done.
	// Gives us more flexibility than http.Get. Also allows us to easily add other http requests in the future

	var responseData T
	defer fmt.Printf("returning %v\n", responseData)

	// check first if data is stored in cache. we can return this and avoid http request
	body, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(body, &responseData); err != nil {
			return responseData, fmt.Errorf("unmarshalling response body: %w", err)
		}
		return responseData, nil
	}

	body, err := c.requestFromURL(url)
	if err != nil {
		return responseData, err
	}

	// now we have the data as slice of bytes, store in cache before unmarshallig and returning to caller
	cache.Add(url, body)

	// now unmarshal and return
	if err := json.Unmarshal(body, &responseData); err != nil {
		return responseData, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return responseData, nil
}
