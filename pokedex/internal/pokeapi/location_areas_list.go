package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (c *Client) GetLocationData(url string) (LocationResponse, error) {
	// Typical pattern:
	//     Create a client with timeouts.
	//     Build a request (optionally attach headers, context).
	//     Do the request with client.Do(req).
	//     Always close resp.Body when done.
	// Gives us more flexibility than http.Get. Also allows us to easily add other http requests in the future

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationResponse{}, fmt.Errorf("creating request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return LocationResponse{}, fmt.Errorf("fetching location area data: %w", err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			fmt.Println("Error closing response body: ", err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		return LocationResponse{}, fmt.Errorf("response failed with status code %d: %s", res.StatusCode, string(body))
	}
	if err != nil {
		return LocationResponse{}, fmt.Errorf("reading response body: %w", err)
	}

	var responseData LocationResponse
	if err = json.Unmarshal(body, &responseData); err != nil {
		return LocationResponse{}, fmt.Errorf("unmarshalling response body: %w", err)
	}

	return responseData, nil
}
