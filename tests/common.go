package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/app"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/joho/godotenv"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
)

type auth interface {
	authorizeUser(client http.Client, apiScheme config.Endpoints) (model.User, []*http.Cookie)
	authorizeUserWithChars(client http.Client, apiScheme config.Endpoints, numberCharLimit int) ([]*http.Cookie, []model.Character)
	authorizeUserWithCharsAndSelect(client http.Client, apiScheme config.Endpoints, numberCharLimit int) []*http.Cookie
}

type readAndReq interface {
	doRequest(client http.Client, method string, url string, body io.Reader, cookie []*http.Cookie) (*http.Response, []*http.Cookie)
	readRespBody(resp *http.Response) []byte
	unmarshalInt(data []byte) int
	unmarshalChar(data []byte) model.Character
	unmarshalChars(data []byte) []model.Character
	unmarshalAvailableItems(data []byte) model.Items
	cheapestOrder(availableItems model.Items) model.Items
}

type templates struct {
	auth
	readAndReq
}

func NewTemplates(log logger.Logger) *templates {
	return &templates{
		auth:       NewAuthorization(log),
		readAndReq: NewReadAndRequest(log),
	}
}

func NewAuthorization(log logger.Logger) *authorization {
	return &authorization{
		log: log,
	}
}

func NewReadAndRequest(log logger.Logger) *readAndRequest {
	return &readAndRequest{
		log: log,
	}
}

func loadApiScheme(log logger.Logger) config.Endpoints {
	var apiScheme config.Endpoints
	resultJson := <-app.Ch
	if err := json.Unmarshal(resultJson, &apiScheme); err != nil {
		log.Errorf("%s: %s", log.CallInfoStr(), err.Error())
	}
	return apiScheme
}

func loadEnv(log logger.Logger) (config.LogMsg, config.MsgToUser, config.CharacterConfig) {
	re := regexp.MustCompile(`^(.*` + "ArchaicReverie" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	godotenv.Load(string(rootPath) + `/.env`)

	app.YmlInfrastructureFilename = os.Getenv("CONFIG_INFRASTRUCTURE")
	app.YmlGameConfig = os.Getenv("CONFIG_GAME")
	app.YmlMsgFilename = os.Getenv("CONFIG_MSG")

	configs := config.NewConfig(os.Getenv("CONFIG_GAME"), os.Getenv("CONFIG_INFRASTRUCTURE"), os.Getenv("CONFIG_MSG"), log)

	logMsg, msgToUser, _ := configs.MsgConfigs.InitMsg()
	charConfig := configs.GameConfigs.InitCharConfig()

	return logMsg, msgToUser, charConfig
}

func reqBody(log logger.Logger, v interface{}) *bytes.Buffer {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		log.Panicf("%s: %s", log.CallInfoStr(), err.Error())
	}
	return bytes.NewBuffer(jsonStr)
}

func generateTestUser() model.User {
	return model.User{
		Username: randStr(7),
		Password: randStr(7),
	}
}

func generateChar(ownedId int) model.Character {
	return model.Character{
		OwnerId: ownedId,
		Name:    randStr(4),
		Growth:  randInt(150, 180),
		Weight:  randInt(40, 90),
	}
}

func generateAction() model.Action {
	return model.Action{
		InAction: "jump",
		Jump: model.Jump{
			SquatDepth:   randIntWithExceptions(-1, 2, 0),
			ArmAmplitude: randIntWithExceptions(-1, 2, 0),
			BodyTilt:     randIntWithExceptions(-1, 2, 0),
			RunUp:        randIntWithExceptions(-1, 2, 0),
		},
	}
}

func randIntWithExceptions(min, max int, exc ...int) int {
	result := randInt(min, max)

	if len(exc) == 1 {
		for result == exc[0] {
			result = randInt(min, max)
		}
	} else {
		for i := 0; i < len(exc); i++ {
			if result == exc[i] {
				result = randInt(min, max)
				i = 0
			}
		}
	}

	return result
}

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func randStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
