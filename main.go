package main

import (
	"time"

	"github.com/LoreviQ/PokeGo/internal/pokeAPI"
)

func main() {
	pokeClient := pokeAPI.NewClient(5 * time.Second)
	config := &config{
		Client: pokeClient,
	}
	startRepl(config)
}
