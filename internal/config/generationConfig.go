package config

import "github.com/Aserose/ArchaicReverie/pkg/logger"

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

func InitGenerationConfig(filename string, log logger.Logger) GenerationConfig {
	var genConfig GenerationConfig

	unmarshalYaml(filename, log, &genConfig)

	return genConfig
}
