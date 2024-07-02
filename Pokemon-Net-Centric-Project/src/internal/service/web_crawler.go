package service

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"pokecat_pokebat/internal/model"
	"strings"
)

func ScrapePokedex() error {
	pokedex := []model.Pokemon{}

	c := colly.NewCollector()

	c.OnHTML("table tr", func(e *colly.HTMLElement) {
		var pokemon model.Pokemon
		e.ForEach("td", func(index int, element *colly.HTMLElement) {
			text := strings.TrimSpace(element.Text)
			switch index {
			case 0:
				fmt.Sscanf(text, "%d", &pokemon.No)
			case 1:
				pokemon.Image = element.ChildAttr("a img", "src")
			case 2:
				pokemon.Name = strings.TrimSpace(element.ChildText("a"))
			case 3:
				fmt.Sscanf(text, "%d", &pokemon.Exp)
			case 4:
				fmt.Sscanf(text, "%d", &pokemon.HP)
			case 5:
				fmt.Sscanf(text, "%d", &pokemon.Attack)
			case 6:
				fmt.Sscanf(text, "%d", &pokemon.Defense)
			case 7:
				fmt.Sscanf(text, "%d", &pokemon.SpAttack)
			case 8:
				fmt.Sscanf(text, "%d", &pokemon.SpDefense)
			case 9:
				fmt.Sscanf(text, "%d", &pokemon.Speed)
			case 10:
				fmt.Sscanf(text, "%d", &pokemon.TotalEvs)
			}
		})
		if pokemon.Name != "" {
			pokedex = append(pokedex, pokemon)
		}
	})

	c.Visit("http://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_effort_value_yield")

	file, err := json.MarshalIndent(pokedex, "", "  ")
	if err != nil {
		return err
	}

	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	return os.WriteFile("data/pokedex.json", file, 0644)
}
