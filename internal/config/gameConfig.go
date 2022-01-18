package config

import "github.com/Aserose/ArchaicReverie/pkg/logger"

type GameConfig struct {
	gameFilename string
	log          logger.Logger
}

type CharacterConfig struct {
	Conversion struct {
		BaseHP         int `yaml:"baseHP"`
		AddHP          int `yaml:"addHP"`
		BaseMP         int `yaml:"baseMP"`
		AddMP          int `yaml:"addMP"`
		BaseAmountCoin int `yaml:"baseAmountCoin"`
	} `yaml:"conversion"`

	Restriction struct {
		NumberCharLimit     int `yaml:"numberCharLimit"`
		MinCharWeight       int `yaml:"minCharWeight"`
		MaxCharWeight       int `yaml:"maxCharWeight"`
		MinCharGrowth       int `yaml:"minCharGrowth"`
		MaxCharGrowth       int `yaml:"maxCharGrowth"`
		WeaponInventorySize int `yaml:"weaponInventorySize"`
	} `yaml:"restriction"`
}

type GenerationConfig struct {
	GenerationEnemy struct {
		MinQuantityOfEnemies int `yaml:"minQuantityOfEnemies"`
		MaxQuantityOfEnemies int `yaml:"maxQuantityOfEnemies"`
		MinClassOfEnemy      int `yaml:"minClassOfEnemy"`
		MaxClassOfEnemy      int `yaml:"maxClassOfEnemy"`
	} `yaml:"generationEnemy"`
	GenerationTypeOfTask struct {
		Location          int `yaml:"location"`
		LocationWithEnemy int `yaml:"locationWithEnemy"`
	} `yaml:"typeOfTask"`
}

func NewGameConfig(gameFilename string, log logger.Logger) *GameConfig {
	return &GameConfig{
		gameFilename: gameFilename,
		log:          log,
	}
}

func (r GameConfig) InitGenerationConfig() GenerationConfig {
	var genConfig GenerationConfig

	unmarshalYaml(r.gameFilename, r.log, &genConfig)

	return genConfig
}

func (r GameConfig) InitCharConfig() CharacterConfig {
	var charConfig CharacterConfig

	unmarshalYaml(r.gameFilename, r.log, &charConfig)

	return charConfig
}
