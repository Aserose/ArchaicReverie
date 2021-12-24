package api

import (
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
)

type CharacterApi interface {
	CreateCharacter(character model.Character) error
	GetAllCharacters(userId int) []model.Character
	GetOne(userId, charId int) model.Character
	Update(character model.Character) error
	Delete(charId int) error
}

type ServiceApi struct {
	db *repository.DB
	CharacterApi
}

func NewApiService(db *repository.DB) *ServiceApi {
	return &ServiceApi{
		db:           db,
		CharacterApi: NewCharacterData(db),
	}
}
