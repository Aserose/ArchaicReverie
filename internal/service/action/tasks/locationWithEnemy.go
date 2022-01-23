package tasks

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"sort"
)

type locationWithEnemy struct {
	db         *repository.DB
	msgToUser  config.MsgToUser
	Conditions Condition
}

func NewLocationWithEnemy(db *repository.DB, msgToUser config.MsgToUser, conditions Condition) *locationWithEnemy {
	return &locationWithEnemy{
		db:         db,
		msgToUser:  msgToUser,
		Conditions: conditions,
	}
}

func (l locationWithEnemy) Main(character model.Character, action model.Action) (string, model.Character) {
	switch action.InAction {
	case "hit":
		return l.hit(character, action.Hit)
	case "run":
		return l.run(character, action.Run)
	}
	return l.msgToUser.ActionMsg.InvalidCommand, character
}

func (l locationWithEnemy) hit(character model.Character, action model.Hit) (string, model.Character) {
	result := l.calculateCharSumHit(character, action) + l.calculateChallengeHitTotalSum(character)
	strongestEnemy := l.getStrongestEnemies()

	l.Conditions.ResetLocation()
	l.Conditions.ResetEnemy()

	if 1 >= result && result >= -1 {
		return l.msgToUser.ActionMsg.SuccessfulHit, character
	} else {
		character = enemyDamage(character, l.db.Postgres.EventData.GetActionResult(model.ActionResult{Name: "hit"}), strongestEnemy[0])
		return fmt.Sprintf("%s: %s", l.msgToUser.ActionMsg.FailedHit, fmt.Sprintf(l.msgToUser.ActionMsg.RemainHealth, character.RemainHealth)), character
	}

}

func (l locationWithEnemy) run(character model.Character, action model.Run) (string, model.Character) {
	result := l.calculateChallengeRunTotalSum(character) + l.calculateCharSumRun(character, action)
	strongestEnemy := l.getStrongestEnemies()

	l.Conditions.ResetLocation()
	l.Conditions.ResetEnemy()

	if 1 >= result && result >= -1 {
		return l.msgToUser.ActionMsg.SuccessfulEscape, character
	} else {
		character = enemyDamage(character, l.db.Postgres.EventData.GetActionResult(model.ActionResult{Name: "hit"}), strongestEnemy[0])
		return fmt.Sprintf("%s: %s", l.msgToUser.ActionMsg.FailedEscape, fmt.Sprintf(l.msgToUser.ActionMsg.RemainHealth, character.RemainHealth)), character
	}

}

func (l locationWithEnemy) calculateCharSumRun(character model.Character, action model.Run) int {
	convertGrowth, convertWeight := convertCharParams(character)

	return convertGrowth + convertWeight + action.BodyTilt
}

func (l locationWithEnemy) calculateChallengeRunTotalSum(character model.Character) int {
	var (
		fastestEnemy      = l.getFastestEnemy()
		ChallengeTotalSum int
	)

	if len(l.Conditions.GetConditionEnemy()) > 2 {
		ChallengeTotalSum += 1
	} else {
		ChallengeTotalSum -= 1
	}

	if fastestEnemy.Weight < character.Weight {
		ChallengeTotalSum -= 1
	} else {
		ChallengeTotalSum += 1
	}

	if fastestEnemy.Growth > character.Growth {
		ChallengeTotalSum -= 1
	} else {
		ChallengeTotalSum += 1
	}

	return ChallengeTotalSum

}

func (l locationWithEnemy) calculateChallengeHitTotalSum(character model.Character) int {
	var (
		challengeTotalSum     int
		sortEnemiesByStrength = l.getStrongestEnemies()
	)

	challengeTotalSum += l.weaponChallengeTotalSum(sortEnemiesByStrength) + l.calculateChallengeEnemyParam(character)

	return challengeTotalSum

}

func (l locationWithEnemy) getFastestEnemy() model.Enemy {
	var fastestEnemy model.Enemy

	for _, enemy := range l.Conditions.GetConditionEnemy() {
		if enemy.Class <= 2 {
			if enemy.Growth > fastestEnemy.Growth {
				if enemy.Weight < fastestEnemy.Weight || fastestEnemy.Weight == 0 {
					fastestEnemy = enemy
				}
			}
		}
	}

	return fastestEnemy
}

