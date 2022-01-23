package action

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"strings"
)

type recreational struct {
	db         *repository.DB
	menuFood   []model.Food
	weaponList []model.Weapon
	charCfg    config.CharacterConfig
	log        logger.Logger
	msgToUser  config.MsgToUser
}

func NewRecreational(db *repository.DB, charCfg config.CharacterConfig, log logger.Logger, msgToUser config.MsgToUser) *recreational {
	return &recreational{
		db:        db,
		charCfg:   charCfg,
		log:       log,
		msgToUser: msgToUser,
	}
}

func (r *recreational) RecreationalMain(userChar model.Character, orderItems model.Items) (model.Items, model.Character, string) {
	availableItems := model.Items{}
	var status string

	availableItems.Weapons = r.getWeaponList()
	availableItems.Foods = r.getFoodList()

	if orderItems.Weapons != nil {
		status, userChar = r.armament(userChar, orderItems.Weapons[0])
	}
	if orderItems.Foods != nil {
		status, userChar = r.eat(userChar, orderItems.Foods[0])
	}

	return availableItems, userChar, status
}

func (r *recreational) getWeaponList() []model.Weapon {
	if len(r.weaponList) == 0 {
		r.weaponList = r.db.Postgres.GetWeaponAll()
		return r.weaponList
	}
	return nil
}

func (r *recreational) getFoodList() []model.Food {
	if len(r.menuFood) == 0 {
		r.menuFood = r.db.Postgres.GetFoodAll()
		return r.menuFood
	}
	return nil
}

func (r *recreational) eat(character model.Character, order model.Food) (string, model.Character) {
	if order.Price > character.CoinAmount {
		return r.msgToUser.ActionMsg.InvalidSum, character
	}

	if character.RemainHealth == character.ThresholdHealth {
		return r.msgToUser.ActionMsg.NoNeedToRecover, character
	}

	for i, food := range r.menuFood {
		if strings.Contains(food.Name, order.Name) {
			character = r.restoreHealth(character, food.RestoreHp, order.Price)
		} else if i == len(r.menuFood) {
			return r.msgToUser.ActionMsg.InvalidFood, character
		}
	}

	r.menuFood = []model.Food{}

	return ``, character
}

func (r *recreational) armament(character model.Character, order model.Weapon) (string, model.Character) {
	if len(character.Inventory.Weapons) > r.charCfg.Restriction.WeaponInventorySize {
		return fmt.Sprintf(r.msgToUser.ActionMsg.InvalidNumberOfWeapons, r.charCfg.Restriction.WeaponInventorySize), character
	}

	if order.Price > character.CoinAmount {
		return r.msgToUser.ActionMsg.InvalidSum, character
	}

	for i, weapon := range r.weaponList {
		if strings.Contains(weapon.Name, order.Name) {
			character.Inventory.Weapons = append(character.Inventory.Weapons, weapon)
			character.CoinAmount -= weapon.Price
		} else if i == len(r.weaponList) {
			return r.msgToUser.ActionMsg.InvalidWeapon, character
		}
	}

	r.weaponList = []model.Weapon{}

	return ``, character
}

func (r *recreational) restoreHealth(character model.Character, toRestore, price int) model.Character {
	character.RemainHealth += toRestore
	if character.RemainHealth > character.ThresholdHealth {
		character.RemainHealth = character.ThresholdHealth
	}
	character.CoinAmount -= price
	return character
}
