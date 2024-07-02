package model

type Pokemon struct {
	No          int         `json:"no"`
	Image       string      `json:"image"`
	Name        string      `json:"name"`
	Exp         int         `json:"exp"`
	HP          int         `json:"hp"`
	Attack      int         `json:"attack"`
	Defense     int         `json:"defense"`
	SpAttack    int         `json:"sp_attack"`
	SpDefense   int         `json:"sp_defense"`
	Speed       int         `json:"speed"`
	TotalEvs    int         `json:"total_evs"`
	Type        []string    `json:"type"`
	Level       int         `json:"level"`
	Height      string      `json:"height"`
	Weight      string      `json:"weight"`
	CatchRate   string      `json:"catch_rate"`
	GenderRatio string      `json:"gender_ratio"`
	EggGroups   []string    `json:"egg_groups"`
	HatchSteps  int         `json:"hatch_steps"`
	Abilities   []string    `json:"abilities"`
	EVs         string      `json:"evs"`
	Strengths   []Strength  `json:"strengths"`
	Weaknesses  []Weakness  `json:"weaknesses"`
	Evolutions  []Evolution `json:"evolutions"`
	Moves       Moves       `json:"moves"`
}

type Strength struct {
	Type       string `json:"type"`
	Multiplier string `json:"multiplier"`
}

type Weakness struct {
	Type       string `json:"type"`
	Multiplier string `json:"multiplier"`
}

type Evolution struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Level string `json:"level"`
}

type Move struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Class       string `json:"class"`
	Power       string `json:"power"`
	Accuracy    string `json:"accuracy"`
	PP          string `json:"pp"`
	Description string `json:"description"`
}

type Moves struct {
	Natural []Move `json:"natural"`
	Machine []Move `json:"machine"`
	Tutor   []Move `json:"tutor"`
	Egg     []Move `json:"egg"`
}
