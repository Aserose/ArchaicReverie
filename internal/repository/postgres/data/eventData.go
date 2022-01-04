package data

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type PostgresEventData struct {
	db     *sqlx.DB
	log    logger.Logger
	logMsg config.LogMsg
}

func NewEventData(db *sqlx.DB, log logger.Logger, logMsg config.LogMsg) *PostgresEventData {
	return &PostgresEventData{
		db:     db,
		log:    log,
		logMsg: logMsg,
	}
}

func (p *PostgresEventData) GenerateEventLocation() model.Location {
	event := model.Location{}

	query := `SELECT *
    			FROM
			(SELECT p.name "places.name", p.difficulty_movement "places.difficulty_movement" 
				FROM places p ORDER BY random()) places 
			    NATURAL FULL JOIN
			(SELECT t.name "times.name", t.clarity "times.clarity" 
				FROM times t ORDER BY random()) times 
			    NATURAL FULL JOIN
			(SELECT w.name "weathers.name", w.clarity "weathers.clarity", w.difficulty_movement "weathers.difficulty_movement" 
				FROM weathers w ORDER BY random()) weathers
				NATURAL FULL JOIN
			(SELECT o.name "obstacles.name", o.length "obstacles.length", o.height "obstacles.height" 
				FROM obstacles o ORDER BY random()) obstacles`

	if err := p.db.Get(&event, query); err != nil {
		p.log.Panicf(p.logMsg.FormatErr, p.log.CallInfoStr(), fmt.Sprint(p.logMsg.Read, "(generated event)"), err.Error())
	}

	event.TotalSumValues =
		event.Place.DifficultyMovement + event.TimeOfDay.Clarity +
			event.Weather.Clarity + event.Weather.DifficultyMovement + event.Obstacle.Length + event.Obstacle.Height

	return event
}
