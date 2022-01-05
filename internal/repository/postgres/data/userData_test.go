package data

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/internal/repository/postgres"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/joho/godotenv"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestPostgresUserData(t *testing.T) {
	logs := logger.NewLogger()
	logMsg, msgToUser, cfgPostgres,utilitiesStr  := loadEnv(logs)

	db := initPostgresDB(logs, logMsg, cfgPostgres, msgToUser, utilitiesStr.NumberCharacterLimit)

	Convey("setup", t, func() {

		testUser := createTestUser(db)
		So(testUser.Id, ShouldNotBeNil)

		Convey("check", func() {
			id, _ := db.UserData.Check(testUser.Username, testUser.Password)
			So(id, ShouldNotBeZeroValue)

			Convey("deleteAccount", func() {
				if err := db.UserData.DeleteAccount(testUser.Id); err != nil {
					logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
				}

				id, _ := db.UserData.Create(testUser.Username, testUser.Password)
				So(id, ShouldNotBeZeroValue)

				if err := db.UserData.DeleteAccount(testUser.Id); err != nil {
					logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
				}
			})

			Convey("create user with same username", func() {
				_, status := db.UserData.Create(testUser.Username, testUser.Password)

				So(status, ShouldEqual, msgToUser.AuthStatus.BusyUsername)
			})
		})
	})

}

func loadEnv(logs logger.Logger) (config.LogMsg, config.MsgToUser, *config.CfgPostgres,config.UtilitiesStr) {
	re := regexp.MustCompile(`^(.*` + "ArchaicReverie" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	godotenv.Load(string(rootPath) + `/.env`)
	logMsg, msgToUser, utilitiesStr, _ := config.InitStrSet(os.Getenv(`CONFIG_FILE`), logs)
	_, _, cfgPostgres, err := config.Init(os.Getenv("CONFIG_FILE"), logs, logMsg)
	if err != nil {
		logs.Errorf(logMsg.FormatErr, logs.CallInfoStr(), logMsg.InitNoOk, err.Error())
	}

	return logMsg, msgToUser, cfgPostgres,utilitiesStr
}

func initPostgresDB(logs logger.Logger, logMsg config.LogMsg,
	cfgPostgres *config.CfgPostgres, msgToUser config.MsgToUser, numberCharLimit int) *PostgresData {
	return NewPostgresData(postgres.Postgres(cfgPostgres, logs, logMsg,numberCharLimit), msgToUser, logs, logMsg,numberCharLimit)
}

func createTestUser(pqDB *PostgresData) model.User {
	var testUser model.User
	for {
		testUser.Username = RandStr(5)
		testUser.Password = RandStr(5)
		testUser.Id, _ = pqDB.UserData.Create(testUser.Username, testUser.Password)
		if testUser.Id != 0 {
			break
		}
	}

	return testUser
}
func RandStr(length int) string {
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
