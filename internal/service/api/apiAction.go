package api

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"math/rand"
	"strings"
	"time"
)

type ActionScene struct {
	db            *repository.DB
	utilitiesStr  config.UtilitiesStr
	locationScene model.Location
	menuFood      []model.Food
	msgToUser     config.MsgToUser
}

func NewActionScene(db *repository.DB, utilitiesStr config.UtilitiesStr, msgToUser config.MsgToUser) *ActionScene {
	return &ActionScene{
		db:           db,
		utilitiesStr: utilitiesStr,
		msgToUser:    msgToUser,
	}
}

func (a *ActionScene) GenerateScene() string {
	if a.locationScene == (model.Location{}) {

		a.locationScene = a.db.Postgres.EventData.GenerateEventLocation()

		return fmt.Sprintf(a.msgToUser.ActionMsg.LocationFormat,
			a.locationScene.Place.Name, a.locationScene.Weather.Name,
			a.locationScene.TimeOfDay.Name, a.locationScene.Obstacle.Name)
	}

	return ""
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

func (a *ActionScene) Jump(character model.Character, jumpPosition model.Jump) (string, model.Character) {
	if character.RemainHealth<9{
		if character.RemainHealth<0{
			character.RemainHealth=0
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

	a.locationScene = model.Location{}

	if 1 >= result && result >= -1 {
		return a.msgToUser.ActionMsg.JumpOver, character
	} else {
		return fmt.Sprintf("%s %s",a.msgToUser.ActionMsg.JumpFell, fmt.Sprintf(a.msgToUser.ActionMsg.RemainHealth, character.RemainHealth)),
		damage(character, 10)
	}
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
