package service

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"pokecat_pokebat/internal/model"
	"strconv"
	"strings"
)

const minPokemonPerUser = 3

var realWorldNames = []string{"Alice", "Bob", "Charlie", "David", "Emma", "Frank", "Grace", "Henry", "Ivy", "Jack"}

func GenerateRandomPlayerPokemonLists(pokedex *[]model.Pokemon) map[string]*model.PlayerPokemonList {
	playerPokemonLists := make(map[string]*model.PlayerPokemonList)

	fmt.Print("Enter the number of players to generate: ")
	reader := bufio.NewReader(os.Stdin)
	numPlayersStr, _ := reader.ReadString('\n')
	numPlayersStr = strings.TrimSpace(numPlayersStr)
	numPlayers, err := strconv.Atoi(numPlayersStr)
	if err != nil || numPlayers <= 0 {
		fmt.Println("Invalid number of players. Exiting...")
		os.Exit(1)
	}

	fmt.Print("Enter the minimum number of pokemons per player: ")
	minPokemonsStr, _ := reader.ReadString('\n')
	minPokemonsStr = strings.TrimSpace(minPokemonsStr)
	minPokemons, err := strconv.Atoi(minPokemonsStr)
	if err != nil || minPokemons <= 0 {
		fmt.Println("Invalid minimum number of pokemons. Exiting...")
		os.Exit(1)
	}

	fmt.Print("Enter the maximum number of pokemons per player: ")
	maxPokemonsStr, _ := reader.ReadString('\n')
	maxPokemonsStr = strings.TrimSpace(maxPokemonsStr)
	maxPokemons, err := strconv.Atoi(maxPokemonsStr)
	if err != nil || maxPokemons <= 0 || maxPokemons < minPokemons {
		fmt.Println("Invalid maximum number of pokemons. Exiting...")
		os.Exit(1)
	}

	for i := 1; i <= numPlayers; i++ {
		playerName := getRandomRealWorldName()
		numPokemons := rand.Intn(maxPokemons-minPokemons+1) + minPokemons
		playerPokemonLists[playerName] = GeneratePokemonList(pokedex, numPokemons)
	}

	return playerPokemonLists
}

func GenerateManualPlayerPokemonLists(pokedex *[]model.Pokemon) map[string]*model.PlayerPokemonList {
	playerPokemonLists := make(map[string]*model.PlayerPokemonList)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the number of players: ")
	numPlayersStr, _ := reader.ReadString('\n')
	numPlayersStr = strings.TrimSpace(numPlayersStr)
	numPlayers, err := strconv.Atoi(numPlayersStr)
	if err != nil || numPlayers <= 0 {
		fmt.Println("Invalid number of players. Exiting...")
		os.Exit(1)
	}

	for i := 1; i <= numPlayers; i++ {
		fmt.Printf("Enter name for player %d: ", i)
		playerName, _ := reader.ReadString('\n')
		playerName = strings.TrimSpace(playerName)

		numPokemons := minPokemonPerUser
		playerPokemonLists[playerName] = GeneratePokemonList(pokedex, numPokemons)
	}

	return playerPokemonLists
}

func getRandomRealWorldName() string {
	return realWorldNames[rand.Intn(len(realWorldNames))]
}
