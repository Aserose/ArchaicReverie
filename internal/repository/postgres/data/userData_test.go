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
	logMsg, msgToUser, cfgPostgres, charConfig := loadEnv(logs)

	db := initPostgresDB(logs, logMsg, cfgPostgres, msgToUser, charConfig)

	Convey("setup", t, func() {

		testUser := createTestUser(db)
		So(testUser.Id, ShouldNotBeNil)

		Convey("check", func() {
			id, _ := db.UserData.Check(testUser.Username, testUser.Password)
			So(id, ShouldNotBeZeroValue)

			Convey("create user with same username", func() {
				_, status := db.UserData.Create(testUser.Username, testUser.Password)

				So(status, ShouldEqual, msgToUser.AuthStatus.BusyUsername)
			})

			Convey("deleteAccount", func() {
				if err := db.UserData.DeleteAccount(testUser.Id, testUser.Password); err != nil {
					logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
				}

				id, _ := db.UserData.Create(testUser.Username, testUser.Password)
				So(id, ShouldNotBeZeroValue)

				if err := db.UserData.DeleteAccount(testUser.Id, testUser.Password); err != nil {
					logs.Panicf(logMsg.Format, logs.CallInfoStr(), err.Error())
				}
			})
		})
	})
}

func loadEnv(logs logger.Logger) (config.LogMsg, config.MsgToUser, *config.CfgPostgres, config.CharacterConfig) {
	re := regexp.MustCompile(`^(.*` + "ArchaicReverie" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	godotenv.Load(string(rootPath) + `/.env`)

	configs := config.NewConfig(os.Getenv("CONFIG_GAME"),
		os.Getenv("CONFIG_INFRASTRUCTURE"), os.Getenv("CONFIG_MSG"), logs)

	charConfig := configs.GameConfigs.InitCharConfig()
	logMsg, msgToUser, _ := configs.MsgConfigs.InitMsg()
	_, _, cfgPostgres, _, err := configs.InfrastructureConfigs.InitInfrastructureConfigs(logMsg)
	if err != nil {
		logs.Errorf(logMsg.Format, logs.CallInfoStr(), err.Error())
	}

	return logMsg, msgToUser, cfgPostgres, charConfig
}

func initPostgresDB(logs logger.Logger, logMsg config.LogMsg,
	cfgPostgres *config.CfgPostgres, msgToUser config.MsgToUser, charConfig config.CharacterConfig) *PostgresData {
	return NewPostgresData(postgres.Postgres(cfgPostgres, logs, logMsg, charConfig), msgToUser, logs, logMsg, charConfig)
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
