package controller

import (
	"math/rand"
	"pokecat_pokebat/internal/model"
	"pokecat_pokebat/internal/service"
	"time"
)

type WorldController struct {
	WorldService *service.WorldService
	WorldID      int
	stop         chan struct{}
}

func NewWorldController(worldID int) *WorldController {
	return &WorldController{
		WorldService: service.NewWorldService(1000, 1000),
		WorldID:      worldID,
		stop:         make(chan struct{}),
	}
}

func (wc *WorldController) AddPlayer(player *model.Player) {
	wc.WorldService.AddPlayer(player)
}

func (wc *WorldController) SpawnPokemons(pokedexData []model.Pokemon, numPokemons int) {
	if numPokemons <= 0 {
		numPokemons = 50
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for i := 0; i < numPokemons; i++ {
				pokemon := wc.generateRandomPokemon(pokedexData)
				pos := model.Position{X: rand.Intn(1000), Y: rand.Intn(1000)}
				wc.WorldService.SpawnPokemon(pokemon, pos)

				go func(pokemon *model.Pokemon, pos model.Position) {
					time.Sleep(5 * time.Minute)
					wc.WorldService.DespawnPokemon(pos)
				}(pokemon, pos)
			}
		case <-wc.stop:
			return
		}
	}
}

func (wc *WorldController) generateRandomPokemon(pokedexData []model.Pokemon) *model.Pokemon {
	randomIndex := rand.Intn(len(pokedexData))
	selectedPokemon := pokedexData[randomIndex]

	level := rand.Intn(100) + 1
	ev := 0.5 + rand.Float64()*0.5

	return &model.Pokemon{
		No:          selectedPokemon.No,
		Image:       selectedPokemon.Image,
		Name:        selectedPokemon.Name,
		Type:        selectedPokemon.Type,
		Level:       level,
		TotalEvs:    int(ev),
		Exp:         selectedPokemon.Exp,
		HP:          selectedPokemon.HP,
		Attack:      selectedPokemon.Attack,
		Defense:     selectedPokemon.Defense,
		SpAttack:    selectedPokemon.SpAttack,
		SpDefense:   selectedPokemon.SpDefense,
		Speed:       selectedPokemon.Speed,
		Height:      selectedPokemon.Height,
		Weight:      selectedPokemon.Weight,
		CatchRate:   selectedPokemon.CatchRate,
		GenderRatio: selectedPokemon.GenderRatio,
		EggGroups:   selectedPokemon.EggGroups,
		HatchSteps:  selectedPokemon.HatchSteps,
		Abilities:   selectedPokemon.Abilities,
		EVs:         selectedPokemon.EVs,
		Strengths:   selectedPokemon.Strengths,
		Weaknesses:  selectedPokemon.Weaknesses,
		Evolutions:  selectedPokemon.Evolutions,
		Moves:       selectedPokemon.Moves,
	}
}

func (wc *WorldController) StopSpawning() {
	close(wc.stop)
}

func (wc *WorldController) DespawnPokemon(pos model.Position) {
	wc.WorldService.DespawnPokemon(pos)
}
