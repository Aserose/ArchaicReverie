package data

import (
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPostgresEventData(t *testing.T) {
	logs := logger.NewLogger()
	logMsg, msgToUser, cfgPostgres, charConfig := loadEnv(logs)

	db := initPostgresDB(logs, logMsg, cfgPostgres, msgToUser, charConfig)

	Convey("generateLocation", t, func() {
		So(db.EventData.GenerateEventLocation(), ShouldNotBeEmpty)

		Convey("getFood", func() {
			So(db.EventData.GetFood("apple"), ShouldNotBeEmpty)

			Convey("getResult", func() {
				So(db.EventData.GetActionResult(model.ActionResult{Name: "fall"}), ShouldNotBeEmpty)
				logs.Print(db.EventData.GetActionResult(model.ActionResult{Name: "fall"}))

				Convey("generateEnemy", func() {
					So(db.EventData.GenerateEnemy([]model.Enemy{{Class: 1}, {Class: 2}, {Class: 3}}), ShouldNotBeEmpty)
					logs.Print(db.EventData.GenerateEnemy([]model.Enemy{{Class: 1}, {Class: 2}, {Class: 3}}))

					Convey("getWeapon", func() {
						So(db.EventData.GetWeapon("knife"), ShouldNotBeEmpty)
						So(db.EventData.GetWeaponAll(), ShouldNotResemble, []model.Weapon{})
					})
				})
			})
		})
	})
}
