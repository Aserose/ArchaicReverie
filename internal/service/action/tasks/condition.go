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

type conditions struct {
	enemy         []model.Enemy
	locationScene model.Location
	Set
	Get
	Reset
}

func NewCondition() *conditions {
	return &conditions{}
}

func (c *conditions) SetConditionLocation(set model.Location) {
	c.locationScene = set
}
func (c *conditions) SetConditionEnemy(set []model.Enemy) {
	c.enemy = set
}

func (c conditions) GetConditionLocation() model.Location {
	return c.locationScene
}

func (c conditions) GetConditionEnemy() []model.Enemy {
	return c.enemy
}

func (c *conditions) ResetLocation() {
	c.locationScene = model.Location{}
}

func (c *conditions) ResetEnemy() {
	c.enemy = nil
}
