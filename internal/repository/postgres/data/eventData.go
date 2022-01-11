package data

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"time"
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

func (p PostgresEventData) GenerateEventLocation() model.Location {
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
		p.log.Panicf(p.logMsg.Format, p.log.CallInfoStr(), err.Error())
	}

	event.TotalSumValues =
		event.Place.DifficultyMovement + event.TimeOfDay.Clarity +
			event.Weather.Clarity + event.Weather.DifficultyMovement + event.Obstacle.Length + event.Obstacle.Height

	return event
}

func (p PostgresEventData) GetActionResult(actionResult model.ActionResult) model.ActionResult {

	row := p.db.QueryRowx(`SELECT a.name, d.damage_hp "damage_type.damage_hp", d.damage_mp "damage_type.damage_mp"
							FROM action_result a
							JOIN damage_type d ON (a.name = d.name)
							WHERE a.name = $1`, actionResult.Name)

	if err := row.StructScan(&actionResult); err != nil {
		p.log.Errorf(p.logMsg.Format, p.log.CallInfoStr(), err.Error())
	}

	return actionResult
}

func (p PostgresEventData) GetListFood() []model.Food {
	var food []model.Food

	if err := p.db.Select(&food, "SELECT * FROM foods"); err != nil {
		p.log.Errorf(p.logMsg.Format, p.log.CallInfoStr(), err.Error())
	}

	return food
}

func (p PostgresEventData) GenerateEnemy(settingEnemy []model.Enemy) []model.Enemy {
	enemies := []model.Enemy{}
	classWeapProb := map[int][]int{
		1: {90, 70, 0},
		2: {70, 30, 85},
		3: {15, 80, 90},
	}

	for _, enemy := range settingEnemy {
		row := p.db.QueryRowx(`SELECT * FROM
								(SELECT a.name "name" FROM enemy a WHERE class=$1) enemy
								    NATURAL FULL JOIN
								(SELECT b.name "weapon.name",b.sharp "weapon.sharp", b.weight "weapon.weight" FROM weapon b WHERE weapon_class=$2 
									ORDER BY random()) weapon`,
			enemy.Class, weaponSpawnProbability(
				classWeapProb[enemy.Class][0],
				classWeapProb[enemy.Class][1],
				classWeapProb[enemy.Class][2]))

		if err := row.StructScan(&enemy); err != nil {
			p.log.Errorf(p.logMsg.Format, p.log.CallInfoStr(), err.Error())
		}
		enemy.CoinAmount = generateCoins(enemy.Class)
		enemies = append(enemies, enemy)
	}

	return enemies
}

func (p PostgresEventData) GetFood(name string) model.Food {
	var food model.Food

	if err := p.db.Get(&food, `SELECT * FROM foods WHERE name=$1`, name); err != nil {
		p.log.Panicf(p.logMsg.Format, p.log.CallInfoStr(), err.Error())
	}

	return food
}

func generateCoins(class int) int {
	switch class {
	case 3:
		return randInt(0, 10)
	case 2:
		return randInt(5, 15)
	case 1:
		return randInt(10, 20)
	}
	return 0
}

func weaponSpawnProbability(a, b, c int) int {
	weaponSpawnProbability := randInt(0, 100)

	if weaponSpawnProbability > a {
		return 3
	} else if weaponSpawnProbability > b {
		return 2
	} else if weaponSpawnProbability > c {
		return 1
	}
	return 0
}

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
