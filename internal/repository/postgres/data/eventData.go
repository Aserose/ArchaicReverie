package data

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/internal/repository/postgres/data/eventData"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type Enemy interface {
	GenerateEnemy(settingEnemy []model.Enemy) []model.Enemy
}

type Location interface {
	GenerateEventLocation() model.Location
}

type Item interface {
	GetWeaponAll() []model.Weapon
	GetWeapon(name string) model.Weapon
	GetFoodAll() []model.Food
	GetFood(name string) model.Food
}

type PostgresEventData struct {
	db     *sqlx.DB
	log    logger.Logger
	logMsg config.LogMsg
	Enemy
	Location
	Item
}

func NewEventData(db *sqlx.DB, log logger.Logger, logMsg config.LogMsg, cfg config.CharacterConfig) *PostgresEventData {
	return &PostgresEventData{
		db:       db,
		log:      log,
		logMsg:   logMsg,
		Enemy:    eventData.NewEnemy(db, cfg, log),
		Location: eventData.NewLocation(db, log),
		Item:     eventData.NewItem(db, log),
	}
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
