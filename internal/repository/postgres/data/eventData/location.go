package eventData

import (
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type location struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewLocation(db *sqlx.DB, log logger.Logger) *location {
	return &location{
		db:  db,
		log: log,
	}
}

func (l location) GenerateEventLocation() model.Location {
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

	if err := l.db.Get(&event, query); err != nil {
		l.log.Panicf("%s: %s", l.log.CallInfoStr(), err.Error())
	}

	event.TotalSumValues =
		event.Place.DifficultyMovement + event.TimeOfDay.Clarity +
			event.Weather.Clarity + event.Weather.DifficultyMovement + event.Obstacle.Length + event.Obstacle.Height

	return event
}
