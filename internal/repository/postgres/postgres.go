package postgres

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model/scheme"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strings"
)

func Postgres(cfgPostgres *config.CfgPostgres, log logger.Logger, logMsg config.LogMsg, charConfig config.CharacterConfig) *sqlx.DB {

	log.Infof(logMsg.Format, log.PackageAndFileNames(), logMsg.Init)

	db, err := sqlx.Connect(cfgPostgres.DriverName,
		fmt.Sprintf(cfgPostgres.ConnectFormat,
			cfgPostgres.Username,
			cfgPostgres.DBName,
			cfgPostgres.Password,
			cfgPostgres.SSLMode))
	if err != nil {
		log.Panicf(logMsg.Format, log.CallInfoStr(), err.Error())
	}

	createTables(db, log,
		scheme.CreateSchemaUser(charConfig.NumberCharLimit+1),
		scheme.CreateSchemaCharacter(charConfig),
		scheme.SchemaLocation,
		scheme.SchemaFood)

	log.Infof(logMsg.Format, log.PackageAndFileNames(), logMsg.InitOk)

	return db
}

func createTables(db *sqlx.DB, log logger.Logger, schemes ...string) {
	for _, schemeDB := range schemes {
		log.Infof("%s : %s", log.CallInfoStr(),
			strings.Replace(strings.Split(schemeDB, `(`)[0], ` IF NOT EXISTS`, ``, 1))
		db.MustExec(schemeDB)
	}
}
