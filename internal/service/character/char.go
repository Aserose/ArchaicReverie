package character

import (
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
)

func (c characterService) CreateCharacter(character model.Character) (int, error) {
	return c.db.Postgres.CharacterData.Create(character)
}

func (c characterService) GetAllCharacters(userId int) []model.Character {
	return c.db.Postgres.CharacterData.ReadAll(userId)
}

func (c characterService) SelectChar(userId, charId int) model.Character {
	return c.setChar(userId, charId)
}

func (c characterService) setChar(userId, charId int) model.Character {
	selectedChar := c.db.Postgres.CharacterData.ReadOne(userId, charId)

	selectedChar.ThresholdHealth = c.calculateThresholdHP(selectedChar.Weight)
	selectedChar.RemainHealth = selectedChar.ThresholdHealth

	selectedChar.ThresholdEnergy = c.calculateThresholdMP(selectedChar.Weight)
	selectedChar.RemainEnergy = selectedChar.ThresholdEnergy
	selectedChar.CoinAmount = c.charConfig.Conversion.BaseAmountCoin

	return selectedChar
}

func (c characterService) calculateThresholdMP(weight int) int {
	mp := c.charConfig.Conversion.BaseMP
	for i := c.charConfig.Restriction.MinCharWeight; i <= weight; i += c.charConfig.Conversion.AddMP {
		mp -= c.charConfig.Conversion.AddMP
		if mp == c.charConfig.Conversion.MinThresholdMP {
			break
		}
	}
	return mp
}

func (c characterService) calculateThresholdHP(weight int) int {
	hp := c.charConfig.Conversion.BaseHP
	for i := c.charConfig.Restriction.MinCharWeight; i <= weight; i += c.charConfig.Conversion.AddHP {
		hp += c.charConfig.Conversion.AddHP
	}
	return hp
}

func (c characterService) Update(character model.Character) error {
	return c.db.Postgres.CharacterData.Update(character)
}

func (c characterService) Delete(userId, charId int) error {
	return c.db.Postgres.CharacterData.Delete(userId, charId)
}

func (c characterService) DeleteAll(userId int) error {
	return c.db.Postgres.CharacterData.DeleteAll(userId)
}
