package main

import (
	"errors"
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
		locationData, err := config.Client.ExploreLocation(location.Name)
		if err != nil {
			return err
		}
		for _, area := range locationData.Areas {
			fmt.Printf("     - Area - %s\n", area.Name)
		}
	}
	return nil
}

func callExplore(config *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("no area supplied")
	}
	areaData, err := config.Client.ExploreArea(args[0])
	if err != nil {
		return err
	}
	for _, encounter := range areaData.PokemonEncounters {
		fmt.Printf(" - Found Pokemon - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func callCatch(config *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("no pokemon supplied")
	}
	err := config.Client.Catch(args[0])
	return err
}

func callInspect(config *config, args ...string) error {
	if len(args) == 0 {
		if len(config.Client.Pokedex) == 0 {
			return errors.New("you haven't caught any pokemon")
		}
		for pokemon := range config.Client.Pokedex {
			fmt.Printf(" - Has Caught - %s\n", pokemon)
		}
	} else {
		pokemon, ok := config.Client.Pokedex[args[0]]
		if !ok {
			return fmt.Errorf("you haven't caught %s", args[0])
		}
		fmt.Printf(" - Name: %v\n", pokemon.Name)
		fmt.Printf(" - Height: %v\n", pokemon.Height)
		fmt.Printf(" - Weight: %v\n", pokemon.Weight)
		fmt.Print(" - Stats:\n")
		for _, stat := range pokemon.Stats {
			fmt.Printf("     - %v: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Print(" - Types:\n")
		for _, t := range pokemon.Types {
			fmt.Printf("     - %v\n", t.Type.Name)
		}
	}
	return nil
}
