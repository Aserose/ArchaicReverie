package data

import (
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPostgresEventData(t *testing.T) {
	logs := logger.NewLogger()
	logMsg, msgToUser, cfgPostgres := loadEnv(logs)

	db := initPostgresDB(logs, logMsg, cfgPostgres, msgToUser)

	Convey("generateLocation", t, func() {
		So(db.EventData.GenerateEventLocation(), ShouldNotBeEmpty)
	})
}
