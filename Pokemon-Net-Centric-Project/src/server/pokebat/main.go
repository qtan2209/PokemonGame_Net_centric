package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"pokecat_pokebat/controller"
	"pokecat_pokebat/internal/service"
)

func main() {
	playerPokemonDataFile := "data/player_pokemon_list.json"
	playerService := service.NewPlayerService(playerPokemonDataFile)
	battleController := controller.NewBattleController()

	r := mux.NewRouter()

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controller.LoginPlayer(w, r, playerService, playerPokemonDataFile)
	}).Methods("POST")
	r.HandleFunc("/start-battle", func(w http.ResponseWriter, r *http.Request) {
		battleController.StartBattle(w, r, playerService)
	}).Methods("POST")
	r.HandleFunc("/player-pokemons", func(w http.ResponseWriter, r *http.Request) {
		controller.GetPlayerPokemons(w, r, playerService)
	}).Methods("GET")

	port := "8080"
	fmt.Printf("Server is running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
