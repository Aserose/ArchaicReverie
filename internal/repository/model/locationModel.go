package model

type Location struct {
	TotalSumValues          int `json:"-"`
	TotalClarity            int `json:"-"`
	TotalDifficultyMovement int `json:"-"`
	Place                   `json:"places" db:"places"`
	Weather                 `json:"weathers" db:"weathers"`
	TimeOfDay               `json:"times" db:"times"`
	Obstacle                `json:"obstacles" db:"obstacles"`
}

type Place struct {
	Name               string `json:"name" db:"name"`
	DifficultyMovement int    `json:"-" db:"difficulty_movement"`
}

type Weather struct {
	Name               string `json:"name" db:"name"`
	Clarity            int    `json:"-" db:"clarity"`
	DifficultyMovement int    `json:"-" db:"difficulty_movement"`
}

type TimeOfDay struct {
	Name    string `json:"name" db:"name"`
	Clarity int    `json:"-" db:"clarity"`
}

type Obstacle struct {
	Name   string `json:"name" db:"name"`
	Length int    `json:"-" db:"length"`
	Height int    `json:"-" db:"height"`
}
