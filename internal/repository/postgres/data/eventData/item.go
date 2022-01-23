package eventData

import (
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type item struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewItem(db *sqlx.DB, log logger.Logger) *item {
	return &item{
		db:  db,
		log: log,
	}
}

func (i item) GetWeaponAll() []model.Weapon {
	var weapons []model.Weapon

	if err := i.db.Select(&weapons, "SELECT name, price, sharp, weight FROM weapon"); err != nil {
		i.log.Errorf("%s: %s", i.log.CallInfoStr(), err.Error())
	}

	return weapons
}

func (i item) GetWeapon(name string) model.Weapon {
	var weapon model.Weapon

	if err := i.db.Get(&weapon, `SELECT name, price, sharp, weight FROM weapon WHERE name=$1`, name); err != nil {
		i.log.Panicf("%s: %s", i.log.CallInfoStr(), err.Error())
	}

	return weapon
}

func (i item) GetFoodAll() []model.Food {
	var foods []model.Food

	if err := i.db.Select(&foods, "SELECT * FROM foods"); err != nil {
		i.log.Errorf("%s: %s", i.log.CallInfoStr(), err.Error())
	}

	return foods
}

func (i item) GetFood(name string) model.Food {
	var food model.Food

	if err := i.db.Get(&food, `SELECT * FROM foods WHERE name=$1`, name); err != nil {
		i.log.Panicf("%s: %s", i.log.CallInfoStr(), err.Error())
	}

	return food
}
