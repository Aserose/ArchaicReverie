package tests

import (
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
	_, msgToUser,utilitiesStr := loadEnv(logs)
	go app.Start(1)
	apiScheme := loadApiScheme(logs)

	cv.Convey("setup", t, func() {
		var (
			client   = http.Client{}
			cookie   []*http.Cookie
			resp     *http.Response
			numberCharLimit = utilitiesStr.NumberCharacterLimit
			testUser = generateTestUser()
			temp     templates
		)
		cookie = temp.authorization.selectUserChar(client,apiScheme,logs,numberCharLimit)

		cv.Convey("BeginActionScene", func() {
			resp, _ = temp.doRequest(
				client,
				strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[0],
				strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[1],
				reqBody(logs, testUser),
				cookie)

			cv.So(string(temp.readRespBody(resp)), cv.ShouldNotBeEmpty)

			var result string
			for {
				resp, _ = temp.doRequest(
					client,
					strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[0],
					strings.Split(apiScheme.ActionEndpoints.BeginActionScene, " ")[1],
					reqBody(logs, generateAction()),
					cookie)
				result = string(temp.readRespBody(resp))
				if result == msgToUser.ActionMsg.JumpOver {
					cv.So(result, cv.ShouldEqual, msgToUser.ActionMsg.JumpOver)
					break
				}
				cv.So(result, cv.ShouldNotBeEmpty)
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
}

func (a authorization) selectUserChar(client http.Client, apiScheme config.Endpoints, logs logger.Logger,numberCharLimit int) []*http.Cookie {
	cookie,chars := a.authorizeUserWithChars(client,apiScheme,logs,numberCharLimit)
	var temp templates

	_, cookie = temp.doRequest(
		client,
		strings.Split(apiScheme.ApiEndpoints.SelectChar, " ")[0],
		strings.Split(apiScheme.ApiEndpoints.SelectChar, " ")[1],
		reqBody(logs, chars[0]),
		cookie)

	return cookie
}

func (a authorization) authorizeUserWithChars(client http.Client, apiScheme config.Endpoints, logs logger.Logger,numberCharLimit int) ([]*http.Cookie,[]model.Character){ //TODO edit func name
	testUser, cookie := a.authorizeUser(client, apiScheme, logs)

	var (
		chars []model.Character
		resp *http.Response
		temp templates
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
	return cookie,chars
}