package api

import (
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
)

type CharacterData struct {
	db *repository.DB
}

func NewCharacterData(db *repository.DB) *CharacterData {
	return &CharacterData{
		db: db,
	}
}

func (c CharacterData) CreateCharacter(character model.Character) error {
	return c.db.Postgres.CharacterData.Create(character)
}

func (c CharacterData) GetAllCharacters(userId int) []model.Character {
	return c.db.Postgres.CharacterData.ReadAll(userId)
}

func (c CharacterData) GetOne(userId, charId int) model.Character {
	return c.db.Postgres.CharacterData.ReadOne(userId, charId)
}

func (c CharacterData) Update(character model.Character) error {
	return c.db.Postgres.CharacterData.Update(character)
}

func (c CharacterData) Delete(charId int) error {
	return c.db.Postgres.CharacterData.Delete(charId)
}
