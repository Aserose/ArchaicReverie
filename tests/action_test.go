package tests

import (
	"encoding/json"
	"github.com/Aserose/ArchaicReverie/internal/app"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	cv "github.com/smartystreets/goconvey/convey"
	"net/http"
	"strings"
	"testing"
)

func TestAction(t *testing.T) {
	logs := logger.NewLogger()
	_, msgToUser, charConfig := loadEnv(logs)
	go app.Start(1)
	apiScheme := loadApiScheme(logs)

	cv.Convey("setup", t, func() {
		var (
			temp            templates
			numberCharLimit = charConfig.NumberCharLimit
			client          = http.Client{}
			cookie          = temp.authorization.createCharAndSelectChar(client, apiScheme, logs, numberCharLimit)
			resp            *http.Response
			testUser        = generateTestUser()
		)

		cv.Convey("beginActionScene", func() {
			resp, _ = temp.doRequest(
				client,
				strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[0],
				strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[1],
				reqBody(logs, testUser),
				cookie)

			cv.So(string(temp.readRespBody(resp)), cv.ShouldNotBeEmpty)

			var result string

			for {
				resp, updCookie := temp.doRequest(
					client,
					strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[0],
					strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[1],
					reqBody(logs, generateAction()),
					cookie)

				if len(updCookie) > 0 {
					if updCookie[0].Value != "" {
						cookie = updCookie
					}
				}

				result = string(temp.readRespBody(resp))
				logs.Print(result)

				if result == msgToUser.ActionMsg.JumpOver {
					cv.So(result, cv.ShouldEqual, msgToUser.ActionMsg.JumpOver)
					break
				}
				if result == msgToUser.ActionMsg.LowHP {
					cv.So(result,cv.ShouldEqual,"low health points")
					break
				}
				cv.So(result, cv.ShouldNotBeEmpty)
			}

			cv.Convey("getInfoAboutChar", func() {
				resp, _ = temp.doRequest(
					client,
					strings.Split(apiScheme.ActionEndpoints.InfoAboutSelectedChar, " ")[0],
					strings.Split(apiScheme.ActionEndpoints.InfoAboutSelectedChar, " ")[1],
					reqBody(logs, testUser),
					cookie)

				selectedChar := temp.unmarshalChar(temp.readRespBody(resp), logs)
				logs.Printf("%d health points left out of %d", selectedChar.RemainHealth, selectedChar.ThresholdHealth)

				if selectedChar.RemainHealth < selectedChar.ThresholdHealth {
					resp, _ = temp.doRequest(
						client,
						strings.Split(apiScheme.ActionEndpoints.BeginRepast, " ")[0],
						strings.Split(apiScheme.ActionEndpoints.BeginRepast, " ")[1],
						reqBody(logs, testUser),
						cookie)

					selectedFood := temp.unmarshalFood(temp.readRespBody(resp), logs)

					logs.Printf("restore health by %d points with %s", selectedFood.RestoreHp, selectedFood.Name)

					resp, cookie = temp.doRequest(
						client,
						strings.Split(apiScheme.ActionEndpoints.BeginRepast, " ")[0],
						strings.Split(apiScheme.ActionEndpoints.BeginRepast, " ")[1],
						reqBody(logs, selectedFood),
						cookie)

					resp, _ = temp.doRequest(
						client,
						strings.Split(apiScheme.ActionEndpoints.InfoAboutSelectedChar, " ")[0],
						strings.Split(apiScheme.ActionEndpoints.InfoAboutSelectedChar, " ")[1],
						reqBody(logs, testUser),
						cookie)

					selectedChar := temp.unmarshalChar(temp.readRespBody(resp), logs)
					logs.Printf("%d health points left out of %d", selectedChar.RemainHealth, selectedChar.ThresholdHealth)
				}

				cv.Convey("delete", func() {
					resp, cookie = temp.doRequest(
						client,
						strings.Split(apiScheme.AuthEndpoints.DeleteAccount, " ")[0],
						strings.Split(apiScheme.AuthEndpoints.DeleteAccount, " ")[1],
						reqBody(logs, testUser),
						cookie)

					resp, cookie = temp.doRequest(
						client,
						strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[0],
						strings.Split(apiScheme.AuthEndpoints.SignIn, " ")[1],
						reqBody(logs, testUser),
						cookie)

					cv.So(string(temp.readRespBody(resp)), cv.ShouldEqual, msgToUser.AuthStatus.InvalidUsername)
				})
			})
		})
	})
}

func (r readAndRequest) unmarshalFood(data []byte, log logger.Logger) model.Food {
	var (
		foods []model.Food
		food  model.Food
	)

	if err := json.Unmarshal(data, &foods); err != nil {
		log.Errorf("%s:%s", log.CallInfoStr(), err.Error())
	}

	for _, foo := range foods {
		if food.Price == 0 || food.Price > foo.Price {
			food = foo
		}
	}
	return food
}

func (a authorization) createCharAndSelectChar(client http.Client, apiScheme config.Endpoints, logs logger.Logger, numberCharLimit int) []*http.Cookie {
	cookie, chars := a.authorizeUserWithChars(client, apiScheme, logs, numberCharLimit)
	var temp templates
	_, cookie = temp.doRequest(
		client,
		strings.Split(apiScheme.ApiEndpoints.SelectChar, " ")[0],
		strings.Split(apiScheme.ApiEndpoints.SelectChar, " ")[1],
		reqBody(logs, chars[0]),
		cookie)

	return cookie
}

func (a authorization) authorizeUserWithChars(client http.Client, apiScheme config.Endpoints,
	logs logger.Logger, numberCharLimit int) ([]*http.Cookie, []model.Character) {
	testUser, cookie := a.authorizeUser(client, apiScheme, logs)

	var (
		chars []model.Character
		resp  *http.Response
		temp  templates
	)

	for i := 0; i < numberCharLimit; i++ {
		chars = append(chars, generateChar(testUser.Id))
		resp, _ = temp.doRequest(
			client,
			strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[0],
			strings.Split(apiScheme.ApiEndpoints.CreateChar, " ")[1],
			reqBody(logs, chars[i]),
			cookie)

		chars[i].CharId = temp.unmarshalChar(temp.readRespBody(resp), logs).CharId
	}
	return cookie, chars
}