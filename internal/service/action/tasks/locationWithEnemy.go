package tasks

import (
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
)

type LocationWithEnemy struct {
	db         *repository.DB
	Conditions Condition
}

func NewLocationWithEnemy(db *repository.DB, conditions Condition) *LocationWithEnemy {
	return &LocationWithEnemy{
		db:         db,
		Conditions: conditions,
	}
}

func (l LocationWithEnemy) Main(character model.Character, action model.Action) (string, model.Character) {
	switch action.InAction {
	case "hit":
		return l.hit(character, action.Hit)
	case "run":
		return l.run(character, action.Run)
	}
	return "invalid command", character
}

func (l LocationWithEnemy) hit(character model.Character, action model.Hit) (string, model.Character) {

	return "", character
}

func (l LocationWithEnemy) run(character model.Character, action model.Run) (string, model.Character) {

	return "", character
}
