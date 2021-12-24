package model

type User struct {
	Id       int    `json:"-" db:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
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
	NamePlace               string `json:"place"`
	NameTimeOfDay           string `json:"timeOfDay"`
	NameWeather             string `json:"weather"`
	NameObstacle string `json:"obstacle"`
	TotalClarity            int    `json:"totalLightLevel"`
	TotalDifficultyMovement int    `json:"totalDifficultyMovement"`
	Place
	Weather
	TimeOfDay
	Obstacle
}

type Place struct {
	DifficultyMovement int
}

type Weather struct {
	Clarity            int
	DifficultyMovement int
}

type TimeOfDay struct {
	Clarity int
}

type Obstacle struct {
	Length int
	Height int
}
