package service

import (
	"encoding/json"
	"math/rand"
	"os"
	"pokecat_pokebat/internal/model"
	"strconv"
)

func GeneratePokemonList(pokedex *[]model.Pokemon, numPokemon int) *model.PlayerPokemonList {
	playerPokemons := model.PlayerPokemonList{}

	for i := 0; i < numPokemon; i++ {
		evs := 0.5 + rand.Float64()*0.5
		if len(*pokedex) > 0 {
			newPokemon := CapturePokemon(&(*pokedex)[rand.Intn(len(*pokedex))], 1, strconv.FormatFloat(evs, 'f', -1, 64))
			playerPokemons = append(playerPokemons, newPokemon)
		}
	}

	return &playerPokemons
}

func SavePlayerPokemonLists(filename string, playerPokemonLists map[string]*model.PlayerPokemonList) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encodedPlayerPokemonLists := make(map[string][]model.CapturedPokemon)

	for playerName, pokemonList := range playerPokemonLists {
		encodedPlayerPokemonLists[playerName] = *pokemonList
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(encodedPlayerPokemonLists)
	if err != nil {
		return err
	}

	return nil
}
