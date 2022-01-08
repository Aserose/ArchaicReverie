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

type templates struct {
	authorization
	readAndRequest
}

type authorization struct {
}

type readAndRequest struct {
}

func (readAndRequest) doRequest(client http.Client, method string, url string, body io.Reader, cookie []*http.Cookie) (*http.Response, []*http.Cookie) {
	request, _ := http.NewRequest(method, url, body)
	for i := range cookie {
		request.AddCookie(cookie[i])
	}
	resp, _ := client.Do(request)
	return resp, resp.Cookies()
}

func (readAndRequest) readRespBody(resp *http.Response) []byte {
	bodyBytes, _ := io.ReadAll(resp.Body)
	return bodyBytes
}

func (readAndRequest) unmarshalInt(data []byte, log logger.Logger) int {
	var beInt int
	if err := json.Unmarshal(data, &beInt); err != nil {
		log.Errorf("%s: %s: %s", log.CallInfoStr(), err.Error(), string(data))
	}

	return beInt
}

func (readAndRequest) unmarshalChar(data []byte, log logger.Logger) model.Character {
	var char model.Character

	if err := json.Unmarshal(data, &char); err != nil {
		log.Errorf("%s: %s: %s", log.CallInfoStr(), err.Error(), string(data))
	}

	return char
}

func (readAndRequest) unmarshalChars(data []byte, log logger.Logger) []model.Character {
	var chars []model.Character
	if err := json.Unmarshal(data, &chars); err != nil {
		log.Errorf("%s: %s: %s", log.CallInfoStr(), err.Error(), string(data))
	}

	return chars
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
	app.YmlFilename = os.Getenv("CONFIG_FILE")
	logMsg, msgToUser, _, _, charConfig := config.InitStrSet(app.YmlFilename, log)

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
		Username: RandStr(7),
		Password: RandStr(7),
	}
}

func generateChar(ownedId int) model.Character {
	return model.Character{
		OwnerId: ownedId,
		Name:    RandStr(4),
		Growth:  RandInt(150, 180),
		Weight:  RandInt(40, 90),
	}
}

func generateAction() model.Action {
	return model.Action{
		InAction: "jump",
		Jump: model.Jump{
			SquatDepth:   RandIntWithExceptions(-1, 2, 0),
			ArmAmplitude: RandIntWithExceptions(-1, 2, 0),
			BodyTilt:     RandIntWithExceptions(-1, 2, 0),
			RunUp:        RandIntWithExceptions(-1, 2, 0),
		},
	}
}

func RandIntWithExceptions(min,max int,exc ... int) int{
	rand.Seed(time.Now().UnixNano())
	result := RandInt(min,max)

	if len(exc)== 1 {
		for result == exc[0]{
			result = RandInt(min,max)
		}
	} else {
		for i:=0;i<len(exc);i++ {
			if result == exc[i] {
				result = RandInt(min, max)
				i = 0
			}
		}
	}

	return result
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func RandStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
