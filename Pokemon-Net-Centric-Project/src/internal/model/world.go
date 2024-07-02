package model

type World struct {
	Width    int
	Height   int
	Pokemons map[Position]*Pokemon
}
