package service

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"pokecat_pokebat/internal/model"
)

func LoadPokedex(filename string) (*[]model.Pokemon, error) {
	var pokedex []model.Pokemon

	file, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &pokedex)

	return &pokedex, err
}

func LoadPokedexData(filename string) ([]model.Pokemon, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var pokedexData []model.Pokemon
	err = json.Unmarshal(data, &pokedexData)
	if err != nil {
		return nil, err
	}

	return pokedexData, nil
}
