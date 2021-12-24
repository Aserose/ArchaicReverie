package data

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type UserData interface {
	Create(username, password string) (id int, authStatus string)
	Read()
	Check(username, password string) (id int, authStatus string)
	Update()
	Delete()
}

type CharacterData interface {
	Create(character model.Character) error
	ReadAll(userId int) []model.Character
	ReadOne(userId, charId int) model.Character
	Update(character model.Character) error
	Delete(charId int) error
}

type EventData interface {
	//TODO
}

type PostgresData struct {
	db *sqlx.DB
	UserData
	CharacterData
}

func NewPostgresData(db *sqlx.DB, msgToUser config.MsgToUser, log logger.Logger, logMsg config.LogMsg) *PostgresData {
	return &PostgresData{
		db:            db,
		UserData:      NewUserData(db, msgToUser, log, logMsg),
		CharacterData: NewCharacterData(db, log, logMsg),
	}
}
