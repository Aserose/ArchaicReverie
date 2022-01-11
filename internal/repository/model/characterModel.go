package model

type Character struct {
	CharId          int    `json:"charId"`
	OwnerId         int    `json:"ownerId"`
	Name            string `json:"name"`
	Growth          int    `json:"growth"`
	Weight          int    `json:"weight"`
	ThresholdHealth int    `json:"thresholdHp"`
	RemainHealth    int    `json:"remainHp"`
	ThresholdEnergy int    `json:"thresholdMp"`
	RemainEnergy    int    `json:"remainMp"`
	Inventory       `json:"inventory" db:"inventory"`
}
