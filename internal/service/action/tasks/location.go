package tasks

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
)

type location struct {
	db           *repository.DB
	Conditions   Condition
	log          logger.Logger
	msgToUser    config.MsgToUser
	utilitiesStr config.UtilitiesStr
}

func NewLocation(db *repository.DB, log logger.Logger, condition Condition, msgToUser config.MsgToUser, utilitiesStr config.UtilitiesStr) *location {
	return &location{
		db:           db,
		Conditions:   condition,
		log:          log,
		msgToUser:    msgToUser,
		utilitiesStr: utilitiesStr,
	}
}

func (a location) Main(character model.Character, action model.Action) (string, model.Character) {
	return a.jump(character, action.Jump)
}

func (a location) jump(character model.Character, jumpPosition model.Jump) (string, model.Character) {
	actionResult := a.db.Postgres.EventData.GetActionResult(model.ActionResult{Name: "fall"})

	if character.RemainHealth < 9 {
		if character.RemainHealth < 0 {
			character.RemainHealth = 0
		}
		return a.msgToUser.ActionMsg.LowHP, character
	}

	if validateActionJumpPosition(
		jumpPosition.RunUp,
		jumpPosition.BodyTilt,
		jumpPosition.ArmAmplitude,
		jumpPosition.SquatDepth) == false {
		return a.utilitiesStr.BadRequest, character
	}

	calcGrowth, calcWeight := convertCharParams(character)

	sumJumpPosition :=
		jumpPosition.RunUp +
			jumpPosition.BodyTilt +
			jumpPosition.ArmAmplitude +
			jumpPosition.SquatDepth +
			calcGrowth +
			calcWeight +
			randomOutcomeVariable()

	result := a.Conditions.GetConditionLocation().TotalSumValues + sumJumpPosition

	a.Conditions.ResetLocation()

	if 1 >= result && result >= -1 {
		return a.msgToUser.ActionMsg.JumpOver, character
	} else {
		return fmt.Sprintf("%s %s", a.msgToUser.ActionMsg.JumpFell, fmt.Sprintf(a.msgToUser.ActionMsg.RemainHealth, character.RemainHealth)),
			damage(character, actionResult.DamageHP)
	}
}
