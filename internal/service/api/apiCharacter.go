package api

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
)

type CharacterService struct {
	db         *repository.DB
	msgToUser  config.MsgToUser
	charConfig config.CharacterConfig
}

func NewCharacterService(db *repository.DB, msgToUser config.MsgToUser, charConfig config.CharacterConfig) *CharacterService {
	return &CharacterService{
		db:         db,
		msgToUser:  msgToUser,
		charConfig: charConfig,
	}
}

func (c CharacterService) CreateCharacter(character model.Character) (int, error) {
	return c.db.Postgres.CharacterData.Create(character)
}

func (c CharacterService) GetAllCharacters(userId int) []model.Character {
	return c.db.Postgres.CharacterData.ReadAll(userId)
}

func (c CharacterService) SelectChar(userId, charId int) model.Character {
	return c.setChar(userId, charId)
}

func (c CharacterService) setChar(userId, charId int) model.Character {
	selectedChar := c.db.Postgres.CharacterData.ReadOne(userId, charId)

	selectedChar.ThresholdHealth = c.calculateThresholdHP(selectedChar.Weight)
	selectedChar.RemainHealth = selectedChar.ThresholdHealth

	selectedChar.ThresholdEnergy = c.calculateThresholdMP(selectedChar.Weight)
	selectedChar.RemainEnergy = selectedChar.ThresholdEnergy
	selectedChar.CoinAmount = 100

	return selectedChar
}

func (c CharacterService) calculateThresholdMP(weight int) int {
	mp := 100
	for i := c.charConfig.MinCharWeight; i <= weight; i += 10 {
		mp -= 10
	}
	return mp
}

func (c CharacterService) calculateThresholdHP(weight int) int {
	hp := 50
	for i := c.charConfig.MinCharWeight; i <= weight; i += 10 {
		hp += 10
	}
	return hp
}

func (c CharacterService) Update(character model.Character) error {
	return c.db.Postgres.CharacterData.Update(character)
}

func (c CharacterService) Delete(userId, charId int) error {
	return c.db.Postgres.CharacterData.Delete(userId, charId)
}

func (c CharacterService) DeleteAll(userId int) error {
	return c.db.Postgres.CharacterData.DeleteAll(userId)
}
