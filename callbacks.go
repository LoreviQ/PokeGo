package main

import (
	"fmt"
	"os"
)

func callHelp(config config, args ...string) (config, error) {
	fmt.Print("Welcome to PokéGo!\n\nThe available commands are:\n")
	for _, c := range getCliCommands() {
		c.log()
	}
	return config, nil
}

func callExit(config config, args ...string) (config, error) {
	fmt.Print("Thank you for using PokéGO!\n")
	os.Exit(0)
	return config, nil
}

func callMap(config config, args ...string) (config, error) {
	endpoint, err := getEndpoint(config, args)
	if err != nil {
		return config, err
	}
	body, err := getAPI(endpoint)
	if err != nil {
		return config, err
	}
	mapData, err := convertToStruct(body)
	if err != nil {
		return config, err
	}
	config.Next = mapData.Next
	config.Previous = mapData.Previous
	for _, location := range mapData.Results {
		fmt.Printf("%s\n", location.Name)
	}
	return config, nil
}
