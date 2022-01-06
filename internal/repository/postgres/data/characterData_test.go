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
	logMsg, msgToUser, cfgPostgres, charConfig := loadEnv(logs)

	db := initPostgresDB(logs, logMsg, cfgPostgres, msgToUser, charConfig)

	Convey("setup", t, func() {
		var chars []model.Character
		testUser := createTestUser(db)
		chars = append(chars, testCharModel(testUser.Id, 5))
		chars = append(chars, testCharModel(testUser.Id, 10))

		So(testUser.Id, ShouldNotBeZeroValue)

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
				So(db.CharacterData.ReadAll(testUser.Id), ShouldResemble, chars)

				Convey("updateChar", func() {
					updChar := testCharModel(testUser.Id, 15)
					updChar.CharId = chars[len(chars)-1].CharId

					if err := db.CharacterData.Update(updChar); err != nil {
						logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
					}
					Convey("readOneChar", func() {
						So(db.ReadOne(testUser.Id, updChar.CharId), ShouldResemble, updChar)

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

						Convey("deleteAll", func() {
							if err := db.DeleteAll(testUser.Id); err != nil {
								logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
							}

							So(db.CharacterData.ReadAll(testUser.Id), ShouldBeNil)

							if err := db.UserData.DeleteAccount(testUser.Id, testUser.Password); err != nil {
								logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
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
