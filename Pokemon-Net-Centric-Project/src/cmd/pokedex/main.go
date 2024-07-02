package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"pokecat_pokebat/internal/model"
	"pokecat_pokebat/internal/service"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	if _, err := os.Stat("data/pokedex.json"); os.IsNotExist(err) {
		fmt.Println("pokedex.json not found, scraping data...")
		err := service.ScrapePokedexOrg()
		if err != nil {
			fmt.Println("Error scraping pokedex data:", err)
			return
		}
	}

	pokedex, err := service.LoadPokedex("data/pokedex.json")
	if err != nil {
		fmt.Println("Error loading pokedex:", err)
		return
	}

	fmt.Println("Choose generation method:")
	fmt.Println("1. Random")
	fmt.Println("2. Manual Input")
	fmt.Print("Enter your choice: ")
	reader := bufio.NewReader(os.Stdin)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || (choice != 1 && choice != 2) {
		fmt.Println("Invalid choice. Exiting...")
		return
	}

	var playerPokemonLists map[string]*model.PlayerPokemonList

	switch choice {
	case 1:
		playerPokemonLists = service.GenerateRandomPlayerPokemonLists(pokedex)
	case 2:
		playerPokemonLists = service.GenerateManualPlayerPokemonLists(pokedex)
	}

	err = service.SavePlayerPokemonLists("data/player_pokemon_list.json", playerPokemonLists)
	if err != nil {
		fmt.Println("Error saving player's pokemon lists:", err)
	} else {
		fmt.Println("Player's pokemon lists saved successfully.")
	}
}
