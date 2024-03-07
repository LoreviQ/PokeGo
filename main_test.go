package main

import (
	"fmt"
	"testing"
)

func TestGetCommand(t *testing.T) {
	commands := getCliCommands()
	i := 1
	for key, value := range commands {
		t.Run(fmt.Sprintf("Test Case %v", i), func(t *testing.T) {
			command, err := getCommand(key)
			if err != nil {
				t.Errorf("expected to find command")
			}
			if command.name != value.name {
				t.Errorf("wrong command")
			}
		})
		i++
	}
	t.Run(fmt.Sprintf("Test Case %v", i), func(t *testing.T) {
		_, err := getCommand("no-command")
		if err == nil {
			t.Errorf("Expected error to be returned")
		}
	})
}
