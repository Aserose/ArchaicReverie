package data

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
)

const empty = ""

type PostgresUserData struct {
	db        *sqlx.DB
	msgToUser config.MsgToUser
	log       logger.Logger
	logMsg    config.LogMsg
}

func NewUserData(db *sqlx.DB, msgToUser config.MsgToUser, log logger.Logger, logMsg config.LogMsg) *PostgresUserData {
	return &PostgresUserData{
		db:        db,
		msgToUser: msgToUser,
		log:       log,
		logMsg:    logMsg,
	}
}

func (p *PostgresUserData) Create(username, passwordHash string) (id int, authStatus string) {
	var person model.User

	p.db.Get(&person, `SELECT id FROM users WHERE username=$1`, username)

	if person.Id != 0 {
		return 0, p.msgToUser.AuthStatus.BusyUsername
	}

	_, err := p.db.Queryx(`INSERT INTO users (username,password) VALUES ($1,$2)`, username, passwordHash)
	if err != nil {

	}
	p.db.Get(&person, `SELECT id FROM users WHERE username=$1 AND password=$2`, username, passwordHash)

	return person.Id, empty
}

func (p *PostgresUserData) Check(username, passwordHash string) (id int, authStatus string) {
	var person model.User

	p.db.Get(&person, `SELECT username FROM users WHERE username=$1`, username)

	if person.Username == empty {
		return person.Id, p.msgToUser.AuthStatus.InvalidUsername
	}

	p.db.Get(&person, `SELECT * FROM users WHERE username=$1 AND password=$2`, username, passwordHash)

	if person.Id == 0 {
		return person.Id, p.msgToUser.AuthStatus.InvalidPassword
	}

	return person.Id, empty
}

func (p *PostgresUserData) UpdatePassword(userId int, newPassword string) error {
	if _, err := p.db.Exec(`UPDATE users SET password=$1 WHERE id=$2`, newPassword, userId); err != nil {
		p.log.Errorf(p.logMsg.Format, p.log.CallInfoStr(), err.Error())
		return err
	}
	return nil
}

func (p *PostgresUserData) DeleteAccount(userId int, password string) error {
	if _, err := p.db.Exec(`DELETE FROM users WHERE id=$1 AND password=$2`, userId, password); err != nil {
		p.log.Errorf(p.logMsg.Format, p.log.CallInfoStr(), err.Error())
		return err
	}
	return nil
}
