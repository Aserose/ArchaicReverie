package postgres

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model/scheme"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Postgres(cfgPostgres *config.CfgPostgres, log logger.Logger, logMsg config.LogMsg) *sqlx.DB {

	log.Infof(logMsg.Format, log.PackageAndFileNames(), logMsg.Init)

	db, err := sqlx.Connect(cfgPostgres.DriverName,
		fmt.Sprintf(cfgPostgres.ConnectFormat,
			cfgPostgres.Username,
			cfgPostgres.DBName,
			cfgPostgres.Password,
			cfgPostgres.SSLMode))
	if err != nil {
		log.Panicf(logMsg.FormatErr, log.CallInfoStr(), logMsg.InitNoOk, err.Error())
	}

	createTables(db,
		scheme.SchemaUser,
		scheme.SchemaCharacter,
		scheme.SchemaLocation)

	log.Infof(logMsg.Format, log.PackageAndFileNames(), logMsg.InitOk)

	return db
}

func createTables(db *sqlx.DB, schemes ...string) {
	for _, schemeDB := range schemes {
		db.MustExec(schemeDB)
	}
}
