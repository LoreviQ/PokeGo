package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func callHelp() error {
	fmt.Print("Welcome to PokéGo!\n\nThe available commands are:\n")
	for _, c := range getCliCommands() {
		c.log()
	}
	return nil
}

func callExit() error {
	fmt.Print("Thank you for using PokéGO!\n")
	os.Exit(0)
	return nil
}

func callMap() error {
	res, err := http.Get("https://pokeapi.co/api/v2/location/")
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", body)
	return nil
}
