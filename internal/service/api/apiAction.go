package api

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	wr "github.com/mroth/weightedrand"
	"math/rand"
	"strings"
	"time"
)

type ActionScene struct {
	db            *repository.DB
	utilitiesStr  config.UtilitiesStr
	locationScene model.Location
	enemy         []model.Enemy
	menuFood      []model.Food
	log           logger.Logger
	msgToUser     config.MsgToUser
}

func NewActionScene(db *repository.DB, utilitiesStr config.UtilitiesStr, log logger.Logger, msgToUser config.MsgToUser) *ActionScene {
	return &ActionScene{
		db:           db,
		utilitiesStr: utilitiesStr,
		log:          log,
		msgToUser:    msgToUser,
	}
}

func (a *ActionScene) GenerateScene() map[string]interface{} {
	var (
		chooser *wr.Chooser
		err     error
	)

	if a.locationScene == (model.Location{}) && len(a.enemy) == 0 {
		var result map[string]interface{}

		if chooser, err = wr.NewChooser(
			wr.Choice{Item: func() map[string]interface{} {
				a.locationScene = a.db.Postgres.EventData.GenerateEventLocation()
				result = map[string]interface{}{"location": a.locationScene}
				return result
			}(), Weight: 5},

			wr.Choice{Item: func() map[string]interface{} {
				a.locationScene = a.db.Postgres.EventData.GenerateEventLocation()
				a.enemy = a.db.Postgres.EventData.GenerateEnemy(generateSettingEnemy())
				result = map[string]interface{}{"location": a.locationScene, "enemy": a.enemy}
				return result
			}(), Weight: 5}); err != nil {
			a.log.Errorf("%s:%s", a.log.CallInfoStr(), err.Error())
		}

		chooser.Pick()

		return result
	}
	return nil
}

func generateSettingEnemy() []model.Enemy {
	enemies := []model.Enemy{}

	for i := 0; i <= RandInt(1, 4, i); i++ {
		enemies = append(enemies, model.Enemy{Class: RandInt(1, 4, i)})
	}

	return enemies
}

func RandInt(min, max, add int) int {
	rand.Seed(time.Now().UnixNano() + int64(add))
	return rand.Intn(max-min) + min
}

func (a *ActionScene) GetFoodList() []model.Food {
	if len(a.menuFood) == 0 {
		a.menuFood = a.db.Postgres.GetListFood()
		return a.menuFood
	}
	return []model.Food{}
}

func (a *ActionScene) Eat(character model.Character, order model.Food) (string, model.Character) {
	if order.Price > character.CoinAmount {
		return a.msgToUser.ActionMsg.InvalidSum, character
	}
	for _, food := range a.menuFood {
		if strings.Contains(food.Name, order.Name) {
			if character.RemainHealth == character.ThresholdHealth {
				return a.msgToUser.ActionMsg.NoNeedToRecover, character
			}
			character = restoreHealth(character, food.RestoreHp, order.Price)
		} else {
			return a.msgToUser.ActionMsg.InvalidFood, character
		}
	}
	return "", character
}

func restoreHealth(character model.Character, toRestore, price int) model.Character {
	character.RemainHealth += toRestore
	if character.RemainHealth > character.ThresholdHealth {
		character.RemainHealth = character.ThresholdHealth
	}
	character.CoinAmount -= price
	return character
}

func (a *ActionScene) HitOrRun() {
	//TODO event check
}

func (a *ActionScene) Jump(character model.Character, jumpPosition model.Jump) (string, model.Character) {
	//TODO event check

	actionResult := a.db.Postgres.EventData.GetActionResult(model.ActionResult{Name: "fall"})

	if character.RemainHealth < 9 {
		if character.RemainHealth < 0 {
			character.RemainHealth = 0
		}
		return a.msgToUser.ActionMsg.LowHP, character
	}

	if validateActionJumpPosition(
		jumpPosition.RunUp,
		jumpPosition.BodyTilt,
		jumpPosition.ArmAmplitude,
		jumpPosition.SquatDepth) == false {
		return a.utilitiesStr.BadRequest, character
	}

	calcGrowth, calcWeight := convertCharParams(character)

	sumJumpPosition :=
		jumpPosition.RunUp +
			jumpPosition.BodyTilt +
			jumpPosition.ArmAmplitude +
			jumpPosition.SquatDepth +
			calcGrowth +
			calcWeight +
			random()

	result := a.locationScene.TotalSumValues + sumJumpPosition

	a.reset()

	if 1 >= result && result >= -1 {
		return a.msgToUser.ActionMsg.JumpOver, character
	} else {
		return fmt.Sprintf("%s %s", a.msgToUser.ActionMsg.JumpFell, fmt.Sprintf(a.msgToUser.ActionMsg.RemainHealth, character.RemainHealth)),
			damage(character, actionResult.DamageHP)
	}
}

func (a *ActionScene) reset() {
	a.locationScene = model.Location{}
	a.enemy = nil
}

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

func random() int {
	rand.Seed(time.Now().UnixNano())
	min := -1
	max := 1
	return rand.Intn(max-min+1) + min
}
