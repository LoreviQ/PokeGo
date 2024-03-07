package main

import (
	"time"

	"https://github.com/LoreviQ/PokeGo/pokeAPI"
)

func main() {
	pokeClient := pokeAPI.NewClient(5 * time.Second)
	config := &config{
		Client: pokeClient,
	}
	startRepl(config)
}
