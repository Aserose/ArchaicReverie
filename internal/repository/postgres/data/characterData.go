package data

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
)

const CharacterDataPackageName = "CharacterData"

type PostgresCharacterData struct {
	db     *sqlx.DB
	log    logger.Logger
	logMsg config.LogMsg
}

func NewCharacterData(db *sqlx.DB, log logger.Logger, logMsg config.LogMsg) *PostgresCharacterData {
	return &PostgresCharacterData{
		db:     db,
		log:    log,
		logMsg: logMsg,
	}
}

func (p *PostgresCharacterData) Create(character model.Character) error {
	_, err := p.db.Queryx("INSERT INTO characters (ownerId,name,growth,weight) VALUES ($1,$2,$3,$4)",
		character.OwnerId, character.Name, character.Growth, character.Weight)
	if err != nil {
		p.log.Errorf(p.logMsg.FormatErr, CharacterDataPackageName, p.logMsg.Create, err.Error())
		return err
	}
	return nil
}

func (p *PostgresCharacterData) ReadAll(userId int) []model.Character {
	var characters []model.Character

	if err := p.db.Select(&characters, "SELECT * FROM characters WHERE ownerId=$1", userId); err != nil {
		p.log.Errorf(p.logMsg.FormatErr, CharacterDataPackageName, p.logMsg.Read, err.Error())
	}

	return characters
}

func (p *PostgresCharacterData) ReadOne(userId, charId int) model.Character {
	var character model.Character

	p.db.Get(&character, "SELECT * FROM characters WHERE ownerId=$1 AND charId=$2", userId, charId)

	return character
}

func (p PostgresCharacterData) Update(character model.Character) error {
	_, err := p.db.Queryx("UPDATE characters SET name=$1,growth=$2,weight=$3 WHERE charId=$5 ",
		character.Name, character.Growth, character.Weight, character.CharId)
	if err != nil {
		p.log.Errorf(p.logMsg.FormatErr, CharacterDataPackageName, p.logMsg.Update, err.Error())
		return err
	}

	return nil
}

func (p PostgresCharacterData) Delete(charId int) error {
	_, err := p.db.Queryx("DELETE FROM characters WHERE charId=$1", charId)
	if err != nil {
		p.log.Errorf(p.logMsg.FormatErr, CharacterDataPackageName, p.logMsg.Delete, err.Error())
		return err
	}

	return nil
}
