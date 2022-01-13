package service

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/internal/service/action"
	"github.com/Aserose/ArchaicReverie/internal/service/authorization"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
)

type Authorization interface {
	SignIn(username, password string) (string, int)
	SignUp(username, password string) (string, int)
	UpdateToken(userId int, character model.Character) string
	UpdatePassword(username, password, newPassword string) string
	Verification(token string) (int, model.Character, error)
	DeleteAccount(username, password string) string
}

type Character interface {
	CreateCharacter(character model.Character) (int, error)
	GetAllCharacters(userId int) []model.Character
	SelectChar(userId, charId int) model.Character
	Update(character model.Character) error
	Delete(userId, charId int) error
	DeleteAll(userId int) error
}

type Action interface {
	GenerateScene() map[string]interface{}
	Action(character model.Character, action model.Action) (string, model.Character)
	GetFoodList() []model.Food
	Eat(character model.Character, order model.Food) (string, model.Character)
}

type Service struct {
	Authorization
	Character
	Action
}

func NewService(db *repository.DB, utilitiesStr config.UtilitiesStr, cfgServices *config.CfgServices,
	msgToUser config.MsgToUser, log logger.Logger, logMsg config.LogMsg, charConfig config.CharacterConfig, genConfig config.GenerationConfig) *Service {
	return &Service{
		Authorization: authorization.NewServiceAuthorization(db, cfgServices, log, logMsg, msgToUser),
		Character:     action.NewCharacterService(db, msgToUser, charConfig),
		Action:        action.NewActionScene(db, genConfig, utilitiesStr, log, msgToUser),
	}
}
