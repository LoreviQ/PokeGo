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
	callback    func() error
}

func (c *cliCommand) log() error {
	blankSpace := strings.Repeat(" ", 9-len(c.name))
	_, ok := fmt.Printf("%s%s- %s\n", c.name, blankSpace, c.description)
	return ok
}

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
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("PokéGO > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		input = strings.Replace(input, "\n", "", -1)
		words := strings.Split(input, " ")
		command, err := getCommand(words[0])
		if err == nil {
			command.callback()
		} else {
			fmt.Println(err)
		}
	}
}
