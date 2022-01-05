package data

import (
	"errors"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type PostgresCharacterData struct {
	db     *sqlx.DB
	log    logger.Logger
	logMsg config.LogMsg
	numberCharLimit int
}

func NewCharacterData(db *sqlx.DB, log logger.Logger, logMsg config.LogMsg, numberCharLimit int) *PostgresCharacterData {
	return &PostgresCharacterData{
		db:     db,
		log:    log,
		logMsg: logMsg,
		numberCharLimit: numberCharLimit,
	}
}

func (p PostgresCharacterData) Create(character model.Character) (int, error) {
	var (
		charId    = 0
		numOfChar = len(p.ReadAll(character.OwnerId))
	)

	if numOfChar >= p.numberCharLimit {
		return charId, errors.New(p.logMsg.CharLimitOutErr)
	}

	if 145 > character.Growth || character.Growth > 200 {
		if 40 > character.Weight || character.Weight > 120 {
			return charId, errors.New(p.logMsg.CharGrowthAndWeightOutErr)
		}
		return charId, errors.New(p.logMsg.CharGrowthOutErr)
	}
	if 40 > character.Weight || character.Weight > 120 {
		return charId, errors.New(p.logMsg.CharWeightOutErr)
	}

	if err := p.db.Get(&charId, "INSERT INTO characters (ownerId,name,growth,weight) VALUES ($1,$2,$3,$4) RETURNING charId",
		character.OwnerId, character.Name, character.Growth, character.Weight); err != nil {
		return charId, err
	}

	return charId, nil
}

func (p PostgresCharacterData) ReadAll(userId int) []model.Character {
	var characters []model.Character

	if err := p.db.Select(&characters, "SELECT * FROM characters WHERE ownerId=$1", userId); err != nil {
		p.log.Errorf(p.logMsg.FormatErr, p.log.CallInfoStr(), p.logMsg.Read, err.Error())
	}

	return characters
}

func (p PostgresCharacterData) ReadOne(userId, charId int) model.Character {
	var character model.Character

	p.db.Get(&character, "SELECT * FROM characters WHERE ownerId=$1 AND charId=$2", userId, charId)

	return character
}

func (p PostgresCharacterData) Update(character model.Character) error {
	if 145 > character.Growth || character.Growth > 200 {
		if 40 > character.Weight || character.Weight > 120 {
			return errors.New(p.logMsg.CharGrowthAndWeightOutErr)
		}
		return errors.New(p.logMsg.CharGrowthOutErr)
	}
	if 40 > character.Weight || character.Weight > 120 {
		return errors.New(p.logMsg.CharWeightOutErr)
	}

	_, err := p.db.Queryx("UPDATE characters SET name=$1,growth=$2,weight=$3 WHERE charId=$4 AND ownerId=$5 ",
		character.Name, character.Growth, character.Weight, character.CharId, character.OwnerId)
	if err != nil {
		p.log.Errorf(p.logMsg.FormatErr, p.log.CallInfoStr(), p.logMsg.Update, err.Error())
		return err
	}

	return nil
}

func (p PostgresCharacterData) Delete(userId, charId int) error {
	_, err := p.db.Query("DELETE FROM characters WHERE charId=$1 AND ownerId=$2", charId, userId)
	if err != nil {
		p.log.Errorf(p.logMsg.FormatErr, p.log.CallInfoStr(), p.logMsg.Delete, err.Error())
		return err
	}

	return nil
}

func (p PostgresCharacterData) DeleteAll(userId int) error {
	_, err := p.db.Query("DELETE FROM characters WHERE ownerId=$1", userId)
	if err != nil {
		p.log.Errorf(p.logMsg.FormatErr, p.log.CallInfoStr(), p.logMsg.Delete, err.Error())
		return err
	}
	return nil
}
