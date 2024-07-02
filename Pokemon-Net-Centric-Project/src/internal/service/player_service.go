package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"pokecat_pokebat/internal/model"
	"time"
)

type PlayerService struct {
	Players               []*model.Player
	PlayerPokemonDataFile string
	PlayerPokemonMap      map[string][]*model.CapturedPokemon
}

func NewPlayerService(playerPokemonDataFile string) *PlayerService {
	ps := &PlayerService{
		Players:               []*model.Player{},
		PlayerPokemonDataFile: playerPokemonDataFile,
		PlayerPokemonMap:      make(map[string][]*model.CapturedPokemon),
	}
	ps.loadPlayerList()
	return ps
}

func (ps *PlayerService) loadPlayerList() {
	data, err := ioutil.ReadFile(ps.PlayerPokemonDataFile)
	if err != nil {
		log.Fatalf("Failed to read player Pokemon data file: %v", err)
	}

	err = json.Unmarshal(data, &ps.PlayerPokemonMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal player Pokemon data: %v", err)
	}

	for name, capturedPokemons := range ps.PlayerPokemonMap {
		player := &model.Player{
			Name:     name,
			Position: model.Position{X: rand.Intn(1000), Y: rand.Intn(1000)},
			Pokemons: capturedPokemons,
		}
		ps.Players = append(ps.Players, player)
	}
}

func (ps *PlayerService) LoadPlayerList(filename string) map[string][]*model.CapturedPokemon {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read player data file: %v", err)
	}

	playerPokemonMap := make(map[string][]*model.CapturedPokemon)

	err = json.Unmarshal(data, &playerPokemonMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal player data: %v", err)
	}

	return playerPokemonMap
}

func (ps *PlayerService) CreatePlayer(name string, initialPosition model.Position) *model.Player {
	player := &model.Player{
		Name:     name,
		Pokemons: []*model.CapturedPokemon{},
		Position: initialPosition,
	}
	ps.Players = append(ps.Players, player)
	return player
}

func (ps *PlayerService) CatchPokemon(player *model.Player, pokemon *model.Pokemon) {
	capturedPokemon := &model.CapturedPokemon{
		No:          pokemon.No,
		Image:       pokemon.Image,
		Name:        pokemon.Name,
		Level:       1,
		Exp:         0,
		EVs:         "0",
		HP:          pokemon.HP,
		Attack:      pokemon.Attack,
		Defense:     pokemon.Defense,
		SpAttack:    pokemon.SpAttack,
		SpDefense:   pokemon.SpDefense,
		Speed:       pokemon.Speed,
		TotalEvs:    pokemon.TotalEvs,
		Type:        pokemon.Type,
		Height:      pokemon.Height,
		Weight:      pokemon.Weight,
		CatchRate:   pokemon.CatchRate,
		GenderRatio: pokemon.GenderRatio,
		EggGroups:   pokemon.EggGroups,
		HatchSteps:  pokemon.HatchSteps,
		Abilities:   pokemon.Abilities,
		Strengths:   pokemon.Strengths,
		Weaknesses:  pokemon.Weaknesses,
		Evolutions:  pokemon.Evolutions,
		Moves:       pokemon.Moves,
	}
	player.Pokemons = append(player.Pokemons, capturedPokemon)

	ps.SavePlayerPokemons(player)
}

func (ps *PlayerService) SavePlayerPokemons(player *model.Player) {
	ps.PlayerPokemonMap[player.Name] = player.Pokemons
	log.Println(ps.PlayerPokemonMap[player.Name])
	data, err := json.MarshalIndent(ps.PlayerPokemonMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal player Pokémon data: %v", err)
	}

	err = ioutil.WriteFile(ps.PlayerPokemonDataFile, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write player Pokémon data file: %v", err)
	}
}

func (ps *PlayerService) IsValidMove(player *model.Player, worldService *WorldService) bool {
	pos := player.Position
	return pos.X >= worldService.MinX() && pos.X <= worldService.MaxX() &&
		pos.Y >= worldService.MinY() && pos.Y <= worldService.MaxY()
}

func (ps *PlayerService) MovePlayer(player *model.Player, direction string, worldService *WorldService) {
	switch direction {
	case "up":
		ps.MoveUp(player)
	case "down":
		ps.MoveDown(player)
	case "left":
		ps.MoveLeft(player)
	case "right":
		ps.MoveRight(player)
	}

	log.Printf("Player %s moved %s to position %+v\n", player.Name, direction, player.Position) // Log player movement

	pos := player.Position
	if worldService.HasPokemon(pos) {
		pokemon := worldService.CapturePokemon(pos)
		ps.CatchPokemon(player, pokemon)
	}
}

func (ps *PlayerService) AutoMovePlayer(player *model.Player, worldService *WorldService, delay time.Duration, broadcast chan<- string) *model.Pokemon {
	directions := []func(*model.Player){
		ps.MoveUp,
		ps.MoveDown,
		ps.MoveLeft,
		ps.MoveRight,
	}
	directionNames := []string{"up", "down", "left", "right"}

	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			moveIndex := rand.Intn(len(directions))
			move := directions[moveIndex]
			direction := directionNames[moveIndex]

			tempPlayer := *player
			move(&tempPlayer)

			if ps.IsValidMove(&tempPlayer, worldService) {
				move(player)
				logMessage := fmt.Sprintf("Player %s moved %s to position %+v", player.Name, direction, player.Position)
				log.Println(logMessage)

				broadcast <- logMessage

				pos := player.Position
				if worldService.HasPokemon(pos) {
					pokemon := worldService.CapturePokemon(pos)
					broadcast <- fmt.Sprintf("Player %s found a Pokémon: %s. Do you want to catch it? (yes/no): ", player.Name, pokemon.Name)
					return pokemon
				}
			} else {
				log.Printf("Player %s tried to move %s but would move outside the world boundaries.\n", player.Name, direction)
			}
		}
	}
}

func (ps *PlayerService) GetPlayerByName(name string) *model.Player {
	for _, player := range ps.Players {
		if player.Name == name {
			return player
		}
	}
	return nil
}

func (ps *PlayerService) SavePlayerList(filename string, playerPokemonMap map[string][]*model.CapturedPokemon) error {
	data, err := json.MarshalIndent(playerPokemonMap, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PlayerService) MoveUp(player *model.Player) {
	player.Position.Y++
}

func (ps *PlayerService) MoveDown(player *model.Player) {
	player.Position.Y--
}

func (ps *PlayerService) MoveLeft(player *model.Player) {
	player.Position.X--
}

func (ps *PlayerService) MoveRight(player *model.Player) {
	player.Position.X++
}
