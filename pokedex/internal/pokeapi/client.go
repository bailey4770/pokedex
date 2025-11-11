// Package pokeapi provides a client for interacting with the public PokeAPI.
package pokeapi

import (
	"net/http"
	"time"
)

// Client wraps the standard http.Client for the PokeAPI. Allows adding methods and better error handling
type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client{
		http.Client{
			Timeout: timeout,
		},
	}
}
