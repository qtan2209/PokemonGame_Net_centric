package service

import (
	"pokecat_pokebat/internal/model"
)

func CapturePokemon(pokedex *model.Pokemon, level int, evs string) model.CapturedPokemon {
	return model.CapturedPokemon{
		No:          pokedex.No,
		Image:       pokedex.Image,
		Name:        pokedex.Name,
		Level:       level,
		Exp:         0,
		EVs:         evs,
		HP:          pokedex.HP,
		Attack:      pokedex.Attack,
		Defense:     pokedex.Defense,
		SpAttack:    pokedex.SpAttack,
		SpDefense:   pokedex.SpDefense,
		Speed:       pokedex.Speed,
		TotalEvs:    pokedex.TotalEvs,
		Type:        pokedex.Type,
		Height:      pokedex.Height,
		Weight:      pokedex.Weight,
		CatchRate:   pokedex.CatchRate,
		GenderRatio: pokedex.GenderRatio,
		EggGroups:   pokedex.EggGroups,
		HatchSteps:  pokedex.HatchSteps,
		Abilities:   pokedex.Abilities,
		Strengths:   pokedex.Strengths,
		Weaknesses:  pokedex.Weaknesses,
		Evolutions:  pokedex.Evolutions,
		Moves:       pokedex.Moves,
	}
}
