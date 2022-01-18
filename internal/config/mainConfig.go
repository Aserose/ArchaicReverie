package config

import (
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type GameConfigs interface {
	InitCharConfig() CharacterConfig
	InitGenerationConfig() GenerationConfig
}

type MsgConfigs interface {
	InitMsg() (LogMsg, MsgToUser, UtilitiesStr)
}

type InfrastructureConfigs interface {
	InitInfrastructureConfigs(logMsg LogMsg) (*CfgServer, *CfgServices,
		*CfgPostgres, Endpoints, error)
}

type Config struct {
	GameConfigs
	MsgConfigs
	InfrastructureConfigs
}

func NewConfig(gameFilename, infrastructureFilename, msgFilename string, log logger.Logger) *Config {
	return &Config{
		GameConfigs:           NewGameConfig(gameFilename, log),
		MsgConfigs:            NewMsgConfig(msgFilename, log),
		InfrastructureConfigs: NewInfrastructureConfig(infrastructureFilename, log),
	}
}

func unmarshalYaml(filename string, log logger.Logger, outs ...interface{}) {
	ymlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("yamlFile.Get err   #%v ", err)
	}

	for _, out := range outs {
		err = yaml.Unmarshal(ymlFile, out)
		if err != nil {
			log.Errorf("error: %v", err)
		}
	}
}

func readEnv(log logger.Logger, logMsg LogMsg, cfgs ...interface{}) {
	for _, cfg := range cfgs {
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			log.Errorf(logMsg.Format, log.PackageAndFileNames(), err.Error())
		}
	}
}
