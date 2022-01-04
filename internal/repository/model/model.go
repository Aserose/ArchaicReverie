package model

type User struct {
	Id                 int    `json:"-" db:"id"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	NumberOfCharacters int    `json:"numberOfCharacters"`
}

type Character struct {
	CharId  int    `json:"charId"`
	OwnerId int    `json:"ownerId"`
	Name    string `json:"name"`
	Growth  int    `json:"growth"`
	Weight  int    `json:"weight"`
}

type Action struct {
	InAction string `json:"inAction"`
	Jump     `json:"jump"`
}

type Jump struct {
	SquatDepth   int `json:"squatDepth"`
	ArmAmplitude int `json:"armAmplitude"`
	BodyTilt     int `json:"bodyTilt"`
	RunUp        int `json:"runUp"`
}

type Location struct {
	TotalSumValues          int `json:"totalSumValues"`
	TotalClarity            int `json:"totalLightLevel"`
	TotalDifficultyMovement int `json:"totalDifficultyMovement"`
	Place                   `json:"places" db:"places"`
	Weather                 `json:"weathers" db:"weathers"`
	TimeOfDay               `json:"times" db:"times"`
	Obstacle                `json:"obstacles" db:"obstacles"`
}

type Place struct {
	Name               string `json:"name" db:"name"`
	DifficultyMovement int    `json:"difficulty_movement" db:"difficulty_movement"`
}

type Weather struct {
	Name               string `json:"name" db:"name"`
	Clarity            int    `json:"clarity" db:"clarity"`
	DifficultyMovement int    `json:"difficulty_movement" db:"difficulty_movement"`
}

type TimeOfDay struct {
	Name    string `json:"name" db:"name"`
	Clarity int    `json:"clarity" db:"clarity"`
}

type Obstacle struct {
	Name   string `json:"name" db:"name"`
	Length int    `json:"length" db:"length"`
	Height int    `json:"height" db:"height"`
}
