package postgres

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model/scheme"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const PostgresPackageName = "postgres"

func Postgres(cfgPostgres *config.CfgPostgres, log logger.Logger, logMsg config.LogMsg) *sqlx.DB {
	log.Infof(logMsg.Format, PostgresPackageName, logMsg.Init)

	db, err := sqlx.Connect(cfgPostgres.DriverName,
		fmt.Sprintf(cfgPostgres.ConnectFormat,
			cfgPostgres.Username,
			cfgPostgres.DBName,
			cfgPostgres.Password,
			cfgPostgres.SSLMode))
	if err != nil {
		log.Panicf(logMsg.FormatErr, PostgresPackageName, logMsg.InitNoOk, err.Error())
	}

	db.MustExec(scheme.SchemaUser)
	db.MustExec(scheme.SchemaCharacter)

	log.Infof(logMsg.Format, PostgresPackageName, logMsg.InitOk)

	return db
}
