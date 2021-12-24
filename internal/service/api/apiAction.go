package api

import (
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
)

type ActionScene struct {
	db *repository.DB
}

func NewActionScene(db *repository.DB) *ActionScene {
	return &ActionScene{
		db: db,
	}
}

func (a ActionScene) Jump(character model.Character, jumpPosition model.Jump) string {

	//TODO

	_, sum := a.exampleSceneLocation()

	sumChar := character.Growth-character.Weight
	sumJumpPosition := jumpPosition.RunUp*(jumpPosition.BodyTilt+jumpPosition.ArmAmplitude+jumpPosition.SquatDepth)

	cj := sumChar+sumJumpPosition

	if cj > sum+15 || cj < sum-15{
		return "You fell"
	}

	return "You jumped over "

	//TODO
}


func (a ActionScene) exampleSceneLocation() (location model.Location, exampleSum int) {

	//TODO

	var (Location model.Location
	Night model.TimeOfDay
	 Fog model.Weather
	 Swamp model.Place
	 Puddle model.Obstacle)


	Swamp.DifficultyMovement = -50

	Night.Clarity = 0

	Fog.Clarity = 0
	Fog.DifficultyMovement = 0

	Puddle.Height = -5
	Puddle.Length = 10

	sum := Swamp.DifficultyMovement+Night.Clarity+Fog.Clarity+Fog.DifficultyMovement+Puddle.Height+Puddle.Length

	Location.Obstacle = Puddle
	Location.Weather = Fog
	Location.TimeOfDay = Night
	Location.Place = Swamp

	return Location, sum

	//TODO
}