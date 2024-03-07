package pokeAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type Client struct {
	httpClient http.Client
}

type APImapData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) getEndpoint(Next, Previous string, args []string) (string, error) {
	forward := true
	if len(args) != 0 && (args[0] == "-back" || args[0] == "-b") {
		forward = false
	}
	var endpoint string
	if forward {
		if Next == "" {
			endpoint = "https://pokeapi.co/api/v2/location/"
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

func (c *Client) getAPI(endpoint string) ([]byte, error) {
	res, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) convertToStruct(body []byte) (APImapData, error) {
	var JSON APImapData
	err := json.Unmarshal([]byte(body), &JSON)
	if err != nil {
		var zeroVal APImapData
		return zeroVal, err
	}
	return JSON, nil
}

func (c *Client) GetLocations(Next, Previous string, args []string) (APImapData, error) {
	var zeroVal APImapData
	endpoint, err := c.getEndpoint(Next, Previous, args)
	if err != nil {
		return zeroVal, err
	}
	body, err := c.getAPI(endpoint)
	if err != nil {
		return zeroVal, err
	}
	mapData, err := c.convertToStruct(body)
	if err != nil {
		return zeroVal, err
	}
	return mapData, nil
}
