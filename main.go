package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args ...string) error
}

func (c *cliCommand) log() error {
	blankSpace := strings.Repeat(" ", 9-len(c.name))
	_, ok := fmt.Printf("%s%s- %s\n", c.name, blankSpace, c.description)
	return ok
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
			description: "Exit the Pokedex",
			callback:    callExit,
		},
		"map": {
			name:        "map",
			description: "Displays 20 locations",
			callback:    callMap,
		},
	}
}

func main() {
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
		err = command.callback(words[1:]...)
		if err != nil {
			fmt.Println(err)
		}
	}
}
