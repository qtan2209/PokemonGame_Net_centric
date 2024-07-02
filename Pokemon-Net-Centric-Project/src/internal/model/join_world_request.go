package model

type JoinWorldRequest struct {
	PlayerName    string `json:"player_name"`
	WorldID       int    `json:"world_id"`
	Mode          string `json:"mode"`
	Direction     string `json:"direction,omitempty"`
	AutoMoveDelay int    `json:"auto_move_delay,omitempty"`
}
