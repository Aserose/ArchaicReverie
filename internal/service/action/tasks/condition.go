package tasks

import (
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
)

type Condition interface {
	Set
	Get
	Reset
}

type Set interface {
	SetConditionLocation(model.Location)
	SetConditionEnemy([]model.Enemy)
}

type Get interface {
	GetConditionLocation() model.Location
	GetConditionEnemy() []model.Enemy
}

type Reset interface {
	ResetLocation()
	ResetEnemy()
}

type Conditions struct {
	enemy         []model.Enemy
	locationScene model.Location
	Set
	Get
	Reset
}

func NewCondition() *Conditions {
	return &Conditions{}
}

func (c *Conditions) SetConditionLocation(set model.Location) {
	c.locationScene = set
}
func (c *Conditions) SetConditionEnemy(set []model.Enemy) {
	c.enemy = set
}

func (c Conditions) GetConditionLocation() model.Location {
	return c.locationScene
}

func (c Conditions) GetConditionEnemy() []model.Enemy {
	return c.enemy
}

func (c *Conditions) ResetLocation() {
	c.locationScene = model.Location{}
}

func (c *Conditions) ResetEnemy() {
	c.enemy = nil
}
