package tasks

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"strings"
)

type Recreational struct {
	db        *repository.DB
	menuFood  []model.Food
	log       logger.Logger
	msgToUser config.MsgToUser
}

func NewRecreational(db *repository.DB, log logger.Logger, msgToUser config.MsgToUser) *Recreational {
	return &Recreational{
		db:        db,
		log:       log,
		msgToUser: msgToUser,
	}
}

func (r *Recreational) GetFoodList() []model.Food {
	if len(r.menuFood) == 0 {
		r.menuFood = r.db.Postgres.GetListFood()
		return r.menuFood
	}
	return []model.Food{}
}

func (r *Recreational) Eat(character model.Character, order model.Food) (string, model.Character) {
	if order.Price > character.CoinAmount {
		return r.msgToUser.ActionMsg.InvalidSum, character
	}

	if character.RemainHealth == character.ThresholdHealth {
		return r.msgToUser.ActionMsg.NoNeedToRecover, character
	}

	for _, food := range r.menuFood {
		if strings.Contains(food.Name, order.Name) {
			character = r.restoreHealth(character, food.RestoreHp, order.Price)
		} else {
			return r.msgToUser.ActionMsg.InvalidFood, character
		}
	}
	return "", character
}

func (r Recreational) restoreHealth(character model.Character, toRestore, price int) model.Character {
	character.RemainHealth += toRestore
	if character.RemainHealth > character.ThresholdHealth {
		character.RemainHealth = character.ThresholdHealth
	}
	character.CoinAmount -= price
	return character
}
