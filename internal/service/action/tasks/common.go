package tasks

import (
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"math/rand"
	"time"
)

func convertWeaponParams(weapon model.Weapon) int {
	var convertDamage int

	if weapon.Sharp > weapon.Weight {
		convertDamage = 1
	} else {
		convertDamage = -1
	}

	if weapon.Sharp > 2 {
		convertDamage += 1
	} else {
		convertDamage -= 1
	}

	if weapon.Weight > 2 {
		convertDamage -= 1
	} else {
		convertDamage += 1
	}

	return convertDamage

}

func enemyDamage(char model.Character, result model.ActionResult, enemy model.Enemy) model.Character {
	weaponVar, addWeaponDamageHP := convertWeaponParams(enemy.Weapons), 0

	switch weaponVar < 0 {
	case true:
		for i, j := weaponVar, 0; i <= 0; i, j = i+1, j+3 {
			addWeaponDamageHP += j
		}
	case false:
		for i, j := weaponVar, 0; i >= 0; i, j = i-1, j+3 {
			addWeaponDamageHP += j
		}
	}

	char.RemainHealth -= result.DamageHP + addWeaponDamageHP

	return char
}

func damage(char model.Character, result model.ActionResult) model.Character {
	char.RemainHealth -= result.DamageHP
	char.RemainEnergy -= result.DamageMP

	return char
}

func convertCharParams(character model.Character) (int, int) {
	var weaponWeight, convertGrowth, convertWeight int

	for _, weapon := range character.Weapons {
		weaponWeight += weapon.Weight
	}

	if character.Growth >= 170 {
		convertGrowth = 1
	} else {
		convertGrowth = -1
	}

	if (character.Weight + weaponWeight) <= 85 {
		convertWeight = 1
	} else {
		convertWeight = -1
	}

	return convertGrowth, convertWeight
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
