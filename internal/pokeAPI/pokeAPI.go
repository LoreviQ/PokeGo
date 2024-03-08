package pokeAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
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
	Pokedex    map[string]Pokemon
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

type Pokemon struct {
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

type species struct {
	Name        string `json:"name"`
	CaptureRate int    `json:"capture_rate"`
}

func NewClient(timeout, interval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache:   pokeCache.NewCache(interval),
		Pokedex: map[string]Pokemon{},
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

func endpointToJSON[T any](c *Client, endpoint string, _ T) (T, error) {
	var zeroVal T
	var JSON T
	body, err := c.getFromAPI(endpoint)
	if err != nil {
		return zeroVal, err
	}
	err = json.Unmarshal([]byte(body), &JSON)
	if err != nil {
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
	mapData, err := endpointToJSON(c, endpoint, zeroVal)
	if err != nil {
		return zeroVal, err
	}
	return mapData, nil
}

func (c *Client) ExploreLocation(location string) (APIlocationData, error) {
	var zeroVal APIlocationData
	endpoint := baseURL + "/location/" + location
	locationData, err := endpointToJSON(c, endpoint, zeroVal)
	if err != nil {
		return zeroVal, err
	}
	return locationData, nil
}

func (c *Client) ExploreArea(area string) (APIareaData, error) {
	var zeroVal APIareaData
	endpoint := baseURL + "/location-area/" + area
	areaData, err := endpointToJSON(c, endpoint, zeroVal)
	if err != nil {
		return zeroVal, err
	}
	return areaData, nil
}

func (c *Client) Catch(pokemon string) error {
	var zeroVal species
	endpoint := baseURL + "/pokemon-species/" + pokemon
	pokemonSpecies, err := endpointToJSON(c, endpoint, zeroVal)
	if err != nil {
		return err
	}
	if rand.IntN(255) < pokemonSpecies.CaptureRate {
		fmt.Print("Caught!\n")
		_, ok := c.Pokedex[pokemon]
		if !ok {
			c.addToPokedex(pokemon)
		}
	} else {
		fmt.Print("Failed to catch\n")
	}
	return nil
}

func (c *Client) addToPokedex(pokemon string) error {
	var zeroVal Pokemon
	endpoint := baseURL + "/pokemon/" + pokemon
	pokemonData, err := endpointToJSON(c, endpoint, zeroVal)
	if err != nil {
		return err
	}
	c.Pokedex[pokemon] = pokemonData
	fmt.Print("Added to dex\n")
	return nil
}
