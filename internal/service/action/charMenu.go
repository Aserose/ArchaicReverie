package action

import (
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
)

type characterMenu struct {
}

func NewCharacterMenu() *characterMenu {
	return &characterMenu{}
}

func (c *characterMenu) CharMenu(userCharacter model.Character, items model.Items) model.Character {
	if items.Weapons != nil {
		return c.discardWeapons(userCharacter, items.Weapons)
	}

	return userCharacter
}

func (c characterMenu) discardWeapons(character model.Character, weapons []model.Weapon) model.Character {

	for i, charWeapon := range character.Inventory.Weapons {
		for _, discardWeapon := range weapons {
			if charWeapon.Name == discardWeapon.Name {
				character.Inventory.CoinAmount += (discardWeapon.Price)/2
				character.Inventory.Weapons[i] = model.Weapon{}
			}
		}
	}

	return character
}
