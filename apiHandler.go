package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type APImapData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func getEndpoint(config config, args []string) (string, error) {
	forward := true
	if len(args) != 0 && (args[0] == "-back" || args[0] == "-b") {
		forward = false
	}
	var endpoint string
	if forward {
		if config.Next == "" {
			endpoint = "https://pokeapi.co/api/v2/location/"
		} else {
			endpoint = config.Next
		}
	} else {
		if config.Previous == "" {
			return "", errors.New("cannot go backwards from the start")
		} else {
			endpoint = config.Previous
		}
	}
	return endpoint, nil
}

func getAPI(endpoint string) ([]byte, error) {
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

func convertToStruct(body []byte) (APImapData, error) {
	var JSON APImapData
	err := json.Unmarshal([]byte(body), &JSON)
	if err != nil {
		var zeroVal APImapData
		return zeroVal, err
	}
	return JSON, nil
}
