package model

import (
	"fmt"
	"strings"
)

type Player struct {
	Name          string
	Position      Position
	Surrendered   bool
	WorldID       int
	Pokemons      []*CapturedPokemon
	ActivePokemon *CapturedPokemon
}

func (p *Player) String() string {
	pokemonNames := make([]string, len(p.Pokemons))
	for i, pokemon := range p.Pokemons {
		pokemonNames[i] = pokemon.Name
	}

	pokemonsStr := fmt.Sprintf("[%s]",
		strings.Join(pokemonNames, ", "))

	return fmt.Sprintf("Name: %s Position: %+v Pokemons: %s",
		p.Name, p.Position, pokemonsStr)
}
