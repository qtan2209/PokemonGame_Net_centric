package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"os"
	"pokecat_pokebat/internal/model"
	"strconv"
	"strings"
	"time"
)

func ScrapePokedexOrg() error {
	var pokedex []model.Pokemon

	l := launcher.New().Headless(true).MustLaunch()
	browser := rod.New().ControlURL(l).MustConnect()
	defer browser.MustClose()

	for i := 1; i <= 120; i++ {
		url := fmt.Sprintf("https://pokedex.org/#/pokemon/%d", i)
		fmt.Println("Generated URL:", url)

		page := browser.MustPage(url).MustWaitLoad()

		// Wait for the JavaScript to render the content
		time.Sleep(3 * time.Second) // Adjust sleep time as needed

		var pokemon model.Pokemon

		pokemon.Name = strings.TrimSpace(page.MustElement(".detail-panel-header").MustText())

		elements := page.MustElements(".detail-national-id span")
		for _, el := range elements {
			fmt.Sscanf(el.MustText(), "#%d", &pokemon.No)
		}

		elements = page.MustElements(".detail-types .monster-type")
		for _, el := range elements {
			pokemon.Type = append(pokemon.Type, strings.TrimSpace(el.MustText()))
		}

		elements = page.MustElements(".detail-stats-row")
		for _, el := range elements {
			stat := strings.TrimSpace(el.MustElement("span:nth-child(1)").MustText())
			value := strings.TrimSpace(el.MustElement("span:nth-child(2) .stat-bar-fg").MustText())

			var val int
			fmt.Sscanf(value, "%d", &val)
			switch stat {
			case "HP":
				pokemon.HP = val
			case "Attack":
				pokemon.Attack = val
			case "Defense":
				pokemon.Defense = val
			case "Speed":
				pokemon.Speed = val
			case "Sp Atk":
				pokemon.SpAttack = val
			case "Sp Def":
				pokemon.SpDefense = val
			}
		}

		elements = page.MustElements(".monster-minutia")
		for _, el := range elements {
			info := strings.TrimSpace(el.MustText())
			if strings.Contains(info, "Height:") && strings.Contains(info, "Weight:") {
				fmt.Sscanf(info, "Height:%s mWeight:%s kg", &pokemon.Height, &pokemon.Weight)
			}
			if strings.Contains(info, "Catch Rate:") && strings.Contains(info, "Gender Ratio:") {
				catchRate := strings.TrimSpace(strings.Split(info, "Catch Rate:")[1])
				pokemon.CatchRate = strings.Split(strings.Split(catchRate, "Gender Ratio:")[0], "%")[0]
				pokemon.GenderRatio = strings.TrimSpace(strings.Split(info, "Gender Ratio:")[1])
			}
			if strings.Contains(info, "Egg Groups:") && strings.Contains(info, "Hatch Steps:") {
				parts := strings.Split(info, "Hatch Steps:")
				eggGroupsPart := strings.TrimSpace(strings.Split(parts[0], "Egg Groups:")[1])
				groupsList := strings.Split(eggGroupsPart, ",")
				for _, group := range groupsList {
					pokemon.EggGroups = append(pokemon.EggGroups, strings.TrimSpace(group))
				}
				var hatchSteps int
				fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &hatchSteps)
				pokemon.HatchSteps = hatchSteps
			}
			if strings.Contains(info, "Abilities:") && strings.Contains(info, "EVs:") {
				abilities := strings.Split(strings.Split(info, "Abilities:")[1], "EVs:")[0]
				abilitiesList := strings.Split(abilities, ",")
				for _, ability := range abilitiesList {
					pokemon.Abilities = append(pokemon.Abilities, strings.TrimSpace(ability))
				}
				pokemon.EVs = strings.TrimSpace(strings.Split(info, "EVs:")[1])
			}
		}

		elements = page.MustElements(".when-attacked-row")
		for _, el := range elements {
			var weakness model.Weakness
			var strength model.Strength

			types := el.MustElements("span.monster-type")
			multipliers := el.MustElements("span.monster-multiplier")

			for i, typ := range types {
				typeText := strings.TrimSpace(typ.MustText())

				if i >= len(multipliers) {
					fmt.Println("Error: No corresponding multiplier for type", typeText)
					continue
				}

				multiplier := strings.TrimSpace(multipliers[i].MustText())
				multiplier = strings.ReplaceAll(multiplier, "x", "")
				multiplierVal, err := strconv.ParseFloat(multiplier, 64)
				if err != nil {
					fmt.Println("Error parsing multiplier:", err)
					continue
				}

				if multiplierVal > 1 {
					weakness.Type = typeText
					weakness.Multiplier = multiplier
					pokemon.Weaknesses = append(pokemon.Weaknesses, weakness)
				} else if multiplierVal < 1 {
					strength.Type = typeText
					strength.Multiplier = multiplier
					pokemon.Strengths = append(pokemon.Strengths, strength)
				}
			}
		}

		elements = page.MustElements(".evolution-row")
		for _, el := range elements {
			labelText := el.MustElement(".evolution-label span").MustText()
			if labelText != "" {
				parts := strings.Split(labelText, " evolves into ")
				if len(parts) == 2 {
					from := strings.TrimSpace(parts[0])
					toAndLevel := strings.Split(parts[1], " at level ")
					if len(toAndLevel) == 2 {
						to := strings.TrimSpace(toAndLevel[0])
						level := strings.TrimSpace(strings.Split(toAndLevel[1], ".")[0])
						evolution := model.Evolution{
							From:  from,
							To:    to,
							Level: level,
						}
						pokemon.Evolutions = append(pokemon.Evolutions, evolution)
					}
				}
			}
		}

		elements = page.MustElements(".monster-moves")
		for _, monsterMoves := range elements {
			moveCategories := monsterMoves.MustElements(".moves-subtitle")
			for _, movesSubtitle := range moveCategories {
				moveType := strings.TrimSpace(movesSubtitle.MustText())

				var moves *[]model.Move
				switch moveType {
				case "Natural Moves":
					moves = &pokemon.Moves.Natural
				case "Machine Moves":
					moves = &pokemon.Moves.Machine
				case "Tutor Moves":
					moves = &pokemon.Moves.Tutor
				case "Egg Moves":
					moves = &pokemon.Moves.Egg
				}

				moveRows := monsterMoves.MustElements(".moves-row")
				for _, moveRow := range moveRows {
					move := extractMove(moveRow)
					*moves = append(*moves, move)
				}
			}
		}

		if pokemon.Name != "" {
			pokedex = append(pokedex, pokemon)
		}

		page.MustClose()
	}

	file, err := json.MarshalIndent(pokedex, "", "  ")
	if err != nil {
		return err
	}

	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	return os.WriteFile("data/pokedex.json", file, 0644)
}

func extractMove(moveRow *rod.Element) model.Move {
	move := model.Move{}

	move.Name = moveRow.MustElement("span:nth-child(2)").MustText()
	move.Type = moveRow.MustElement("span:nth-child(3)").MustText()
	move.Description = moveRow.MustElement(".move-description").MustText()

	stats := moveRow.MustElements(".moves-row-stats span")
	for _, stat := range stats {
		statText := stat.MustText()
		if strings.Contains(statText, "Power") {
			move.Power = strings.TrimSpace(strings.Split(statText, ":")[1])
		} else if strings.Contains(statText, "Acc") {
			move.Accuracy = strings.TrimSpace(strings.Split(statText, ":")[1])
		} else if strings.Contains(statText, "PP") {
			move.PP = strings.TrimSpace(strings.Split(statText, ":")[1])
		}
	}

	return move
}
