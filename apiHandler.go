package main

import (
	"encoding/json"
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
