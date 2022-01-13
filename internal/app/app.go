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
	YmlFilename           = "configs/config.yml"
	YmlGenerationFilename = "configs/generationConfig.yml"
	Ch                    = make(chan []byte)
)

func Start(mode int) {
	log := logger.NewLogger()

	logMsg, msgToUser, utilitiesStr, endpoints, charConfig := config.InitStrSet(YmlFilename, log)

	cfgServer, cfgServices, cfgPostgres, err := config.Init(YmlFilename, log, logMsg)
	if err != nil {
		log.Errorf(logMsg.Format, log.CallInfoStr(), err.Error())
	}

	cfgGeneration := config.InitGenerationConfig(YmlGenerationFilename, log)

	postgresData := data.NewPostgresData(postgres.Postgres(cfgPostgres, log, logMsg, charConfig), msgToUser, log, logMsg, charConfig)

	db := repository.NewDB(postgresData)

	services := service.NewService(db, utilitiesStr, cfgServices, msgToUser, log, logMsg, charConfig, cfgGeneration)

	handlers := handler.NewHandler(services, utilitiesStr, msgToUser, log, logMsg)

	servers := server.Server{}

	if mode == 1 {
		getApiScheme(endpoints, handlers.Routes(endpoints).Routes(), log)
	}

	if err := servers.Start(cfgServer.Port, handlers.Routes(endpoints), log, logMsg); err != nil {
		log.Errorf(logMsg.Format, log.CallInfoStr(), err.Error())
	}
}

func getApiScheme(endpoints config.Endpoints, info gin.RoutesInfo, log logger.Logger) {
	var toTest map[string]map[string]string

	inrec, err := json.Marshal(endpoints)
	if err != nil {
		log.Errorf("%s %s", log.CallInfoStr(), err.Error())
	}
	if err := json.Unmarshal(inrec, &toTest); err != nil {
		log.Errorf("%s %s", log.CallInfoStr(), err.Error())
	}

	for i := 0; i <= len(info)-1; i++ {
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
	}
	resultJson, _ := json.Marshal(toTest)

	Ch <- resultJson
	close(Ch)
}
