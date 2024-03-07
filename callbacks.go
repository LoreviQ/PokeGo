package main

import (
	"fmt"
	"os"
)

func callHelp(config *config, args ...string) error {
	fmt.Print("Welcome to PokéGo!\n\nThe available commands are:\n")
	for _, c := range getCliCommands() {
		c.log()
	}
	return nil
}

func callExit(config *config, args ...string) error {
	fmt.Print("Thank you for using PokéGO!\n")
	os.Exit(0)
	return nil
}

func callMap(config *config, args ...string) error {
	mapData, err := config.Client.GetLocations(config.Next, config.Previous, args)
	if err != nil {
		return err
	}
	config.Next = mapData.Next
	config.Previous = mapData.Previous
	for _, location := range mapData.Results {
		fmt.Printf(" - Location - %s\n", location.Name)
	}
	return nil
}

/*
func callExplore(config *config, args ...string) error {
	if len(args) != 0 {
		if args[0] == "-back" || args[0] == "-b" {
			forward = false
		} else {
			return "", fmt.Errorf("'%s' is an invalid arg", args[0])
		}
	}
	return nil
}
*/
