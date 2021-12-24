package app

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/postgres"
	"github.com/Aserose/ArchaicReverie/internal/repository/postgres/data"
	"github.com/Aserose/ArchaicReverie/internal/server"
	"github.com/Aserose/ArchaicReverie/internal/server/handler"
	"github.com/Aserose/ArchaicReverie/internal/service"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
)

const (
	ymlFilename = "configs/config.yml"
)

func Start() {
	log := logger.NewLogger()

	logMsg, msgToUser := config.InitStrSet(ymlFilename, log)

	cfgServer, cfgServices, cfgPostgres, err := config.Init(ymlFilename, log, logMsg)
	if err != nil {
		log.Errorf(logMsg.FormatErr, config.ConfigPackageName, logMsg.InitNoOk, err.Error())
	}

	postgresData := data.NewPostgresData(postgres.Postgres(cfgPostgres, log, logMsg), msgToUser, log, logMsg)

	db := repository.NewDB(postgresData)

	services := service.NewService(db, cfgServices, log, logMsg)

	handlers := handler.NewHandler(services, msgToUser, log, logMsg)
	servers := server.Server{}

	if err := servers.Start(cfgServer.Port, handlers.Routes(), log, logMsg); err != nil {
		log.Errorf(logMsg.FormatErr, server.ServerPackageName, logMsg.InitNoOk, err.Error())
	}

}
