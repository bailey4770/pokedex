// Package pokeapi provides a client for interacting with the public PokeAPI.
package pokeapi

import (
	"fmt"
	"io"
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

func (c *Client) requestFromURL(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("creating request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("fetching location area data: %w", err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			fmt.Println("Error closing response body: ", err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		return []byte{}, fmt.Errorf("response failed with status code %d: %s", res.StatusCode, string(body))
	}
	if err != nil {
		return []byte{}, fmt.Errorf("reading response body: %w", err)
	}

	return body, nil
}
