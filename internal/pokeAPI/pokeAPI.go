package pokeAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/LoreviQ/PokeGo/internal/pokeCache"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type Client struct {
	httpClient http.Client
	cache      pokeCache.Cache
}

type APImapData struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type APIlocationData struct {
	Areas []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"areas"`
}

type APIareaData struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func NewClient(timeout, interval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokeCache.NewCache(interval),
	}
}

func getEndpoint(Next, Previous string, args []string) (string, error) {
	forward := true
	if len(args) != 0 {
		if args[0] == "-back" || args[0] == "-b" {
			forward = false
		} else {
			return "", fmt.Errorf("'%s' is an invalid arg", args[0])
		}
	}
	var endpoint string
	if forward {
		if Next == "" {
			endpoint = baseURL + "/location/?offset=0&limit=20"
		} else {
			endpoint = Next
		}
	} else {
		if Previous == "" {
			return "", errors.New("cannot go backwards from the start")
		} else {
			endpoint = Previous
		}
	}
	return endpoint, nil
}

func (c *Client) getFromAPI(endpoint string) ([]byte, error) {
	body, err := c.cache.Get(endpoint)
	if err == nil {
		return body, nil
	}
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	body, err = io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return nil, err
	}
	err = c.cache.Add(endpoint, body)
	if err != nil {
		return body, err
	}
	return body, nil
}

func convertJsonToStruct[T any](body []byte, _ T) (T, error) {
	var JSON T
	err := json.Unmarshal([]byte(body), &JSON)
	if err != nil {
		var zeroVal T
		return zeroVal, err
	}
	return JSON, nil
}

func (c *Client) GetLocations(Next, Previous string, args []string) (APImapData, error) {
	var zeroVal APImapData
	endpoint, err := getEndpoint(Next, Previous, args)
	if err != nil {
		return zeroVal, err
	}
	body, err := c.getFromAPI(endpoint)
	if err != nil {
		return zeroVal, err
	}
	mapData, err := convertJsonToStruct(body, zeroVal)
	if err != nil {
		return zeroVal, err
	}
	return mapData, nil
}

func (c *Client) ExploreLocation(location string) (APIlocationData, error) {
	var zeroVal APIlocationData
	endpoint := baseURL + "/location/" + location
	body, err := c.getFromAPI(endpoint)
	if err != nil {
		return zeroVal, err
	}
	locationData, err := convertJsonToStruct(body, zeroVal)
	if err != nil {
		return zeroVal, err
	}
	return locationData, nil
}
