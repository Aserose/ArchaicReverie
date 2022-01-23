package eventData

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"strings"
	"time"
)

type enemy struct {
	db  *sqlx.DB
	cfg config.CharacterConfig
	log logger.Logger
}

func NewEnemy(db *sqlx.DB, cfg config.CharacterConfig, log logger.Logger) *enemy {
	return &enemy{
		db:  db,
		cfg: cfg,
		log: log,
	}
}

func (e enemy) GenerateEnemy(settingEnemy []model.Enemy) []model.Enemy {
	enemies := []model.Enemy{}

	// class 1 enemy: {probability of weapon class 1, weapon class 2, weapon class 3},
	// class 2 enemy: {weapon class 1, weapon class 2, weapon class 3},
	//...
	classWeapProb := map[int][]int{
		1: {90, 70, 0},
		2: {70, 30, 85},
		3: {15, 80, 90},
	}

	for _, enemy := range settingEnemy {
		row := e.db.QueryRowx(`SELECT * FROM
								(SELECT a.name "name" FROM enemy a WHERE class=$1) enemy
								    NATURAL FULL JOIN
								(SELECT b.name "weapon.name", b.weapon_class "weapon.weapon_class", b.sharp "weapon.sharp", b.weight "weapon.weight" FROM weapon b WHERE weapon_class=$2 
									ORDER BY random()) weapon`,
			enemy.Class, weaponSpawnProbability(
				classWeapProb[enemy.Class][0],
				classWeapProb[enemy.Class][1],
				classWeapProb[enemy.Class][2]))

		if err := row.StructScan(&enemy); err != nil {
			if strings.Contains(err.Error(), "converting NULL to string") {
				continue
			} else {
				e.log.Errorf("%s: %s", e.log.CallInfoStr(), err.Error())
			}
		}

		if enemy.Weapons == (model.Weapon{}) {
			enemy.Weapons.Name = "fists"
		}

		enemies = append(enemies, generateEnemyParams(enemy, e.cfg))

	}

	return enemies
}

func generateEnemyParams(enemy model.Enemy, cfg config.CharacterConfig) model.Enemy {
	enemy.Growth = randInt(cfg.Restriction.MinCharGrowth, cfg.Restriction.MaxCharGrowth)
	enemy.Weight = randInt(cfg.Restriction.MinCharWeight, cfg.Restriction.MaxCharWeight) + enemy.Weapons.Weight
	enemy.CoinAmount = generateEnemyCoins(enemy.Class)
	enemy.RemainHealth = calculateThresholdHP(enemy.Weight, cfg)
	enemy.RemainEnergy = calculateThresholdMP(enemy.Weight, cfg)

	return enemy
}

func calculateThresholdMP(weight int, cfg config.CharacterConfig) int {
	mp := cfg.Conversion.BaseMP
	for i := cfg.Restriction.MinCharWeight; i <= weight; i += cfg.Conversion.AddMP {
		mp -= cfg.Conversion.AddMP
		if mp == 10 {
			break
		}
	}
	return mp
}

func calculateThresholdHP(weight int, cfg config.CharacterConfig) int {
	hp := cfg.Conversion.BaseHP
	for i := cfg.Restriction.MinCharWeight; i <= weight; i += cfg.Conversion.AddHP {
		hp += cfg.Conversion.AddHP
	}
	return hp
}

func generateEnemyCoins(class int) int {
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
