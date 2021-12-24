package service

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/internal/service/api"
	"github.com/Aserose/ArchaicReverie/internal/service/authorization"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
)

type Authorization interface {
	SignIn(username, password string) string
	SignUp(username, password string) string
	UpdateToken(userId int, character model.Character) string
	Verification(token string) (int, model.Character, error)
}

type Character interface {
	CreateCharacter(character model.Character) error
	GetAllCharacters(userId int) []model.Character
	GetOne(userId, charId int) model.Character
	Update(character model.Character) error
	Delete(charId int) error
}

type Service struct {
	Authorization
	Character
	//TODO
}

func NewService(db *repository.DB, cfgServices *config.CfgServices, log logger.Logger, logMsg config.LogMsg) *Service {
	return &Service{
		Authorization: authorization.NewServiceAuthorization(db, cfgServices, log, logMsg),
		Character:     api.NewApiService(db),
	}
}
