package controller

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"pokecat_pokebat/internal/model"
	"pokecat_pokebat/internal/service"
	"sync"
)

var (
	mu sync.Mutex
)

func RegisterPlayer(w http.ResponseWriter, r *http.Request, playerService *service.PlayerService, playerPokemonDataFile string) {
	var request struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if playerService.GetPlayerByName(request.Name) != nil {
		http.Error(w, "Player already registered", http.StatusConflict)
		return
	}

	player := playerService.CreatePlayer(request.Name, model.Position{X: rand.Intn(1000), Y: rand.Intn(1000)})

	playerPokemonMap := playerService.LoadPlayerList(playerPokemonDataFile)
	playerPokemonMap[request.Name] = player.Pokemons

	err = playerService.SavePlayerList(playerPokemonDataFile, playerPokemonMap)
	if err != nil {
		http.Error(w, "Failed to save player data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LoginPlayer(w http.ResponseWriter, r *http.Request, playerService *service.PlayerService, playerPokemonDataFile string) {
	var request struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	player := playerService.GetPlayerByName(request.Name)

	if player == nil {
		http.Error(w, "Player not found", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
