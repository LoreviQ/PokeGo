package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/LoreviQ/PokeGo/internal/pokeAPI"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *config, args ...string) error
}

func (c *cliCommand) log() error {
	blankSpace := strings.Repeat(" ", 9-len(c.name))
	_, err := fmt.Printf("%s%s- %s\n", c.name, blankSpace, c.description)
	return err
}

type config struct {
	Client   pokeAPI.Client
	Next     string
	Previous string
}

func startRepl(config *config) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("PokéGO > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		input = strings.Replace(input, "\n", "", -1)
		words := strings.Split(input, " ")
		command, err := getCommand(words[0])
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = command.callback(config, words[1:]...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getCommand(word string) (cliCommand, error) {
	command, ok := getCliCommands()[word]
	if ok {
		return command, nil
	} else {
		var zeroVal cliCommand
		return zeroVal, fmt.Errorf("%s is an invalid command", word)
	}
}

func getCliCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    callHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit PokéGO",
			callback:    callExit,
		},
		"map": {
			name:        "map",
			description: "Displays 20 locations. Repeated uses displays additional locations. '-b' or -'back' can be used to show the previous 20 locations",
			callback:    callMap,
		},
		"explore": {
			name:        "explore",
			description: "Lists all Pokemon in an area. 'Area' must be supplied as the first arg",
			callback:    callExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon. 'Pokemon' name must be supplied as the first arg",
			callback:    callCatch,
		},
	}
}
