package tasks

import (
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"math/rand"
	"time"
)

func damage(char model.Character, damage int) model.Character {
	char.RemainHealth -= damage
	return char
}

func convertCharParams(character model.Character) (int, int) {
	var calcGrowth, calcWeight int

	if character.Growth >= 170 {
		calcGrowth = 1
	} else {
		calcGrowth = -1
	}

	if character.Weight <= 85 {
		calcWeight = 1
	} else {
		calcWeight = -1
	}

	return calcGrowth, calcWeight
}

func validateActionJumpPosition(jumpPosition ...int) bool {
	for _, v := range jumpPosition {
		if v > 1 || v < -1 {
			return false
		}
	}
	return true
}

func randomOutcomeVariable() int {
	rand.Seed(time.Now().UnixNano())
	min := -1
	max := 1
	return rand.Intn(max-min+1) + min
}
