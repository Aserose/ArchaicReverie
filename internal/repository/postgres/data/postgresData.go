package data

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type UserData interface {
	Create(username, password string) (id int, authStatus string)
	Check(username, password string) (id int, authStatus string)
	UpdatePassword(userId int, newPassword string) error
	DeleteAccount(userId int, password string) error
}

type CharacterData interface {
	Create(character model.Character) (int, error)
	ReadAll(userId int) []model.Character
	ReadOne(userId, charId int) model.Character
	Update(character model.Character) error
	Delete(userId, charId int) error
	DeleteAll(userId int) error
}

type EventData interface {
	GenerateEventLocation() model.Location
	GetListFood() []model.Food
	GetFood(name string) model.Food
	GetActionResult(actionResult model.ActionResult) model.ActionResult
	GenerateEnemy(settingEnemy []model.Enemy) []model.Enemy
}

type PostgresData struct {
	db *sqlx.DB
	UserData
	CharacterData
	EventData
}

func NewPostgresData(db *sqlx.DB, msgToUser config.MsgToUser, log logger.Logger, logMsg config.LogMsg, charConfig config.CharacterConfig) *PostgresData {
	return &PostgresData{
		db:            db,
		UserData:      NewUserData(db, msgToUser, log, logMsg),
		CharacterData: NewCharacterData(db, log, logMsg, charConfig),
		EventData:     NewEventData(db, log, logMsg),
	}
}
