package action

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/internal/service/action/tasks"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	wr "github.com/mroth/weightedrand"
	"math/rand"
	"time"
)

type Recreational interface {
	GetFoodList() []model.Food
	Eat(character model.Character, order model.Food) (string, model.Character)
}

type Location interface {
	Main(character model.Character, action model.Action) (string, model.Character)
}

type LocationWithEnemy interface {
	Main(character model.Character, action model.Action) (string, model.Character)
}

type ActionScene struct {
	db           *repository.DB
	genConfig    config.GenerationConfig
	utilitiesStr config.UtilitiesStr
	log          logger.Logger
	msgToUser    config.MsgToUser
	Conditions   tasks.Condition
	Location
	LocationWithEnemy
	Recreational
}

func NewActionScene(db *repository.DB, genConfig config.GenerationConfig,
	utilitiesStr config.UtilitiesStr, log logger.Logger, msgToUser config.MsgToUser) *ActionScene {

	cond := tasks.NewCondition()

	return &ActionScene{
		db:                db,
		genConfig:         genConfig,
		utilitiesStr:      utilitiesStr,
		log:               log,
		msgToUser:         msgToUser,
		Conditions:        cond,
		Location:          tasks.NewLocation(db, log, cond, msgToUser, utilitiesStr),
		LocationWithEnemy: tasks.NewLocationWithEnemy(db, cond),
		Recreational:      tasks.NewRecreational(db, log, msgToUser),
	}
}

func (a *ActionScene) GenerateScene() map[string]interface{} {
	var (
		chooser *wr.Chooser
		err     error
	)

	if !a.inActive() {
		var result map[string]interface{}

		if chooser, err = wr.NewChooser(
			wr.Choice{Item: 1, Weight: uint(a.genConfig.GenerationTypeOfTask.Location)},
			wr.Choice{Item: 2, Weight: uint(a.genConfig.GenerationTypeOfTask.LocationWithEnemy)}); err != nil {
			a.log.Errorf("%s:%s", a.log.CallInfoStr(), err.Error())
		}

		pick := chooser.Pick().(int)

		switch pick {
		case 1:
			a.Conditions.SetConditionLocation(a.db.Postgres.EventData.GenerateEventLocation())
			result = map[string]interface{}{"location": a.Conditions.GetConditionLocation()}
		case 2:
			a.Conditions.SetConditionLocation(a.db.Postgres.EventData.GenerateEventLocation())
			a.Conditions.SetConditionEnemy(a.db.Postgres.EventData.GenerateEnemy(a.generateSettingEnemy()))
			result = map[string]interface{}{"location": a.Conditions.GetConditionLocation(), "enemy": a.Conditions.GetConditionEnemy()}
		}
		return result
	}

	return nil
}

func (a ActionScene) inActive() bool {
	if a.Conditions.GetConditionEnemy() != nil ||
		a.Conditions.GetConditionLocation() != (model.Location{}) {
		return true
	}

	return false
}

func (a ActionScene) Action(character model.Character, action model.Action) (string, model.Character) {
	if a.Conditions.GetConditionEnemy() != nil {
		return a.LocationWithEnemy.Main(character, action)
	} else {
		return a.Location.Main(character, action)
	}
}

func (a ActionScene) generateSettingEnemy() []model.Enemy {
	enemies := []model.Enemy{}

	for i := 0; i <= RandInt(a.genConfig.GenerationEnemy.MinQuantityOfEnemies, a.genConfig.GenerationEnemy.MaxClassOfEnemy, i); i++ {
		enemies = append(
			enemies, model.Enemy{
				Class: RandInt(a.genConfig.GenerationEnemy.MinClassOfEnemy,
					a.genConfig.GenerationEnemy.MaxClassOfEnemy, i),
			})
	}

	return enemies
}

func RandInt(min, max, add int) int {
	rand.Seed(time.Now().UnixNano() + int64(add))
	return rand.Intn(max-min) + min
}
