package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"pokecat_pokebat/internal/model"
	"pokecat_pokebat/internal/service"
)

type BattleController struct {
	Logs []string
}

func NewBattleController() *BattleController {
	return &BattleController{
		Logs: []string{},
	}
}

func (bc *BattleController) StartBattle(w http.ResponseWriter, r *http.Request, playerService *service.PlayerService) {
	var request struct {
		Player1          string   `json:"player1"`
		Player2          string   `json:"player2"`
		SelectedPokemons []string `json:"selectedPokemons"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	player1 := playerService.GetPlayerByName(request.Player1)
	player2 := playerService.GetPlayerByName(request.Player2)

	if player1 == nil || player2 == nil {
		http.Error(w, "Both players must be registered", http.StatusBadRequest)
		return
	}

	player1.Pokemons = bc.selectPokemons(player1, request.SelectedPokemons)
	player2.Pokemons = bc.selectRandomPokemons(player2)

	if len(player1.Pokemons) < 3 || len(player2.Pokemons) < 3 {
		http.Error(w, "Both players must have at least 3 Pokémon", http.StatusBadRequest)
		return
	}

	bc.LogBattleStart(player1, player2)

	player1.ActivePokemon = player1.Pokemons[0]
	player2.ActivePokemon = player2.Pokemons[0]

	firstPlayer, secondPlayer := determineFirstPlayer(player1, player2)
	bc.LogTurnOrder(firstPlayer, secondPlayer)

	winner := bc.executeBattle(firstPlayer, secondPlayer)
	bc.LogBattleEnd(winner, player1, player2)

	opponentPokemons := []string{}
	for _, p := range player2.Pokemons {
		opponentPokemons = append(opponentPokemons, p.Name)
	}

	response := struct {
		Winner           string   `json:"winner"`
		Logs             []string `json:"logs"`
		OpponentPokemons []string `json:"opponent_pokemons"`
	}{
		Winner:           winner.Name,
		Logs:             bc.Logs,
		OpponentPokemons: opponentPokemons,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func (bc *BattleController) selectPokemons(player *model.Player, selectedPokemons []string) []*model.CapturedPokemon {
	selected := []*model.CapturedPokemon{}
	for _, pokemonName := range selectedPokemons {
		for _, pokemon := range player.Pokemons {
			if pokemon.Name == pokemonName {
				selected = append(selected, pokemon)
			}
		}
	}
	return selected
}

func (bc *BattleController) selectRandomPokemons(player *model.Player) []*model.CapturedPokemon {
	rand.Shuffle(len(player.Pokemons), func(i, j int) {
		player.Pokemons[i], player.Pokemons[j] = player.Pokemons[j], player.Pokemons[i]
	})
	if len(player.Pokemons) > 3 {
		return player.Pokemons[:3]
	}
	return player.Pokemons
}

func (bc *BattleController) executeBattle(firstPlayer, secondPlayer *model.Player) *model.Player {
	for {
		if firstPlayer.ActivePokemon != nil {
			bc.performTurn(firstPlayer, secondPlayer)
			if bc.isBattleOver(firstPlayer, secondPlayer) {
				return firstPlayer
			}
		}
		if secondPlayer.ActivePokemon != nil {
			bc.performTurn(secondPlayer, firstPlayer)
			if bc.isBattleOver(firstPlayer, secondPlayer) {
				return secondPlayer
			}
		}
	}
}

func determineFirstPlayer(player1, player2 *model.Player) (firstPlayer, secondPlayer *model.Player) {
	if player1.ActivePokemon.Speed > player2.ActivePokemon.Speed {
		return player1, player2
	} else if player2.ActivePokemon.Speed > player1.ActivePokemon.Speed {
		return player2, player1
	} else {
		if rand.Intn(2) == 0 {
			return player1, player2
		}
		return player2, player1
	}
}

func (bc *BattleController) performTurn(attacker, defender *model.Player) {
	damage := bc.calculateDamage(attacker.ActivePokemon, defender.ActivePokemon)
	bc.LogAttack(attacker, defender, damage, "normal attack")
	defender.ActivePokemon.HP -= damage

	if defender.ActivePokemon.HP <= 0 {
		bc.LogFaint(defender)
		defender.ActivePokemon = bc.chooseActivePokemon(defender)
		if defender.ActivePokemon != nil {
			bc.LogSwitch(defender)
		}
	}
}

func (bc *BattleController) calculateDamage(attacker, defender *model.CapturedPokemon) int {
	damage := attacker.Attack - defender.Defense
	if damage < 0 {
		return 0
	}
	return damage
}

func (bc *BattleController) performAttack(attacker, defender *model.Player) {
	isSpecialAttack := rand.Intn(2) == 0
	var damage int
	var attackType string

	if isSpecialAttack {
		damage = calculateSpecialDamage(attacker.ActivePokemon, defender.ActivePokemon)
		attackType = "special attack"
	} else {
		damage = calculateNormalDamage(attacker.ActivePokemon, defender.ActivePokemon)
		attackType = "normal attack"
	}

	if damage < 0 {
		damage = 0
	}

	defender.ActivePokemon.HP -= damage

	bc.LogAttack(attacker, defender, damage, attackType)
}

func calculateNormalDamage(attacker, defender *model.CapturedPokemon) int {
	damage := attacker.Attack - defender.Defense
	return damage
}

func calculateSpecialDamage(attacker, defender *model.CapturedPokemon) int {
	damage := attacker.SpAttack - defender.SpDefense
	return damage
}

func GetPlayerPokemons(w http.ResponseWriter, r *http.Request, playerService *service.PlayerService) {
	playerName := r.URL.Query().Get("name")
	player := playerService.GetPlayerByName(playerName)
	if player == nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(player.Pokemons)
}

func (bc *BattleController) chooseActivePokemon(player *model.Player) *model.CapturedPokemon {
	for _, pokemon := range player.Pokemons {
		if pokemon.HP > 0 {
			return pokemon
		}
	}
	return nil
}

func (bc *BattleController) isBattleOver(player1, player2 *model.Player) bool {
	return bc.areAllPokemonsFainted(player1) || bc.areAllPokemonsFainted(player2)
}

func (bc *BattleController) areAllPokemonsFainted(player *model.Player) bool {
	for _, pokemon := range player.Pokemons {
		if pokemon.HP > 0 {
			return false
		}
	}
	return true
}

func (bc *BattleController) LogBattleStart(player1, player2 *model.Player) {
	bc.Logs = append(bc.Logs, fmt.Sprintf("Battle started between %s and %s", player1.Name, player2.Name))
}

func (bc *BattleController) LogTurnOrder(firstPlayer, secondPlayer *model.Player) {
	bc.Logs = append(bc.Logs, fmt.Sprintf("%s's %s will go first", firstPlayer.Name, firstPlayer.ActivePokemon.Name))
}

func (bc *BattleController) LogSwitch(player *model.Player) {
	bc.Logs = append(bc.Logs, fmt.Sprintf("%s switched to %s", player.Name, player.ActivePokemon.Name))
}

func (bc *BattleController) LogAttack(attacker, defender *model.Player, damage int, attackType string) {
	bc.Logs = append(bc.Logs, fmt.Sprintf("%s's %s used a %s and dealt %d damage to %s's %s",
		attacker.Name, attacker.ActivePokemon.Name, attackType, damage, defender.Name, defender.ActivePokemon.Name))
}

func (bc *BattleController) LogFaint(player *model.Player) {
	bc.Logs = append(bc.Logs, fmt.Sprintf("%s's %s has fainted", player.Name, player.ActivePokemon.Name))
}

func (bc *BattleController) LogBattleEnd(winner, player1, player2 *model.Player) {
	var loserName string
	if winner.Name == player1.Name {
		loserName = player2.Name
	} else {
		loserName = player1.Name
	}
	bc.Logs = append(bc.Logs, fmt.Sprintf("Battle ended. %s won against %s", winner.Name, loserName))
}
