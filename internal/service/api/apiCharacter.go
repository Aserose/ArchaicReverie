package api

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
)

type CharacterService struct {
	db        *repository.DB
	msgToUser config.MsgToUser
}

func NewCharacterService(db *repository.DB, msgToUser config.MsgToUser) *CharacterService {
	return &CharacterService{
		db:        db,
		msgToUser: msgToUser,
	}
}

func (c CharacterService) CreateCharacter(character model.Character) (int, error) {
	return c.db.Postgres.CharacterData.Create(character)
}

func (c CharacterService) GetAllCharacters(userId int) []model.Character {
	return c.db.Postgres.CharacterData.ReadAll(userId)
}

func (c CharacterService) GetOne(userId, charId int) model.Character {
	return c.db.Postgres.CharacterData.ReadOne(userId, charId)
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
