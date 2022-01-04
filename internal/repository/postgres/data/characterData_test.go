package data

import (
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
	"time"
)

func TestPostgresCharacterData(t *testing.T) {
	logs := logger.NewLogger()
	logMsg, msgToUser, cfgPostgres := loadEnv(logs)

	db := initPostgresDB(logs, logMsg, cfgPostgres, msgToUser)

	Convey("setup", t, func() {
		var chars []model.Character
		ownerId := createTestUser(db).Id
		chars = append(chars, testCharModel(ownerId, 5))
		chars = append(chars, testCharModel(ownerId, 10))

		So(ownerId, ShouldNotBeZeroValue)

		Convey("createCharacter", func() {

			for i, char := range chars {
				id, err := createCharacter(db, char)
				if err != nil {
					logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
				}
				So(id, ShouldNotBeZeroValue)
				chars[i].CharId = id
			}

			Convey("readAllChar", func() {
				So(db.CharacterData.ReadAll(ownerId), ShouldResemble, chars)

				Convey("updateChar", func() {
					updChar := testCharModel(ownerId, 15)
					updChar.CharId = chars[len(chars)-1].CharId

					if err := db.CharacterData.Update(updChar); err != nil {
						logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
					}
					Convey("readOneChar", func() {
						So(db.ReadOne(ownerId, updChar.CharId), ShouldResemble, updChar)

						Convey("deleteAll", func() {
							if err := db.DeleteAll(ownerId); err != nil {
								logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
							}

							So(db.CharacterData.ReadAll(ownerId), ShouldBeNil)

							if err := db.UserData.DeleteAccount(ownerId); err != nil {
								logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
							}
						})
						Convey("growth/weight out range", func() {
							updChar.Growth = 210
							updChar.Weight = 35

							if err := db.CharacterData.Update(updChar); err != nil {
								switch err.Error() {
								case logMsg.CharGrowthAndWeightOutErr:
									So(err.Error(), ShouldEqual, logMsg.CharGrowthAndWeightOutErr)
								default:
									logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
								}
							}
						})
					})
				})
			})
		})
	})
}

func createCharacter(db *PostgresData, char model.Character) (int, error) {
	charId, err := db.CharacterData.Create(char)
	if err != nil {
		return 0, err
	}
	return charId, nil
}

func testCharModel(ownedId, n int) model.Character {
	return model.Character{
		OwnerId: ownedId,
		Name:    RandStr(4),
		Growth:  RandInt(145, 180) + n,
		Weight:  RandInt(40, 90) + n,
	}
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
