package data

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
	"log"
)

const (
	empty               = ""
	userDataPackageName = "CharacterData"
)

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

func (p *PostgresUserData) Create(username, password string) (id int, authStatus string) {

	var person model.User

	p.db.Get(&person, "SELECT id FROM users WHERE username=$1", username)

	if person.Id != 0 {
		return 0, p.msgToUser.AuthStatus.BusyUsername
	}

	_, err := p.db.Queryx("INSERT INTO users (username,password) VALUES ($1,$2)", username, password)
	if err != nil {

	}

	p.db.Get(&person, "SELECT id FROM users WHERE username=$1 AND password=$2", username, password)

	return person.Id, empty
}

func (p *PostgresUserData) Check(username, password string) (id int, authStatus string) {
	var person model.User

	p.db.Get(&person, "SELECT username FROM users WHERE username=$1", username)

	if person.Username == empty {
		return person.Id, p.msgToUser.AuthStatus.InvalidUsername
	}

	p.db.Get(&person, "SELECT * FROM users WHERE username=$1 AND password=$2", username, password)

	if person.Id == 0 {
		return person.Id, p.msgToUser.AuthStatus.InvalidPassword
	}

	return person.Id, empty
}

func (p *PostgresUserData) Read() {

	var per model.User

	rows, err := p.db.Queryx("SELECT * FROM users")
	if err != nil {
		p.log.Errorf(p.logMsg.FormatErr, userDataPackageName, p.logMsg.Read, err.Error())
	}

	for rows.Next() {
		err := rows.StructScan(&per)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%#v\n", per)
	}
}

func (p *PostgresUserData) Update() {
	//TODO
}

func (p *PostgresUserData) Delete() {
	//TODO
}
