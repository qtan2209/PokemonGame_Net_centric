package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"pokecat_pokebat/controller"
	"pokecat_pokebat/internal/model"
	"pokecat_pokebat/internal/service"
)

var (
	playerService         *service.PlayerService
	worldController       *controller.WorldController
	pokedexData           []model.Pokemon
	playerPokemonDataFile string
	worldDataFile         string
	mu                    sync.Mutex
	activeWorlds          = make(map[int]*controller.WorldController)
	upgrader              = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan string)
)

func init() {
	rand.Seed(time.Now().UnixNano())
	var err error
	playerPokemonDataFile = filepath.Join("data", "player_pokemon_list.json")
	worldDataFile = filepath.Join("data", "world_data.json")
	pokedexData, err = service.LoadPokedexData("data/pokedex.json")
	if err != nil {
		fmt.Printf("Failed to load Pokedex data: %v\n", err)
		os.Exit(1)
	}
	playerService = service.NewPlayerService(playerPokemonDataFile)
}

func createWorld(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	worldID := rand.Int()
	worldController = controller.NewWorldController(worldID)
	activeWorlds[worldID] = worldController

	log.Printf("World created with ID: %d\n", worldID)

	go worldController.SpawnPokemons(pokedexData, 100)

	worldData := map[string]int{"worldID": worldID}
	file, _ := json.MarshalIndent(worldData, "", "  ")
	_ = os.WriteFile(worldDataFile, file, 0644)

	response := struct {
		Message      string         `json:"message"`
		WorldID      int            `json:"world_id"`
		ActiveWorlds map[int]string `json:"active_worlds"`
	}{
		Message:      fmt.Sprintf("World with ID %d created successfully!", worldID),
		WorldID:      worldID,
		ActiveWorlds: getActiveWorldsMap(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getActiveWorlds(w http.ResponseWriter, r *http.Request) {
	activeWorldIDs := getActiveWorldsMap()
	err := json.NewEncoder(w).Encode(activeWorldIDs)
	if err != nil {
		http.Error(w, "Failed to get active worlds", http.StatusInternalServerError)
		return
	}
}

func getActiveWorldsMap() map[int]string {
	activeWorldIDs := make(map[int]string)
	for id := range activeWorlds {
		activeWorldIDs[id] = fmt.Sprintf("World %d", id)
	}
	return activeWorldIDs
}

func handleClientConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			delete(clients, conn)
			break
		}

		var request model.JoinWorldRequest
		err = json.Unmarshal(msg, &request)
		if err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		err = handleJoinWorld(request, conn)
		if err != nil {
			log.Println("Error handling join world:", err)
			continue
		}
	}
}

func handleJoinWorld(request model.JoinWorldRequest, conn *websocket.Conn) error {
	mu.Lock()
	defer mu.Unlock()

	player := playerService.GetPlayerByName(request.PlayerName)
	fmt.Println(player)
	fmt.Println(player)
	if player == nil {
		return fmt.Errorf("player not found")
	}

	worldController, exists := activeWorlds[request.WorldID]
	if !exists {
		return fmt.Errorf("world not found")
	}

	switch request.Mode {
	case "manual":
		movePlayerManually(player, request.Direction, worldController.WorldService)
	case "auto":
		startAutoMove(player, worldController.WorldService, time.Duration(request.AutoMoveDelay)*time.Millisecond, conn)
	default:
		return fmt.Errorf("invalid mode")
	}

	return nil
}

func startAutoMove(player *model.Player, worldService *service.WorldService, delay time.Duration, conn *websocket.Conn) {
	go func() {
		ticker := time.NewTicker(delay)
		defer ticker.Stop()

		for range ticker.C {
			pokemon := playerService.AutoMovePlayer(player, worldService, delay, broadcast)
			log.Printf("Player %s moved automatically.\n", player.Name)

			if pokemon != nil {
				message := fmt.Sprintf("Player %s moved automatically and found a Pokémon: %s. Do you want to catch it? (yes/no): ", player.Name, pokemon.Name)
				log.Println("Broadcasting message:", message)

				if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("Error writing message:", err)
					return
				}

				_, userResponse, err := conn.ReadMessage()
				if err != nil {
					log.Println("Error reading user response:", err)
				}

				if strings.Contains(string(userResponse), "yes") {
					playerService.CatchPokemon(player, pokemon)
					log.Printf("Player %s caught the Pokémon %s!\n", player.Name, pokemon.Name)
				} else {
					log.Printf("Player %s decided not to catch the Pokémon %s.\n", player.Name, pokemon.Name)
				}

			}
		}
	}()
}

func movePlayerManually(player *model.Player, direction string, worldService *service.WorldService) {
	playerService.MovePlayer(player, direction, worldService)
	log.Printf("Player %s moved %s.\n", player.Name, direction)
}

func broadcastMessages() {
	for {
		message := <-broadcast

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("Error broadcasting message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controller.LoginPlayer(w, r, playerService, playerPokemonDataFile)
	}).Methods("POST")
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controller.RegisterPlayer(w, r, playerService, playerPokemonDataFile)
	}).Methods("POST")
	r.HandleFunc("/createWorld", createWorld).Methods("POST")
	r.HandleFunc("/getActiveWorlds", getActiveWorlds).Methods("GET")
	r.HandleFunc("/ws", handleClientConnection)

	go broadcastMessages()

	log.Fatal(http.ListenAndServe(":8080", r))
}
