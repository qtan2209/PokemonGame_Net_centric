package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"pokecat_pokebat/internal/model"
	"sync"
)

type WorldService struct {
	World *model.World
	mutex sync.RWMutex
}

func NewWorldService(width, height int) *WorldService {
	return &WorldService{
		World: &model.World{
			Width:    width,
			Height:   height,
			Pokemons: make(map[model.Position]*model.Pokemon),
		},
	}
}

func (ws *WorldService) MinX() int {
	return 0
}

func (ws *WorldService) MaxX() int {
	return ws.World.Width - 1
}

func (ws *WorldService) MinY() int {
	return 0
}

func (ws *WorldService) MaxY() int {
	return ws.World.Height - 1
}

func (ws *WorldService) AddPlayer(player *model.Player) {
	// Add player to the world (e.g., to a players map or list)
}

func (ws *WorldService) PlayerHasPokemons() bool {
	for _, pokemon := range ws.World.Pokemons {
		if pokemon != nil {
			return true
		}
	}
	return false
}

func (ws *WorldService) HasPokemon(pos model.Position) bool {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	_, exists := ws.World.Pokemons[pos]
	return exists
}

func (ws *WorldService) CapturePokemon(pos model.Position) *model.Pokemon {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	pokemon := ws.World.Pokemons[pos]
	delete(ws.World.Pokemons, pos)
	return pokemon
}

func (ws *WorldService) SpawnPokemon(pokemon *model.Pokemon, pos model.Position) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	ws.World.Pokemons[pos] = pokemon
}

func (ws *WorldService) DespawnPokemon(pos model.Position) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	delete(ws.World.Pokemons, pos)
}

func (ws *WorldService) LoadPokedexData() ([]model.Pokemon, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %v", err)
	}

	pokedexFilePath := filepath.Join(workingDir, "../../data/pokedex.json")

	data, err := ioutil.ReadFile(pokedexFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pokedex data file: %v", err)
	}

	var pokemonList []model.Pokemon
	if err := json.Unmarshal(data, &pokemonList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pokedex data: %v", err)
	}

	return pokemonList, nil
}
