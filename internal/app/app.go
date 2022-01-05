package app

import (
	"encoding/json"
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository"
	"github.com/Aserose/ArchaicReverie/internal/repository/postgres"
	"github.com/Aserose/ArchaicReverie/internal/repository/postgres/data"
	"github.com/Aserose/ArchaicReverie/internal/server"
	"github.com/Aserose/ArchaicReverie/internal/server/handler"
	"github.com/Aserose/ArchaicReverie/internal/service"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

var (
	YmlFilename = "configs/config.yml"
	Ch          = make(chan []byte)
)

func Start(mode int) {
	log := logger.NewLogger()

	logMsg, msgToUser, utilitiesStr, endpoints := config.InitStrSet(YmlFilename, log)

	cfgServer, cfgServices, cfgPostgres, err := config.Init(YmlFilename, log, logMsg)
	if err != nil {
		log.Errorf(logMsg.FormatErr, log.CallInfoStr(), logMsg.InitNoOk, err.Error())
	}

	postgresData := data.NewPostgresData(postgres.Postgres(cfgPostgres, log, logMsg,utilitiesStr.NumberCharacterLimit), msgToUser, log, logMsg,utilitiesStr.NumberCharacterLimit)

	db := repository.NewDB(postgresData)

	services := service.NewService(db, utilitiesStr, cfgServices, msgToUser, log, logMsg)

	handlers := handler.NewHandler(services, utilitiesStr, msgToUser, log, logMsg)

	servers := server.Server{}

	if mode == 1 {
		testApiScheme(endpoints, handlers.Routes(endpoints).Routes())
	}

	if err := servers.Start(cfgServer.Port, handlers.Routes(endpoints), log, logMsg); err != nil {
		log.Errorf(logMsg.FormatErr, log.CallInfoStr(), logMsg.InitNoOk, err.Error())
	}
}

func testApiScheme(endpoints config.Endpoints, info gin.RoutesInfo) {
	var toTest map[string]map[string]string

	inrec, _ := json.Marshal(endpoints)
	json.Unmarshal(inrec, &toTest)

	for i := 0; i <= len(info); i++ {
		for field, val := range toTest {
			for field2, val2 := range val {
				switch len(strings.Split(info[i].Path, "/")) {
				case 3:
					if strings.Contains(val2, strings.Split(info[i].Path, "/")[2]) {
						toTest[field][field2] = fmt.Sprintf(`%s %s/%s%s`, info[i].Method, os.Getenv("APP_URL"),
							strings.Split(info[i].Path, "/")[1], val2)
					}
				case 4:
					if strings.Contains(val2, strings.Split(info[i].Path, "/")[3]) {
						toTest[field][field2] = fmt.Sprintf(`%s %s/%s/%s%s`, info[i].Method, os.Getenv("APP_URL"),
							strings.Split(info[i].Path, "/")[1],
							strings.Split(info[i].Path, "/")[2], val2)
					}
				}

			}
		}
		if i == 13 {
			break
		}
	}
	resultJson, _ := json.Marshal(toTest)
	Ch <- resultJson
	close(Ch)
}
