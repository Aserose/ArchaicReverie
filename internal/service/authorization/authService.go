package authorization

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
)

type serviceAuthorization struct {
	db          *repository.DB
	cfgServices *config.CfgServices
	log         logger.Logger
	logMsg      config.LogMsg
	msgToUser   config.MsgToUser
}

func NewServiceAuthorization(db *repository.DB, cfgServices *config.CfgServices,
	log logger.Logger, logMsg config.LogMsg, msgToUser config.MsgToUser) *serviceAuthorization {
	return &serviceAuthorization{
		db:          db,
		cfgServices: cfgServices,
		log:         log,
		logMsg:      logMsg,
		msgToUser:   msgToUser,
	}
}