func (l locationWithEnemy) getStrongestEnemies() []model.Enemy {
	sortEnemiesByStrength := l.Conditions.GetConditionEnemy()

	sort.SliceStable(sortEnemiesByStrength, func(i, j int) bool {
		return sortEnemiesByStrength[i].Class < sortEnemiesByStrength[j].Class
	})
	sort.SliceStable(sortEnemiesByStrength, func(i, j int) bool {
		return (sortEnemiesByStrength[i].Weight + sortEnemiesByStrength[i].Growth) > (sortEnemiesByStrength[j].Weight + sortEnemiesByStrength[j].Growth)
	})

	return sortEnemiesByStrength
}

func (l locationWithEnemy) getStrongestEnemyIndex() []int {
	var indexNumber int
	if indexNumber > len(l.Conditions.GetConditionEnemy()) {
		indexNumber = 2
	} else {
		indexNumber = 1
	}

	strongestEnemiesIndex := make([]int, indexNumber)
	sortedEnemiesByStrength := l.getStrongestEnemies()

	for i, enemy := range l.Conditions.GetConditionEnemy() {
		for j := 0; j < indexNumber; j++ {
			if enemy == sortedEnemiesByStrength[j] {
				strongestEnemiesIndex[j] = i
			}
		}
	}

	return strongestEnemiesIndex
}

func (l locationWithEnemy) weaponChallengeTotalSum(strongestEnemy []model.Enemy) int {
	var weaponChallengeSum int

	sort.Slice(strongestEnemy, func(i, j int) bool {
		return strongestEnemy[i].Weapons.Class > strongestEnemy[j].Weapons.Class
	})

	if strongestEnemy[0].Weapons.Sharp > strongestEnemy[0].Weapons.Weight {
		weaponChallengeSum = 1
	} else {
		weaponChallengeSum = -1
	}

	if strongestEnemy[0].Weapons.Sharp > 2 {
		weaponChallengeSum += 1
	} else {
		weaponChallengeSum -= 1
	}

	return weaponChallengeSum
}

func (l locationWithEnemy) calculateChallengeEnemyParam(character model.Character) int {
	enemy := l.getStrongestEnemies()[0]

	if (enemy.Growth + enemy.Weight) > (character.Growth + character.Weight) {
		return -1
	} else {
		return 1
	}

}

func (l locationWithEnemy) calculateCharParamVariable(character model.Character) int {
	charParamW, charParamH := convertCharParams(character)

	if charParamW == 1 && charParamH == 1 {
		return 1
	}
	if charParamW == -1 && charParamH == -1 {
		return -1
	}
	return randomOutcomeVariable()
}

func (l locationWithEnemy) calculateCharWeaponVariable(character model.Character) int {
	var (
		playerSum    int
		weaponParams []int
	)

	if character.Inventory.Weapons != nil {
		if character.Inventory.Weapons[0].Sharp > character.Inventory.Weapons[0].Weight {
			weaponParams = append(weaponParams, -1)
		} else {
			weaponParams = append(weaponParams, 1)
		}
		if character.Inventory.Weapons[0].Sharp > 2 {
			weaponParams = append(weaponParams, -1)
		} else {
			weaponParams = append(weaponParams, 1)
		}
	}

	for i, w := range weaponParams {
		var plus, minus int
		if w == 1 {
			plus++
		} else {
			minus++
		}
		if i == len(weaponParams)-1 {
			if plus > minus {
				playerSum = 1
			} else {
				playerSum = -1
			}
		}
	}

	return playerSum
}

func (l locationWithEnemy) calculateCharSumHit(character model.Character, action model.Hit) int {
	var playerSum int

	for i, strongestEnemyIndex := range l.getStrongestEnemyIndex() {
		if action.Target == strongestEnemyIndex {
			playerSum = 1
		} else if i == len(l.getStrongestEnemyIndex())-1 {
			playerSum = -1
		}
	}

	return playerSum + l.calculateCharWeaponVariable(character) + l.calculateCharParamVariable(character)
}
