package main

import (
	"fmt"
	"os"
)

func callHelp(args ...string) error {
	fmt.Print("Welcome to PokéGo!\n\nThe available commands are:\n")
	for _, c := range getCliCommands() {
		c.log()
	}
	return nil
}

func callExit(args ...string) error {
	fmt.Print("Thank you for using PokéGO!\n")
	os.Exit(0)
	return nil
}

func callMap(args ...string) error {
	body, err := getAPI("https://pokeapi.co/api/v2/location/")
	if err != nil {
		return err
	}
	locations, err := convertToStruct(body)
	if err != nil {
		return err
	}
	fmt.Printf("\n\n json object:::: %v\n", locations)
	return nil
}
